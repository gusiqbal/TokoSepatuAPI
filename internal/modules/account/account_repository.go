package account

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAccountRepository interface {
	CreateAccount(ctx context.Context, request *RegisterUserRequest) error
	GetUserByUserName(ctx context.Context, username string, password string) (*User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) error
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

	if err := a.db.WithContext(ctx).Where("user_name = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Username Does not exist!") // Pesan disamarkan demi keamanan
		}
	}

	return &user, nil
}

func (a *AccountRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	var user User

	if err := a.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User Does not exist!")
		}
	}

	return &user, nil
}

func (a *AccountRepository) UpdateUser(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) error {
	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = *req.PhoneNumber
	}
	if len(updates) == 0 {
		return nil
	}
	return a.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(updates).Error
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
		Level:         "user",
		CreatedAt:     now,
		LastUpdatedAt: now,
	}

	return a.db.WithContext(ctx).Create(&newUser).Error
}
