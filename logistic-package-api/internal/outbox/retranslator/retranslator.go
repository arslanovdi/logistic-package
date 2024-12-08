// Package retranslator get events from database (consumer) and send to kafka (producer)
package retranslator

import (
	"context"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/consumer"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/producer"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/repo"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/sender"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"log/slog"
	"time"
)

// Retranslator считывает события из БД и отправляет в кафку
type Retranslator interface {
	Start(topic string)
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
func NewRetranslator(Repo repo.EventRepo, Sender sender.EventSender) Retranslator {

	log := slog.With("func", "retranslator.NewRetranslator")

	globalcfg := config.GetConfigInstance()

	events := make(chan []model.PackageEvent)
	unlocks := make(chan int64)
	removes := make(chan int64)
	dbconsumer := consumer.NewDbConsumer(
		Repo,
		events,
		unlocks,
		removes,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(globalcfg.Database.QueryTimeout))
	defer cancel()

	err := Repo.UnlockAll(ctx) // Разблокируем сообщения. Такие могут появиться если ретранслятор был завершен без graceful.
	if err != nil {
		log.Error("error unlock all", err)
	}

	kafkaproducer := producer.NewProducer(
		Sender,
		Repo,
		events,
		unlocks,
		removes,
	)

	log.Info("Retranslator created")

	return &retranslator{
		events:   events,
		unlocks:  unlocks,
		removes:  removes,
		consumer: dbconsumer,
		producer: kafkaproducer,
	}
}

// Start запуск ретранслятора (outbox pattern)
func (r *retranslator) Start(topic string) {

	log := slog.With("func", "retranslator.Start")

	r.producer.Start(topic)
	r.consumer.Start()

	log.Info("Retranslator started")
}

// Stop останавливает пул воркеров
func (r *retranslator) Stop() {

	log := slog.With("func", "retranslator.Stop")

	r.consumer.Stop()
	r.producer.Stop()

	close(r.events)
	close(r.unlocks)
	close(r.removes)

	log.Info("Retranslator stopped")
}
