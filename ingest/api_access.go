package main

import (
	"github.com/dtcookie/prometheus-metrics/clients/dynatrace"
	"github.com/dtcookie/prometheus-metrics/clients/prometheus"
)

// APIAccess abstracts away the two REST Clients to be used
type APIAccess struct {
	Dynatrace  dynatrace.Client
	Prometheus prometheus.Client
}

// NewAPIAccess initializes the REST clients of a new API Access
func NewAPIAccess(config Config) *APIAccess {
	return &APIAccess{
		Dynatrace:  dynatrace.NewClient(config.Dynatrace.EnvironmentID, config.Dynatrace.BaseURL, config.Dynatrace.Token),
		Prometheus: prometheus.NewClient(config.Prometheus.Host(), config.Prometheus.Port()),
	}
}

// GetHosts queries for currently known hosts from the Dynatrace environment
func (aa *APIAccess) GetHosts() (dynatrace.Hosts, error) {
	return aa.Dynatrace.GetHosts()
}

// GetMetrics queries Prometheus for all available metric names
func (aa *APIAccess) GetMetrics() ([]string, error) {
	return aa.Prometheus.GetMetrics()
}

// GetMeasurements queries for available measurements for the given metric name
func (aa *APIAccess) GetMeasurements(metric string) (*prometheus.MetricResultData, error) {
	return aa.Prometheus.GetMeasurements(metric)
}

// IngestMetrics sends the given payload as POST body to the Dynatrace environments Metric Ingest API
func (aa *APIAccess) IngestMetrics(data []byte) error {
	return aa.Dynatrace.IngestMetrics(data)
}
