// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "users" (telegram_id)
VALUES (?)
RETURNING id, telegram_id
`

func (q *Queries) CreateUser(ctx context.Context, telegramID string) (Users, error) {
	row := q.db.QueryRowContext(ctx, createUser, telegramID)
	var i Users
	err := row.Scan(&i.ID, &i.TelegramID)
	return i, err
}

const createUserWallet = `-- name: CreateUserWallet :one
INSERT INTO "crypto_wallets" (user_id, name)
VALUES (?, ?)
RETURNING id, user_id, name, created_at
`

type CreateUserWalletParams struct {
	UserID int64
	Name   string
}

func (q *Queries) CreateUserWallet(ctx context.Context, arg CreateUserWalletParams) (CryptoWallets, error) {
	row := q.db.QueryRowContext(ctx, createUserWallet, arg.UserID, arg.Name)
	var i CryptoWallets
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, telegram_id
FROM "users"
WHERE id = ?
LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (Users, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i Users
	err := row.Scan(&i.ID, &i.TelegramID)
	return i, err
}

const getUserWallets = `-- name: GetUserWallets :many
SELECT id, user_id, name, created_at
FROM "crypto_wallets"
WHERE user_id = ?
ORDER BY created_at DESC
`

func (q *Queries) GetUserWallets(ctx context.Context, userID int64) ([]CryptoWallets, error) {
	rows, err := q.db.QueryContext(ctx, getUserWallets, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CryptoWallets
	for rows.Next() {
		var i CryptoWallets
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}