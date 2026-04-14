package order

import (
	"learnapirest/internal/middleware"

	"github.com/gin-gonic/gin"
)

func OrderRouter(app *gin.Engine, s *OrderService, secret []byte) {
	ctrl := NewOrderController(s)

	api := app.Group("/orders")
	api.Use(middleware.JWTAuth(secret))

	api.POST("", ctrl.CreateOrder)
	api.GET("", ctrl.GetOrderHistory)
	api.GET("/:id", ctrl.GetOrderDetail)
}
