package transaction

import (
	"learnapirest/internal/middleware"

	"github.com/gin-gonic/gin"
)

func PaymentRouter(app *gin.Engine, s *PaymentService, secret []byte) {
	ctrl := NewPaymentController(s)

	api := app.Group("/payments")

	// Stripe webhook must be unauthenticated (called by Stripe servers)
	api.POST("/webhook", ctrl.HandleWebhook)

	// Authenticated payment endpoints
	auth := api.Group("")
	auth.Use(middleware.JWTAuth(secret))
	auth.POST("/checkout", ctrl.CreateCheckout)
	auth.GET("/status/:orderId", ctrl.GetPaymentStatus)
}
