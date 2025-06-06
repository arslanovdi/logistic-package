package packaging

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/arslanovdi/logistic-package/telegram_bot/internal/general"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Get обработка команды /get бота
func (c *Commander) Get(message *tgbotapi.Message) {
	log := slog.With("func", "Commander.Get")

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
		c.errorResponseCommand(message, fmt.Sprintf("wrong args %v", args))
		log.Debug("wrong args", slog.String("args", args), slog.String("error", err.Error()))
		return
	}

	pkg, err := c.packageService.Get(id)
	if err != nil {
		log.Error("fail to get product", slog.Uint64("id", id), slog.String("error", err.Error()))
		if errors.Is(err, general.ErrNotFound) {
			c.errorResponseCommand(message, fmt.Sprintf("Package with id: %d not found.\n", id))
			return
		}
		c.errorResponseCommand(message, "fail to get product\n")
		return
	}

	// successful response
	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		pkg.String(),
	)

	_, err = c.bot.Send(msg)
	if err != nil {
		log.Error("error sending reply message to chat", slog.String("error", err.Error()))
	}

	log.Debug("get package", slog.Any("pkg", pkg))
}
