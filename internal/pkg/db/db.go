// Package db provides a database connection
package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
	"github.com/zeusWPI/scc/pkg/config"
)

// DB represents a database connection
type DB struct {
	Pool    *pgxpool.Pool
	Queries *sqlc.Queries
}

// New creates a new database connection
func New() (*DB, error) {
	pgConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, err
	}

	pgConfig.ConnConfig.Host = config.GetDefaultString("db.host", "localhost")
	pgConfig.ConnConfig.Port = config.GetDefaultUint16("db.port", 5432)
	pgConfig.ConnConfig.Database = config.GetDefaultString("db.database", "scc")
	pgConfig.ConnConfig.User = config.GetDefaultString("db.user", "postgres")
	pgConfig.ConnConfig.Password = config.GetDefaultString("db.password", "postgres")

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.TODO()); err != nil {
		return nil, err
	}

	queries := sqlc.New(pool)

	return &DB{Pool: pool, Queries: queries}, nil
}
