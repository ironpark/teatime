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
    recipe_path = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
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