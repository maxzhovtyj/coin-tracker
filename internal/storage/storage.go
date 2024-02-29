package storage

import (
	"github.com/maxzhovtyj/coin-tracker/internal/storage/sqlc"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
)

type Storage struct {
	User
	Wallet
}

type User interface {
	Create(telegramID int64) (db.User, error)
}

type Wallet interface {
	Create(telegramID int64, wallet string) error
}

func New(conn db.DBTX) *Storage {
	return &Storage{
		User:   sqlc.NewUserStorage(conn),
		Wallet: sqlc.NewWalletStorage(conn),
	}
}
