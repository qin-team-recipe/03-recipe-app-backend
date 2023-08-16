-- name: GetChefRecipes :many
SELECT
    *
FROM
    v_recipe
WHERE
    chef_id = @chef_id
ORDER BY
    created_at DESC;

-- name: DeleteChefRecipe :one
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

-- name: SearchChefRecipe :many
SELECT
    id,
    chef_id,
    name,
    servings,
    image_url,
    introduction,
    access_level,
    created_at,
    updated_at,
    num_fav
FROM
    recipe
WHERE
    access_level = 1
AND
    chef_id IS NOT NULL
AND (
    name &@~ @query
OR
    introduction &@~ @query
)
ORDER BY
    pgroonga_score(tableoid, ctid) DESC,
    num_fav DESC;
