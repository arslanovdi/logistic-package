package packaging

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/general"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"strings"
	"time"
)

// Edit обработка команды /edit бота
func (c *Commander) Edit(message *tgbotapi.Message) {

	log := slog.With("func", "Commander.Edit")

	args := message.CommandArguments()

	pkg := model.Package{}

	var err error
	// Обработка опционального параметра Weight
	switch strings.Count(args, " ") {
	case 1:
		_, err = fmt.Sscanf(args, "%d %s", &pkg.ID, &pkg.Title)
		if err != nil {
			log.Info("wrong args", slog.Any("args", args), slog.String("error", err.Error()))
			err = fmt.Errorf("wrong args %v", args)
		}
	case 2:
		pkg.Weight = sql.NullInt64{
			Int64: 0,
			Valid: true,
		}
		_, err = fmt.Sscanf(args, "%d %s %d", &pkg.ID, &pkg.Title, pkg.Weight.Int64)
		if err != nil {
			log.Info("wrong args", slog.Any("args", args), slog.String("error", err.Error()))
			err = fmt.Errorf("wrong args %v", args)
		}
	default:
		log.Info("wrong args count", slog.Any("args", args))
		err = fmt.Errorf("wrong args %v", args)
	}

	if err != nil {
		c.errorResponseCommand(message, err.Error())
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
