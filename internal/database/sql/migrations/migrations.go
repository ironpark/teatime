package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

// RunMigrations runs all database migrations
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

// ResetDatabase drops all tables and runs migrations from scratch
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

// GetMigrationStatus returns the current migration status
func GetMigrationStatus(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.Status(db, ".")
}

// MigrateUp runs one migration up
func MigrateUp(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.UpByOne(db, ".")
}

// MigrateDown runs one migration down
func MigrateDown(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.Down(db, ".")
}

// GetVersion returns the current migration version
func GetVersion(db *sql.DB) (int64, error) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return 0, err
	}

	return goose.GetDBVersion(db)
}
