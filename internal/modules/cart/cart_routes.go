package cart

import (
	"learnapirest/internal/middleware"

	"github.com/gin-gonic/gin"
)

func CartRouter(app *gin.Engine, s *CartService, secret []byte) {
	ctrl := NewCartController(s)

	api := app.Group("/cart")
	api.Use(middleware.JWTAuth(secret))

	api.GET("", ctrl.GetCart)
	api.POST("/items", ctrl.AddItem)
	api.PUT("/items/:id", ctrl.UpdateItem)
	api.DELETE("/items/:id", ctrl.RemoveItem)
	api.DELETE("", ctrl.ClearCart)
}
