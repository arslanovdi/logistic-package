package test

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/retranslator"
	"github.com/arslanovdi/logistic-package/logistic-package-api/mocks"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	if err := config.ReadConfigYML("../config_retranslator_local.yml"); err != nil {
		slog.Warn("Failed init configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func TestRetranslator(t *testing.T) {
	t.Parallel()
	cfg := config.GetConfigInstance()

	t.Run("consume one message; send it to topic; delete in repo", func(t *testing.T) {
		t.Parallel()
		event := model.PackageEvent{
			ID:        55,
			PackageID: 1,
			Type:      1,
			Status:    0,
			Payload:   nil,
			TraceID:   nil,
		}

		repoMock := mocks.NewEventRepo(t)
		senderMock := mocks.NewEventSender(t)

		repoMock.EXPECT().UnlockAll(mock.AnythingOfType("*context.timerCtx")).Return(nil)
		r := retranslator.NewRetranslator(repoMock, senderMock)

		repoMock.EXPECT().Lock(mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("int")).Return([]model.PackageEvent{event}, nil)
		senderMock.EXPECT().Send(mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*model.PackageEvent"), "test_topic").Return(nil)
		repoMock.EXPECT().Remove(mock.AnythingOfType("*context.timerCtx"), []int64{55}).Return(nil)
		repoMock.AssertNotCalled(t, "Unlock", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("[]int64]"))

		r.Start("test_topic")

		time.Sleep(time.Second * time.Duration(cfg.Outbox.Ticker+1)) // Ждем пока сработает тикер, по которому происходит чтение из БД. 1 секунды должно быть достаточно для отправки 1 сообщения.

		r.Stop()
	})

	t.Run("consume one message; error on send it to topic; unlock in repo", func(t *testing.T) {
		t.Parallel()
		event := model.PackageEvent{
			ID:        55,
			PackageID: 1,
			Type:      1,
			Status:    0,
			Payload:   nil,
			TraceID:   nil,
		}

		repoMock := mocks.NewEventRepo(t)
		senderMock := mocks.NewEventSender(t)

		repoMock.EXPECT().UnlockAll(mock.AnythingOfType("*context.timerCtx")).Return(nil)
		r := retranslator.NewRetranslator(repoMock, senderMock)

		repoMock.EXPECT().Lock(mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("int")).Return([]model.PackageEvent{event}, nil)
		senderMock.EXPECT().Send(mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*model.PackageEvent"), "test_topic").Return(kafka.NewError(kafka.ErrMsgTimedOut, "dont send", false))
		repoMock.EXPECT().Unlock(mock.AnythingOfType("*context.timerCtx"), []int64{55}).Return(nil)
		repoMock.AssertNotCalled(t, "Remove", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("[]int64"))

		r.Start("test_topic")

		time.Sleep(time.Second * time.Duration(cfg.Outbox.Ticker+1)) // Ждем пока сработает тикер, по которому происходит чтение из БД. 1 секунды должно быть достаточно для отправки 1 сообщения.

		r.Stop()
	})
}
