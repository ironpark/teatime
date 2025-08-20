package database

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/ironpark/teatime/internal/database/sql/migrations"

	_ "modernc.org/sqlite"
)

// Client wraps database connections with query methods and transaction support.
// It embeds Queries for direct access to generated database operations.
type Client struct {
	*Queries
	DB *sql.DB
}

// OpenOption defines configuration options for opening a database connection.
type OpenOption func(*openConfig)

type openConfig struct {
	inMemory bool
	dbPath   string
}

// WithInMemory configures the database to use in-memory SQLite mode.
// Useful for testing and temporary data storage.
func WithInMemory() OpenOption {
	return func(c *openConfig) {
		c.inMemory = true
	}
}

// Open creates a new database client with the specified options.
// It automatically runs migrations and configures SQLite with foreign keys enabled.
func Open(ctx context.Context, dbPath string, opts ...OpenOption) (*Client, error) {
	config := &openConfig{
		dbPath: dbPath,
	}

	for _, opt := range opts {
		opt(config)
	}

	var connStr string
	if config.inMemory {
		connStr = ":memory:?cache=shared&_fk=1"
		log.Println("Open in-memory database")
	} else {
		connStr = config.dbPath
		if !strings.HasSuffix(connStr, "?cache=shared&_fk=1") {
			connStr += "?cache=shared&_fk=1"
		}
		log.Println("Open database", config.dbPath)
	}

	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := migrations.RunMigrations(db); err != nil {
		db.Close()
		return nil, err
	}

	return &Client{
		Queries: New(db),
		DB:      db,
	}, nil
}

// OpenInMemory is a convenience function for creating an in-memory database.
// It's equivalent to calling Open with WithInMemory() option.
func OpenInMemory() (*Client, error) {
	return Open(context.Background(), "", WithInMemory())
}

// WithTx executes a function within a database transaction.
// The transaction is automatically rolled back if the function returns an error.
func (c *Client) WithTx(ctx context.Context, fn func(ctx context.Context, queries *Queries) error) error {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	queries := c.Queries.WithTx(tx)
	if err := fn(ctx, queries); err != nil {
		return err
	}
	return tx.Commit()
}

// Close closes the database connection.
func (c *Client) Close() error {
	return c.DB.Close()
}

// GetMigrationVersion returns the current database schema migration version.
func (c *Client) GetMigrationVersion() (int64, error) {
	return migrations.GetVersion(c.DB)
}

// GetMigrationStatus displays the current database migration status and history.
func (c *Client) GetMigrationStatus() error {
	return migrations.GetMigrationStatus(c.DB)
}

// ResetDatabase drops all tables and re-runs all migrations from the beginning.
// WARNING: This will destroy all existing data.
func (c *Client) ResetDatabase() error {
	return migrations.ResetDatabase(c.DB)
}
