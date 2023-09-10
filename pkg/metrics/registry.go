package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// BuildRegistry builds a new prometheus.Registry
func BuildRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}
