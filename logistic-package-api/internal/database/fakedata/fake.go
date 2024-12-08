// Package fakedata - создает фэйковые записи в БД
package fakedata

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/service"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/brianvoe/gofakeit/v7"
	"time"
)

// Generate - создает count фэйковых записей в БД
func Generate(count int, repo service.Repo) {
	for i := 0; i < count; i++ {
		id, err := repo.Create(context.Background(), model.Package{
			Title: gofakeit.ProductName(),
			Weight: sql.NullInt64{
				Int64: int64(gofakeit.Uint32()),
				Valid: gofakeit.Bool(),
			},
			Created: gofakeit.DateRange(time.Now().AddDate(0, 0, -2), time.Now()),
			Updated: sql.NullTime{
				Time:  gofakeit.DateRange(time.Now().AddDate(0, 0, -1), time.Now()),
				Valid: gofakeit.Bool(),
			},
		})
		if err != nil {
			fmt.Printf("Generate id: %d, err: %s ", id, err.Error())
		}
	}
}
