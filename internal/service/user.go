package service

import (
	"github.com/maxzhovtyj/coin-tracker/internal/storage"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
)

type UserService struct {
	db storage.User
}

func NewUserService(db storage.User) *UserService {
	return &UserService{
		db: db,
	}
}

func (u *UserService) Create(telegramID int64) (db.User, error) {
	return u.db.Create(telegramID)
}
