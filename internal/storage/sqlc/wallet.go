package sqlc

import (
	"context"
	"errors"
	"github.com/mattn/go-sqlite3"
	"github.com/maxzhovtyj/coin-tracker/internal/storage/models"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
	"time"
)

type WalletStorage struct {
	q *db.Queries
}

func (w WalletStorage) Create(telegramID int64, wallet string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	_, err := w.q.CreateUserWallet(ctx, db.CreateUserWalletParams{
		UserID: telegramID,
		Name:   wallet,
	})
	if err != nil {
		if dbErr, ok := err.(sqlite3.Error); ok && errors.Is(dbErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return models.ErrConstraintUnique
		}

		return err
	}

	return nil
}

func NewWalletStorage(conn db.DBTX) *WalletStorage {
	return &WalletStorage{
		q: db.New(conn),
	}
}
