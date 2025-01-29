// Package consumer принимает события из базы данных и отправляет их в канал
package consumer

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/metrics"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/repo"
	"github.com/arslanovdi/logistic-package/pkg/model"
)

// Consumer читает из базы данных события в n потоков и отправляет их в канал
type Consumer struct {
	repo           repo.EventRepo              // Интерфейс работы с БД
	batchSize      int                         // Кол-во записей получаемых из БД за один запрос
	stop           chan struct{}               // Канал для остановки
	events         chan<- []model.PackageEvent // Канал для передачи событий
	eventsCount    int                         // Кол-во обрабатываемых событий, race condition быть не должно.
	unlocks        chan int64                  // Канал передачи PackageEventID для снятия блокировки
	removes        chan int64                  // Канал передачи PackageEventID для удаления
	unlockEvents   []int64                     // Слайс с PackageEventID, для удаления из БД
	removeEvents   []int64                     // Слайс с PackageEventID, для удаления из БД
	tick           time.Duration               // Интервал между запросами в БД
	queryTimeout   time.Duration               // Таймаут операций с БД
	inWork         bool                        // Флаг обработки батча
	processingTime time.Time                   // Время начала обработки батча
	wg             *sync.WaitGroup
}

// NewDbConsumer конструктор
func NewDbConsumer(
	r repo.EventRepo,
	events chan<- []model.PackageEvent,
	unlocks, removes chan int64,
) *Consumer {
	log := slog.With("func", "consumer.NewDbConsumer")

	cfg := config.GetConfigInstance()

	wg := &sync.WaitGroup{}
	stop := make(chan struct{})

	unlockEvents := make([]int64, 0, cfg.Outbox.BatchSize) // слайс с PackageEventID, для удаления из БД
	removeEvents := make([]int64, 0, cfg.Outbox.BatchSize) // слайс с PackageEventID, для удаления из БД

	c := Consumer{
		events:       events,
		repo:         r,
		batchSize:    cfg.Outbox.BatchSize,
		tick:         time.Second * time.Duration(cfg.Outbox.Ticker),
		wg:           wg,
		stop:         stop,
		unlockEvents: unlockEvents,
		removeEvents: removeEvents,
		unlocks:      unlocks,
		removes:      removes,
		queryTimeout: time.Second * time.Duration(cfg.Database.QueryTimeout),
		inWork:       false,
	}

	// собирает батч событий для удаления из БД
	go func() {
		for event := range c.removes { // Получаем PackageEventID
			c.removeEvents = append(c.removeEvents, event)

			// Если кол-во отправленных в кафку событий равно кол-ву обработанных событий. Это произойдет только в одной из горутин.
			if c.eventsCount == len(c.removeEvents)+len(c.unlockEvents) {
				c.batchProcessing()
			}
		}
	}()

	// собирает батч событий для разблокировки повторной отправки в кафку
	go func() {
		for event := range c.unlocks { // Получаем PackageEventID
			c.unlockEvents = append(c.unlockEvents, event)

			// Если кол-во отправленных в кафку событий равно кол-ву обработанных событий. Это произойдет только в одной из горутин.
			if c.eventsCount == len(c.removeEvents)+len(c.unlockEvents) {
				c.batchProcessing()
			}
		}
	}()

	log.Info("db consumer created")

	return &c
}

// Start starts DB consumer
func (c *Consumer) Start() {
	log := slog.With("func", "consumer.Start")

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		ticker := time.NewTicker(c.tick) // тикер

		for {
			select {
			case <-ticker.C:
				if c.inWork {
					continue
				}

				ctx, cancel := context.WithTimeout(context.Background(), c.queryTimeout)

				events, err := c.repo.Lock(ctx, c.batchSize) // берем события из базы, отсортированные по PackageID

				cancel()

				if err != nil {
					log.Error("Error getting events from db", slog.String("error", err.Error()))
					c.inWork = false
					continue
				}
				if len(events) == 0 {
					c.inWork = false
					continue
				}

				c.wg.Add(1) // wait for batchProcessing
				c.inWork = true
				c.processingTime = time.Now()

				c.eventsCount = len(events)

				log.Debug("got events from db", slog.Int("count", len(events)))

				metrics.RetranslatorEvents.Add(float64(len(events))) // метрика, кол-во обрабатываемых событий, прибавляем к счетчику

				packageID := events[0].PackageID
				index := 0
				for i := 0; i < len(events); i++ {
					if packageID == events[i].PackageID {
						continue
					}
					c.events <- events[index:i] // передает события в канал с разбивкой по PackageID

					log.Debug("send event to channel", slog.Any("PackageID", packageID), slog.Int("event count", len(events[index:i])))

					index = i
					packageID = events[i].PackageID
				}
				c.events <- events[index:] // события последнего packageID

			case <-c.stop:
				return
			}
		}
	}()
}

// Stop consumer
func (c *Consumer) Stop() {
	log := slog.With("func", "consumer.Stop")

	close(c.stop)
	c.wg.Wait()

	log.Info("db consumer stopped")
}

// batchProcessing удаление из БД отправленных в кафку событий
// разблокировка в БД не отправленных событий
func (c *Consumer) batchProcessing() {
	defer c.wg.Done()

	log := slog.With("func", "consumer.batchProcessing")

	ctx, cancel := context.WithTimeout(context.Background(), c.queryTimeout)
	defer cancel()

	if len(c.removeEvents) > 0 {
		// Обработка батча завершена, удаляем события из БД
		err := c.repo.Remove(ctx, c.removeEvents)
		if err != nil { // TODO Будут повторные отправки, процессинг должен быть идемпотентным. Либо ретраить до упора.
			log.Error("Ошибка при удалении события из БД", slog.String("error", err.Error()))
		}
		log.Debug("Remove event from package_events", slog.Int("events count", len(c.removeEvents)))
		c.removeEvents = c.removeEvents[:0]
		c.inWork = false
	}

	if len(c.unlockEvents) > 0 {
		// Обработка батча завершена, снимаем блокировку
		err := c.repo.Unlock(ctx, c.unlockEvents)
		//goland:noinspection GoLinter
		if err != nil {
			log.Error("Ошибка при снятии блокировки с события в БД", slog.String("error", err.Error()))
			os.Exit(1) //nolint:gocritic  // TODO нужно ретраить, либо перезапускать сервис. Иначе может нарушиться порядок событий
		}
		log.Debug("Unlock unset events", slog.Int("events count", len(c.unlockEvents)))
		c.unlockEvents = c.unlockEvents[:0]
		c.inWork = false
	}

	metrics.ProcessingTime.Set(time.Since(c.processingTime).Seconds())
}
