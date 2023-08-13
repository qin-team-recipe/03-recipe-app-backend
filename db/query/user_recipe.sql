-- name: GetUsrRecipes :many
SELECT
    *
FROM
    v_recipe
WHERE
    usr_id = @usr_id
ORDER BY
    created_at DESC;

-- name: DeleteUserRecipe :one
DELETE FROM
    recipe
WHERE
    usr_id = @usr_id
AND
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
