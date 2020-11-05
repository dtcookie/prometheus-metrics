package dynatrace

import (
	"strings"
)

// Host has no documentation
type Host struct {
	EntityID       string `json:"entityId"`
	DisplayName    string `json:"displayName"`
	DiscoveredName string `json:"discoveredName"`
}

func unFqnd(hostname string) string {
	idx := strings.Index(hostname, ".")
	if idx == -1 {
		return hostname
	}
	return hostname[0:idx]
}

// Matches has no documentation
func (h *Host) Matches(hostname string) bool {
	if h.DiscoveredName == hostname {
		return true
	}
	if unFqnd(hostname) == unFqnd(h.DiscoveredName) {
		return true
	}
	return false
}

// Hosts has no documentation
type Hosts []Host

// Match has no documentation
func (h Hosts) Match(hostname string) *Host {
	for _, host := range h {
		if host.Matches(hostname) {
			return &host
		}
	}
	return nil
}
