package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/general"
	"github.com/arslanovdi/logistic-package/pkg/ctxutil"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

// Delete - delete package by id in database
func (r *Repo) Delete(ctx context.Context, id uint64) error {

	log := slog.With("func", "postgres.Delete")

	traceid := ""
	span := trace.SpanContextFromContext(ctx)
	if span.IsSampled() {
		traceid = span.TraceID().String()
	}

	query, args, err1 := psql.Delete("package").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err1 != nil {
		return fmt.Errorf("postgres.Delete: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	pi := psql.Insert("package_events")

	if span.IsSampled() {
		pi = pi.Columns("package_id", "type", "traceid").
			Values(id, model.Removed, traceid)
	} else {
		pi = pi.Columns("package_id", "type").
			Values(id, model.Removed)
	}

	queryEvent, argsEvent, err2 := pi.ToSql()

	if err2 != nil {
		return fmt.Errorf("postgres.Delete: %w", err2)
	}

	log.Debug("queryEvent", slog.String("query", queryEvent), slog.Any("args", argsEvent))

	ctx = ctxutil.Detach(ctx)

	err3 := pgx.BeginFunc(ctx, r.dbpool, func(tx pgx.Tx) error {

		tag, err := tx.Exec(ctx, query, args...)

		if err != nil {
			return err
		}

		if tag.RowsAffected() == 0 { // Получаем количество обновленных строк
			return general.ErrNotFound
		}

		_, err = tx.Exec(ctx, queryEvent, argsEvent...)

		if err != nil {
			return err
		}

		return nil
	})

	if err3 != nil {
		return fmt.Errorf("postgres.Delete: %w", err3)
	}

	log.Debug("Package deleted", slog.Any("id", id))

	return nil
}
