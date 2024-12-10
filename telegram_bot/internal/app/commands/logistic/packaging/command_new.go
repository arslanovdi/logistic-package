package packaging

import (
	"database/sql"
	"fmt"
	"github.com/arslanovdi/logistic-package/pkg/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"strings"
	"time"
)

// New обработка команды /new бота
func (c *Commander) New(message *tgbotapi.Message) {

	log := slog.With("func", "Commander.New")

	args := message.CommandArguments()

	pkg := model.Package{}

	var err error
	// Обработка опционального параметра Weight
	switch strings.Count(args, " ") {
	case 0:
		_, err = fmt.Sscanf(args, "%s", &pkg.Title)
		if err != nil {
			log.Info("wrong args", slog.Any("args", args), slog.String("error", err.Error()))
			err = fmt.Errorf("wrong args %v", args)
		}
	case 1:
		pkg.Weight = sql.NullInt64{
			Int64: 0,
			Valid: true,
		}
		_, err = fmt.Sscanf(args, "%s %d", &pkg.Title, pkg.Weight.Int64)
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

	pkg.Created = time.Now()

	id, err := c.packageService.Create(&pkg)

	if err != nil {
		c.errorResponseCommand(message, fmt.Sprintf("Fail to create package with title %v", pkg.Title))
		log.Error("fail to create package", slog.String("package", pkg.String()), slog.String("error", err.Error()))
		return
	}

	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		fmt.Sprintf("package %v created with id: %d", pkg.Title, id),
	)

	_, err = c.bot.Send(msg)
	if err != nil {
		log.Error("error sending reply message to chat", slog.String("error", err.Error()))
	}

	log.Debug("Package created", slog.Uint64("id", id), slog.String("package", pkg.String()))
}
