package product

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductController struct {
	SepatuService IProductService
}

func NewProductController(sepatuService IProductService) *ProductController {
	return &ProductController{
		SepatuService: sepatuService,
	}
}

func (s *ProductController) CreateSepatu(ginc *gin.Context) {
	var input CreateProductRequest
	if err := ginc.ShouldBindJSON(&input); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.SepatuService.CreateSepatu(ginc.Request.Context(), &input)
	if err != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginc.JSON(http.StatusOK, gin.H{"data": input})
}

func (s *ProductController) GetSepatu(ginc *gin.Context) {
	sepatus, err := s.SepatuService.GetSepatu(ginc.Request.Context())
	if err != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginc.JSON(http.StatusOK, gin.H{"data": sepatus})
}

func (s *ProductController) GetSepatuByID(ginc *gin.Context) {
	id, err := uuid.Parse(ginc.Param("id"))
	if err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	detail, err := s.SepatuService.GetSepatuByID(ginc.Request.Context(), id)
	if err != nil {
		ginc.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	ginc.JSON(http.StatusOK, gin.H{"data": detail})
}

func (s *ProductController) DeleteSepatu(ginc *gin.Context) {
	input := new(string)
	if err := ginc.ShouldBindJSON(&input); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.SepatuService.DeleteSepatu(ginc.Request.Context(), input)
	if err != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginc.JSON(http.StatusOK, "Data has been successfully deleted")
}

func (s *ProductController) UpdateSepatu(ginc *gin.Context) {
	var input UpdateProductRequest
	if errbind := ginc.ShouldBindJSON(&input); errbind != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": errbind.Error()})
		return
	}

	id, errparse := uuid.Parse(*input.ID)
	if errparse != nil {
		log.Println(id)
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": "id is invalid"})
		return
	}

	errupdate := s.SepatuService.UpdateSepatu(ginc.Request.Context(), &input, id)
	if errupdate != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": errupdate.Error()})
		return
	}

	ginc.JSON(http.StatusOK, gin.H{"message": "Data has been successfully updated", "data": input})
}

func (s *ProductController) LikeProduct(ginc *gin.Context) {
	var req LikeProductRequest

	if err := ginc.ShouldBindJSON(&req); err != nil {
		ginc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errLike := s.SepatuService.LikeProduct(ginc.Request.Context(), &req); errLike != nil {
		ginc.JSON(http.StatusInternalServerError, gin.H{"error": errLike.Error()})
		return
	}

	ginc.JSON(http.StatusOK, gin.H{"message": "product liked"})
}
