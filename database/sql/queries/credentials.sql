-- name: GetCredential :one
SELECT * FROM credentials
WHERE id = ? LIMIT 1;

-- name: GetCredentialByName :one
SELECT * FROM credentials
WHERE name = ? LIMIT 1;

-- name: ListCredentials :many
SELECT 
    id, name, type, description, storage_type, 
    created_at, updated_at, last_used_at
FROM credentials
ORDER BY name ASC;

-- name: ListCredentialsByType :many
SELECT 
    id, name, type, description, storage_type, 
    created_at, updated_at, last_used_at
FROM credentials
WHERE type = ?
ORDER BY name ASC;

-- name: CreateCredential :one
INSERT INTO credentials (
    id, name, type, description, keychain_key, 
    encrypted_data, salt, storage_type
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateCredential :one
UPDATE credentials
SET 
    name = ?,
    type = ?,
    description = ?,
    keychain_key = ?,
    encrypted_data = ?,
    salt = ?,
    storage_type = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateCredentialLastUsed :exec
UPDATE credentials
SET last_used_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteCredential :exec
DELETE FROM credentials
WHERE id = ?;

-- name: CheckCredentialNameExists :one
SELECT COUNT(*) as count FROM credentials
WHERE name = ?;

-- name: GetCredentialForUse :one
SELECT * FROM credentials
WHERE name = ?
LIMIT 1;

-- name: SearchCredentialsByName :many
SELECT 
    id, name, type, description, storage_type, 
    created_at, updated_at, last_used_at
FROM credentials
WHERE name LIKE ?
ORDER BY name ASC;

-- name: GetUnusedCredentials :many
SELECT 
    id, name, type, description, storage_type, 
    created_at, updated_at, last_used_at
FROM credentials
WHERE last_used_at IS NULL 
   OR last_used_at < datetime('now', '-90 days')
ORDER BY last_used_at ASC NULLS FIRST;