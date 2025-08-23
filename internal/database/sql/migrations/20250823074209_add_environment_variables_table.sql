-- +goose Up
-- +goose StatementBegin

-- Environment variables table for non-sensitive configuration data
CREATE TABLE environment_variables (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    value TEXT NOT NULL,
    description TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for performance
CREATE INDEX idx_environment_variables_name ON environment_variables(name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_environment_variables_name;
DROP TABLE IF EXISTS environment_variables;

-- +goose StatementEnd
