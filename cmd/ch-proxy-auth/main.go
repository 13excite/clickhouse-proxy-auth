package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
	flag.StringVar(&configPath, "config", "", "path to the cfg file")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	cfg := config.Config{}
	cfg.Defaults()

	cfg.ReadConfigFile(configPath)

	if err := config.InitLogger(&cfg); err != nil {
		log.Fatalf("could not configure logger: %v", err)
	}
	logger := zap.S().With("package", "cmd")

	// create a new prom metric server with default registry
	promRegistry := metrics.BuildRegistry()
	metricsServer := metrics.NewMetricServer(cfg.MetricsPort, promRegistry)

	// create the mux
	mux := router.New(cfg, promRegistry)

	// create a http server
	server := http.Server{

		Addr:              fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort),
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
	logger.Infow("starting server", "host", cfg.ServerHost, "port", cfg.ServerPort)

	err := group.Wait()
	if err != nil {
		logger.Errorf("waitgroup returned an error: %w", err)
	}
}
