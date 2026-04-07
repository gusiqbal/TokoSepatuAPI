package product

import (
	"learnapirest/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRouter(app *gin.Engine, s *ProductService, secret []byte) {
	sepatuCtrl := NewProductController(s)

	api := app.Group("/sepatu")
	api.Use(middleware.JWTAuth(secret))

	api.POST("", sepatuCtrl.CreateSepatu)
	api.GET("", sepatuCtrl.GetSepatu)

	api.PUT("/:id", sepatuCtrl.UpdateSepatu)
	api.DELETE("/:id", sepatuCtrl.DeleteSepatu)

	api.POST("like", sepatuCtrl.LikeProduct)
}
