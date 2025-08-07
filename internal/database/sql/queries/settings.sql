-- name: GetSettings :one
SELECT * FROM settings LIMIT 1;

-- name: CreateDefaultSettings :exec
INSERT INTO settings (id, theme, auto_start, language) 
VALUES (1, 'system', 0, 'en')
ON CONFLICT(id) DO NOTHING;

-- name: GetOrCreateSettings :one
INSERT INTO settings (id, theme, auto_start, language) 
VALUES (1, 'system', 0, 'en')
ON CONFLICT(id) DO UPDATE SET
    id = excluded.id
RETURNING *;

-- name: UpdateTheme :exec
UPDATE settings 
SET theme = ?, updated_at = CURRENT_TIMESTAMP 
WHERE id = 1;

-- name: UpdateAutoStart :exec
UPDATE settings 
SET auto_start = ?, updated_at = CURRENT_TIMESTAMP 
WHERE id = 1;

-- name: UpdateLanguage :exec
UPDATE settings 
SET language = ?, updated_at = CURRENT_TIMESTAMP 
WHERE id = 1;

-- name: UpdateSettings :exec
UPDATE settings 
SET theme = ?, 
    auto_start = ?, 
    language = ?,
    updated_at = CURRENT_TIMESTAMP 
WHERE id = 1;