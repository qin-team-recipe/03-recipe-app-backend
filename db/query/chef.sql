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

-- name: GetChef :one
SELECT
    id,
    name,
    image_url,
    profile,
    link,
    created_at,
    updated_at,
    num_recipe,
    num_follower
FROM
    v_chef
WHERE
    id = @id;

-- name: CreateChef :one
SELECT
    id,
    name,
    image_url,
    profile,
    link,
    created_at,
    updated_at,
    num_recipe,
    num_follower
FROM
    insert_chef(@data);

-- name: UpdateChef :one
SELECT
    id,
    name,
    image_url,
    profile,
    link,
    created_at,
    updated_at,
    num_recipe,
    num_follower
FROM
    update_chef(@id, @data);

-- name: DeleteChef :one
DELETE FROM
    chef
WHERE
    id = @id
RETURNING
    id,
    name,
    image_url,
    profile,
    created_at,
    updated_at,
    num_recipe,
    num_follower;

-- name: SearchChef :many
SELECT
    id,
    name,
    image_url,
    profile,
    created_at,
    updated_at,
    num_recipe,
    num_follower
FROM
    chef
WHERE
    name &@~ @query
OR
    profile &@~ @query
ORDER BY
    pgroonga_score(tableoid, ctid) DESC,
    num_follower DESC;
