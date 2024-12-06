package packaging

import (
	"errors"
	"fmt"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

// Delete обработка команды /delete бота
func (c *Commander) Delete(message *tgbotapi.Message) {

	log := slog.With("func", "Commander.Delete")

	args := message.CommandArguments()

	id := uint64(0)
	_, err := fmt.Sscanf(args, "%d", &id)
	if err != nil {
		c.errorResponseCommand(message, fmt.Sprintf("wrong args %v\n", args))
		log.Info("wrong args", slog.Any("args", args), slog.String("error", err.Error()))
		return
	}

	err = c.packageService.Delete(id)
	if err != nil {
		log.Error("fail to delete package", slog.Uint64("id", id), slog.String("error", err.Error()))
		if errors.Is(err, model.ErrNotFound) {
			c.errorResponseCommand(message, "Package not found")
			return
		}
		c.errorResponseCommand(message, fmt.Sprintf("Fail to delete package with id %d", id))
		return
	}

	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		fmt.Sprintf("Package with id: %d deleted", id),
	)

	_, err = c.bot.Send(msg)
	if err != nil {
		log.Error("error sending reply message to chat", slog.String("error", err.Error()))
	}

	log.Debug("Package deleted", slog.Uint64("id", id))
}
