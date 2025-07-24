package router

import (
	"github.com/gin-gonic/gin"
	"learnapirest/controller"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/sepatu/create", controller.CreateSepatu)
	router.GET("/sepatu/get", controller.GetSepatu)
	router.POST("/sepatu/delete", controller.DeleteSepatu)
	router.POST("/sepatu/update", controller.UpdateSepatu)

	return router
}
