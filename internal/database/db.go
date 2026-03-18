package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Connect opens a PostgreSQL connection pool using the provided connection URL,
// verifies connectivity with a Ping, and configures sensible pool defaults.
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("database: failed to open connection: %w", err)
	}

	// Verify the connection is reachable before returning.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database: failed to ping database: %w", err)
	}

	// Configure connection pool settings.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
