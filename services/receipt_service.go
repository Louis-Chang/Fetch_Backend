package services

import (
	"math"
	"receipt-processor/models"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ReceiptService struct {
	Receipts map[string]models.Receipt
}

func NewReceiptService() *ReceiptService {
	return &ReceiptService{
		Receipts: make(map[string]models.Receipt),
	}
}

func (s *ReceiptService) CalculatePoints(receipt models.Receipt) int {
	points := 0

	// Rule 1: Points for retailer name
	alphanumeric := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(alphanumeric.FindAllString(receipt.Retailer, -1))

	// Rule 2: Round dollar amount
	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == float64(int(total)) {
		points += 50
	}

	// Rule 3: Multiple of 0.25
	if math.Mod(total*100, 25) == 0 {
		points += 25
	}

	// Rule 4: Every two items
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: Description length multiple of 3
	for _, item := range receipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	// Rule 6: Total > 10.00 (LLM rule)
	if total > 10.0 {
		points += 5
	}

	// Rule 7: Odd day
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 == 1 {
		points += 6
	}

	// Rule 8: Time between 2:00 and 4:00 PM
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	hour := purchaseTime.Hour()
	if hour >= 14 && hour < 16 {
		points += 10
	}

	return points
}
