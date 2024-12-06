package postgres

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/model"
	"github.com/arslanovdi/logistic-package/pkg/ctxutil"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

// List - Get packages from database. start index - offset, count - limit
func (r *Repo) List(ctx context.Context, offset uint64, limit uint64) ([]model.Package, error) {

	log := slog.With("func", "postgres.List")

	query, args, err1 := psql.Select("*").
		From("package").
		Where(sq.GtOrEq{"id": offset}).
		//Where(sq.Lt{"id": offset + limit}).
		OrderBy("id ASC").
		Limit(limit).
		ToSql()

	if err1 != nil {
		return nil, fmt.Errorf("postgres.List: %w", err1)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	ctx = ctxutil.Detach(ctx)

	rows, _ := r.dbpool.Query(ctx, query, args...) // Ошибка игнорируется, так как она обрабатывается в CollectRows
	defer rows.Close()

	var packages []model.Package
	packages, err2 := pgx.CollectRows(rows, pgx.RowToStructByName[model.Package]) // десериализуем в слайс структуру
	if err2 != nil {
		if errors.Is(err2, pgx.ErrNoRows) {
			log.Debug("no rows found")
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("postgres.List: %w", err2)
	}

	log.Debug("packages listed", slog.Uint64("offset", offset), slog.Uint64("limit", limit))

	return packages, nil
}
