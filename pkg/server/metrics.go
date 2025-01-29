// Package server - http серверы
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsServer - http сервер для метрик
type MetricsServer struct {
	server *http.Server
	config *MetricsConfig
}

// MetricsConfig - конфигурация http сервера
type MetricsConfig struct {
	Host string
	Port int
	Path string
}

// NewMetricsServer returns http server for metrics
func NewMetricsServer(cfg *MetricsConfig) *MetricsServer {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	mux := http.DefaultServeMux
	mux.Handle(cfg.Path, promhttp.Handler())

	metrics := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: ReadHeaderTimeout,
	}

	return &MetricsServer{
		server: metrics,
		config: cfg,
	}
}

// Start - запуск http сервера
func (s *MetricsServer) Start() {
	log := slog.With("func", "MetricsServer.Start")

	metricsAddr := fmt.Sprintf("%s:%v", s.config.Host, s.config.Port)

	go func() {
		log.Info("Metrics server is running", slog.String("address", metricsAddr))
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed running metric server", slog.String("error", err.Error()))

			os.Exit(1) // приложение завершается с ошибкой, при ошибке запуска сервера обрабатывающего запросы Prometheus
		}
	}()
}

// Stop - остановка http сервера
func (s *MetricsServer) Stop(ctx context.Context) {
	log := slog.With("func", "MetricsServer.Stop")

	if err := s.server.Shutdown(ctx); err != nil {
		log.Error("MetricsServer.Shutdown", slog.String("error", err.Error()))
	} else {
		log.Info("MetricsServer shut down correctly")
	}
}
