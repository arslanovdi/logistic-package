package fake

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/brianvoe/gofakeit/v7"
)

// Генерирует запрос создания пакета
func (f *Faker) eCreate() {
	log := slog.With("func", "fake.eCreate")
	pkg := model.Package{}

	pkg.Weight = sql.NullInt64{
		Int64: gofakeit.Int64(),
		Valid: gofakeit.Bool(),
	}

	pkg.Title = gofakeit.ProductName()
	pkg.Created = gofakeit.DateRange(time.Now().AddDate(0, -2, 0), time.Now())

	_, err := f.pkgService.Create(&pkg)
	if err != nil {
		log.Error("FAKE fail to create package", "error", err)
	}

	f.counter++ // счетчик пакетов
}
