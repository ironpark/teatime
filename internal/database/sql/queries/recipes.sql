-- name: GetRecipe :one
SELECT * FROM recipes
WHERE id = ? LIMIT 1;

-- name: ListRecipes :many
SELECT * FROM recipes
ORDER BY created_at DESC;

-- name: CreateRecipe :one
INSERT INTO recipes (
    id, name, description, recipe_path
) VALUES (
    ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateRecipe :one
UPDATE recipes
SET 
    name = ?,
    description = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateRecipeByPath :one
UPDATE recipes
SET 
    name = ?,
    description = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE recipe_path = ?
RETURNING *;

-- name: DeleteRecipe :exec
DELETE FROM recipes
WHERE id = ?;

-- name: SearchRecipesByName :many
SELECT * FROM recipes
WHERE name LIKE ?
ORDER BY created_at DESC;

-- name: GetRecipeByPath :one
SELECT * FROM recipes
WHERE recipe_path = ?
LIMIT 1;

-- name: ExistsRecipeByPath :one
SELECT COUNT(*) = 1 FROM recipes
WHERE recipe_path = ?;

-- name: GetRecipeIdByPath :one
SELECT id FROM recipes
WHERE recipe_path = ?;