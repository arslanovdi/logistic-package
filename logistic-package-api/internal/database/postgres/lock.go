package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/general"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

// Lock заблокировать в БД n записей событий
func (r *Repo) Lock(ctx context.Context, n int) ([]model.PackageEvent, error) {

	log := slog.With("func", "postgres.Lock")

	// squirell не поддерживает CTE
	// использую "github.com/huandu/go-sqlbuilder", тоже не без замечаний.

	/*
		WITH lockedevents AS (
		    WITH eventstosend AS (
		        SELECT id
		        FROM package_events
		        WHERE (status <> $1 OR status IS NULL)
		        ORDER BY id ASC
		        LIMIT 10
		        )
		    UPDATE package_events
		        SET status = $2
		        FROM eventstosend
		        WHERE package_events.id = eventstosend.id
		        RETURNING package_events.*
		)
		SELECT *
		FROM lockedevents
		ORDER BY package_id ASC, id ASC
	*/

	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL

	eventsBuilder := sqlbuilder.Select("id").From("package_events")
	eventsBuilder.Where(eventsBuilder.Or(eventsBuilder.NotEqual("status", model.Locked),
		eventsBuilder.IsNull("status"))).
		OrderBy("id ASC").
		Limit(n)

	cteu := sqlbuilder.With(
		sqlbuilder.CTETable("eventstosend").As(eventsBuilder),
	)

	ub := sqlbuilder.NewUpdateBuilder()
	ub.With(cteu)
	ub.Update("package_events").
		Set(ub.Assign("status", model.Locked)).
		Where("package_events.id = eventstosend.id").
		SQL("RETURNING package_events.*")

	ctes := sqlbuilder.With(sqlbuilder.CTETable("lockedevents").As(ub))
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*").
		With(ctes).
		OrderBy("package_id ASC, id ASC")

	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL) // сборка запроса для PostgreSQL

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	var events []model.PackageEvent

	err2 := pgx.BeginTxFunc(ctx, r.dbpool, pgx.TxOptions{ /*IsoLevel: "serializable"*/ }, func(tx pgx.Tx) error {

		rows, _ := tx.Query(ctx, query, args...) // Ошибка игнорируется, так как она обрабатывается в CollectRowsr

		defer rows.Close()

		var err error
		events, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.PackageEvent]) // десериализуем в слайс структуру
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Debug("no rows found")
				return general.ErrNotFound
			}
			return err
		}

		return nil
	})

	if err2 != nil {
		return nil, fmt.Errorf("postgres.Lock: %w", err2)
	}

	return events, nil
}
