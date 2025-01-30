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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	h.service.Receipts[id] = receipt

	c.JSON(http.StatusOK, models.ReceiptResponse{ID: id})
}

func (h *ReceiptHandler) GetPoints(c *gin.Context) {
	id := c.Param("id")

	receipt, exists := h.service.Receipts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	points := h.service.CalculatePoints(receipt)
	c.JSON(http.StatusOK, models.PointsResponse{Points: points})
}
