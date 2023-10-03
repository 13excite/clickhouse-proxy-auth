// Package router assembles the router from packages under pkg/api
package router

import (
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/13excite/clickhouse-proxy-auth/pkg/api"
	"github.com/13excite/clickhouse-proxy-auth/pkg/config"
	"github.com/13excite/clickhouse-proxy-auth/pkg/metrics"

	"github.com/go-chi/chi"
)

// New creates a new router
func New(conf config.Config, registry *prometheus.Registry) *chi.Mux {
	mux := chi.NewRouter()
	// health endpoint
	mux.Use(middleware.Heartbeat("/healthz"))
	// inject request-id
	mux.Use(middleware.RequestID)

	mux.Group(func(groupRouter chi.Router) {
		groupRouter.Use(middleware.NoCache)
		groupRouter.Use(loggerHTTPMiddlewareDefault(conf.IgnorePaths))
		if registry != nil {
			groupRouter.Use(metrics.BuildRequestMiddleware(registry))
		}
		groupRouter.Mount("/", api.NewHandler(conf.NetACLClusters, conf.HostToCluster))
	})
	return mux
}
