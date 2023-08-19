-- name: CreateFavoriteRecipe :one
INSERT INTO favoring 
(
    recipe_id,
    usr_id
)
VALUES
(
    @recipe_id,
    @usr_id
)
RETURNING 
    *;

-- name: DeleteFavoriteRecipe :one
DELETE FROM
    favoring
WHERE
    recipe_id = @recipe_id
AND
    usr_id = @usr_id
RETURNING
    *;


-- name: ExistsFavoriteRecipe :one
SELECT EXISTS (
    SELECT
        1
    FROM
        favoring
    WHERE
        recipe_id = @recipe_id
    AND
        usr_id = @usr_id
);

-- name: ListFavoriteRecipe :many
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
    recipe.access_level = 1
    AND
        EXISTS (
            SELECT
                1
            FROM
                favoring
            WHERE
                recipe_id = recipe.id
            AND
                favoring.usr_id = @usr_id
        );