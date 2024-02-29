package service

import (
	"errors"
	"github.com/maxzhovtyj/coin-tracker/internal/storage"
	"github.com/maxzhovtyj/coin-tracker/internal/storage/models"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserService struct {
	db storage.User
}

func NewUserService(db storage.User) *UserService {
	return &UserService{
		db: db,
	}
}

func (u *UserService) Create(telegramID int64) (db.Users, error) {
	user, err := u.db.Create(telegramID)
	if err != nil {
		if errors.Is(err, models.ErrConstraintUnique) {
			return db.Users{}, ErrUserAlreadyExists
		}

		return db.Users{}, err
	}

	return user, nil
}
