package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"learnapirest/model"
	"learnapirest/service"
	"log"
	"net/http"
)

func CreateSepatu(ginc *gin.Context) {
	var input model.Sepatu
	if err := ginc.ShouldBindJSON(&input); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.CreateSepatu(ginc.Request.Context(), &input)
	if err != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ginc.JSON(http.StatusOK, gin.H{"data": input})
}

func GetSepatu(ginc *gin.Context) {
	sepatus, err := service.GetSepatu(ginc.Request.Context())
	if err != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ginc.JSON(http.StatusOK, gin.H{"data": sepatus})
}

func DeleteSepatu(ginc *gin.Context) {
	var input model.DeleteSepatu
	if err := ginc.ShouldBindJSON(&input); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.DeleteSepatu(ginc.Request.Context(), &input)
	if err != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ginc.JSON(http.StatusOK, "Data has been successfully deleted")
}

func UpdateSepatu(ginc *gin.Context) {
	var input model.UpdateSepatu
	if errbind := ginc.ShouldBindJSON(&input); errbind != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": errbind.Error()})
		return
	}
	id, errparse := uuid.Parse(input.ID.String())
	if errparse != nil {
		log.Println(id)
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": "id is invalid"})
		return
	}
	errupdate := service.UpdateSepatu(ginc.Request.Context(), &input, id)
	if errupdate != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": errupdate.Error()})
	}

	ginc.JSON(http.StatusOK, gin.H{"message": "Data has been successfully updated", "data": input})

}
