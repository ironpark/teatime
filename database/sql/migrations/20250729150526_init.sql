-- +goose Up
-- +goose StatementBegin

CREATE TABLE recipes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    recipe_path TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Executions table for workflow execution history
CREATE TABLE executions (
    id TEXT PRIMARY KEY,
    recipe_id TEXT NOT NULL,
    FOREIGN KEY (recipe_id) REFERENCES recipes(id),
    status TEXT NOT NULL CHECK (status IN ('running', 'success', 'failed', 'cancelled')),
    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME,
    duration_ms INTEGER,
    error_message TEXT,
    input_data TEXT, -- JSON string
    output_data TEXT, -- JSON string
    execution_context TEXT -- JSON string for additional context
);

-- Credentials table for secure storage of API keys and secrets
CREATE TABLE credentials (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL, -- 'keychain', 'encrypted', 'env'
    description TEXT,
    keychain_key TEXT, -- OS keychain reference key
    encrypted_data BLOB, -- AES encrypted credential data
    salt BLOB, -- Salt for encryption
    storage_type TEXT NOT NULL CHECK (storage_type IN ('keychain', 'database', 'environment')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at DATETIME
);

-- Indexes for performance
CREATE INDEX idx_executions_workflow_path ON executions(workflow_path);
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
DROP INDEX IF EXISTS idx_executions_workflow_path;

DROP TABLE IF EXISTS credentials;
DROP TABLE IF EXISTS executions;

-- +goose StatementEnd
