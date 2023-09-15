package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/13excite/clickhouse-proxy-auth/pkg/config"
	"github.com/13excite/clickhouse-proxy-auth/pkg/metrics"
	"github.com/13excite/clickhouse-proxy-auth/pkg/router"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "path to the cfgig file")
}

func main() {
	ctx := context.Background()

	cfg := config.Config{}
	cfg.Defaults()

	logger := zap.S().With("package", "cmd")

	// create a new prom metric server with default registry
	metricsServer := metrics.NewMetricServer(cfg.MetricsPort, nil)

	// create the mux
	mux := router.New(cfg)

	// create a http server
	server := http.Server{
		Addr:              cfg.ServerHost + ":" + cfg.ServerPort,
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	// create the Group
	group, _ := errgroup.WithContext(ctx)
	// run HTTP server
	group.Go(
		server.ListenAndServe,
	)
	// run prometheus server
	group.Go(
		metricsServer.ListenAndServe,
	)

	err := group.Wait()
	if err != nil {
		logger.Errorf("waitgroup returned an error: %w", err)
	}
}
