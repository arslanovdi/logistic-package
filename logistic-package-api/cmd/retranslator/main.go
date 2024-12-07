// Сервис для пересылки событий из базы данных в кафку (outbox pattern)
package main

import (
	"context"
	"errors"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/database"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/database/postgres"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/retranslator"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/sender"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/server"
	"github.com/arslanovdi/logistic-package/pkg/logger"
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

	//configFile = "logistic-package-api/config_local.yml"
	configFile = "logistic-package-api/config.yml"
)

func main() {
	logger.InitializeLogger(slog.LevelDebug)

	log := slog.With("func", "Retranslator.main")

	if err := config.ReadConfigYML(configFile); err != nil {
		log.Warn("Failed init configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}
	cfg := config.GetConfigInstance()

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

	ctxTrace, cancelTrace := context.WithTimeout(context.Background(), starttimeout)
	defer cancelTrace()

	trace, err1 := tracer.NewTracer(ctxTrace, cfg.Project.Name+" "+"Retranslator", cfg.Jaeger.Host+cfg.Jaeger.Port)
	if err1 != nil {
		log.Warn("Failed to init tracer", slog.String("error", err1.Error()))
		os.Exit(1)
	}

	isReady := &atomic.Value{}
	isReady.Store(false)
	statusServer := server.NewStatusServer(isReady)
	statusServer.Start()

	go func() { // TODO отсечка статус сервера
		time.Sleep(starttimeout)
		isReady.Store(true)
		log.Info("The service is ready to accept requests")
	}()

	/*metricsServer := server.NewMetricsServer()	TODO metrics
	metricsServer.Start()*/

	dbpool := database.MustGetPgxPool(context.Background())

	repo := postgres.NewPostgresRepo(dbpool)

	kafka := sender.MustNewKafkaSender()

	OutboxRetranslator := retranslator.NewRetranslator(repo, kafka)
	OutboxRetranslator.Start(cfg.Kafka.Topic)

	cancel() // отменяем контекст запуска приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

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

	OutboxRetranslator.Stop()

	//metricsServer.Stop(ctxShutdown)

	err := kafka.Close()
	if err != nil {
		log.Error("Failed to close Kafka producer", slog.String("error", err.Error()))
	}

	if err4 := trace.Shutdown(ctxShutdown); err4 != nil {
		log.Error("Error shutting down tracer provider", slog.String("error", err4.Error()))
	}

	dbpool.Close()

	statusServer.Stop(ctxShutdown)

	log.Info("Application stopped")
}
