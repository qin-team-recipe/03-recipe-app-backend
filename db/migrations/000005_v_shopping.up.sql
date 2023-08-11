DROP VIEW if exists v_shopping_list;

CREATE VIEW v_shopping_list AS
SELECT
    shopping_list.id,
    shopping_list.usr_id,
    shopping_list.recipe_id,
    CASE
    WHEN shopping_list.recipe_id IS NOT NULL THEN
        recipe.name
    ELSE
        '(メモリスト)'
    END AS recipe_name,
    CASE
    WHEN shopping_list.recipe_id IS NOT NULL THEN
        CASE
        WHEN recipe.chef_id IS NOT NULL THEN
            (SELECT
                name
            FROM
                chef
            WHERE
                chef.id = recipe.chef_id)
        ELSE
            NULL::TEXT
        END
    ELSE
        '(メモリスト)'
    END AS chef_name,
    CASE
    WHEN recipe_id IS NOT NULL THEN
        CASE
        WHEN recipe.usr_id IS NOT NULL THEN
            (SELECT
                name
            FROM
                usr
            WHERE
                usr.id = recipe.usr_id)
        ELSE
            NULL::TEXT
        END
    ELSE
        '(メモリスト)'
    END AS general_chef_name,
    shopping_list.description,
    shopping_list.is_fair_copy,
    shopping_list.created_at,
    shopping_list.updated_at,
    COALESCE(
        (
            SELECT JSONB_AGG(
                JSONB_BUILD_OBJECT('id', id) ||
                JSONB_BUILD_OBJECT('ingredientId', ingredient_id) ||
                JSONB_BUILD_OBJECT('name', name) ||
                JSONB_BUILD_OBJECT('supplement', supplement)
--                 JSONB_BUILD_OBJECT('createdAt', created_at) ||
--                 JSONB_BUILD_OBJECT('updatedAt', updated_at)
            )
            FROM
            (
                SELECT
                    *
                FROM
                    shopping_item
                WHERE
                    shopping_list_id = shopping_list.id
                ORDER BY idx
            ) AS ingre
        ),
        TO_JSONB(ARRAY[]::INTEGER[])
    ) AS item,
    shopping_list.r_idx
FROM
    shopping_list
LEFT OUTER JOIN
    recipe
ON
    shopping_list.recipe_id = recipe.id;
