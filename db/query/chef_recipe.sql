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
