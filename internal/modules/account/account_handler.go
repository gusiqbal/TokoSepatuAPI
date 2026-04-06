package account

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	IAccountService IAccountService
}

func NewAccountController(IAccountService IAccountService) *AccountController {
	return &AccountController{
		IAccountService: IAccountService,
	}
}

func (a *AccountController) CreateAccount(ginc *gin.Context) {
	var input RegisterUserRequest
	if err := ginc.ShouldBindJSON(input); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.IAccountService.CreateAccount(ginc, &input); err != nil {
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

	token, err := a.IAccountService.Login(ginc, input.Username, input.Password)

	if err != nil {
		if strings.ToLower(err.Error()) == "password does not match!" {
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
