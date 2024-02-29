package sqlc

import db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"

type WalletStorage struct {
	q *db.Queries
}

func NewWalletStorage(conn db.DBTX) *WalletStorage {
	return &WalletStorage{
		q: db.New(conn),
	}
}
