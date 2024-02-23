package service

import (
	"errors"
	"forum/internal/models"
	"forum/internal/storage"
)

type User interface {
	GetUserByToken(token string) (models.User, error)
	UpdateUserName(id int, username string) error
}

type UserService struct {
	storage *storage.Storage
}

func NewUserService(storage *storage.Storage) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (u *UserService) GetUserByToken(token string) (models.User, error) {
	return u.storage.User.GetUserByToken(token)
}

func (u *UserService) UpdateUserName(id int, username string) error {
	if username == "" {
		return errors.New(" username field not found (empty username)")
	}
	if len(username) > 35 {
		return errors.New(" username should be shorter than 36 symbols")
	}
	check, err := u.storage.User.CheckUserByNameEmail(username, username)
	if err != nil {
		return errors.New(" error func CheckUserByNameEmail")
	}
	if check {
		return errors.New(" Sorry, but you cant use this Username, try other username")
	}
	return u.storage.User.UpdateUserName(id, username)
}
