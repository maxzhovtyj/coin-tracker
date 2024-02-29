package sqlc

import (
	"context"
	"errors"
	"github.com/mattn/go-sqlite3"
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
	"time"
)

type UserStorage struct {
	q *db.Queries
}

func NewUserStorage(conn db.DBTX) *UserStorage {
	return &UserStorage{
		q: db.New(conn),
	}
}

func (u *UserStorage) Create(telegramID int64) (db.User, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	user, err := u.q.CreateUser(ctx, telegramID)
	if err != nil {
		if dbErr, ok := err.(sqlite3.Error); ok && errors.Is(dbErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return db.User{}, models.ErrUserAlreadyExists
		}

		return db.User{}, err
	}

	return user, nil
}
