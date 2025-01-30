package handlers

import (
	"net/http"
	"receipt-processor/models"
	"receipt-processor/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReceiptHandler struct {
	service *services.ReceiptService
}

func NewReceiptHandler(service *services.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{service: service}
}

func (h *ReceiptHandler) ProcessReceipt(c *gin.Context) {
	var receipt models.Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please verify input. The receipt is invalid.",
		})
		return
	}

	if err := h.service.ValidateReceipt(receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please verify input. The receipt is invalid.",
		})
		return
	}

	id := uuid.New().String()
	h.service.Receipts[id] = receipt

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *ReceiptHandler) GetPoints(c *gin.Context) {
	id := c.Param("id")

	receipt, exists := h.service.Receipts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No receipt found for that ID.",
		})
		return
	}

	points := h.service.CalculatePoints(receipt)
	c.JSON(http.StatusOK, gin.H{"points": points})
}
