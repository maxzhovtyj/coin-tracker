package service

import (
	"fmt"
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	"github.com/maxzhovtyj/coin-tracker/internal/storage"
	"github.com/maxzhovtyj/coin-tracker/pkg/binance"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
	"strconv"
)

type WalletService struct {
	db  storage.Wallet
	api binance.API
}

func (w *WalletService) NewTransaction(wallet int64, amount, price float64) error {
	_, err := w.db.CreateWalletRecord(wallet, amount, price)
	if err != nil {
		return err
	}

	return nil
}

func NewWalletService(db storage.Wallet, api binance.API) *WalletService {
	return &WalletService{
		db:  db,
		api: api,
	}
}

func (w *WalletService) All(telegramID int64) ([]db.CryptoWallet, error) {
	return w.db.All(telegramID)
}

func (w *WalletService) Create(telegramID int64, walletName string) error {
	info, err := w.api.Info(walletName)
	if err != nil {
		return err
	}

	if info == nil {
		return fmt.Errorf("unknown coin")
	}

	return w.db.Create(telegramID, walletName)
}

func (w *WalletService) Get(telegramID int64, name string) (models.Wallet, error) {
	storageWallet, err := w.db.Get(telegramID, name)
	if err != nil {
		return models.Wallet{}, err
	}

	coinInfo, err := w.api.Info(storageWallet.Name)
	if err != nil {
		return models.Wallet{}, err
	}

	price, err := strconv.ParseFloat(coinInfo.Price, 64)
	if err != nil {
		return models.Wallet{}, err
	}

	return models.Wallet{
		Id:      storageWallet.ID,
		UserID:  storageWallet.UserID,
		Name:    storageWallet.Name,
		Price:   price,
		Amount:  storageWallet.Amount,
		Balance: storageWallet.Amount * price,
	}, nil
}
