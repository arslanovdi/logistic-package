// Package postgres - Postgres implementation of service.Repo and repo.EventRepo
package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar) // Squirell плэйсхолдер для Postgres

// Repo - Postgres implementation of service.Repo and repo.EventRepo
type Repo struct {
	dbpool *pgxpool.Pool
}

// NewPostgresRepo returns Postgres implementation of service.Repo and repo.EventRepo
func NewPostgresRepo(dbpool *pgxpool.Pool) *Repo {
	return &Repo{
		dbpool: dbpool,
	}
}
