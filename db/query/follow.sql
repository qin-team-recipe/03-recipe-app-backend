-- name: CreateFollowChef :one
INSERT INTO following_chef
(
    chef_id,
    usr_id
)
VALUES
(
    @chef_id,
    @usr_id
)
RETURNING
    *;

-- name: DeleteFollowChef :one
DELETE FROM
    following_chef
WHERE
    chef_id = @chef_id
AND
    usr_id = @usr_id
RETURNING
    *;

-- name: ExistsFollowChef :one
SELECT EXISTS (
    SELECT
        *
    FROM
        following_chef
    WHERE
        chef_id = @chef_id
    AND
        usr_id = @usr_id
);

-- name: GetFollowChef :many
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
    EXISTS (
        SELECT
            1
        FROM
            following_chef
        WHERE
            usr_id = @usr_id
        AND
            chef_id = chef.id
    )
ORDER BY
    num_follower DESC;
