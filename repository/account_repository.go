package repository

import (
	"context"
	"errors"
	"learnapirest/helpers"
	"learnapirest/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAccountRepository interface {
	CreateAccount(ctx context.Context, account *model.User) error
	GetUserByUserName(ctx context.Context, username string, password string) (*model.User, error)
}

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (a *AccountRepository) GetUserByUserName(ctx context.Context, username string, password string) (*model.User, error) {
	account := new(model.User) // Asumsi nama struct Anda adalah model.Account

	if err := a.db.WithContext(ctx).Where("user_name = ?", username).First(account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Username Doesnt Exist!") // Pesan disamarkan demi keamanan
		}
	}

	err := helpers.VerifyPassword(account.PasswordHash, password)
	if err != nil {
		return nil, errors.New("Password does not match!")
	}

	return account, nil
}

func (a *AccountRepository) CreateAccount(ctx context.Context, account *model.User) error {
	var existing model.User

	err := a.db.WithContext(ctx).
		Where("email = ?", account.Email).
		First(&existing).Error

	if err == nil {
		return errors.New("email already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	now := time.Now().Unix()
	account.ID = uuid.New()
	account.CreatedAt = now
	account.LastUpdatedAt = now

	return a.db.WithContext(ctx).Create(account).Error
}
