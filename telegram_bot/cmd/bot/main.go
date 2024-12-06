// Телеграм бот для управления логистикой пакетов
package main

import (
	"context"
	"errors"
	routerPkg "github.com/arslanovdi/logistic-package/telegram_bot/internal/app/router"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/config"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/fake"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/grpc"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/logger"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/service"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/tracer"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const level = slog.LevelDebug // log level

func main() {
	logger.InitializeLogger(level) // slog logger
	log := slog.With("func", "main")

	startCtx, cancel := context.WithTimeout(context.Background(), time.Minute) // контекст запуска приложения
	defer cancel()
	go func() {
		<-startCtx.Done()
		if errors.Is(startCtx.Err(), context.DeadlineExceeded) { // приложение зависло при запуске
			log.Warn("Application startup time exceeded")
			os.Exit(1)
		}
	}()

	err1 := config.ReadConfigYML("config.yml")
	if err1 != nil {
		log.Warn("Failed to read config", slog.String("error", err1.Error()))
		os.Exit(1)
	}

	// TODO move to config
	_ = godotenv.Load()

	token, found := os.LookupEnv("TOKEN")
	if !found {
		log.Warn("environment variable TOKEN not found in .env")
		os.Exit(1)
	}

	ctxTrace, cancelTrace := context.WithCancel(context.Background())
	defer cancelTrace()
	trace, err2 := tracer.New(ctxTrace)
	if err2 != nil {
		log.Warn("Failed to init tracer", slog.String("error", err2.Error()))
		os.Exit(1)
	}

	grpcClient := grpc.NewGrpcClient()
	packageService := service.NewPackageService(grpcClient)

	bot, err3 := tgbotapi.NewBotAPI(token)
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
			slog.Info("Graceful shutdown")
			grpcClient.Close()
			if err := trace.Shutdown(ctxTrace); err != nil {
				log.Error("Error shutting down tracer provider", slog.String("error", err.Error()))
			}
			slog.Info("Application stopped")
			return
		}
	}
}
