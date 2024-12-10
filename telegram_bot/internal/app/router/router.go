// Package router - роутер для обработки сообщений телеграм бота
package router

import (
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/app/commands/logistic"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/service"
	"log/slog"
	"runtime/debug"

	"github.com/arslanovdi/logistic-package/telegram_bot/internal/app/path"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const myDomain = "logistic"

// LogisticCommander - интерфейс представляет методы для обработки команд и кнопок телеграм бота
type LogisticCommander interface {
	HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath)
	HandleCommand(callback *tgbotapi.Message, commandPath path.CommandPath)
}

// Router - роутер для обработки сообщений телеграм бота
type Router struct {
	bot               *tgbotapi.BotAPI
	logisticCommander LogisticCommander // Экземпляр интерфейса обрабатывающий сообщения телеграм бота
}

// New конструктор
func New(
	bot *tgbotapi.BotAPI,
	pkgService *service.LogisticPackageService,
) *Router {
	return &Router{
		bot:               bot,
		logisticCommander: logistic.NewLogisticCommander(bot, pkgService),
	}
}

// HandleUpdate обработка сообщений телеграм бота
func (c *Router) HandleUpdate(update *tgbotapi.Update) {

	log := slog.With("func", "Router.HandleUpdate")

	defer func() { // recover, чтобы бот не умер
		if panicValue := recover(); panicValue != nil {
			log.Warn("recovered from panic", slog.Any("panic value", panicValue), slog.String("stack", string(debug.Stack())))
		}
	}()

	switch {
	case update.CallbackQuery != nil:
		c.handleCallback(update.CallbackQuery) // обработка кнопок
	case update.Message != nil:
		c.handleMessage(update.Message) // обработка сообщений
	}
}

// handleCallback обработка нажатия кнопок
func (c *Router) handleCallback(callback *tgbotapi.CallbackQuery) {

	log := slog.With("func", "Router.handleCallback")

	callbackPath, err := path.ParseCallback(callback.Data)
	if err != nil {
		log.Info("error parsing callback data", slog.String("data", callback.Data), slog.String("error", err.Error()))
		return
	}

	switch callbackPath.Domain {
	case myDomain:
		c.logisticCommander.HandleCallback(callback, callbackPath)
	default:
		log.Info("unknown domain", slog.String("domain", callbackPath.Domain))
	}
}

// handleMessage обработка команд
func (c *Router) handleMessage(msg *tgbotapi.Message) {

	log := slog.With("func", "Router.handleMessage")

	if !msg.IsCommand() {
		c.showCommandFormat(msg)
		return
	}

	commandPath, err := path.ParseCommand(msg.Command())
	if err != nil {
		log.Error("error parsing command", slog.String("command", msg.Command()), slog.String("error", err.Error()))
		return
	}

	switch commandPath.Domain {
	case myDomain:
		c.logisticCommander.HandleCommand(msg, commandPath)
	default:
		log.Info("unknown domain", slog.String("domain", commandPath.Domain))
	}
}

// showCommandFormat выдача в бот сообщения с форматом команд
func (c *Router) showCommandFormat(inputMessage *tgbotapi.Message) {

	log := slog.With("func", "Router.showCommandFormat")

	outputMsg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Command format: /{command}__logistic__package")

	_, err := c.bot.Send(outputMsg)
	if err != nil {
		log.Error("error sending reply message to chat", slog.String("error", err.Error()))
	}
}
