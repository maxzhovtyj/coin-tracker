package service

import (
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	"github.com/maxzhovtyj/coin-tracker/internal/storage"
	"github.com/maxzhovtyj/coin-tracker/pkg/binance"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
)

type Service struct {
	User
	Wallet
}

type User interface {
	Create(telegramID int64) (db.User, error)
}

type Wallet interface {
	Get(telegramID int64, wallet string) (models.Wallet, error)
	Create(telegramID int64, wallet string) error
	All(telegramID int64) ([]db.CryptoWallet, error)
}

func New(storage *storage.Storage, api binance.API) *Service {
	return &Service{
		User:   NewUserService(storage.User),
		Wallet: NewWalletService(storage.Wallet, api),
	}
}
