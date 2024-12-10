// Package fake эмуляция работы пользователя телеграм бота
package fake

import (
	"crypto/rand"
	"database/sql"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/service"
	"github.com/brianvoe/gofakeit/v7"
	"log/slog"
	"math"
	"math/big"
	"time"
)

// % запросов
const (
	create = 30
	del    = 10
	get    = 20
	list   = 20
	update = 20
)

var counter int64

// genInt генерация случайного числа от 0 до max при помощи пакета crypto/rand
func genInt(maximum int64) int64 {
	log := slog.With("func", "fake.genInt")
	if maximum <= 0 {
		maximum = math.MaxInt64 - 1
	}
	rnd, err := rand.Int(rand.Reader, big.NewInt(maximum))
	if err != nil {
		log.Error("FAKE fail to generate random number", slog.String("error", err.Error()))
		return 0
	}
	return rnd.Int64()
}

// Emulate Эмуляция работы пользователя телеграм бота
func Emulate(d int64, pkgService *service.LogisticPackageService) {
	for {
		duration := genInt(d)
		time.Sleep(time.Duration(duration) * time.Millisecond) // от 200ms до 1200 мс на одну операцию

		rnd := genInt(100) // 100%
		switch {
		case rnd < create: // create % созданий пакета
			eCreate(pkgService)

		case rnd < create+del: // del % удаления пакета
			if counter <= 2 {
				continue
			}
			eDelete(pkgService)

		case rnd < create+del+get: // get% получения пакета

			if counter <= 2 {
				continue
			}
			eGet(pkgService)

		case rnd < create+del+get+list: // list % получения списка пакетов
			if counter < 10 {
				continue
			}
			eList(pkgService)

		case rnd < create+del+get+list+update: // update % обновления пакета

			if counter <= 2 {
				continue
			}
			eUpdate(pkgService)
		}
	}
}

func eCreate(pkgService *service.LogisticPackageService) {
	log := slog.With("func", "fake.eCreate")
	pkg := model.Package{}

	pkg.Weight = sql.NullInt64{
		Int64: gofakeit.Int64(),
		Valid: gofakeit.Bool(),
	}

	pkg.Title = gofakeit.ProductName()
	pkg.Created = gofakeit.DateRange(time.Now().AddDate(0, -2, 0), time.Now())

	_, err := pkgService.Create(&pkg)
	if err != nil {
		log.Error("FAKE fail to create package", "error", err)
	}
	counter++
}

func eDelete(pkgService *service.LogisticPackageService) {
	log := slog.With("func", "fake.eDelete")

	id := genInt(counter) + 1

	err := pkgService.Delete(uint64(id))
	if err != nil {
		log.Error("FAKE fail to delete package", "error", err)
	}
	counter--
}

func eGet(pkgService *service.LogisticPackageService) {
	log := slog.With("func", "fake.eGet")
	id := genInt(counter) + 1

	_, err := pkgService.Get(uint64(id))
	if err != nil {
		log.Error("FAKE fail to get package", "id", id)
	}
}

func eList(pkgService *service.LogisticPackageService) {
	log := slog.With("func", "fake.eList")
	offset := genInt(counter/2) + 1
	limit := genInt(counter - offset)
	_, err := pkgService.List(uint64(offset), uint64(limit))
	if err != nil {
		log.Error("FAKE fail to list package", "error", err)
	}
}

func eUpdate(pkgService *service.LogisticPackageService) {
	log := slog.With("func", "fake.eUpdate")
	pkg := model.Package{}

	pkg.ID = uint64(genInt(counter) + 1)

	pkg.Weight = sql.NullInt64{
		Int64: gofakeit.Int64(),
		Valid: gofakeit.Bool(),
	}

	if gofakeit.Bool() {
		pkg.Title = gofakeit.ProductName()
	}

	pkg.Updated = sql.NullTime{Time: time.Now(), Valid: true}

	err := pkgService.Update(&pkg)
	if err != nil {
		log.Error("FAKE fail to update package", "error", err)
	}
}
