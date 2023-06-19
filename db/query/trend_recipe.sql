-- name: ListTrendRecipe :many
WITH
history AS (
SELECT
    SUM(CASE WHEN is_fav THEN 1 ELSE 0 END) - SUM(CASE WHEN is_fav THEN 0 ELSE 1 END) AS score,
    recipe_id
FROM
    fav_history
GROUP BY
    recipe_id
)
SELECT
    history.recipe_id,
    recipe.title,
    recipe.comment,
    recipe.image_url,
    (SELECT
         COUNT(1)
     FROM
         favoring
     WHERE
         favoring.recipe_id = history.recipe_id) AS num_fav,
    history.score
FROM
    history
INNER JOIN
    recipe
ON
    history.recipe_id = recipe.id
ORDER BY
    score DESC;

-- name: FakeListTrendRecipe :many
WITH
RECURSIVE generate_index (ii) AS (
    SELECT 0
    UNION ALL
    SELECT ii + 1
    FROM generate_index
)
SELECT
    nanoid() AS recipe_id,
    '' AS title,
    '' AS comment,
    'https://source.unsplash.com/random/300x300?v=1' AS image_url,
    0 AS num_fav,
    0 AS score
FROM
    generate_index
LIMIT 10;
