// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"time"
)

type CryptoWallet struct {
	ID        int64
	UserID    int64
	Name      string
	CreatedAt time.Time
}

type Transaction struct {
	ID        int64
	WalletID  int64
	Amount    int64
	CreatedAt time.Time
}

type User struct {
	ID         int64
	TelegramID int64
	CreatedAt  time.Time
}
