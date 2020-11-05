package main

import "strings"

// MetricFilter includes or excludes metrics provided by Prometheus
type MetricFilter struct {
	Includes []string `json:"includes,omitempty"`
	Excludes []string `json:"excludes,omitempty"`
}

// NewMetricFilter initializes a new Metric Filter
func NewMetricFilter() *MetricFilter {
	return &MetricFilter{Includes: []string{}, Excludes: []string{}}
}

// Matches returns true if the given metric name should be included based on the filter configuration, false otherwise
func (mf *MetricFilter) Matches(metricName string) bool {
	isIncluded := false
	for _, include := range mf.Includes {
		if strings.HasPrefix(include, "*") {
			if strings.HasSuffix(include, "*") {
				if strings.Contains(metricName, include[1:len(include)-2]) {
					isIncluded = true
					break
				}
			} else {
				if strings.HasSuffix(metricName, include[1:len(include)-1]) {
					isIncluded = true
					break
				}
			}
		} else if strings.HasSuffix(include, "*") {
			if strings.HasPrefix(metricName, include[0:len(include)-1]) {
				isIncluded = true
				break
			}
		}
	}
	if !isIncluded {
		return false
	}
	for _, exclude := range mf.Excludes {
		if strings.HasPrefix(exclude, "*") {
			if strings.HasSuffix(exclude, "*") {
				if strings.Contains(metricName, exclude[1:len(exclude)-2]) {
					return false
				}
			} else {
				if strings.HasSuffix(metricName, exclude[1:len(exclude)-1]) {
					return false
				}
			}
		} else if strings.HasSuffix(exclude, "*") {
			if strings.HasPrefix(metricName, exclude[0:len(exclude)-1]) {
				return false
			}
		}
	}
	return true
}
