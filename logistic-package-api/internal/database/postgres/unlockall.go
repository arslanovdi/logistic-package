package postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

// UnlockAll разблокировать в БД все записи
// используется при инициализации ретранслятора
func (r *Repo) UnlockAll(ctx context.Context) error {

	log := slog.With("func", "postgres.UnlockAll")

	query, args, err1 := psql.Update("package_events").
		Set("status", model.Unlocked).
		Where(sq.NotEq{"status": model.Unlocked}).
		Suffix("RETURNING id").
		ToSql()

	if err1 != nil {
		return err1
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	rows, _ := r.dbpool.Query(ctx, query, args...)
	defer rows.Close()

	events, err2 := pgx.CollectRows(rows, pgx.RowTo[int64]) // десериализуем в слайс
	if err2 != nil {
		return err2
	}

	if len(events) > 0 {
		log.Warn("Found Locked Events", slog.Any("event_id", events))
	}

	return nil
}
