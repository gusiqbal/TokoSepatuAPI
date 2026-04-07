package account

import (
	"context"
	"learnapirest/helpers"
	"learnapirest/internal/config"
	"net/http"
)

type IAccountService interface {
	CreateAccount(ctx context.Context, account *RegisterUserRequest) error
	Login(ctx context.Context, username string, password string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
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

func (a *AccountService) Login(ctx context.Context, username string, password string) (*TokenResponse, error) {

	userData, err := a.repo.GetUserByUserName(ctx, username, password)
	if err != nil {
		return nil, err
	}

	err = helpers.VerifyPassword(userData.PasswordHash, password)
	if err != nil {
		return nil, helpers.NewError(http.StatusUnauthorized, "Invalid username or password")
	}

	accessToken, _ := helpers.GenerateAccessToken(userData.ID, a.conf.JWTSecret)
	RefreshToken, _ := helpers.GenerateRefreshToken(userData.ID, a.conf.JWTSecret)
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: RefreshToken,
	}, nil

}

func (a *AccountService) Logout(ctx context.Context, refreshTokenString string) error {
	// 1. Verifikasi apakah token ini memang valid dan milik sistem kita
	_, err := helpers.VerifyJWT(refreshTokenString, a.conf.JWTSecret)
	if err != nil {
		// Jika token sudah tidak valid/kedaluwarsa, anggap saja proses logout sukses
		// (karena tujuannya memang membuat token tidak bisa dipakai)
		return nil
	}

	// 2. LOGIKA ENTERPRISE (Opsional/Next Step):
	// Di aplikasi berskala besar, kamu harus menghapus token ini dari Database atau Redis.
	// err = s.repo.DeleteRefreshToken(ctx, refreshTokenString)
	// if err != nil {
	//     return apperror.New(http.StatusInternalServerError, "Gagal memproses logout")
	// }

	return nil
}

func (a *AccountService) RefreshToken(ctx context.Context, refresToken string) (*TokenResponse, error) {
	userID, err := helpers.VerifyJWT(refresToken, a.conf.JWTSecret)

	if err != nil {
		return nil, helpers.NewError(http.StatusUnauthorized, "Refresh token is not valid!")
	}

	user, err := a.repo.GetUserByID(ctx, userID)

	if err != nil {
		return nil, helpers.NewError(http.StatusUnauthorized, "User not found, Session not valid!")
	}

	newAccessToken, err := helpers.GenerateAccessToken(user.ID, a.conf.JWTSecret)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "Failed to make access token!")
	}

	newRefreshToken, err := helpers.GenerateRefreshToken(user.ID, a.conf.JWTSecret)

	return &TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
