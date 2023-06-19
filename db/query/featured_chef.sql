-- name: ListFeaturedChef :many
WITH
history AS (
SELECT
    SUM(CASE WHEN is_follow THEN 1 ELSE 0 END) - SUM(CASE WHEN is_follow THEN 0 ELSE 1 END) AS score,
    chef_id
FROM
    follow_history
WHERE
    CURRENT_TIMESTAMP - INTERVAL '3 days' <= created_at
GROUP BY
    chef_id
)
SELECT
    history.chef_id,
    chef.name,
    chef.image_url,
    (SELECT
         COUNT(1)
     FROM
         following
     WHERE
         following.chef_id = history.chef_id) AS num_follower,
    history.score
FROM
    history
INNER JOIN
    chef
ON
    history.chef_id = chef.id
ORDER BY
    score DESC;

-- name: FakeListFeaturedChef :many
WITH
RECURSIVE generate_index (ii) AS (
    SELECT 0
    UNION ALL
    SELECT ii + 1
    FROM generate_index
)
SELECT
    nanoid() AS chef_id,
    '' AS name,
    'https://source.unsplash.com/random/300x300?v=1' AS image_url,
    0 AS num_follower,
    0 AS score
FROM
    generate_index
LIMIT 10;
