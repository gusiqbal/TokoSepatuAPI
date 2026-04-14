package account

import (
	"learnapirest/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AccountRouter(app *gin.Engine, a *AccountService, secret []byte) {
	accountCtrl := NewAccountController(a)

	api := app.Group("/account")

	api.POST("/create", accountCtrl.CreateAccount)
	api.POST("/login", accountCtrl.Login)
	api.POST("/logout", accountCtrl.Logout)
	api.POST("/refresh", accountCtrl.RefreshToken)

	// Protected profile routes
	profile := api.Group("/profile")
	profile.Use(middleware.JWTAuth(secret))
	profile.GET("", accountCtrl.GetProfile)
	profile.PUT("", accountCtrl.UpdateProfile)
}
