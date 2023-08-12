-- name: ListShoppingList :many
SELECT
    id,
    usr_id,
    recipe_id,
    r_idx,
    recipe_name,
    chef_name,
    general_chef_name,
    description,
    is_fair_copy,
    created_at,
    updated_at,
    item
FROM
    v_shopping_list
WHERE
    usr_id = @usr_id
ORDER BY
    r_idx DESC;

-- name: GetShoppingList :one
SELECT
    id,
    usr_id,
    recipe_id,
    r_idx,
    recipe_name,
    chef_name,
    general_chef_name,
    description,
    is_fair_copy,
    created_at,
    updated_at,
    item
FROM
    v_shopping_list
WHERE
    usr_id = @usr_id
AND
    recipe_id = @recipe_id;

-- name: CreateShoppingList :one
INSERT INTO shopping_list (
    usr_id,
    recipe_id,
    r_idx,
    description,
    is_fair_copy
)
VALUES
(
    @usr_id,
    @recipe_id,
    (SELECT
        COALESCE(MAX(r_idx) + 1, 1)
    FROM
        shopping_list
    WHERE
        usr_id = @usr_id
    ),
    @description,
    @is_fair_copy
)
RETURNING
    id,
    usr_id,
    recipe_id,
    r_idx,
    description,
    is_fair_copy,
    created_at,
    updated_at;

-- name: InnerCreateShoppingItem :one
INSERT INTO shopping_item (
    shopping_list_id,
    ingredient_id,
    idx,
    name,
    supplement
)
VALUES
(
    @shopping_list_id,
    @ingredient_id,
    @idx,
    @name,
    @supplement
)
RETURNING
    id,
    ingredient_id,
    name,
    supplement;

-- name: UpdateShoppingList :one
UPDATE shopping_list SET
    r_idx        = @r_idx,
    description  = @description,
    is_fair_copy = @is_fair_copy
WHERE
    usr_id = @usr_id
AND
    id = @id
RETURNING
    id,
    usr_id,
    recipe_id,
    r_idx,
    description,
    is_fair_copy,
    created_at,
    updated_at;

-- name: InnerUpdateShoppingItem :one
UPDATE shopping_item SET
    idx           = @idx,
    name          = @name,
    supplement    = @supplement
WHERE
    (ingredient_id = @ingredient_id OR @ingredient_id IS NULL)
AND
    shopping_list_id = @shopping_list_id
AND
    id = @id
RETURNING
    id,
    ingredient_id,
    name,
    supplement;

-- name: InnerDeleteNotAnyShoppingItem :exec
DELETE FROM
    shopping_item
WHERE
    shopping_list_id = @shopping_list_id
AND
    NOT (id = ANY (@id::UUID[]));

-- name: DeleteShoppingList :one
DELETE FROM
    shopping_list
WHERE
    usr_id = @usr_id
AND
    id = @id
RETURNING
    id,
    usr_id,
    recipe_id,
    r_idx,
    description,
    is_fair_copy,
    created_at,
    updated_at;
