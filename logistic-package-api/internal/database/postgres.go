// Package database - contains all database functions
package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
Postgresql Error Codes
https://www.postgresql.org/docs/16/errcodes-appendix.html
*/

// MustGetPgxPool get pgxpool or os.Exit(1)
func MustGetPgxPool(ctx context.Context) *pgxpool.Pool {
	log := slog.With("func", "database.MustGetPgxPool")

	dbpool, err1 := NewPgxPool(ctx)
	if err1 != nil {
		log.Warn("Failed init postgres", slog.String("error", err1.Error()))
		os.Exit(1)
	}

	return dbpool
}

// NewPgxPool create new pgxpool of connections to postgres
func NewPgxPool(ctx context.Context) (*pgxpool.Pool, error) {
	log := slog.With("func", "database.NewPgxPool")

	cfg := config.GetConfigInstance()

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v", // sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		// cfg.Database.Ssl,
	)

	// Эти параметры можно также задать в DSN
	/* дефолтные значения:
	   pool_max_conn_lifetime = time.Hour
	   pool_max_conn_idle_time = time.Minute * 30
	   pool_health_check_period = time.Minute
	   pool_max_conns = greater of 4 or runtime.NumCPU() если ядер больше 4	*/

	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("database.NewPgxPool: %w", err)
	}

	pgxConfig.ConnConfig.Tracer = otelpgx.NewTracer() // Добавляем OpenTelemetry трассировку для postgresql

	dbpool, err1 := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err1 != nil {
		log.Warn("Error connecting to the database", slog.String("error", err1.Error()))
		return nil, fmt.Errorf("database.NewPgxPool: %w", err1)
	}

	err2 := dbpool.Ping(ctx) // эта команда заменяет acquire + ping
	if err2 != nil {
		log.Warn("Could not ping database", slog.String("error", err2.Error()))
		return nil, fmt.Errorf("database.NewPgxPool: %w", err2)
	}

	log.Info("successfully connected to database")
	return dbpool, nil
}
