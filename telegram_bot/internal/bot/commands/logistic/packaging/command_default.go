package packaging

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Default ответ при неверной команде
func (c *Commander) Default(message *tgbotapi.Message) {
	log := slog.With("func", "Commander.Default")

	log.Debug("wrong command", slog.String("username", message.From.UserName), slog.String("message", message.Text))

	msg := tgbotapi.NewMessage(message.Chat.ID, "You wrote: "+message.Text)

	_, err := c.bot.Send(msg)
	if err != nil {
		log.Error("error sending reply message to chat", slog.String("error", err.Error()))
	}
}
