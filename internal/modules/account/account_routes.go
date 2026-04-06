package account

import (
	"github.com/gin-gonic/gin"
)

func AccountRouter(app *gin.Engine, a *AccountService, secret []byte) {
	accountCtrl := NewAccountController(a)

	api := app.Group("/account")

	api.POST("/create", accountCtrl.CreateAccount)
	api.POST("/login", accountCtrl.Login)
}
