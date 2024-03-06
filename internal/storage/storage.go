package storage

import (
	"github.com/maxzhovtyj/coin-tracker/internal/storage/sqlc"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
)

type Storage struct {
	User
	Wallet
	Subscription
}

type User interface {
	Create(telegramID int64) (db.User, error)
}

type Wallet interface {
	Get(telegramID, wallet int64) (db.CryptoWallet, error)
	Create(telegramID int64, wallet string) error
	All(telegramID int64) ([]db.CryptoWallet, error)
	CreateWalletRecord(walletID int64, amount, price float64) (db.Transaction, error)
	GetTransactions(walletID int64) ([]db.Transaction, error)
	Delete(userID, walletID int64) error
}

type Subscription interface {
	All() ([]db.Subscription, error)
	UserSubscriptions(uid int64) ([]db.Subscription, error)
	Create(uid int64, subscriptionType, data, interval string) (db.Subscription, error)
	UpdateLastNotifiedAt(id int64) error
}

func New(conn db.DBTX) *Storage {
	return &Storage{
		User:         sqlc.NewUserStorage(conn),
		Wallet:       sqlc.NewWalletStorage(conn),
		Subscription: sqlc.NewSubscriptionStorage(conn),
	}
}
