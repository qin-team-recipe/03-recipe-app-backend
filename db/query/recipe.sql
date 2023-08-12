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

-- name: GetChefRecipes :many
SELECT
    *
FROM
    v_recipe
WHERE
    chef_id = @chef_id
ORDER BY
    created_at DESC;

-- name: GetUsrRecipes :many
SELECT
    *
FROM
    v_recipe
WHERE
    usr_id = @usr_id
ORDER BY
    created_at DESC;

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

-- name: DeleteRecipe :one
DELETE FROM
    recipe
WHERE
    id = @id
RETURNING
    id,
    chef_id,
    usr_id,
    name,
    servings,
    image_url,
    introduction,
    link,
    access_level,
    created_at,
    updated_at,
    num_fav;
