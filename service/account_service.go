package service

import (
	"context"
	"learnapirest/config"
	"learnapirest/helpers"
	"learnapirest/model"
	"learnapirest/repository"
)

type IAccountService interface {
	CreateAccount(ctx context.Context, account *model.User) error
	Login(ctx context.Context, username string, password string) (string, error)
}

type AccountService struct {
	repo repository.AccountRepository
	conf config.Config
}

func NewAccountService(repo repository.AccountRepository, config config.Config) *AccountService {
	return &AccountService{
		repo: repo,
		conf: config,
	}
}

func (a *AccountService) CreateAccount(ctx context.Context, account *model.User) error {
	hashed, err := helpers.HashPassword(account.PasswordHash)
	if err != nil {
		return err
	}
	account.PasswordHash = hashed

	return a.repo.CreateAccount(ctx, account)
}

func (a *AccountService) Login(ctx context.Context, username string, password string) (string, error) {
	userData, err := a.repo.GetUserByUserName(ctx, username, password)
	if err != nil {
		return "", err
	}

	token, _ := helpers.GenerateJWT(userData.ID, a.conf.JWTSecret)

	return token, nil

}
