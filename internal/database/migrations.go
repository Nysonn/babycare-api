package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

// RunMigrations applies all pending goose migrations from the provided directory.
func RunMigrations(db *sql.DB, migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("database: set goose dialect: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("database: apply migrations: %w", err)
	}

	return nil
}
