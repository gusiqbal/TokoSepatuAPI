package controller

import (
	"learnapirest/model"
	"learnapirest/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	IAccountService service.IAccountService
}

func NewAccountController(accountService service.IAccountService) *AccountController {
	return &AccountController{
		IAccountService: accountService,
	}
}

func (a *AccountController) CreateAccount(ginc *gin.Context) {
	input := new(model.User)
	if err := ginc.ShouldBindJSON(input); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.IAccountService.CreateAccount(ginc, input); err != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginc.JSON(http.StatusOK, gin.H{"message": "Account has been created", "username": input.UserName})

}

func (a *AccountController) Login(ginc *gin.Context) {
	input := new(model.LoginInput)

	if err := ginc.ShouldBindJSON(input); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.IAccountService.Login(ginc, input.Username, input.Password)

	if err != nil {
		if err.Error() == "username atau password salah" {
			ginc.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ginc.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan pada server"})
		return
	}

	ginc.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
	})
}
