package packaging

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/general"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Edit обработка команды /edit бота
func (c *Commander) Edit(message *tgbotapi.Message) {
	log := slog.With("func", "Commander.Edit")

	args := message.CommandArguments()

	pkg := model.Package{}

	var err error
	// Обработка опционального параметра Weight
	switch strings.Count(args, " ") {
	case 1: // без веса
		_, err = fmt.Sscanf(args, "%d %s", &pkg.ID, &pkg.Title)
		if err != nil {
			log.Debug("wrong args", slog.String("args", args), slog.String("error", err.Error()))
			err = fmt.Errorf("wrong args %v", args)
		}
	case 2: // с весом
		pkg.Weight = sql.NullInt64{
			Int64: 0,
			Valid: true,
		}

		_, err = fmt.Sscanf(args, "%d %s %d", &pkg.ID, &pkg.Title, &pkg.Weight.Int64)
		if err != nil {
			log.Debug("wrong args", slog.String("args", args), slog.String("error", err.Error()))
			err = fmt.Errorf("wrong args %v", args)
			break
		}

		if pkg.Weight.Int64 < 0 { // weight must be positive
			log.Debug("wrong args", slog.String("args", args), slog.String("error", "weight must be positive"))
			err = fmt.Errorf("weight must be positive %v", args)
			break
		}
	default:
		log.Debug("wrong args count", slog.String("args", args))

		err = fmt.Errorf("wrong args %v", args)
	}

	if err != nil {
		c.errorResponseCommand(message, err.Error()) // return error to telegram client
		return
	}

	pkg.Updated = sql.NullTime{Time: time.Now(), Valid: true}

	err = c.packageService.Update(&pkg)
	if err != nil {
		log.Error("fail to edit package", slog.Uint64("id", pkg.ID), slog.String("error", err.Error()))
		if errors.Is(err, general.ErrNotFound) {
			c.errorResponseCommand(message, "Package not found")
			return
		}
		c.errorResponseCommand(message, fmt.Sprintf("Fail to edit package with id %d", pkg.ID))
		return
	}

	// successful response
	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		fmt.Sprintf("Package with id: %d updated", pkg.ID),
	)

	_, err = c.bot.Send(msg)
	if err != nil {
		log.Error("error sending reply message to chat", slog.String("error", err.Error()))
	}

	log.Debug("Package updated", slog.Any("package", pkg))
}
