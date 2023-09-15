package router

import (
	"github.com/go-chi/chi/middleware"

	"github.com/13excite/clickhouse-proxy-auth/pkg/api"
	"github.com/13excite/clickhouse-proxy-auth/pkg/config"

	"github.com/go-chi/chi"
)

// New creates a new router
func New(conf config.Config) *chi.Mux {
	mux := chi.NewRouter()
	// health endpoint
	mux.Use(middleware.Heartbeat("/healthz"))
	// inject request-id
	mux.Use(middleware.RequestID)

	mux.Group(func(groupRouter chi.Router) {
		groupRouter.Use(middleware.NoCache)
		groupRouter.Use(loggerHTTPMiddlewareDefault(conf.IgnorePaths))

		groupRouter.Mount("/auth", api.NewHandler(conf.NetAclClusters, conf.HostToCluster))

	})
	return mux

}
