-- name: ListFeaturedChef :many
WITH
history AS (
    SELECT
        SUM(CASE WHEN is_follow THEN 1 ELSE 0 END) - SUM(CASE WHEN is_follow THEN 0 ELSE 1 END) AS score,
        chef_id
    FROM
        follow_chef_history
    WHERE
        CURRENT_TIMESTAMP - INTERVAL '3 days' <= created_at
    GROUP BY
        chef_id
)
SELECT
    history.chef_id,
    chef.name,
    chef.image_url,
    chef.num_follower,
    history.score
FROM
    history
INNER JOIN
    chef
ON
    history.chef_id = chef.id
ORDER BY
    score DESC
LIMIT @lim;
