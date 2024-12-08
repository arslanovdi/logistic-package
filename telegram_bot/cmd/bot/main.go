// Телеграм бот для управления логистикой пакетов
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/arslanovdi/logistic-package/pkg/logger"
	"github.com/arslanovdi/logistic-package/pkg/server"
	"github.com/arslanovdi/logistic-package/pkg/tracer"
	routerPkg "github.com/arslanovdi/logistic-package/telegram_bot/internal/app/router"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/config"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/fake"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/grpc"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func main() {
	logger.InitializeLogger(slog.LevelDebug) // slog logger

	version := flag.String("version", "dev", "Defines the version of the service")
	commitHash := flag.String("commitHash", "-", "Defines the commit hash of the service")
	configFile := flag.String("config", "telegram_bot/config_local.yml", "Defines the config file of the service")
	flag.Parse()

	log := slog.With("func", "main")

	err1 := config.ReadConfigYML(*configFile)
	if err1 != nil {
		log.Warn("Failed to read config", slog.String("error", err1.Error()))
		os.Exit(1)
	}
	cfg := config.GetConfigInstance()

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

	go func() { // TODO отсечка статус сервера
		time.Sleep(starttimeout)
		isReady.Store(true)
		log.Info("The service is ready to accept requests")
	}()

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

	grpcClient := grpc.NewGrpcClient()
	packageService := service.NewPackageService(grpcClient)

	bot, err3 := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err3 != nil {
		log.Warn("Failed to create new bot", slog.String("error", err3.Error()))
		os.Exit(1)
	}

	log.Info("Telegram bot authorized on account ", slog.String("account", bot.Self.UserName))

	// Uncomment if you want debugging
	// bot.Debug = true

	u := tgbotapi.UpdateConfig{
		Timeout: 60,
	}

	updates := bot.GetUpdatesChan(u) // получаем канал обновлений телеграм бота

	routerHandler := routerPkg.New(bot, packageService) // Создаем обработчик телегрм бота

	go fake.Emulate(2000, packageService) // запускаем эмуляцию пользователей телеграм бота

	cancel() // отменяем контекст запуска приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT) // подписываем канал на сигналы завершения процесса
	for {
		select {
		case update := <-updates:
			routerHandler.HandleUpdate(update)
		case <-stop:
			log.Info("Graceful shutdown")

			ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.Project.ShutdownTimeout))
			defer cancelShutdown()
			go func() {
				<-ctxShutdown.Done()
				log.Warn("Application shutdown time exceeded")
				os.Exit(1)
			}()

			isReady.Store(false)

			metricsServer.Stop(ctxShutdown)

			grpcClient.Close()
			if err := trace.Shutdown(ctxTrace); err != nil {
				log.Error("Error shutting down tracer provider", slog.String("error", err.Error()))
			}

			statusServer.Stop(ctxShutdown)

			log.Info("Application stopped")
			return
		}
	}
}
