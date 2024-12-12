// Package fake - эмуляция работы пользователя телеграм бота
package fake

import (
	"crypto/rand"
	"log/slog"
	"math"
	"math/big"
	"time"

	"github.com/arslanovdi/logistic-package/telegram_bot/internal/service"
)

// % генерируемых запросов
const (
	create = 30
	del    = 10
	get    = 20
	list   = 20
	update = 20
)

const listCount = 10 // Кол-во package выдаваемое за 1 раз

type Faker struct {
	counter    uint64 // Счетчик созданных пакетов
	pkgService *service.LogisticPackageService
	duration   int64
	stop       chan struct{}
}

// genInt генерация случайного числа от 0 до max при помощи пакета crypto/rand
func genInt(maximum uint64) uint64 {
	log := slog.With("func", "fake.genInt")

	if maximum > math.MaxInt64 { // проверка на переполнение
		maximum = math.MaxInt64
	}

	if maximum <= 0 { // проверка на отрицательное значение
		maximum = 1
	}

	rnd, err := rand.Int(rand.Reader, big.NewInt(int64(maximum))) //nolint:gosec // переполнение исключается предыдущим условием
	if err != nil {
		log.Error("FAKE fail to generate random number", slog.String("error", err.Error()))
		return 0
	}

	return rnd.Uint64()
}

// NewFaker конструктор
func NewFaker(duration int64, pkgService *service.LogisticPackageService) Faker {
	return Faker{
		duration:   duration,
		pkgService: pkgService,
		stop:       make(chan struct{}),
	}
}

// Emulate запускает эмуляцию
func (f *Faker) Emulate() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(f.duration))

	go func() {
		for {
			select {
			case <-f.stop:
				return
			case <-ticker.C:
				rnd := genInt(uint64(create + del + get + list + update))

				switch {
				case rnd < create: // create % созданий пакета
					f.eCreate()
				case rnd < create+del: // del % удаления пакета
					f.eDelete()
				case rnd < create+del+get: // get % получения пакета
					f.eGet()
				case rnd < create+del+get+list: // list % получения списка пакетов
					f.eList()
				case rnd < create+del+get+list+update: // update % обновления пакета
					f.eUpdate()
				}
			}
		}
	}()
}

// Stop останавливает эмуляцию
func (f *Faker) Stop() {
	if f.stop == nil {
		return
	}

	close(f.stop)
}
