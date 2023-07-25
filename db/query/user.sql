-- name: ExistsUser :one
SELECT EXISTS (SELECT * FROM usr WHERE email = $1);
