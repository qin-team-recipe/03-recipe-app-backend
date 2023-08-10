-- name: ExistsUser :one
SELECT EXISTS (SELECT * FROM usr WHERE email = $1);

-- name: CreateUser :one
INSERT INTO usr (
  email, name, auth_server, auth_userinfo
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, name, email;

-- name: GetUserId :one
SELECT id FROM usr WHERE email = $1;

-- name: GetSelf :one
SELECT
    id,
    email,
    name,
    image_url,
    profile,
    link,
    created_at,
    updated_at,
    num_recipe
FROM
    v_usr
WHERE
    email = @email;

-- name: GetUser :one
SELECT
    id,
    email,
    name,
    image_url,
    profile,
    link,
    created_at,
    updated_at,
    num_recipe
FROM
    v_usr
WHERE
    id = @id;

-- name: UpdateUser :one
SELECT
    id,
    email,
    name,
    image_url,
    profile,
    link,
    created_at,
    updated_at,
    num_recipe
FROM
    update_usr(@email, @data);
