// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: user.sql

package db

import (
	"context"

	dto "github.com/aopontann/gin-sqlc/db/dto"
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

const deleteUser = `-- name: DeleteUser :one
DELETE FROM
    usr
WHERE
    email = $1
RETURNING
    id,
    email,
    name,
    image_url,
    profile,
    created_at,
    updated_at,
    num_recipe
`

type DeleteUserRow struct {
	ID        pgtype.UUID        `json:"id"`
	Email     string             `json:"email"`
	Name      string             `json:"name"`
	ImageUrl  pgtype.Text        `json:"imageUrl"`
	Profile   pgtype.Text        `json:"profile"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
	NumRecipe int32              `json:"numRecipe"`
}

func (q *Queries) DeleteUser(ctx context.Context, email string) (DeleteUserRow, error) {
	row := q.db.QueryRow(ctx, deleteUser, email)
	var i DeleteUserRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.ImageUrl,
		&i.Profile,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.NumRecipe,
	)
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

const getSelf = `-- name: GetSelf :one
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
    email = $1
`

type GetSelfRow struct {
	ID        pgtype.UUID          `json:"id"`
	Email     string               `json:"email"`
	Name      string               `json:"name"`
	ImageUrl  pgtype.Text          `json:"imageUrl"`
	Profile   pgtype.Text          `json:"profile"`
	Link      dto.ChefLinkArrayDto `json:"link"`
	CreatedAt pgtype.Timestamptz   `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz   `json:"updatedAt"`
	NumRecipe int32                `json:"numRecipe"`
}

func (q *Queries) GetSelf(ctx context.Context, email string) (GetSelfRow, error) {
	row := q.db.QueryRow(ctx, getSelf, email)
	var i GetSelfRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.ImageUrl,
		&i.Profile,
		&i.Link,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.NumRecipe,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
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
    id = $1
`

type GetUserRow struct {
	ID        pgtype.UUID          `json:"id"`
	Email     string               `json:"email"`
	Name      string               `json:"name"`
	ImageUrl  pgtype.Text          `json:"imageUrl"`
	Profile   pgtype.Text          `json:"profile"`
	Link      dto.ChefLinkArrayDto `json:"link"`
	CreatedAt pgtype.Timestamptz   `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz   `json:"updatedAt"`
	NumRecipe int32                `json:"numRecipe"`
}

func (q *Queries) GetUser(ctx context.Context, id pgtype.UUID) (GetUserRow, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i GetUserRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.ImageUrl,
		&i.Profile,
		&i.Link,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.NumRecipe,
	)
	return i, err
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

const updateUser = `-- name: UpdateUser :one
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
    update_usr($1, $2)
`

type UpdateUserParams struct {
	Email string `json:"email"`
	Data  []byte `json:"data"`
}

type UpdateUserRow struct {
	ID        pgtype.UUID          `json:"id"`
	Email     string               `json:"email"`
	Name      string               `json:"name"`
	ImageUrl  pgtype.Text          `json:"imageUrl"`
	Profile   pgtype.Text          `json:"profile"`
	Link      dto.ChefLinkArrayDto `json:"link"`
	CreatedAt pgtype.Timestamptz   `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz   `json:"updatedAt"`
	NumRecipe int32                `json:"numRecipe"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (UpdateUserRow, error) {
	row := q.db.QueryRow(ctx, updateUser, arg.Email, arg.Data)
	var i UpdateUserRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.ImageUrl,
		&i.Profile,
		&i.Link,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.NumRecipe,
	)
	return i, err
}
