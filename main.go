package main

import (
	"receipt-processor/handlers"
	"receipt-processor/services"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	receiptService := services.NewReceiptService()
	receiptHandler := handlers.NewReceiptHandler(receiptService)

	router.POST("/receipts/process", receiptHandler.ProcessReceipt)
	router.GET("/receipts/:id/points", receiptHandler.GetPoints)

	router.Run(":8080")
}
