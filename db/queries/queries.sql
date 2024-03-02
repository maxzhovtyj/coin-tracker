-- name: GetUser :one
SELECT *
FROM users
WHERE id = ?
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (telegram_id)
VALUES (?)
RETURNING *;

-- name: GetUserWallets :many
SELECT *
FROM crypto_wallets
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: GetUserWallet :one
SELECT *
FROM crypto_wallets
WHERE user_id = ? AND name = ?
ORDER BY created_at DESC;

-- name: CreateUserWallet :one
INSERT INTO crypto_wallets (user_id, name)
VALUES (?, ?)
RETURNING *;

-- name: CreateTransaction :one
INSERT INTO transactions (wallet_id, amount, price) VALUES (?, ?, ?) RETURNING *;

-- name: UpdateWalletBalance :one
UPDATE crypto_wallets SET amount = amount + ?  WHERE id = ? RETURNING *;
