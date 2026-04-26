package transaction

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentController struct {
	service IPaymentService
}

func NewPaymentController(service IPaymentService) *PaymentController {
	return &PaymentController{service: service}
}

func (pc *PaymentController) CreateCheckout(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	resp, err := pc.service.CreateCheckout(c.Request.Context(), userID, orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (pc *PaymentController) HandleWebhook(c *gin.Context) {
	if err := pc.service.HandleWebhook(c.Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (pc *PaymentController) GetPaymentStatus(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	status, err := pc.service.GetPaymentStatus(c.Request.Context(), userID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": status})
}

func userIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}
