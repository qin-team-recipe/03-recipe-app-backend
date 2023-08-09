DROP FUNCTION if exists update_usr;


CREATE OR REPLACE FUNCTION update_usr(
    email TEXT,
    data JSONB
)
    RETURNS v_usr
    LANGUAGE plpgsql
    volatile
AS
$$
DECLARE
    updating_link type_chef_link[];
    rec v_usr%ROWTYPE;
BEGIN
    updating_link = ARRAY[]::type_chef_link[];
    FOR i IN 0..JSONB_ARRAY_LENGTH(data->'link') - 1 LOOP
        updating_link = ARRAY_APPEND(updating_link,
            ROW(data->'link'->i->>'label',
                data->'link'->i->>'url')::type_chef_link);
    END LOOP;

    UPDATE usr SET
        name          = data->>'name',
        image_url     = data->>'imageUrl',
        profile       = data->>'profile',
        link          = updating_link
--         auth_server   = data->>'authServer',
--         auth_userinfo = data->'authUserinfo'
    WHERE
        usr.email = update_usr.email
    RETURNING
        id,
        usr.email,
        name,
        image_url,
        profile,
        TO_JSONB(link) AS link,
        auth_server,
        auth_userinfo,
        created_at,
        updated_at,
        num_recipe
    INTO
        rec;

    RETURN rec;
END
$$;
