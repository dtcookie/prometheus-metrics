package prometheus

import (
	"encoding/json"
	"fmt"

	"github.com/dtcookie/prometheus-metrics/clients/rest"
)

// Client has no documentation
type Client interface {
	GetMetrics() ([]string, error)
	GetMeasurements(metric string) (*MetricResultData, error)
}

// NewClient has no documentation
func NewClient(host string, port int) Client {
	return &client{Host: host, Port: port, Client: rest.NewClient()}
}

type client struct {
	Host   string
	Port   int
	Client rest.Client
}

// GetMetrics queries Prometheus for all available metric names
func (pc *client) GetMetrics() ([]string, error) {
	var valuesResult ValuesResult

	err := pc.Client.Get(fmt.Sprintf("http://%s:%d/api/v1/label/__name__/values", pc.Host, pc.Port), func(status int, message string, body []byte) error {
		switch status {
		case 200:
			return json.Unmarshal(body, &valuesResult)
		default:
			return fmt.Errorf("[%d] %s", status, message)
		}

	})
	return valuesResult.Data, err
}

// GetMeasurements queries for available measurements for the given metric name
func (pc *client) GetMeasurements(metric string) (*MetricResultData, error) {
	var metricResult MetricResult
	err := pc.Client.Get(fmt.Sprintf("http://%s:%d/api/v1/query?query=%s", pc.Host, pc.Port, metric), func(status int, message string, body []byte) error {
		switch status {
		case 200:
			return json.Unmarshal(body, &metricResult)
		default:
			return fmt.Errorf("[%d] %s", status, message)
		}

	})
	return &metricResult.Data, err
}
