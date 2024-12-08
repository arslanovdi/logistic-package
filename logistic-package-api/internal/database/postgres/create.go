package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arslanovdi/logistic-package/pkg/ctxutil"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

// Create - create new package in database
func (r *Repo) Create(ctx context.Context, pkg model.Package) (*uint64, error) {

	log := slog.With("func", "postgres.Create")

	traceid := ""
	span := trace.SpanContextFromContext(ctx)
	if span.IsSampled() {
		traceid = span.TraceID().String()
	}

	query, args, err1 := psql.Insert("package").
		Columns("weight", "title", "created").
		Values(pkg.Weight, pkg.Title, pkg.Created).
		Suffix("RETURNING id").
		ToSql()
	if err1 != nil {
		return nil, fmt.Errorf("postgres.Create: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	ctx = ctxutil.Detach(ctx)

	err2 := pgx.BeginFunc(ctx, r.dbpool, func(tx pgx.Tx) error { // Запускаем транзакцию
		err := tx.QueryRow(ctx, query, args...).Scan(&pkg.ID)
		//err := r.dbpool.QueryRow(ctx, query, args...).Scan(&pkg.ID)
		if err != nil {
			return err
		}
		pkgJSON, err := json.Marshal(pkg)
		if err != nil {
			return err
		}

		pi := psql.Insert("package_events")

		if span.IsSampled() {
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

		_, err = tx.Exec(ctx, queryEvent, argsEvent...)

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
