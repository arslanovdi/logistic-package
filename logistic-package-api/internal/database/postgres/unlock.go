package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/model"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/server"
	"log/slog"
)

// Unlock разблокировать в БД n записей событий
func (r *Repo) Unlock(ctx context.Context, eventID []int64) error {

	log := slog.With("func", "postgres.Unlock")

	query, args, err1 := psql.Update("package_events").
		Set("status", model.Unlocked).
		Where(sq.Eq{"id": eventID}).
		ToSql()

	if err1 != nil {
		return fmt.Errorf("postgres.Unlock: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	tag, err2 := r.dbpool.Exec(ctx, query, args...)
	if err2 != nil {
		return fmt.Errorf("postgres.Unlock: %w", err2)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	server.RetranslatorEvents.Sub(float64(len(eventID))) // TODO метрика, кол-во обрабатываемых событий, убавляем счетчик

	return nil
}
