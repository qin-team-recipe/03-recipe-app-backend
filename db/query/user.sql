-- name: ExistsUser :one
SELECT EXISTS (SELECT * FROM usr WHERE email = $1);

-- name: CreateUser :one
INSERT INTO usr (
  email, name, auth_server, auth_userinfo
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, name, email;