package router

import (
	"learnapirest/controller"
	"learnapirest/middleware"
	"learnapirest/service"

	"github.com/gin-gonic/gin"
)

func SepatuRouter(app *gin.Engine, s *service.SepatuService, secret []byte) {
	sepatuCtrl := controller.NewSepatuController(s)

	api := app.Group("/sepatu")
	api.Use(middleware.JWTAuth(secret))

	api.POST("", sepatuCtrl.CreateSepatu)
	api.GET("", sepatuCtrl.GetSepatu)

	api.PUT("/:id", sepatuCtrl.UpdateSepatu)
	api.DELETE("/:id", sepatuCtrl.DeleteSepatu)
}
