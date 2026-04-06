package account

import (
	"context"
	"errors"
	"learnapirest/helpers"
	"learnapirest/internal/config"
)

type IAccountService interface {
	CreateAccount(ctx context.Context, account *RegisterUserRequest) error
	Login(ctx context.Context, username string, password string) (string, error)
}

type AccountService struct {
	repo *AccountRepository
	conf *config.Config
}

func NewAccountService(repo *AccountRepository, config *config.Config) *AccountService {
	return &AccountService{
		repo: repo,
		conf: config,
	}
}

func (a *AccountService) CreateAccount(ctx context.Context, request *RegisterUserRequest) error {
	hashed, err := helpers.HashPassword(request.Password)
	if err != nil {
		return err
	}

	request.Password = hashed

	return a.repo.CreateAccount(ctx, request)
}

func (a *AccountService) Login(ctx context.Context, username string, password string) (string, error) {

	userData, err := a.repo.GetUserByUserName(ctx, username, password)
	if err != nil {
		return "", err
	}

	err = helpers.VerifyPassword(userData.PasswordHash, password)
	if err != nil {
		return "", errors.New("Invalid username or password")
	}

	token, _ := helpers.GenerateJWT(userData.ID, a.conf.JWTSecret)

	return token, nil

}
