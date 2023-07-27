-- name: CreateRecipe :one
SELECT
    *
FROM
    insert_recipe(@data);
