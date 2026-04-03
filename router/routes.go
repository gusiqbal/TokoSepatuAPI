package router

import (
	"learnapirest/controller"
	"learnapirest/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(s *service.SepatuService) *gin.Engine {
	sepatu := controller.NewSepatuController(s)

	router := gin.Default()

	router.POST("/sepatu/create", sepatu.CreateSepatu)
	router.GET("/sepatu/get", sepatu.GetSepatu)
	router.POST("/sepatu/delete", sepatu.DeleteSepatu)
	router.POST("/sepatu/update", sepatu.UpdateSepatu)

	return router
}
