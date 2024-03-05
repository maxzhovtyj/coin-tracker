package service

import (
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	"github.com/maxzhovtyj/coin-tracker/internal/storage"
	"github.com/maxzhovtyj/coin-tracker/pkg/binance"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
	"time"
)

type Service struct {
	User
	Wallet
	Subscription
}

type User interface {
	Create(telegramID int64) (db.User, error)
}

type Wallet interface {
	Get(telegramID int64, wallet string) (models.Wallet, error)
	Create(telegramID int64, wallet string) error
	All(telegramID int64) ([]db.CryptoWallet, error)
	NewTransaction(wallet int64, amount, price float64) error
	NetWorth(telegramID int64) (models.UserNetWorth, error)
	GetTransactions(wallet int64) ([]models.Transaction, error)
}

type Subscription interface {
	All() ([]models.Subscription, error)
	UserSubscriptions(uid int64) ([]models.Subscription, error)
	NewCoinSubscription(uid int64, coinName string, interval time.Duration) error
	Notified(id int64) error
	CoinTicker(coin, windowSize string) (binance.SymbolTicker, error)
}

func New(storage *storage.Storage, api binance.API) *Service {
	return &Service{
		User:         NewUserService(storage.User),
		Wallet:       NewWalletService(storage.Wallet, api),
		Subscription: NewSubscriptionService(storage.Subscription, api),
	}
}
