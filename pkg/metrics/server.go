package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewMetricServer builds srv for prometheus metrics
func NewMetricServer(port int, registry *prometheus.Registry) *http.Server {
	mux := chi.NewRouter()
	var handler http.Handler
	if registry != nil {
		handler = promhttp.InstrumentMetricHandler(
			registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
		)
	} else {
		handler = promhttp.Handler()
	}

	mux.Handle("/metrics", handler)

	metricsServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	return metricsServer
}
