package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

// Create - create new package in database
func (r *Repo) Create(ctx context.Context, pkg *model.Package) (*uint64, error) {
	log := slog.With("func", "postgres.Create")

	traceid := ""
	span := trace.SpanContextFromContext(ctx)
	if span.IsSampled() {
		traceid = span.TraceID().String()
		log.With("trace_id", traceid) // insert traceid to log
	}

	// сборка первого запроса - query
	query, args, err1 := psql.Insert("package").
		Columns("weight", "title", "created").
		Values(pkg.Weight, pkg.Title, pkg.Created).
		Suffix("RETURNING id").
		ToSql()
	if err1 != nil {
		return nil, fmt.Errorf("postgres.Create: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	ctx = context.WithoutCancel(ctx) // Отвязать контекст, нужно завершить операцию, даже если клиент отвалился

	err2 := pgx.BeginFunc(ctx, r.dbpool, func(tx pgx.Tx) error { // Запуск транзакции, автоматический rollback при ошибке
		err := tx.QueryRow(ctx, query, args...).Scan(&pkg.ID) // выполнить первый запрос
		if err != nil {
			return err
		}

		// сборка второго запроса - queryEvent
		pkgJSON, err := json.Marshal(pkg)
		if err != nil {
			return err
		}

		pi := psql.Insert("package_events")

		if span.IsSampled() { // insert traceid to package_events if it exists
			pi = pi.Columns("package_id", "type", "payload", "traceid").
				Values(pkg.ID, model.Created, pkgJSON, traceid)
		} else {
			pi = pi.Columns("package_id", "type", "payload").
				Values(pkg.ID, model.Created, pkgJSON)
		}

		queryEvent, argsEvent, err := pi.ToSql()
		if err != nil {
			return err
		}

		log.Debug("queryEvent", slog.String("query", queryEvent), slog.Any("args", argsEvent))

		_, err = tx.Exec(ctx, queryEvent, argsEvent...) // выполнить второй запрос
		if err != nil {
			return err
		}

		return nil
	})

	if err2 != nil {
		return nil, fmt.Errorf("postgres.Create: %w", err2)
	}

	log.Debug("package created", slog.String("package", pkg.String()))

	return &pkg.ID, nil
}
