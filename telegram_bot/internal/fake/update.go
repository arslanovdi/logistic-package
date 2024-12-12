package fake

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/brianvoe/gofakeit/v7"
)

// Генерирует запрос изменения данных в пакете
func (f *Faker) eUpdate() {
	if f.counter <= 2 { // отбрасываем запросы, пока не создано хотя-бы 2 пакета
		return
	}

	log := slog.With("func", "fake.eUpdate")
	pkg := model.Package{}

	pkg.ID = genInt(f.counter) + 1 // меняет случайный пакет

	pkg.Weight = sql.NullInt64{
		Int64: gofakeit.Int64(),
		Valid: gofakeit.Bool(), // случайное изменение поля weight
	}

	if gofakeit.Bool() { // случайное изменение поля title
		pkg.Title = gofakeit.ProductName()
	}

	pkg.Updated = sql.NullTime{Time: time.Now(), Valid: true}

	err := f.pkgService.Update(&pkg)
	if err != nil {
		log.Error("FAKE fail to update package", "error", err)
	}
}
