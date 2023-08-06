DROP TYPE if exists type_vrecipe_ingredient CASCADE;

CREATE TYPE type_vrecipe_ingredient AS (
    id UUID,
    name TEXT,
    supplement TEXT
);


DROP VIEW if exists v_recipe CASCADE;

CREATE VIEW v_recipe AS
SELECT
    id,
    chef_id,
    usr_id,
    name,
    servings,
    COALESCE(
        (
            SELECT JSONB_AGG(
                JSONB_BUILD_OBJECT('id', id) ||
--                 JSONB_BUILD_OBJECT('recipeId', recipe_id) ||
--                 JSONB_BUILD_OBJECT('idx', idx) ||
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
                    ingredient
                WHERE
                    recipe_id = recipe.id
                ORDER BY idx
            ) AS ingre
        ),
        TO_JSONB(ARRAY[]::INTEGER[])
    ) AS ingredient,
    TO_JSONB(method) AS method,
    image_url,
    introduction,
    link,
    access_level,
    created_at,
    updated_at,
    num_fav
FROM
    recipe;


DROP FUNCTION if exists insert_recipe CASCADE;

CREATE OR REPLACE FUNCTION insert_recipe(
    data JSONB
)
    RETURNS v_recipe
    LANGUAGE plpgsql
    volatile
AS
$$
DECLARE
    inserting_method type_recipe_method[];
    rec v_recipe%ROWTYPE;
    ingredient_rec ingredient%ROWTYPE;
    ingredient_array type_vrecipe_ingredient[];
BEGIN
    inserting_method = ARRAY[]::type_recipe_method[];
    FOR i IN 0..JSONB_ARRAY_LENGTH(data->'method') - 1 LOOP
        inserting_method = ARRAY_APPEND(inserting_method,
            ROW(data->'method'->i->>'html',
                data->'method'->i->>'supplement')::type_recipe_method);
    END LOOP;

    INSERT INTO recipe
    (
        chef_id,
        usr_id,
        name,
        servings,
        method,
        image_url,
        introduction,
        link,
        access_level
    )
    VALUES
    (
        (data->>'chefId')::UUID,
        (data->>'usrId')::UUID,
        data->>'name',
        (data->'servings')::INTEGER,
        inserting_method,
        data->>'imageUrl',
        data->>'introduction',
        (SELECT ARRAY_AGG(value) FROM JSONB_ARRAY_ELEMENTS_TEXT(data->'link')),
        (data->'accessLevel')::INTEGER
    )
    RETURNING
        id,
        chef_id,
        usr_id,
        name,
        servings,
        NULL,
        TO_JSONB(method) AS method,
        image_url,
        introduction,
        link,
        access_level,
        created_at,
        updated_at,
        num_fav
    INTO
        rec;

    ingredient_array = ARRAY[]::type_vrecipe_ingredient[];
    FOR i IN 0..JSONB_ARRAY_LENGTH(data->'ingredient') - 1 LOOP
        INSERT INTO ingredient
        (
            recipe_id,
            idx,
            name,
            supplement
        )
        VALUES
        (
            rec.id,
            i + 1,
            data->'ingredient'->i->>'name',
            data->'ingredient'->i->>'supplement'
        )
        RETURNING
            *
        INTO
            ingredient_rec;

        ingredient_array = ARRAY_APPEND(ingredient_array,
            ROW(ingredient_rec.id,
                ingredient_rec.name,
                ingredient_rec.supplement)::type_vrecipe_ingredient);
    END LOOP;
    rec.ingredient = TO_JSONB(ingredient_array);

    RETURN rec;
END
$$;


DROP FUNCTION if exists update_recipe CASCADE;

CREATE OR REPLACE FUNCTION update_recipe(
    id UUID,
    data JSONB
)
    RETURNS v_recipe
    LANGUAGE plpgsql
    volatile
AS
$$
DECLARE
    ingredient_id UUID;
    ingredient_ids UUID[];
    updating_method type_recipe_method[];
    rec v_recipe%ROWTYPE;
    ingredient_rec ingredient%ROWTYPE;
    ingredient_array type_vrecipe_ingredient[];
BEGIN
    ingredient_ids = ARRAY[]::UUID[];
    FOR i IN 0..JSONB_ARRAY_LENGTH(data->'ingredient') - 1 LOOP
        ingredient_id = (data->'ingredient'->i->>'id')::UUID;
        IF ingredient_id IS NOT NULL THEN
            ingredient_ids = ARRAY_APPEND(ingredient_ids, ingredient_id);
        END IF;
    END LOOP;

    DELETE FROM
        ingredient
    WHERE
        recipe_id = update_recipe.id
    AND
        NOT (ingredient.id = ANY (ingredient_ids));

    updating_method = ARRAY[]::type_recipe_method[];
    FOR i IN 0..JSONB_ARRAY_LENGTH(data->'method') - 1 LOOP
        updating_method = ARRAY_APPEND(updating_method,
            ROW(data->'method'->i->>'html',
                data->'method'->i->>'supplement')::type_recipe_method);
    END LOOP;

    UPDATE recipe SET
--         chef_id      = (data->>'chefId')::UUID,
--         usr_id       = (data->>'usrId')::UUID,
        name         = data->>'name',
        servings     = (data->'servings')::INTEGER,
        method       = updating_method,
        image_url    = data->>'imageUrl',
        introduction = data->>'introduction',
        link         = (SELECT ARRAY_AGG(value) FROM JSONB_ARRAY_ELEMENTS_TEXT(data->'link')),
        access_level = (data->'accessLevel')::INTEGER
    WHERE
        recipe.id = update_recipe.id
    RETURNING
        recipe.id,
        chef_id,
        usr_id,
        name,
        servings,
        NULL,
        TO_JSONB(method) AS method,
        image_url,
        introduction,
        link,
        access_level,
        created_at,
        updated_at,
        num_fav
    INTO
        rec;

    ingredient_array = ARRAY[]::type_vrecipe_ingredient[];
    FOR i IN 0..JSONB_ARRAY_LENGTH(data->'ingredient') - 1 LOOP
        ingredient_id = (data->'ingredient'->i->>'id')::UUID;
        IF ingredient_id IS NOT NULL THEN
            UPDATE ingredient SET
                idx        = i + 1,
                name       = data->'ingredient'->i->>'name',
                supplement = data->'ingredient'->i->>'supplement'
            WHERE
                recipe_id = update_recipe.id
            AND
                ingredient.id = ingredient_id
            RETURNING
                *
            INTO
                ingredient_rec;
        ELSE
            INSERT INTO ingredient
            (
                recipe_id,
                idx,
                name,
                supplement
            )
            VALUES
            (
                update_recipe.id,
                i + 1,
                data->'ingredient'->i->>'name',
                data->'ingredient'->i->>'supplement'
            )
            RETURNING
                *
            INTO
                ingredient_rec;
        END IF;

        ingredient_array = ARRAY_APPEND(ingredient_array,
            ROW(ingredient_rec.id,
                ingredient_rec.name,
                ingredient_rec.supplement)::type_vrecipe_ingredient);
    END LOOP;
    rec.ingredient = TO_JSONB(ingredient_array);

    RETURN rec;
END
$$;
