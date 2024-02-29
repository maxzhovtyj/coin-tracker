package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
	"time"
)

type WalletStorage struct {
	q *db.Queries
}

func NewWalletStorage(conn db.DBTX) *WalletStorage {
	return &WalletStorage{
		q: db.New(conn),
	}
}

func (w *WalletStorage) Get(telegramID int64, wallet string) (db.CryptoWallet, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	userWallet, err := w.q.GetUserWallet(ctx, db.GetUserWalletParams{
		UserID: telegramID,
		Name:   wallet,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.CryptoWallet{}, models.ErrWalletNotFound
		}

		return db.CryptoWallet{}, err
	}

	return userWallet, nil
}

func (w *WalletStorage) All(telegramID int64) ([]db.CryptoWallet, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	return w.q.GetUserWallets(ctx, telegramID)
}

func (w *WalletStorage) Create(telegramID int64, wallet string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	_, err := w.q.CreateUserWallet(ctx, db.CreateUserWalletParams{
		UserID: telegramID,
		Name:   wallet,
	})
	if err != nil {
		if dbErr, ok := err.(sqlite3.Error); ok && errors.Is(dbErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return models.ErrWalletAlreadyExists
		}

		return err
	}

	return nil
}
