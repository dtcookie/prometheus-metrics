package dynatrace

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dtcookie/prometheus-metrics/clients/rest"
)

// Client has no documentation
type Client interface {
	GetHosts() (Hosts, error)        // GetHosts queries for currently known hosts from the Dynatrace environment
	IngestMetrics(data []byte) error // IngestMetrics sends the given payload as POST body to the Dynatrace environments Metric Ingest API
}

// NewClient creates a new API client targeted to a Dynatrace environment
// If baseURL is empty it is assumed that it is a SaaS environment (<environmentid>.live.dynatrace.com)
// The specified token is being used for authentication
func NewClient(environmentID string, baseURL string, token string) Client {
	if len(baseURL) == 0 {
		return &client{Endpoint: &saasEndpoint{EnvironmentID: environmentID}, Client: rest.NewClient().SetHeader("Authorization", "Api-Token "+token)}
	}
	return &client{Endpoint: &managedEndpoint{EnvironmentID: environmentID, BaseURL: baseURL}, Client: rest.NewClient().SetHeader("Authorization", "Api-Token "+token)}
}

type client struct {
	Endpoint endpoint
	Client   rest.Client
}

// GetHosts queries for currently known hosts from the Dynatrace environment
func (sc *client) GetHosts() (Hosts, error) {
	var hosts []Host

	err := sc.Client.Get(sc.Endpoint.Format("/api/v1/entity/infrastructure/hosts"), func(status int, message string, body []byte) error {
		switch status {
		case 200:
			return json.Unmarshal(body, &hosts)
		default:
			var errResponse errorResponse
			if err := json.Unmarshal(body, &errResponse); err != nil {
				message = strings.TrimSpace(message)
				if len(message) == 0 {
					message = "no error message available"
				}
				return fmt.Errorf("[%d] %s", status, message)
			}
			return errResponse
		}

	})
	return hosts, err
}

// IngestMetrics sends the given payload as POST body to the Dynatrace environments Metric Ingest API
func (sc *client) IngestMetrics(data []byte) error {
	sc.Client.Trace(true)
	return sc.Client.Post(sc.Endpoint.Format("/api/v2/metrics/ingest"), data, func(status int, message string, body []byte) error {
		var iRes ingestResponse
		switch status {
		case 202:
			return json.Unmarshal(body, &iRes)
		case 400:
			if err := json.Unmarshal(body, &iRes); err != nil {
				return err
			}
			return iRes
		default:
			var errResponse errorResponse
			if err := json.Unmarshal(body, &errResponse); err != nil {
				message = strings.TrimSpace(message)
				if len(message) == 0 {
					message = "no error message available"
				}
				return fmt.Errorf("[%d] %s", status, message)
			}
			return errResponse
		}
	})
}

// https://siz65484.live.dynatrace.com/api/v2/metrics/ingest
