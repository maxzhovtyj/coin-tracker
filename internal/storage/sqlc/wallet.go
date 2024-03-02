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
	dbRaw *sql.DB
	q     *db.Queries
}

func NewWalletStorage(conn db.DBTX) *WalletStorage {
	dbRaw, ok := conn.(*sql.DB)
	if !ok {
		panic("can't get raw db connection")
	}

	return &WalletStorage{
		dbRaw: dbRaw,
		q:     db.New(conn),
	}
}

func (w *WalletStorage) CreateTransaction(withTx *db.Queries, walletID int64, amount, price float64) (db.Transaction, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	tr, err := withTx.CreateTransaction(ctx, db.CreateTransactionParams{
		WalletID: walletID,
		Amount:   amount,
		Price:    price,
	})
	if err != nil {
		return db.Transaction{}, err
	}

	return tr, nil
}

func (w *WalletStorage) CreateWalletRecord(walletID int64, amount, price float64) (db.Transaction, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	tx, err := w.dbRaw.Begin()
	if err != nil {
		return db.Transaction{}, err
	}
	defer func() {
		err = tx.Rollback()
		if err != nil {
			return
		}
	}()

	qTx := w.q.WithTx(tx)

	tr, err := w.CreateTransaction(qTx, walletID, amount, price)
	if err != nil {
		return db.Transaction{}, err
	}

	_, err = qTx.UpdateWalletBalance(ctx, db.UpdateWalletBalanceParams{
		ID:     walletID,
		Amount: amount,
	})
	if err != nil {
		return db.Transaction{}, err
	}

	return tr, tx.Commit()
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
