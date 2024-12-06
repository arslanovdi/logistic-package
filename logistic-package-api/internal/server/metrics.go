package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// GRPCNotFoundCounter - счетчик не найденных запросов
var GRPCNotFoundCounter = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: "logistic",
	Subsystem: "package_api",
	Name:      "grpc_not_found",
	Help:      "Total gRPC not found calls",
})

// CRUDCounter - счетчик CRUD запросов
var CRUDCounter = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: "logistic",
	Subsystem: "package_api",
	Name:      "crud",
	Help:      "Total CRUD calls",
})

// GRPC2 - гистограмма времени выполнения gRPC запросов
var GRPC2 = promauto.NewHistogram(prometheus.HistogramOpts{
	Namespace: "logistic",
	Subsystem: "package_api",
	Name:      "grpc2",
	Help:      "grpc2 calls",
}, // []string{"method"}, // метка для метрики, для каждой метки будет свой график
)

// RetranslatorEvents - счетчик событий которые сейчас отправляются в кафку
var RetranslatorEvents = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "logistic",
	Subsystem: "package_api",
	Name:      "retranslator",
	Help:      "Retranslator events in work",
})

// MetricsServer - http сервер для метрик
type MetricsServer struct {
	server *http.Server
}

// NewMetricsServer returns http server for metrics
func NewMetricsServer() *MetricsServer {

	cfg := config.GetConfigInstance()

	addr := fmt.Sprintf("%s:%d", cfg.Metrics.Host, cfg.Metrics.Port)

	mux := http.DefaultServeMux
	mux.Handle(cfg.Metrics.Path, promhttp.Handler())

	metrics := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: time.Second * 5,
	}

	return &MetricsServer{
		server: metrics,
	}
}

// Start - запуск http сервера
func (s *MetricsServer) Start() { //(cancelFunc context.CancelFunc) {

	log := slog.With("func", "MetricsServer.Start")

	cfg := config.GetConfigInstance()

	metricsAddr := fmt.Sprintf("%s:%v", cfg.Metrics.Host, cfg.Metrics.Port)

	go func() {
		log.Info("Metrics server is running", slog.String("address", metricsAddr))
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed running metrics server", slog.String("error", err.Error()))

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
