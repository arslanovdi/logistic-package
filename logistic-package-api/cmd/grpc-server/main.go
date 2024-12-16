// Core домен
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/cache"

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/database/postgres"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/service"
	"github.com/arslanovdi/logistic-package/pkg/logger"
	"github.com/arslanovdi/logistic-package/pkg/tracer"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/database"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/server"
	pkgserver "github.com/arslanovdi/logistic-package/pkg/server"
)

const (
	startTimeout = 5 * time.Second
)

func main() {
	logger.InitializeLogger(slog.LevelDebug)

	version := flag.String("version", "dev", "Defines the version of the service")
	commitHash := flag.String("commitHash", "-", "Defines the commit hash of the service")
	// дефолтный конфиг файл для локального запуска
	configFile := flag.String("config", "logistic-package-api/config_local.yml", "Defines the config file of the service")
	migration := flag.Bool("migration", true, "Defines the migration start option") // TODO set default to false
	flag.Parse()

	log := slog.With("func", "main")

	if err := config.ReadConfigYML(*configFile); err != nil { // чтение конфигурации, в докере подставляется свой конфиг
		log.Warn("Failed init configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}
	cfg := config.GetConfigInstance()

	cfg.Project.Version = *version // загружаем флагами, полученные из командной строки
	cfg.Project.CommitHash = *commitHash

	if cfg.Project.Debug {
		logger.SetLogLevel(slog.LevelDebug)
	} else {
		logger.SetLogLevel(slog.LevelInfo)
	}

	startCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(cfg.Project.StartupTimeout)) // контекст запуска приложения
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

	ctxTrace, cancelTrace := context.WithTimeout(context.Background(), startTimeout) // контекст запуска grpc экспортера в jaeger
	defer cancelTrace()

	trace, err1 := tracer.NewTracer(
		ctxTrace,
		cfg.Project.Instance,
		cfg.Jaeger.Host+cfg.Jaeger.Port)
	if err1 != nil {
		log.Warn("Failed to init tracer", slog.String("error", err1.Error()))
		os.Exit(1)
	}

	isReady := &atomic.Bool{} // состояние приложения
	isReady.Store(false)

	go func() { // TODO отсечка статус сервера
		time.Sleep(startTimeout)
		isReady.Store(true)
		log.Info("The service is ready to accept requests")
	}()

	statusServer := pkgserver.NewStatusServer(
		isReady,
		&pkgserver.StatusConfig{
			Host:          cfg.Status.Host,
			Port:          cfg.Status.Port,
			LivenessPath:  cfg.Status.LivenessPath,
			ReadinessPath: cfg.Status.ReadinessPath,
			VersionPath:   cfg.Status.VersionPath,
		},
		&pkgserver.ProjectInfo{
			Name:        cfg.Project.Name,
			Debug:       cfg.Project.Debug,
			Environment: cfg.Project.Environment,
			Version:     cfg.Project.Version,
			CommitHash:  cfg.Project.CommitHash,
			Instance:    cfg.Project.Instance,
		},
	)
	statusServer.Start()

	metricsServer := pkgserver.NewMetricsServer(
		&pkgserver.MetricsConfig{
			Host: cfg.Metrics.Host,
			Port: cfg.Metrics.Port,
			Path: cfg.Metrics.Path,
		})
	metricsServer.Start()

	dbpool := database.MustGetPgxPool(context.Background())

	repo := postgres.NewPostgresRepo(dbpool)           // интерфейс работы с БД
	redis := cache.NewRedis(repo)                      // интерфейс работы с кэшем
	packageService := service.NewPackageService(redis) // интерфейс работы с пакетами

	if *migration { // миграция параметром из командной строки
		log.Info("Migration started")
		if err := goose.Up(stdlib.OpenDBFromPool(dbpool), // получаем соединение с базой данных из пула
			cfg.Database.Migrations); err != nil {
			log.Warn("Migration failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}

	grpcServer := server.NewGrpcServer(packageService)
	grpcServer.Start()

	gatewayServer := server.NewGatewayServer() // grpc-gateway
	gatewayServer.Start()

	cancel() // отменяем контекст запуска приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Info("Graceful shutdown")

	ctxShutdown, cancelShutdown := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(cfg.Project.ShutdownTimeout)) // контекст останова приложения
	defer cancelShutdown()
	go func() {
		<-ctxShutdown.Done()
		log.Warn("Application shutdown time exceeded")
		os.Exit(1)
	}()

	isReady.Store(false)

	if err3 := grpcServer.Stop(); err3 != nil {
		log.Error("Failed to stop gRPC server", slog.String("error", err3.Error()))
	}

	metricsServer.Stop(ctxShutdown)
	gatewayServer.Stop(ctxShutdown)
	if err4 := trace.Shutdown(ctxShutdown); err4 != nil {
		log.Error("Error shutting down tracer provider", slog.String("error", err4.Error()))
	}

	redis.Close()
	dbpool.Close()

	statusServer.Stop(ctxShutdown)

	log.Info("Application stopped correctly")
}
