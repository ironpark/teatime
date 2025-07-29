-- name: CreateExecution :one
INSERT INTO executions (
    id, workflow_path, workflow_name, status, started_at, input_data, execution_context
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetExecution :one
SELECT * FROM executions WHERE id = ?;

-- name: GetExecutionsByWorkflow :many
SELECT * FROM executions 
WHERE workflow_path = ?
ORDER BY started_at DESC
LIMIT ?;

-- name: GetExecutionsByStatus :many
SELECT * FROM executions 
WHERE status = ?
ORDER BY started_at DESC
LIMIT ?;

-- name: GetRecentExecutions :many
SELECT * FROM executions 
ORDER BY started_at DESC 
LIMIT ?;

-- name: UpdateExecutionStatus :exec
UPDATE executions 
SET status = ?, finished_at = ?, duration_ms = ?, error_message = ?, output_data = ?
WHERE id = ?;

-- name: UpdateExecutionFinished :exec
UPDATE executions 
SET status = ?, finished_at = CURRENT_TIMESTAMP, duration_ms = ?, output_data = ?
WHERE id = ?;

-- name: UpdateExecutionError :exec
UPDATE executions 
SET status = 'failed', finished_at = CURRENT_TIMESTAMP, duration_ms = ?, error_message = ?
WHERE id = ?;

-- name: DeleteExecution :exec
DELETE FROM executions WHERE id = ?;

-- name: DeleteOldExecutions :exec
DELETE FROM executions 
WHERE started_at < datetime('now', '-' || ? || ' days');

-- name: CountExecutionsByStatus :one
SELECT COUNT(*) FROM executions WHERE status = ?;

-- name: GetExecutionStats :one
SELECT 
    COUNT(*) as total_count,
    COUNT(CASE WHEN status = 'success' THEN 1 END) as success_count,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_count,
    COUNT(CASE WHEN status = 'running' THEN 1 END) as running_count,
    AVG(CASE WHEN duration_ms IS NOT NULL THEN duration_ms END) as avg_duration_ms
FROM executions
WHERE workflow_path = ?;

-- Credentials queries
-- name: CreateCredential :one
INSERT INTO credentials (
    id, name, type, description, keychain_key, encrypted_data, salt, storage_type
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetCredential :one
SELECT * FROM credentials WHERE id = ?;

-- name: GetCredentialByName :one
SELECT * FROM credentials WHERE name = ?;

-- name: ListCredentials :many
SELECT id, name, type, description, storage_type, created_at, updated_at, last_used_at 
FROM credentials
ORDER BY name;

-- name: UpdateCredential :exec
UPDATE credentials 
SET type = ?, description = ?, keychain_key = ?, encrypted_data = ?, salt = ?, storage_type = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdateCredentialLastUsed :exec
UPDATE credentials 
SET last_used_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteCredential :exec
DELETE FROM credentials WHERE id = ?;

-- name: GetCredentialsByType :many
SELECT * FROM credentials 
WHERE type = ?
ORDER BY name;

-- name: GetCredentialsByStorageType :many
SELECT * FROM credentials 
WHERE storage_type = ?
ORDER BY name;

-- name: CheckCredentialExists :one
SELECT COUNT(*) FROM credentials WHERE name = ?;