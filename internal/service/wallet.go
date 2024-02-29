package service

import (
	"errors"
	"fmt"
	"github.com/maxzhovtyj/coin-tracker/internal/storage"
	"github.com/maxzhovtyj/coin-tracker/internal/storage/models"
	"github.com/maxzhovtyj/coin-tracker/pkg/binance"
)

var ErrWalletAlreadyExists = errors.New("wallet already exists")

type WalletService struct {
	db  storage.Wallet
	api binance.API
}

func (w *WalletService) Create(telegramID int64, walletName string) error {
	info, err := w.api.Info(walletName)
	if err != nil {
		return err
	}

	if info == nil {
		return fmt.Errorf("unknown coin")
	}

	err = w.db.Create(telegramID, walletName)
	if err != nil {
		if errors.Is(err, models.ErrConstraintUnique) {
			return ErrWalletAlreadyExists
		}

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
