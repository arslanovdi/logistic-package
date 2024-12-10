package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/general"
	"github.com/arslanovdi/logistic-package/pkg/ctxutil"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

// Update - update package by id in database
func (r *Repo) Update(ctx context.Context, pkg *model.Package) error {

	log := slog.With("func", "postgres.Update")

	traceid := ""
	span := trace.SpanContextFromContext(ctx)
	if span.IsSampled() {
		traceid = span.TraceID().String()
	}

	query, args, err1 := psql.Update("package").
		Set("weight", pkg.Weight).
		Set("title", pkg.Title).
		Set("updated", pkg.Updated).
		Where(sq.Eq{"id": pkg.ID}).
		Suffix("RETURNING created, removed").
		ToSql()
	if err1 != nil {
		return fmt.Errorf("postgres.Update: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	ctx = ctxutil.Detach(ctx)

	err2 := pgx.BeginFunc(ctx, r.dbpool, func(tx pgx.Tx) error {

		err := tx.QueryRow(ctx, query, args...).Scan(&pkg.Created, &pkg.Removed)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return general.ErrNotFound
			}
			return err
		}

		pkgJSON, err := json.Marshal(pkg)
		if err != nil {
			return err
		}

		pi := psql.Insert("package_events")
		if span.IsSampled() {
			pi = pi.Columns("package_id", "type", "payload", "traceid").
				Values(pkg.ID, model.Updated, pkgJSON, traceid)
		} else {
			pi = pi.Columns("package_id", "type", "payload").
				Values(pkg.ID, model.Updated, pkgJSON)
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
		return fmt.Errorf("postgres.Update: %w", err2)
	}

	log.Debug("package updated", slog.String("package", pkg.String()))

	return nil
}
