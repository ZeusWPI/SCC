// Package db provides a database connection
package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
)

// DB represents a database connection
type DB struct {
	DB      *sql.DB
	Queries *sqlc.Queries
}

// New creates a new database connection
func New() (*DB, error) {
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	queries := sqlc.New(db)

	return &DB{DB: db, Queries: queries}, nil
}
