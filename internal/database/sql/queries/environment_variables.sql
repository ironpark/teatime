-- Environment Variables queries

-- name: CreateEnvironmentVariable :one
INSERT INTO environment_variables (id, name, value, description)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetEnvironmentVariable :one
SELECT * FROM environment_variables 
WHERE id = ?;

-- name: GetEnvironmentVariableByName :one
SELECT * FROM environment_variables 
WHERE name = ?;

-- name: ListEnvironmentVariables :many
SELECT * FROM environment_variables 
ORDER BY name;

-- name: UpdateEnvironmentVariable :one
UPDATE environment_variables 
SET name = ?, value = ?, description = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteEnvironmentVariable :exec
DELETE FROM environment_variables 
WHERE id = ?;

-- name: CheckEnvironmentVariableNameExists :one
SELECT COUNT(*) FROM environment_variables 
WHERE name = ?;

-- name: SearchEnvironmentVariablesByName :many
SELECT * FROM environment_variables 
WHERE name LIKE ?
ORDER BY name;