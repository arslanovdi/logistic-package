// Package producer get messages from channel and send from Sender (Kafka is default)
package producer

import (
	"context"
	"log/slog"
	"sync"

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/metrics"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/sender"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"go.opentelemetry.io/otel/trace"
)

// Producer читает из канала событий и отправляет в sender, в n потоков
type Producer struct {
	sender  sender.EventSender          // Интерфейс для отправки сообщений (kafka)
	n       int                         // Кол-во потоков
	stop    chan struct{}               // Канал для остановки
	events  <-chan []model.PackageEvent // Канал для получения событий
	unlocks chan int64                  // Канал для снятия блокировки с событий в БД (отправка в kafka неудачная)
	removes chan int64                  // Канал для удаления отправленных в кафку событий из БД
	wg      *sync.WaitGroup
}

// NewProducer конструктор
func NewProducer(
	s sender.EventSender,
	events <-chan []model.PackageEvent,
	unlocks, removes chan int64,
) *Producer {
	log := slog.With("func", "producer.NewProducer")

	cfg := config.GetConfigInstance()

	wg := &sync.WaitGroup{}
	stop := make(chan struct{})

	log.Info("kafka producer created")

	return &Producer{
		n:       cfg.Outbox.ProducerCount,
		sender:  s,
		events:  events,
		wg:      wg,
		stop:    stop,
		unlocks: unlocks,
		removes: removes,
	}
}

// Start запуск продюсера
func (p *Producer) Start(topic string) {
	log := slog.With("func", "producer.Start")

	for i := 0; i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case event := <-p.events: // получаем слайс событий по конкретному PackageID и отправляем их в kafka
					for j := 0; j < len(event); j++ {
						ctx := context.Background()  // контекст без трассировки
						if event[j].TraceID != nil { // Есть информация о root TraceID
							traceid, err := trace.TraceIDFromHex(*event[j].TraceID)
							if err != nil {
								log.Error("TraceID error", slog.String("error", err.Error()))
							}

							ctx = trace.ContextWithRemoteSpanContext( // контекст с трассировкой
								context.Background(),
								trace.NewSpanContext(trace.SpanContextConfig{
									TraceID:    traceid,
									TraceFlags: trace.FlagsSampled,
								}),
							)
						}

						err := p.sender.Send(ctx, &event[j], topic) // Поочередно отправляем события
						if err != nil {
							log.Error("Send event", slog.String("error", err.Error()))

							p.unlocks <- event[j].ID // Снять блокировку с события в БД, для повторной отправки, так как отправка в кавку неудачная

							for l := j + 1; l < len(event); l++ {
								p.unlocks <- event[l].ID // Снять блокировку со всех последующих событий PackageID, чтобы не нарушать последовательность событий
							}

							break
						}

						p.removes <- event[j].ID // Удалить событие из БД, так как оно обработано и отправлено в кавку
					}
					metrics.RetranslatorEvents.Add(-1 * float64(len(event))) // метрика, кол-во обрабатываемых событий

				case <-p.stop:
					return
				}
			}
		}()
	}
}

// Stop останавливает пул и ждет, пока все задачи будут выполнены.
// dbConsumer должен быть остановлен.
func (p *Producer) Stop() {
	log := slog.With("func", "producer.Stop")

	close(p.stop)
	p.wg.Wait()

	log.Debug("Producer stopped")
}
