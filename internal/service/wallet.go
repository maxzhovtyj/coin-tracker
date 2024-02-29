package service

import "github.com/maxzhovtyj/coin-tracker/internal/storage"

type WalletService struct {
	db storage.Wallet
}

func NewWalletService(db storage.Wallet) *WalletService {
	return &WalletService{
		db: db,
	}
}
