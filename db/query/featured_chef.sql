-- name: ListFeaturedChef :many
WITH
history AS (
SELECT
    SUM(CASE WHEN is_follow THEN 1 ELSE 0 END) - SUM(CASE WHEN is_follow THEN 0 ELSE 1 END) AS kpi_featured,
    chef_id
FROM
    follow_history
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
    history.kpi_featured
FROM
    history
INNER JOIN
    chef
ON
    history.chef_id = chef.id
ORDER BY
    kpi_featured DESC;

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
    0 AS kpi_featured
FROM
    generate_index
LIMIT 10;
