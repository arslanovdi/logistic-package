// Package fake эмуляция работы пользователя телеграм бота
package fake

import (
	"crypto/rand"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/model"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/service"
	"github.com/brianvoe/gofakeit/v7"
	"log/slog"
	"math/big"
	"time"
)

const (
	create = 30
	del    = 10
	get    = 20
	list   = 20
	update = 20
)

var counter uint64

// genInt генерация случайного числа от 0 до max при помощи пакета crypto/rand
func genInt(max uint64) uint64 {
	rnd, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		slog.Error("FAKE fail to generate random number", slog.String("error", err.Error()))
		return 0
	}
	return rnd.Uint64()
}

// Emulate эмуляция работы пользователя телеграм бота
func Emulate(d uint64, pkgService *service.LogisticPackageService) {
	for {
		time.Sleep(time.Duration(genInt(d)) * time.Millisecond) // от 200ms до 1200 мс на одну операцию)

		rnd := genInt(100) // 100%
		switch {
		case rnd < create: // create % созданий пакета
			pkg := model.Package{}

			if gofakeit.Bool() {
				pkg.Weight = new(uint64)
				*pkg.Weight = gofakeit.Uint64()
			}

			pkg.Title = gofakeit.ProductName()
			pkg.Created = gofakeit.DateRange(time.Now().AddDate(0, -2, 0), time.Now())

			_, err := pkgService.Create(pkg)
			if err != nil {
				slog.Error("FAKE fail to create package", "error", err)
			}
			counter++
		case rnd < create+del: // del % удаления пакета

			if counter <= 2 {
				continue
			}

			id := genInt(counter) + 1

			err := pkgService.Delete(id)
			if err != nil {
				slog.Error("FAKE fail to delete package", "error", err)
			}
			counter--
		case rnd < create+del+get: // get% получения пакета

			if counter <= 2 {
				continue
			}
			id := genInt(counter) + 1

			_, err := pkgService.Get(id)
			if err != nil {
				slog.Error("FAKE fail to get package", "id", id)
			}
		case rnd < create+del+get+list: // list % получения списка пакетов
			if counter < 10 {
				continue
			}
			offset := genInt(counter/2) + 1
			limit := genInt(counter - offset)
			_, err := pkgService.List(offset, limit)
			if err != nil {
				slog.Error("FAKE fail to list package", "error", err)
			}

		case rnd < create+del+get+list+update: // update % обновления пакета

			if counter <= 2 {
				continue
			}

			pkg := model.Package{}

			pkg.ID = genInt(counter) + 1

			if gofakeit.Bool() {
				pkg.Weight = new(uint64)
				*pkg.Weight = gofakeit.Uint64()
			}

			if gofakeit.Bool() {
				pkg.Title = gofakeit.ProductName()
			}

			pkg.Updated = new(time.Time)
			*pkg.Updated = time.Now()

			err := pkgService.Update(pkg)
			if err != nil {
				slog.Error("FAKE fail to update package", "error", err)
			}
		}
	}
}
