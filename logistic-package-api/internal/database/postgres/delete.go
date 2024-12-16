package postgres

import (
	"context"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/general"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

// Delete - delete package by id in database
func (r *Repo) Delete(ctx context.Context, id uint64) error {
	log := slog.With("func", "postgres.Delete")

	traceid := ""
	span := trace.SpanContextFromContext(ctx)
	if span.IsSampled() {
		traceid = span.TraceID().String()
		log.With("trace_id", traceid) // insert traceid to log
	}

	// сборка первого запроса - query
	query, args, err1 := psql.Delete("package").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err1 != nil {
		return fmt.Errorf("postgres.Delete: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	// сборка второго запроса - queryEvent
	pi := psql.Insert("package_events")

	if span.IsSampled() { // insert traceid to package_events if it exists
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

	ctx = context.WithoutCancel(ctx) // Отвязать контекст, нужно завершить операцию, даже если клиент отвалился

	err3 := pgx.BeginFunc(ctx, r.dbpool, func(tx pgx.Tx) error { // Запуск транзакции, автоматический rollback при ошибке
		tag, err := tx.Exec(ctx, query, args...) // выполнить первый запрос
		if err != nil {
			return err
		}

		if tag.RowsAffected() == 0 { // Получаем количество затронутых строк
			return general.ErrNotFound
		}

		_, err = tx.Exec(ctx, queryEvent, argsEvent...) // выполнить второй запрос
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
