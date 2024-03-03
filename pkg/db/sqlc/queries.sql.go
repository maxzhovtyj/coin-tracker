// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package db

import (
	"context"
)

const createSubscription = `-- name: CreateSubscription :one
INSERT INTO subscriptions (type, user_id, data, notify_interval) VALUES (?, ?, ?, ?) RETURNING id, type, user_id, data, notify_interval, last_notified_at
`

type CreateSubscriptionParams struct {
	Type           string
	UserID         int64
	Data           string
	NotifyInterval string
}

func (q *Queries) CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) (Subscription, error) {
	row := q.db.QueryRowContext(ctx, createSubscription,
		arg.Type,
		arg.UserID,
		arg.Data,
		arg.NotifyInterval,
	)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.UserID,
		&i.Data,
		&i.NotifyInterval,
		&i.LastNotifiedAt,
	)
	return i, err
}

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions (wallet_id, amount, price) VALUES (?, ?, ?) RETURNING id, wallet_id, amount, price, created_at
`

type CreateTransactionParams struct {
	WalletID int64
	Amount   float64
	Price    float64
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction, arg.WalletID, arg.Amount, arg.Price)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.WalletID,
		&i.Amount,
		&i.Price,
		&i.CreatedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (telegram_id)
VALUES (?)
RETURNING id, telegram_id, created_at
`

func (q *Queries) CreateUser(ctx context.Context, telegramID int64) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, telegramID)
	var i User
	err := row.Scan(&i.ID, &i.TelegramID, &i.CreatedAt)
	return i, err
}

const createUserWallet = `-- name: CreateUserWallet :one
INSERT INTO crypto_wallets (user_id, name)
VALUES (?, ?)
RETURNING id, user_id, name, amount, created_at
`

type CreateUserWalletParams struct {
	UserID int64
	Name   string
}

func (q *Queries) CreateUserWallet(ctx context.Context, arg CreateUserWalletParams) (CryptoWallet, error) {
	row := q.db.QueryRowContext(ctx, createUserWallet, arg.UserID, arg.Name)
	var i CryptoWallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getSubscriptions = `-- name: GetSubscriptions :many
SELECT id, type, user_id, data, notify_interval, last_notified_at
FROM subscriptions
`

func (q *Queries) GetSubscriptions(ctx context.Context) ([]Subscription, error) {
	rows, err := q.db.QueryContext(ctx, getSubscriptions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Subscription
	for rows.Next() {
		var i Subscription
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.UserID,
			&i.Data,
			&i.NotifyInterval,
			&i.LastNotifiedAt,
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

const getUser = `-- name: GetUser :one
SELECT id, telegram_id, created_at
FROM users
WHERE id = ?
LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(&i.ID, &i.TelegramID, &i.CreatedAt)
	return i, err
}

const getUserSubscription = `-- name: GetUserSubscription :one
SELECT id, type, user_id, data, notify_interval, last_notified_at
FROM subscriptions
WHERE user_id = ? AND type = ?
`

type GetUserSubscriptionParams struct {
	UserID int64
	Type   string
}

func (q *Queries) GetUserSubscription(ctx context.Context, arg GetUserSubscriptionParams) (Subscription, error) {
	row := q.db.QueryRowContext(ctx, getUserSubscription, arg.UserID, arg.Type)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.UserID,
		&i.Data,
		&i.NotifyInterval,
		&i.LastNotifiedAt,
	)
	return i, err
}

const getUserSubscriptions = `-- name: GetUserSubscriptions :many
SELECT id, type, user_id, data, notify_interval, last_notified_at
FROM subscriptions
WHERE user_id = ?
`

func (q *Queries) GetUserSubscriptions(ctx context.Context, userID int64) ([]Subscription, error) {
	rows, err := q.db.QueryContext(ctx, getUserSubscriptions, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Subscription
	for rows.Next() {
		var i Subscription
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.UserID,
			&i.Data,
			&i.NotifyInterval,
			&i.LastNotifiedAt,
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

const getUserWallet = `-- name: GetUserWallet :one
SELECT id, user_id, name, amount, created_at
FROM crypto_wallets
WHERE user_id = ? AND name = ?
ORDER BY created_at DESC
`

type GetUserWalletParams struct {
	UserID int64
	Name   string
}

func (q *Queries) GetUserWallet(ctx context.Context, arg GetUserWalletParams) (CryptoWallet, error) {
	row := q.db.QueryRowContext(ctx, getUserWallet, arg.UserID, arg.Name)
	var i CryptoWallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getUserWallets = `-- name: GetUserWallets :many
SELECT id, user_id, name, amount, created_at
FROM crypto_wallets
WHERE user_id = ?
ORDER BY created_at DESC
`

func (q *Queries) GetUserWallets(ctx context.Context, userID int64) ([]CryptoWallet, error) {
	rows, err := q.db.QueryContext(ctx, getUserWallets, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CryptoWallet
	for rows.Next() {
		var i CryptoWallet
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Amount,
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

const updateLastNotifiedAt = `-- name: UpdateLastNotifiedAt :one
UPDATE subscriptions SET last_notified_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING id, type, user_id, data, notify_interval, last_notified_at
`

func (q *Queries) UpdateLastNotifiedAt(ctx context.Context, id int64) (Subscription, error) {
	row := q.db.QueryRowContext(ctx, updateLastNotifiedAt, id)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.UserID,
		&i.Data,
		&i.NotifyInterval,
		&i.LastNotifiedAt,
	)
	return i, err
}

const updateWalletBalance = `-- name: UpdateWalletBalance :one
UPDATE crypto_wallets SET amount = amount + ?  WHERE id = ? RETURNING id, user_id, name, amount, created_at
`

type UpdateWalletBalanceParams struct {
	Amount float64
	ID     int64
}

func (q *Queries) UpdateWalletBalance(ctx context.Context, arg UpdateWalletBalanceParams) (CryptoWallet, error) {
	row := q.db.QueryRowContext(ctx, updateWalletBalance, arg.Amount, arg.ID)
	var i CryptoWallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}
