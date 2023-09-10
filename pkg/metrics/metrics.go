// Package metrics provides http prometheus metrics server
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func buildHTTPRequestCounterCollector() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ch_proxy_auth_http_requests_total",
		Help: "Count of HTTP requests",
	}, []string{"method", "path", "statuscode"})
}
