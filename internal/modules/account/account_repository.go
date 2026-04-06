package account

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAccountRepository interface {
	CreateAccount(ctx context.Context, account *User) error
	GetUserByUserName(ctx context.Context, username string, password string) (*User, error)
}

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (a *AccountRepository) GetUserByUserName(ctx context.Context, username string, password string) (*User, error) {
	var user User

	if err := a.db.WithContext(ctx).Where("user_name = ?", username).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Username Doesnt Exist!") // Pesan disamarkan demi keamanan
		}
	}

	return &user, nil
}

func (a *AccountRepository) CreateAccount(ctx context.Context, request *RegisterUserRequest) error {
	var existing User

	err := a.db.WithContext(ctx).
		Where("email = ?", request.Email).
		First(&existing).Error

	if err == nil {
		return errors.New("email already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	now := time.Now().Unix()

	newUser := User{
		ID:            uuid.New(),
		UserName:      request.UserName,
		PasswordHash:  request.Password,
		Email:         request.Email,
		PhoneNumber:   request.PhoneNumber,
		CreatedAt:     now,
		LastUpdatedAt: now,
	}

	return a.db.WithContext(ctx).Create(newUser).Error
}
