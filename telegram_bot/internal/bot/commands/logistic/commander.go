// Package logistic пакет для обработки команд логистики телеграм бота
package logistic

import (
	"log/slog"

	"github.com/arslanovdi/logistic-package/telegram_bot/internal/bot/commands/logistic/packaging"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/service"

	"github.com/arslanovdi/logistic-package/telegram_bot/internal/bot/path"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const mySubDomain = "package"

// PackageCommander интерфейс обрабатывающий команды логистики телеграм бота
type PackageCommander interface {
	// HandleCallback process button events from telegram
	HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath)
	// HandleCommand process messages from telegram
	HandleCommand(message *tgbotapi.Message, commandPath path.CommandPath)
}

// Commander структура обработчика команд логистики телеграм бота
type Commander struct {
	bot              *tgbotapi.BotAPI
	packageCommander PackageCommander
}

// NewLogisticCommander конструктор
func NewLogisticCommander(bot *tgbotapi.BotAPI, pkgService *service.LogisticPackageService) *Commander {
	return &Commander{
		bot:              bot,
		packageCommander: packaging.NewCommander(bot, pkgService),
	}
}

// HandleCallback обработка нажатия кнопок в телеграм боте
func (c *Commander) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	log := slog.With("func", "LogisticCommander.HandleCallback")

	switch callbackPath.Subdomain {
	case mySubDomain:
		c.packageCommander.HandleCallback(callback, callbackPath)
	default:
		log.Info("unknown subdomain", slog.String("subdomain", callbackPath.Subdomain))
	}
}

// HandleCommand обработка команд в телеграм боте
func (c *Commander) HandleCommand(msg *tgbotapi.Message, commandPath path.CommandPath) {
	log := slog.With("func", "LogisticCommander.HandleCommand")

	switch commandPath.Subdomain {
	case mySubDomain:
		c.packageCommander.HandleCommand(msg, commandPath)
	default:
		log.Info("unknown subdomain", slog.String("subdomain", commandPath.Subdomain))
	}
}
