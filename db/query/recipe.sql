-- name: ListTrendRecipe :many
WITH
history AS (
    SELECT
        SUM(CASE WHEN is_fav THEN 1 ELSE 0 END) - SUM(CASE WHEN is_fav THEN 0 ELSE 1 END) AS score,
        recipe_id
    FROM
        fav_history
    WHERE
        chef_id IS NOT NULL
    AND
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
    recipe.access_level = 1
AND
    history.recipe_id = recipe.id
ORDER BY
    score DESC
LIMIT @lim;

-- name: GetRecipe :one
SELECT
    *
FROM
    v_recipe
WHERE
    id = @id;

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

-- name: ListRecipe :many
SELECT
    id,
    chef_id,
    name,
    servings,
    image_url,
    introduction,
    created_at,
    updated_at,
    num_fav
FROM
    recipe
WHERE
    access_level = 1
AND
    chef_id IS NOT NULL
ORDER BY
    num_fav DESC;
