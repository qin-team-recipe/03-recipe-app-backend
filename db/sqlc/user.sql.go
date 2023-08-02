// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: user.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO usr (
  email, name, auth_server, auth_userinfo
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, name, email
`

type CreateUserParams struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	AuthServer   string `json:"authServer"`
	AuthUserinfo []byte `json:"authUserinfo"`
}

type CreateUserRow struct {
	ID    pgtype.UUID `json:"id"`
	Name  string      `json:"name"`
	Email string      `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Email,
		arg.Name,
		arg.AuthServer,
		arg.AuthUserinfo,
	)
	var i CreateUserRow
	err := row.Scan(&i.ID, &i.Name, &i.Email)
	return i, err
}

const existsUser = `-- name: ExistsUser :one
SELECT EXISTS (SELECT id, email, name, image_url, profile, link, auth_server, auth_userinfo, created_at, updated_at, num_recipe FROM usr WHERE email = $1)
`

func (q *Queries) ExistsUser(ctx context.Context, email string) (bool, error) {
	row := q.db.QueryRow(ctx, existsUser, email)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getUserId = `-- name: GetUserId :one
SELECT id FROM usr WHERE email = $1
`

func (q *Queries) GetUserId(ctx context.Context, email string) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, getUserId, email)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}
