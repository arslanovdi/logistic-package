package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/arslanovdi/logistic-package/events/internal/config"
	"github.com/arslanovdi/logistic-package/events/internal/consumer"
	"github.com/arslanovdi/logistic-package/events/internal/process"
	"github.com/arslanovdi/logistic-package/pkg/logger"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/arslanovdi/logistic-package/pkg/server"
	"github.com/arslanovdi/logistic-package/pkg/tracer"
	"log/slog"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	starttimeout = 5 * time.Second
)

type PackageConsumer interface {
	Run(topic string, handler func(key string, msg model.PackageEvent, offset int64)) error
	Close()
}

func main() {
	logger.InitializeLogger(slog.LevelDebug)

	flag.Parse()
	version := flag.String("version", "dev", "Defines the version of the service")
	commitHash := flag.String("commitHash", "-", "Defines the commit hash of the service")
	configFile := flag.String("configFile", "events/config_local.yml", "Defines the config file of the service")

	if err := config.ReadConfigYML(*configFile); err != nil {
		slog.Warn("Failed init configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}
	cfg := config.GetConfigInstance()

	log := slog.With("func", cfg.Project.Instance+".main")

	cfg.Project.Version = *version
	cfg.Project.CommitHash = *commitHash

	if cfg.Project.Debug {
		logger.SetLogLevel(slog.LevelDebug)
	} else {
		logger.SetLogLevel(slog.LevelInfo)
	}

	startCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.Project.StartupTimeout)) // контекст запуска приложения
	defer cancel()
	go func() {
		<-startCtx.Done()
		if errors.Is(startCtx.Err(), context.DeadlineExceeded) { // приложение зависло при запуске
			log.Warn("Application startup time exceeded")
			os.Exit(1)
		}
	}()

	log.Info(fmt.Sprintf("Starting service %s", cfg.Project.Name),
		slog.String("version", cfg.Project.Version),
		slog.String("commitHash", cfg.Project.CommitHash),
		slog.Bool("debug", cfg.Project.Debug),
		slog.String("environment", cfg.Project.Environment),
		slog.String("instance", cfg.Project.Instance),
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ctxTrace, cancelTrace := context.WithTimeout(context.Background(), starttimeout)
	defer cancelTrace()

	trace, err := tracer.NewTracer(
		ctxTrace,
		cfg.Project.Instance,
		cfg.Jaeger.Host+cfg.Jaeger.Port)
	if err != nil {
		log.Warn("Failed to init tracer", slog.String("error", err.Error()))
		os.Exit(1)
	}

	isReady := &atomic.Value{}
	isReady.Store(false)

	statusServer := server.NewStatusServer(
		isReady,
		server.StatusConfig{
			Host:          cfg.Status.Host,
			Port:          cfg.Status.Port,
			LivenessPath:  cfg.Status.LivenessPath,
			ReadinessPath: cfg.Status.ReadinessPath,
			VersionPath:   cfg.Status.VersionPath,
		},
		server.ProjectInfo{
			Name:        cfg.Project.Name,
			Debug:       cfg.Project.Debug,
			Environment: cfg.Project.Environment,
			Version:     cfg.Project.Version,
			CommitHash:  cfg.Project.CommitHash,
			Instance:    cfg.Project.Instance,
		},
	)
	statusServer.Start()

	mcfg := server.MetricsConfig{
		Host: cfg.Metrics.Host,
		Port: cfg.Metrics.Port,
		Path: cfg.Metrics.Path,
	}
	metricsServer := server.NewMetricsServer(mcfg)
	metricsServer.Start()

	kafka, err := consumer.NewKafkaConsumer()
	if err != nil {
		log.Warn("Failed to init kafka consumer", slog.String("error", err.Error()))
		os.Exit(1)
	}

	go func() {
		err = kafka.Run(cfg.Kafka.Topic, process.PrintPackageEvent)
		if err != nil {
			log.Error("kafka consumer error", slog.String("error", err.Error()))
		}
	}()

	isReady.Store(true)

	cancel() // отменяем контекст запуска приложения
	<-stop
	log.Info("Graceful shutdown")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.Project.ShutdownTimeout))
	defer cancelShutdown()
	go func() {
		<-ctxShutdown.Done()
		log.Warn("Application shutdown time exceeded")
		os.Exit(1)
	}()

	isReady.Store(false)

	kafka.Close()

	if err := trace.Shutdown(ctxShutdown); err != nil {
		log.Error("Error shutting down tracer provider", slog.String("error", err.Error()))
	}

	statusServer.Stop(ctxShutdown)

	metricsServer.Stop(ctxShutdown)

	log.Info("Application stopped")
}
