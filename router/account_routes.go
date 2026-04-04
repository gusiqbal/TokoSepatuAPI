package router

import (
	"learnapirest/controller"
	// "learnapirest/middleware"
	"learnapirest/service"

	"github.com/gin-gonic/gin"
)

func AccountRouter(app *gin.Engine, a *service.AccountService, secret []byte) {
	accountCtrl := controller.NewAccountController(a)

	api := app.Group("/account")

	api.POST("/create", accountCtrl.CreateAccount)
	api.POST("/login", accountCtrl.Login)
}
