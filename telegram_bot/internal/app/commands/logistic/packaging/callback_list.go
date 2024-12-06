package packaging

import (
	"encoding/json"
	"errors"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/model"
	"log/slog"
	"strings"

	"github.com/arslanovdi/logistic-package/telegram_bot/internal/app/path"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CallbackListData структура данных для обработки реакции на нажатие кнопки
type CallbackListData struct {
	Offset int `json:"offset"` // смещение, с которого выводятся записи в телеграм боте
}

// CallbackList обработка реакции на нажатие кнопки
func (c *Commander) CallbackList(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {

	log := slog.With("func", "Commander.CallbackList")

	parsedData := CallbackListData{}
	err := json.Unmarshal([]byte(callbackPath.CallbackData), &parsedData)
	if err != nil {
		c.errorResponseCallback(callback, "внутренняя ошибка")
		log.Error("fail to read json data for type CallbackListData from input string",
			slog.String("input string", callbackPath.CallbackData),
			slog.String("error", err.Error()))
		return
	}

	outputMsgText := strings.Builder{}

	packages, err := c.packageService.List(uint64(parsedData.Offset)+limit, limit) // Запрашиваем клиентов со смещением

	var endOfList bool

	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			c.errorResponseCallback(callback, "packages not found")
			return
		}
		if errors.Is(err, model.ErrEndOfList) {
			endOfList = true
		} else {
			c.errorResponseCallback(callback, "Ошибка получения списка")
			log.Error("fail to get list of packages", slog.String("error", err.Error()))
			return
		}
	}

	outputMsgText.WriteString("These are our packages: \n\n")
	for _, p := range packages {
		outputMsgText.WriteString(p.String())
		outputMsgText.WriteString("\n")
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, outputMsgText.String())

	if !endOfList {
		serializedData, _ := json.Marshal(CallbackListData{ // данные сериализуемые в кнопку
			Offset: parsedData.Offset + limit,
		})
		callbackPath.CallbackData = string(serializedData)

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

	log.Debug("Callback List packages", slog.Uint64("offset", uint64(parsedData.Offset)+limit), slog.Uint64("limit", limit))
}
