package router

import (
	"learnapirest/config"
	"learnapirest/service"

	"learnapirest/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(s *service.SepatuService, a *service.AccountService, config *config.Config) *gin.Engine {
	limiter := middleware.NewRateLimiter(100, time.Minute)
	secret := []byte(config.JWTSecret)
	app := gin.New()

	app.Use(
		middleware.Recovery(),
		limiter.Middleware(),
		middleware.Logger(),
	)

	SepatuRouter(app, s, secret)
	AccountRouter(app, a, secret)

	return app
}
