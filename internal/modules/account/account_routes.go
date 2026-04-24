package account

import (
	"learnapirest/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AccountRouter(app *gin.Engine, a *AccountService, secret []byte) {
	accountCtrl := NewAccountController(a)
	account := app.Group("/accounts")
	account.POST("", accountCtrl.CreateAccount) // create account

	session := app.Group("/sessions")
	session.POST("", accountCtrl.Login)         // login (create session)
	session.DELETE("", accountCtrl.Logout)      // logout (destroy session)
	session.POST("/refresh", accountCtrl.RefreshToken)

	profile := app.Group("/profile")
	profile.Use(middleware.JWTAuth(secret))
	profile.GET("", accountCtrl.GetProfile)
	profile.PUT("", accountCtrl.UpdateProfile)
}
