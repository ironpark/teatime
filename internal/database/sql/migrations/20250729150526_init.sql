-- +goose Up
-- +goose StatementBegin

-- Enable foreign key enforcement
PRAGMA foreign_keys = ON;

CREATE TABLE recipes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    recipe_path TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Executions table for workflow execution history
CREATE TABLE executions (
    id TEXT PRIMARY KEY,
    recipe_id TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('running', 'success', 'failed', 'cancelled')),
    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME,
    duration_ms INTEGER,
    error_message TEXT NOT NULL DEFAULT '',
    input_data TEXT NOT NULL DEFAULT '{}',
    output_data TEXT NOT NULL DEFAULT '{}',
    execution_context TEXT NOT NULL DEFAULT '{}',
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
);


-- Credentials table for secure storage of API keys and secrets
CREATE TABLE credentials (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL,
    description TEXT,
    keychain_key TEXT,
    encrypted_data BLOB,
    salt BLOB,
    storage_type TEXT NOT NULL CHECK (storage_type IN ('keychain', 'database', 'environment')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at DATETIME
);

-- Settings table for application preferences
CREATE TABLE settings (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    theme TEXT NOT NULL DEFAULT 'system' CHECK (theme IN ('system', 'light', 'dark')),
    auto_start INTEGER NOT NULL DEFAULT 0 CHECK (auto_start IN (0, 1)),
    language TEXT NOT NULL DEFAULT 'en',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert default settings row
INSERT INTO settings (id, theme, auto_start, language) VALUES (1, 'system', 0, 'en');

-- Indexes for performance
CREATE INDEX idx_executions_recipe_id ON executions(recipe_id);
CREATE INDEX idx_executions_status ON executions(status);
CREATE INDEX idx_executions_started_at ON executions(started_at);
CREATE INDEX idx_credentials_name ON credentials(name);
CREATE INDEX idx_credentials_type ON credentials(type);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_credentials_type;
DROP INDEX IF EXISTS idx_credentials_name;
DROP INDEX IF EXISTS idx_executions_started_at;
DROP INDEX IF EXISTS idx_executions_status;
DROP INDEX IF EXISTS idx_executions_recipe_id;

DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS credentials;
DROP TABLE IF EXISTS executions;
DROP TABLE IF EXISTS recipes;

-- +goose StatementEnd
