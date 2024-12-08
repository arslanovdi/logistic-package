// Package consumer принимает события из базы данных и отправляет их в канал
package consumer

import (
	"context"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/metrics"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/outbox/repo"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"log/slog"
	"sync"
	"time"
)

// Consumer читает из базы данных события в n потоков и отправляет их в канал
type Consumer struct {
	repo         repo.EventRepo              // интерфейс работы с БД
	batchSize    int                         // кол-во записей получаемых из БД за один запрос
	stop         chan struct{}               // канал для остановки
	events       chan<- []model.PackageEvent // канал для передачи событий
	eventsCount  int                         // кол-во обрабатываемых событий, конкурентный доступ на чтение! Т.к. сначала происходит запись и только потом чтение в потоках, race condition быть не должно.
	unlocks      chan int64                  // канал передачи PackageEventID для снятия блокировки
	removes      chan int64                  // канал передачи PackageEventID для удаления
	unlockEvents []int64                     // слайс с PackageEventID, для удаления из БД
	removeEvents []int64                     // слайс с PackageEventID, для удаления из БД
	tick         time.Duration               // интервал между запросами в БД
	queryTimeout time.Duration               // таймаут операций с БД
	wg           *sync.WaitGroup
	inwork       bool
}

// NewDbConsumer конструктор
func NewDbConsumer(
	repo repo.EventRepo,
	events chan<- []model.PackageEvent,
	unlocks chan int64,
	removes chan int64,
) *Consumer {

	log := slog.With("func", "consumer.NewDbConsumer")

	cfg := config.GetConfigInstance()

	wg := &sync.WaitGroup{}
	stop := make(chan struct{})

	unlockEvents := make([]int64, 0, cfg.Outbox.BatchSize) // слайс с PackageEventID, для удаления из БД
	removeEvents := make([]int64, 0, cfg.Outbox.BatchSize) // слайс с PackageEventID, для удаления из БД

	c := Consumer{
		events:       events,
		repo:         repo,
		batchSize:    cfg.Outbox.BatchSize,
		tick:         time.Second * time.Duration(cfg.Outbox.Ticker),
		wg:           wg,
		stop:         stop,
		unlockEvents: unlockEvents,
		removeEvents: removeEvents,
		unlocks:      unlocks,
		removes:      removes,
		queryTimeout: time.Second * time.Duration(cfg.Database.QueryTimeout),
		inwork:       false,
	}

	// собирает батч событий для удаления из БД
	go func() {
		for {
			select {
			case event, ok := <-c.removes: // Получаем PackageEventID
				if !ok {
					return
				}
				c.removeEvents = append(c.removeEvents, event)

				if c.eventsCount == len(c.removeEvents)+len(c.unlockEvents) { // когда кол-во отправленных в кафку событий станет равно кол-ву обработанных событий. Это произойдет только в одной из горутин.
					c.batchprocessing()
				}
			}
		}
	}()

	// собирает батч событий для разблокировки повторной отправки в кафку
	go func() {
		for {
			select {
			case event, ok := <-c.unlocks: // Получаем PackageEventID
				if !ok {
					return
				}
				c.unlockEvents = append(c.unlockEvents, event)

				if c.eventsCount == len(c.removeEvents)+len(c.unlockEvents) { // когда кол-во отправленных в кафку событий станет равно кол-ву обработанных событий. Это произойдет только в одной из горутин.
					c.batchprocessing()
				}
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
				if c.inwork {
					continue
				}

				ctx, cancel := context.WithTimeout(context.Background(), c.queryTimeout)

				events, err := c.repo.Lock(ctx, c.batchSize) // берем события из базы, отсортированные по .PackageID

				cancel()

				if err != nil {
					log.Error("Error getting events from db", slog.String("error", err.Error()))
					c.inwork = false
					continue
				}
				if len(events) == 0 {
					c.inwork = false
					continue
				}

				c.wg.Add(1) // whait for batchprocessing
				c.inwork = true

				c.eventsCount = len(events)

				log.Debug("got events from db", slog.Int("count", len(events)))

				metrics.RetranslatorEvents.Add(float64(len(events))) // метрика, кол-во обрабатываемых событий, прибавляем к счетчику

				packageid := events[0].PackageID
				index := 0
				for i := 0; i < len(events); i++ {
					if packageid == events[i].PackageID {
						continue
					}
					c.events <- events[index:i] // передает события в канал с разбивкой по PackageID

					log.Debug("send event to channel", slog.Any("PackageID", packageid), slog.Int("event count", len(events[index:i])))

					index = i
					packageid = events[i].PackageID
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

// batchprocessing удаление из БД отправленных в кафку событий, батчем
// разблокировка в БД не отправленных событий, батчем
func (c *Consumer) batchprocessing() {
	defer c.wg.Done()

	log := slog.With("func", "consumer.batchprocessing")

	ctx, cancel := context.WithTimeout(context.Background(), c.queryTimeout)
	defer cancel()

	if len(c.removeEvents) > 0 {
		// Обработка батча завершена, удаляем события из БД
		err := c.repo.Remove(ctx, c.removeEvents)
		if err != nil {
			log.Error("Ошибка при удалении события из БД", slog.String("error", err.Error()))
		}
		log.Debug("Remove event from package_events", slog.Int("events count", len(c.removeEvents)))
		c.removeEvents = c.removeEvents[:0]
		c.inwork = false
	}

	if len(c.unlockEvents) > 0 {
		// Обработка батча завершена, снимаем блокировку
		err := c.repo.Unlock(ctx, c.unlockEvents)
		if err != nil {
			log.Error("Ошибка при снятии блокировки с события в БД", slog.String("error", err.Error()))
		}
		log.Debug("Unlock unset events", slog.Int("events count", len(c.unlockEvents)))
		c.unlockEvents = c.unlockEvents[:0]
		c.inwork = false
	}

}
