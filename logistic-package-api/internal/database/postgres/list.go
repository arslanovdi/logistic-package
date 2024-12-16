package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel/trace"

	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/general"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/jackc/pgx/v5"
)

// List - Get packages from database. start index - offset, count - limit
func (r *Repo) List(ctx context.Context, offset, limit uint64) ([]model.Package, error) {
	log := slog.With("func", "postgres.List")

	span := trace.SpanContextFromContext(ctx)
	if span.IsSampled() {
		log.With("trace_id", span.TraceID().String()) // insert traceid to log
	}

	// сборка запроса - query
	query, args, err1 := psql.Select("*").
		From("package").
		Where(sq.GtOrEq{"id": offset}).
		OrderBy("id ASC").
		Limit(limit).
		ToSql()

	if err1 != nil {
		return nil, fmt.Errorf("postgres.List: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	ctx = context.WithoutCancel(ctx) // Отвязать контекст, нужно завершить операцию, даже если клиент отвалился

	// выполнить запрос
	rows, _ := r.dbpool.Query(ctx, query, args...) //nolint:errcheck // Ошибка игнорируется, так как она обрабатывается в CollectRows
	defer rows.Close()

	var packages []model.Package
	packages, err2 := pgx.CollectRows(rows, pgx.RowToStructByName[model.Package]) // десериализовать в слайс структур
	if err2 != nil {
		if errors.Is(err2, pgx.ErrNoRows) {
			log.Debug("no rows found")
			return nil, general.ErrNotFound
		}
		return nil, fmt.Errorf("postgres.List: %w", err2)
	}

	log.Debug("packages listed", slog.Uint64("offset", offset), slog.Uint64("limit", limit))

	return packages, nil
}
