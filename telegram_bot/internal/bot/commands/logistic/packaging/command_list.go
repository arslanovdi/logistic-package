package packaging

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"github.com/arslanovdi/logistic-package/telegram_bot/internal/general"

	"github.com/arslanovdi/logistic-package/telegram_bot/internal/bot/path"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// List обработка команды /list бота
func (c *Commander) List(message *tgbotapi.Message) {
	log := slog.With("func", "Commander.List")

	outputMsgText := strings.Builder{}
	outputMsgText.WriteString("These are all our packages: \n\n")

	var endOfList bool

	packages, err := c.packageService.List(1, limit)
	if err != nil {
		if errors.Is(err, general.ErrNotFound) {
			c.errorResponseCommand(message, "packages not found")
			return
		}
		if !errors.Is(err, general.ErrEndOfList) {
			c.errorResponseCommand(message, "Ошибка получения списка")
			log.Error("fail to get list of packages", slog.String("error", err.Error()))
			return
		}

		endOfList = true
	}

	for _, p := range packages {
		outputMsgText.WriteString(p.String())
		outputMsgText.WriteString("\n")
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, outputMsgText.String())

	// добавляем кнопку в ответ, если это не конец списка
	if !endOfList {
		serializedData, err1 := json.Marshal(CallbackListData{ // serialized data in button
			Offset: 1,
		})

		if err1 != nil {
			c.errorResponseCommand(message, "Error serializing data")
			log.Error("Error serializing data", slog.String("error", err1.Error()))
			return
		}

		callbackPath := path.CallbackPath{ // собираем структуру кнопки
			Domain:       "logistic",
			Subdomain:    "package",
			CallbackName: "list",
			CallbackData: string(serializedData),
		}

		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup( // добавляем кнопку в ответ
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Next page", callbackPath.String()),
			),
		)
	}

	_, err = c.bot.Send(msg)
	if err != nil {
		log.Error("error sending reply message to chat", slog.String("error", err.Error()))
	}

	log.Debug("Command List packages", slog.Uint64("offset", 1), slog.Uint64("limit", limit))
}
