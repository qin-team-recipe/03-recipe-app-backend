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
