package service

import (
	"errors"
	"forum/internal/models"
	"forum/internal/storage"
	"regexp"
	"time"
	"unicode"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	CreateUser(user models.User) error
	CheckUser(user models.User) (string, time.Time, error)
	DeleteToken(token string) error
	CreateOrLoginByGoogle(user_g models.GoogleLoginUserData) (string, time.Time, error)
	CreateOrLoginByGithub(user_g models.GithubUserData) (string, time.Time, error)
}

type AuthService struct {
	storage *storage.Storage
}

func NewAuthService(storage *storage.Storage) *AuthService {
	return &AuthService{
		storage: storage,
	}
}

func (a *AuthService) CreateUser(user models.User) error {
	if err := validUser(user); err != nil {
		return err
	}

	uniq, err := a.storage.CheckUserByNameEmail(user.Email, user.Username)
	if err != nil {
		return err
	}
	if uniq {
		return errors.New(" Username or Email is already in used! ")
	}

	user.Password, err = generateHashPassword(user.Password)
	if err != nil {
		return err
	}

	return a.storage.Auth.CreateUser(user)
}

func (a *AuthService) CheckUser(user models.User) (string, time.Time, error) {
	password, err := a.storage.GetPasswordByUsername(user.Username)
	if err != nil {
		return "", time.Time{}, errors.New(" There is no user with that name <" + user.Username + "> ")
	}
	if err := compareHashAndPassword(password, user.Password); err != nil {
		return "", time.Time{}, err
	}

	token := uuid.NewGen()
	d, err := token.NewV4()
	if err != nil {
		return "", time.Time{}, err
	}
	expired := time.Now().Add(time.Hour * 12)
	if err := a.storage.SaveToken(d.String(), expired, user.Username); err != nil {
		return "", time.Time{}, err
	}
	return d.String(), expired, nil
}

func generateHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func compareHashAndPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return errors.New("password does not match")
	}
	return nil
}

func (a *AuthService) DeleteToken(token string) error {
	return a.storage.Auth.DeleteToken(token)
}

func validUser(user models.User) error {
	for _, char := range user.Username {
		if char <= 32 || char >= 127 {
			return models.ErrInvalidUserName
		}
	}
	validEmail, err := regexp.MatchString(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, user.Email)
	if err != nil {
		return err
	}
	if !validEmail {
		return models.ErrInvalidEmail
	}
	if len(user.Username) < 6 || len(user.Username) >= 36 {
		return models.ErrInvalidUserName
	}

	if !passIsValid(user.Password) {
		return models.ErrShortPassword
	}
	if user.Password != user.RepeatPassword {
		return models.ErrPasswordDontMatch
	}
	return nil
}

func passIsValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 8 || len(s) <= 20 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
