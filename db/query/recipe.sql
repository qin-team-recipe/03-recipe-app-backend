-- name: ListTrendRecipe :many
WITH
history AS (
    SELECT
        SUM(CASE WHEN is_fav THEN 1 ELSE 0 END) - SUM(CASE WHEN is_fav THEN 0 ELSE 1 END) AS score,
        recipe_id
    FROM
        fav_history
    WHERE
        CURRENT_TIMESTAMP - INTERVAL '3 days' <= created_at
    GROUP BY
        recipe_id
)
SELECT
    history.recipe_id,
    recipe.name,
    recipe.introduction,
    recipe.image_url,
    recipe.num_fav,
    history.score
FROM
    history
INNER JOIN
    recipe
ON
    history.recipe_id = recipe.id
ORDER BY
    score DESC
LIMIT @lim;

-- name: CreateRecipe :one
SELECT
    *
FROM
    insert_recipe(@data);

-- name: UpdateRecipe :one
SELECT
    *
FROM
    update_recipe(@id, @data);