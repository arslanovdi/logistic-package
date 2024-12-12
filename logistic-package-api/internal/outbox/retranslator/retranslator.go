// Package retranslator get events from database (consumer) and send to kafka (producer)
package retranslator

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/consumer"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/producer"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/repo"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/sender"
	"github.com/arslanovdi/logistic-package/pkg/model"
)

// Retranslator считывает события из БД и отправляет в кафку
type Retranslator interface {
	// Start retranslator
	Start(topic string)
	// Stop retranslator
	Stop()
}

type retranslator struct {
	events   chan []model.PackageEvent
	unlocks  chan int64
	removes  chan int64
	consumer *consumer.Consumer
	producer *producer.Producer
}

// NewRetranslator конструктор
func NewRetranslator(r repo.EventRepo, s sender.EventSender) Retranslator {
	log := slog.With("func", "retranslator.NewRetranslator")

	cfg := config.GetConfigInstance()

	events := make(chan []model.PackageEvent)
	unlocks := make(chan int64)
	removes := make(chan int64)
	dbConsumer := consumer.NewDbConsumer(
		r,
		events,
		unlocks,
		removes,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.Database.QueryTimeout))
	defer cancel()

	err := r.UnlockAll(ctx) // Разблокируем сообщения. Такие могут появиться если ретранслятор был завершен без graceful.
	if err != nil {
		log.Error("error unlock all", slog.String("error", err.Error()))
		os.Exit(1) // TODO нужно ретраить, либо перезапускать сервис. Иначе может нарушиться порядок событий
	}

	kafkaProducer := producer.NewProducer(
		s,
		events,
		unlocks,
		removes,
	)

	log.Info("Retranslator created")

	return &retranslator{
		events:   events,
		unlocks:  unlocks,
		removes:  removes,
		consumer: dbConsumer,
		producer: kafkaProducer,
	}
}

// Start retranslator (outbox pattern)
func (r *retranslator) Start(topic string) {
	log := slog.With("func", "retranslator.Start")

	r.producer.Start(topic)
	r.consumer.Start()

	log.Info("Retranslator started")
}

// Stop retranslator
func (r *retranslator) Stop() {
	log := slog.With("func", "retranslator.Stop")

	r.consumer.Stop()
	r.producer.Stop()

	close(r.events)
	close(r.unlocks)
	close(r.removes)

	log.Info("Retranslator stopped")
}
