package router

import (
	"learnapirest/internal/config"
	"learnapirest/internal/modules/account"
	"learnapirest/internal/modules/cart"
	"learnapirest/internal/modules/order"
	"learnapirest/internal/modules/product"

	"learnapirest/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	s *product.ProductService,
	a *account.AccountService,
	c *cart.CartService,
	o *order.OrderService,
	config *config.Config,
) *gin.Engine {
	limiter := middleware.NewRateLimiter(100, time.Minute)
	secret := []byte(config.JWTSecret)
	app := gin.New()

	app.Use(
		middleware.Recovery(),
		limiter.Middleware(),
		middleware.Logger(),
		middleware.CorsSetup(),
	)

	product.ProductRouter(app, s, secret)
	account.AccountRouter(app, a, secret)
	cart.CartRouter(app, c, secret)
	order.OrderRouter(app, o, secret)

	return app
}
