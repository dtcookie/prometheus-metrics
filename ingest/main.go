package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dtcookie/prometheus-metrics/clients/dynatrace"
	"github.com/dtcookie/prometheus-metrics/clients/prometheus"
)

func main() {
	var err error
	var logFile *os.File

	exeName := os.Args[0]
	if strings.HasSuffix(exeName, ".exe") {
		exeName = exeName[:len(exeName)-len(".exe")]
	}

	logFile, err = os.OpenFile(exeName+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()
	wrt := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(wrt)

	config := readConfig()

	api := NewAPIAccess(config)

	for {
		ingestLoop(config, api)
		time.Sleep(50 * time.Second)
	}

}

func ingestLoop(config Config, api *APIAccess) {
	var err error
	var hosts dynatrace.Hosts

	if hosts, err = api.GetHosts(); err != nil {
		log.Fatal(err)
		return
	}

	var metrics []string

	if metrics, err = api.GetMetrics(); err != nil {
		log.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	for _, metricName := range metrics {
		if !config.Metrics.Matches(metricName) {
			continue
		}

		if err = handleMetric(api, hosts, buf, metricName); err != nil {
			log.Println(err)
		}
	}

	data := buf.Bytes()
	if len(data) > 0 {
		if err := api.IngestMetrics(data); err != nil {
			log.Println(err)
		}
	}
}

func handleMetric(api *APIAccess, hosts dynatrace.Hosts, buf *bytes.Buffer, metricName string) error {
	var err error
	var data *prometheus.MetricResultData
	if data, err = api.GetMeasurements(metricName); err != nil {
		return err
	}
	if data.ResultType != "vector" {
		return nil
	}
	for _, record := range data.Result {
		dimensions := record.Dimensions
		name := dimensions["__name__"]
		vInstance := dimensions["instance"]
		instance := instanceOf(vInstance)
		if instance == nil {
			continue
		}
		host := hosts.Match(instance.Host)
		if host == nil {
			continue
		}
		output := fmt.Sprintf("%s,dt.entity.host=%s", name, host.EntityID)
		for k, v := range dimensions {
			if k == "job" {
				continue
			}
			if k == "__name__" {
				continue
			}
			if k == "instance" {
				continue
			}
			output = output + "," + uscore(k) + "=" + quote(v)
		}
		output = fmt.Sprintf("%s %f", output, record.Value.Value)

		fmt.Fprintln(buf, output)
	}
	return nil
}
