// Package migrations provides database schema migration management for Teatime.
//
// This package handles the complete lifecycle of database schema changes:
//   - Executing all pending migrations or individual migration steps
//   - Managing schema versions and tracking migration history
//   - Supporting incremental migrations for development and debugging
//   - Providing database reset functionality for testing environments
//
// All SQL migration files are embedded in the binary for easy deployment.
// The package uses the goose migration tool internally for reliable schema management.
package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
// embedMigrations contains all SQL migration files embedded in the binary.
var embedMigrations embed.FS

// RunMigrations executes all pending database migrations in order.
// It uses embedded SQL files and the goose migration tool.
func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}
	if err := goose.Up(db, "."); err != nil {
		return err
	}

	return nil
}

// ResetDatabase drops all tables and re-runs all migrations from the beginning.
// WARNING: This destroys all existing data in the database.
func ResetDatabase(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Reset(db, "."); err != nil {
		return err
	}

	if err := goose.Up(db, "."); err != nil {
		return err
	}

	return nil
}

// GetMigrationStatus displays the current migration status and history.
// It shows which migrations have been applied and which are pending.
func GetMigrationStatus(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.Status(db, ".")
}

// MigrateUp executes the next pending migration.
// Use this to apply migrations one at a time for testing or debugging.
func MigrateUp(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.UpByOne(db, ".")
}

// MigrateDown reverts the most recent migration.
// WARNING: This may cause data loss if the migration drops tables or columns.
func MigrateDown(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.Down(db, ".")
}

// GetVersion returns the current database schema version number.
// Returns 0 if no migrations have been applied yet.
func GetVersion(db *sql.DB) (int64, error) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return 0, err
	}

	return goose.GetDBVersion(db)
}
