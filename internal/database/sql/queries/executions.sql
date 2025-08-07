-- name: GetExecution :one
SELECT * FROM executions
WHERE id = ? LIMIT 1;

-- name: ListExecutions :many
SELECT * FROM executions
ORDER BY started_at DESC;

-- name: ListExecutionsByRecipe :many
SELECT * FROM executions
WHERE recipe_id = ?
ORDER BY started_at DESC;

-- name: ListExecutionsByStatus :many
SELECT * FROM executions
WHERE status = ?
ORDER BY started_at DESC;

-- name: CreateExecution :one
INSERT INTO executions (
    id, recipe_id, status, input_data, execution_context
) VALUES (
    ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateExecutionStatus :one
UPDATE executions
SET 
    status = ?1,
    finished_at = CASE 
        WHEN ?1 IN ('success', 'failed', 'cancelled') THEN CURRENT_TIMESTAMP 
        ELSE finished_at 
    END,
    duration_ms = CASE 
        WHEN ?1 IN ('success', 'failed', 'cancelled') 
        THEN CAST((julianday(CURRENT_TIMESTAMP) - julianday(started_at)) * 86400000 AS INTEGER)
        ELSE duration_ms 
    END,
    error_message = ?2,
    output_data = ?3
WHERE id = ?4
RETURNING *;

-- name: DeleteExecution :exec
DELETE FROM executions
WHERE id = ?;

-- name: GetRunningExecutions :many
SELECT * FROM executions
WHERE status = 'running'
ORDER BY started_at DESC;

-- name: GetExecutionsByRecipeAndStatus :many
SELECT * FROM executions
WHERE recipe_id = ? AND status = ?
ORDER BY started_at DESC;

-- name: GetRecentExecutions :many
SELECT * FROM executions
WHERE started_at >= datetime('now', '-7 days')
ORDER BY started_at DESC
LIMIT ?;

-- name: GetExecutionStats :one
SELECT 
    COUNT(*) as total_executions,
    COUNT(CASE WHEN status = 'success' THEN 1 END) as successful_executions,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_executions,
    COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled_executions,
    COUNT(CASE WHEN status = 'running' THEN 1 END) as running_executions,
    AVG(CASE WHEN status = 'success' THEN duration_ms END) as avg_success_duration_ms
FROM executions
WHERE recipe_id = ?;

-- name: CleanupOldExecutions :exec
DELETE FROM executions
WHERE finished_at < datetime('now', '-30 days')
AND status IN ('success', 'failed', 'cancelled');