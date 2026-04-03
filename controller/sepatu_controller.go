package controller

import (
	"learnapirest/model"
	"learnapirest/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SepatuController struct {
	SepatuService service.ISepatuService
}

func NewSepatuController(sepatuService service.ISepatuService) *SepatuController {
	return &SepatuController{
		SepatuService: sepatuService,
	}
}

func (s *SepatuController) CreateSepatu(ginc *gin.Context) {
	var input model.Sepatu
	if err := ginc.ShouldBindJSON(&input); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.SepatuService.CreateSepatu(ginc.Request.Context(), &input)
	if err != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ginc.JSON(http.StatusOK, gin.H{"data": input})
}

func (s *SepatuController) GetSepatu(ginc *gin.Context) {
	sepatus, err := s.SepatuService.GetSepatu(ginc.Request.Context())
	if err != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ginc.JSON(http.StatusOK, gin.H{"data": sepatus})
}

func (s *SepatuController) DeleteSepatu(ginc *gin.Context) {
	var input model.DeleteSepatu
	if err := ginc.ShouldBindJSON(&input); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.SepatuService.DeleteSepatu(ginc.Request.Context(), &input)
	if err != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ginc.JSON(http.StatusOK, "Data has been successfully deleted")
}

func (s *SepatuController) UpdateSepatu(ginc *gin.Context) {
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
	errupdate := s.SepatuService.UpdateSepatu(ginc.Request.Context(), &input, id)
	if errupdate != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": errupdate.Error()})
	}

	ginc.JSON(http.StatusOK, gin.H{"message": "Data has been successfully updated", "data": input})

}
