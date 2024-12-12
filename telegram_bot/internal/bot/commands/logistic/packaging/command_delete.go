package packaging

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/arslanovdi/logistic-package/telegram_bot/internal/general"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Delete обработка команды /delete бота
func (c *Commander) Delete(message *tgbotapi.Message) {
	log := slog.With("func", "Commander.Delete")

	args := message.CommandArguments()

	// validate args count
	if strings.Count(args, " ") > 0 {
		c.errorResponseCommand(message, fmt.Sprintf("too many args %v", args))
		log.Debug("wrong args", slog.String("args", args), slog.String("error", "too many args"))
		return
	}

	// validate uint type
	id := uint64(0)
	_, err := fmt.Sscanf(args, "%d", &id)
	if err != nil {
		c.errorResponseCommand(message, fmt.Sprintf("wrong args %v\n", args))
		log.Debug("wrong args", slog.String("args", args), slog.String("error", err.Error()))
		return
	}

	err = c.packageService.Delete(id)
	if err != nil {
		if errors.Is(err, general.ErrNotFound) {
			c.errorResponseCommand(message, "Package not found")
			return
		}
		log.Error("fail to delete package", slog.Uint64("id", id), slog.String("error", err.Error()))
		c.errorResponseCommand(message, fmt.Sprintf("Fail to delete package with id %d", id))
		return
	}

	// successful response
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
