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

-- -- name: ListFavoriteRecipe :many
-- SELECT
--     id,
--     name,
--     image_url,
--     profile,
--     created_at,
--     updated_at,
--     num_recipe,
--     num_follower
-- FROM
--     chef
-- WHERE
--     EXISTS (
--         SELECT
--             1
--         FROM
--             favoring
--         WHERE
--         recipe_id=@recipe_id
--         AND
--         usr_id=@usr_id
--     )
-- ORDER BY
--     num_follower DESC;