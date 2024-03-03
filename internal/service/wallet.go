package service

import (
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	"github.com/maxzhovtyj/coin-tracker/internal/storage"
	"github.com/maxzhovtyj/coin-tracker/pkg/binance"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
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

func (w *WalletService) NetWorth(telegramID int64) (models.UserNetWorth, error) {
	wallets, err := w.db.All(telegramID)
	if err != nil {
		return models.UserNetWorth{}, err
	}

	symbols := make([]string, len(wallets))
	walletsResp := make(map[string]models.Wallet)
	for i, cw := range wallets {
		symbols[i] = cw.Name
		walletsResp[cw.Name] = models.Wallet{
			Id:     cw.ID,
			UserID: cw.UserID,
			Name:   cw.Name,
			Amount: cw.Amount,
		}
	}

	list, err := w.api.CoinsList(symbols...)
	if err != nil {
		return models.UserNetWorth{}, err
	}

	var netWorth models.UserNetWorth
	netWorth.Wallets = make([]models.Wallet, len(list))

	for i, s := range list {
		wr := walletsResp[s.Symbol]

		netWorth.Wallets[i] = models.Wallet{
			Id:      wr.Id,
			UserID:  wr.UserID,
			Name:    wr.Name,
			Price:   s.Price,
			Amount:  wr.Amount,
			Balance: s.Price * wr.Amount,
		}

		netWorth.Balance += netWorth.Wallets[i].Balance
	}

	return netWorth, nil
}

func (w *WalletService) All(telegramID int64) ([]db.CryptoWallet, error) {
	return w.db.All(telegramID)
}

func (w *WalletService) Create(telegramID int64, walletName string) error {
	_, err := w.api.Coin(walletName)
	if err != nil {
		return err
	}

	return w.db.Create(telegramID, walletName)
}

func (w *WalletService) Get(telegramID int64, name string) (models.Wallet, error) {
	storageWallet, err := w.db.Get(telegramID, name)
	if err != nil {
		return models.Wallet{}, err
	}

	price, err := w.api.Coin(storageWallet.Name)
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
