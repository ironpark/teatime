package database

import (
	"context"
	"database/sql"
	"log"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/ironpark/teatime/database/sql/migrations"
)

type Client struct {
	*Queries
	DB *sql.DB
}

// OpenOption defines options for opening a database
type OpenOption func(*openConfig)

type openConfig struct {
	inMemory bool
	dbPath   string
}

// WithInMemory sets the database to use in-memory mode
func WithInMemory() OpenOption {
	return func(c *openConfig) {
		c.inMemory = true
	}
}

// Open creates a new database client with options
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

// OpenInMemory is a convenience function for creating an in-memory database
func OpenInMemory() (*Client, error) {
	return Open(context.Background(), "", WithInMemory())
}

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

func (c *Client) Close() error {
	return c.DB.Close()
}

// GetMigrationVersion returns the current migration version
func (c *Client) GetMigrationVersion() (int64, error) {
	return migrations.GetVersion(c.DB)
}

// GetMigrationStatus prints the current migration status
func (c *Client) GetMigrationStatus() error {
	return migrations.GetMigrationStatus(c.DB)
}

// ResetDatabase drops all tables and runs migrations from scratch
func (c *Client) ResetDatabase() error {
	return migrations.ResetDatabase(c.DB)
}
