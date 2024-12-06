// Package producer get messages from channel and send from Sender (Kafka is default)
package producer

import (
	"context"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/model"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/repo"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/sender"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/server"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"sync"
)

// Producer читает из канала событий и отправляет в sender, в n потоков
type Producer struct {
	sender  sender.EventSender          // Интерфейс для отправки сообщений (kafka)
	repo    repo.EventRepo              // Интерфейс для работы с БД
	n       int                         // кол-во потоков
	stop    chan struct{}               // канал для остановки
	events  <-chan []model.PackageEvent // канал для получения событий
	unlocks chan int64                  // канал для снятия блокировки с событий в БД (отправка в кавку неудачная)
	removes chan int64                  // канал для удаления отправленных в кафку событий из БД
	wg      *sync.WaitGroup
}

// NewProducer конструктор
func NewProducer(
	sender sender.EventSender,
	repo repo.EventRepo,
	events <-chan []model.PackageEvent,
	unlocks chan int64,
	removes chan int64,
) *Producer {

	log := slog.With("func", "producer.NewProducer")

	cfg := config.GetConfigInstance()

	wg := &sync.WaitGroup{}
	stop := make(chan struct{})

	log.Info("kafka producer created")

	return &Producer{
		n:       cfg.Outbox.ProducerCount,
		sender:  sender,
		events:  events,
		repo:    repo,
		wg:      wg,
		stop:    stop,
		unlocks: unlocks,
		removes: removes,
	}
}

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

						ctx := context.Background()
						if event[j].TraceID != nil { // Есть информация о root TraceID
							straceid := *event[j].TraceID
							traceid, err := trace.TraceIDFromHex(straceid)
							if err != nil {
								log.Error("TraceID error", slog.String("error", err.Error()))
							}
							ctx = trace.ContextWithRemoteSpanContext(
								context.Background(),
								trace.NewSpanContext(trace.SpanContextConfig{
									TraceID: traceid,
								}),
							)
						}

						err := p.sender.Send(ctx, &event[j], topic) // Поочередно отправляем события

						if err != nil {
							log.Error("Send event", slog.String("error", err.Error()))

							p.unlocks <- event[j].ID // снимаем блокировку с события в БД, для повторной отправки, т.к. отправка в кавку неудачная

							for l := j + 1; l < len(event); l++ {
								p.unlocks <- event[l].ID // снимаем блокировку со всех последующих событий PackageID, для повторной отправки, чтобы не нарушать последовательность событий
							}

							break
						}

						p.removes <- event[j].ID // удаляем событие из БД, т.к. оно обработано и отправлено в кавку
					}
					server.RetranslatorEvents.Add(-1 * float64(len(event))) // метрика, кол-во обрабатываемых событий

				case <-p.stop:
					return
				}
			}
		}()
	}
}

// Stop останавливает пул и ждет, пока все задачи будут выполнены.
// DBConcumer должен быть остановлен.
func (p *Producer) Stop() {

	log := slog.With("func", "producer.Stop")

	close(p.stop)
	p.wg.Wait()

	log.Debug("Producer stopped")
}
