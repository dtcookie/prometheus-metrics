package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// Config contains configuration settings
type Config struct {
	Dynatrace  Dynatrace    `json:"dynatrace,omitempty"`
	Prometheus Prometheus   `json:"prometheus,omitempty"`
	Metrics    MetricFilter `json:"metrics,omitempty"`
}

func (c *Config) checkValid() error {
	if err := c.Dynatrace.checkValid(); err != nil {
		return err
	}
	if err := c.Prometheus.checkValid(); err != nil {
		return err
	}
	return nil
}

// Dynatrace contains configuration settings for Dynatrace
type Dynatrace struct {
	EnvironmentID string `json:"environment,omitempty"`
	BaseURL       string `json:"baseURL,omitempty"`
	Token         string `json:"token,omitempty"`
}

func (d *Dynatrace) checkValid() error {
	if empty(d.EnvironmentID) {
		return errors.New("no environment specified")
	}
	if empty(d.Token) {
		return errors.New("no API Token specified")
	}
	if empty(d.BaseURL) {
		return nil
	}
	if !strings.HasPrefix(d.BaseURL, "http://") && !strings.HasPrefix(d.BaseURL, "https://") {
		return fmt.Errorf("'%s' doesn't look like a valid Base URL for a Dynatrace Server", d.BaseURL)
	}
	return nil
}

// Prometheus contains configuration settings for Prometheus
type Prometheus string

func (p Prometheus) checkValid() error {
	if empty(p.Host()) {
		return fmt.Errorf("'%s' doesn't look like a valid host name of a Prometheus server", p)
	}
	if p.Port() == 0 {
		return fmt.Errorf("'%s' doesn't look like a valid host name of a Prometheus server", p)
	}
	return nil
}

// Host extracts the host part of the Prometheus setting (host:port)
func (p Prometheus) Host() string {
	s := strings.TrimSpace(string(p))
	idx := strings.Index(s, ":")
	if idx == -1 {
		return s
	}
	result := strings.TrimSpace(s[0:idx])
	if strings.HasPrefix(result, "http://") {
		return result[len("http://"):]
	}
	if strings.HasPrefix(result, "https://") {
		return result[len("https://"):]
	}
	return result
}

// Port extracts the port part of the Prometheus setting (host:port)
func (p Prometheus) Port() int {
	var port int
	s := strings.TrimSpace(string(p))
	idx := strings.Index(s, ":")
	if idx == -1 {
		return 9090
	}
	port, err := strconv.Atoi(s[idx+1:])
	if err != nil {
		return 0
	}
	return port
}

func readConfig() Config {
	var err error
	var settingsFileName string
	var argConfig Config

	baseURL := flag.String("baseurl", "", "The base url of your Managed Dynatrace Server")
	environmentID := flag.String("environment", "", "the environment id of your Dynatrace Tenant")
	token := flag.String("token", "", "an API Token to access the REST API of your Dynatrace Tenant")
	prometheus := flag.String("prometheus", "", "Host and Port of your Prometheus Server")
	pConfigFileName := flag.String("config", "", "a JSON file containing settings")

	flag.Parse()

	setFlag(baseURL, &argConfig.Dynatrace.BaseURL)
	setFlag(environmentID, &argConfig.Dynatrace.EnvironmentID)
	setFlag(token, &argConfig.Dynatrace.Token)
	if !pempty(prometheus) {
		argConfig.Prometheus = Prometheus(*prometheus)
	}
	setFlag(pConfigFileName, &settingsFileName)

	if !empty(settingsFileName) {
		if !fileExists(settingsFileName) {
			fmt.Println(fmt.Sprintf("config file '%s' doesn't exist", settingsFileName))
			os.Exit(1)
		}
	} else {
		settingsFileName = "settings.json"
	}

	var config Config

	if fileExists(settingsFileName) {
		var settingsFile *os.File
		if settingsFile, err = os.Open(settingsFileName); err != nil {
			log.Fatal(fmt.Sprintf("Unable to read settings file '%s'", settingsFileName), err)
			os.Exit(1)
		}
		defer settingsFile.Close()

		var data []byte
		if data, err = ioutil.ReadAll(settingsFile); err != nil {
			log.Fatal(fmt.Sprintf("Unable to read settings file '%s'", settingsFileName), err)
			os.Exit(1)
		}
		if err = json.Unmarshal(data, &config); err != nil {
			log.Fatal(fmt.Sprintf("Invalid settings file '%s'", settingsFileName), err)
			os.Exit(1)
		}
	}

	if len(argConfig.Prometheus) > 0 {
		config.Prometheus = argConfig.Prometheus
	}
	if !empty(argConfig.Dynatrace.BaseURL) {
		config.Dynatrace.BaseURL = argConfig.Dynatrace.BaseURL
	}
	if !empty(argConfig.Dynatrace.EnvironmentID) {
		config.Dynatrace.EnvironmentID = argConfig.Dynatrace.EnvironmentID
	}
	if !empty(argConfig.Dynatrace.Token) {
		config.Dynatrace.Token = argConfig.Dynatrace.Token
	}

	if err = config.checkValid(); err != nil {
		flag.Usage()
		fmt.Println()
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return config

}
