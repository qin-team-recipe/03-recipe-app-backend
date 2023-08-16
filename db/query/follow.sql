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

-- name: ListFollowChef :many
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

-- name: ListFollowChefNewRecipe :many
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
    EXISTS (
        SELECT
            1
        FROM
            following_chef
        WHERE
            following_chef.usr_id = @usr_id
        AND
            following_chef.chef_id = recipe.chef_id
    )
ORDER BY
    created_at DESC;
