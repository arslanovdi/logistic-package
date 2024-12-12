package postgres

import (
	"context"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/general"
)

// Remove удалить из БД n записей событий
func (r *Repo) Remove(ctx context.Context, eventIDs []int64) error {
	log := slog.With("func", "postgres.Remove")

	// сборка запроса - query
	query, args, err1 := psql.Delete("package_events").
		Where(sq.Eq{"id": eventIDs}).
		ToSql()

	if err1 != nil {
		return fmt.Errorf("postgres.Remove: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	tag, err2 := r.dbpool.Exec(ctx, query, args...) // выполнить запрос
	if err2 != nil {
		return fmt.Errorf("postgres.Remove: %w", err2)
	}

	if tag.RowsAffected() == 0 {
		return general.ErrNotFound
	}

	return nil
}
