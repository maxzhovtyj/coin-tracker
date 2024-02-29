package service

import (
	"github.com/maxzhovtyj/coin-tracker/internal/storage"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
)

type Service struct {
	User
	Wallet
}

type User interface {
	Create(telegramID int64) (db.Users, error)
}

type Wallet interface {
}

func New(storage *storage.Storage) *Service {
	return &Service{
		User:   NewUserService(storage.User),
		Wallet: NewWalletService(storage.Wallet),
	}
}
