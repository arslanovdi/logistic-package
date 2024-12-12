package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel/trace"

	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/general"
	"github.com/arslanovdi/logistic-package/pkg/ctxutil"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/jackc/pgx/v5"
)

// Get - get package by id from database
func (r *Repo) Get(ctx context.Context, id uint64) (*model.Package, error) {
	log := slog.With("func", "postgres.Get")

	span := trace.SpanContextFromContext(ctx)
	if span.IsSampled() {
		log.With("trace_id", span.TraceID().String()) // insert traceid to log
	}

	// сборка запроса - query
	query, args, err1 := psql.Select("*").
		From("package").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err1 != nil {
		return nil, fmt.Errorf("postgres.Get: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	ctx = ctxutil.Detach(ctx) // Отвязать таймер в контексте

	// выполнить запрос
	rows, _ := r.dbpool.Query(ctx, query, args...) //nolint:errcheck    // Ошибка игнорируется, так как она обрабатывается в CollectOneRow
	defer rows.Close()

	pkg, err2 := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Package]) // десериализовать в структуру
	if err2 != nil {
		if errors.Is(err2, pgx.ErrNoRows) {
			return nil, general.ErrNotFound
		}
		return nil, fmt.Errorf("postgres.Get: %w", err2)
	}

	log.Debug("Get", slog.Any("pkg", pkg))

	return &pkg, nil
}
