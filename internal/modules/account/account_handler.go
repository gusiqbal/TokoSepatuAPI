package account

import (
	"errors"
	"learnapirest/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	AccountService *AccountService
}

func NewAccountController(AccountService *AccountService) *AccountController {
	return &AccountController{
		AccountService: AccountService,
	}
}

func (a *AccountController) CreateAccount(ginc *gin.Context) {
	var input RegisterUserRequest
	if err := ginc.ShouldBindJSON(input); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.AccountService.CreateAccount(ginc, &input); err != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginc.JSON(http.StatusOK, gin.H{"message": "Account has been created", "username": input.UserName})

}

func (a *AccountController) Login(ginc *gin.Context) {
	input := new(LoginRequest)

	if err := ginc.ShouldBindJSON(input); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.AccountService.Login(ginc, input.Username, input.Password)

	if err != nil {
		var appErr *helpers.AppError

		if errors.As(err, &appErr) {
			ginc.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}

		ginc.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan pada server"})
		return
	}

	ginc.JSON(http.StatusOK, gin.H{
		"message": "Login success",
		"token":   token,
	})
}

func (a *AccountController) Logout(ginc *gin.Context) {

	var LogoutRequest LogoutRequest

	if err := ginc.ShouldBindJSON(LogoutRequest); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request!"})
		return
	}

	if err := a.AccountService.Logout(ginc, LogoutRequest.RefreshToken); err != nil {
		var appErr *helpers.AppError
		if errors.As(err, &appErr) {
			ginc.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan sistem"})
		return
	}

	ginc.JSON(http.StatusOK, gin.H{
		"message": "Logout Success",
	})

}

func (a *AccountController) RefreshToken(ginc *gin.Context) {
	var reqRefreshToken RefreshTokenRequest

	if err := ginc.ShouldBindJSON(reqRefreshToken); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
		return
	}

	tokenResponse, err := a.AccountService.RefreshToken(ginc, reqRefreshToken.RefreshToken)

	if err != nil {
		ginc.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ginc.JSON(http.StatusOK, tokenResponse.RefreshToken)

}
