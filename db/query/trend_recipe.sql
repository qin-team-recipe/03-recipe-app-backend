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

-- name: FakeListTrendRecipe :many
WITH
RECURSIVE generate_index (ii) AS (
    SELECT 0
    UNION ALL
    SELECT ii + 1
    FROM generate_index
)
SELECT
    GEN_RANDOM_UUID() AS recipe_id,
    '' AS name,
    '' AS introduction,
    'https://source.unsplash.com/random/300x300?v=1' AS image_url,
    0 AS num_fav,
    0 AS score
FROM
    generate_index
LIMIT @lim;
