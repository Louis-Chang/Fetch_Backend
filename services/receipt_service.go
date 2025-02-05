package services

import (
	"fmt"
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

func (s *ReceiptService) ValidateReceipt(receipt models.Receipt) error {
	// Validate retailer
	retailerPattern := regexp.MustCompile(`^[\w\s\-&]+$`)
	if !retailerPattern.MatchString(receipt.Retailer) {
		return fmt.Errorf("invalid retailer format")
	}

	// Validate purchaseDate
	if _, err := time.Parse("2006-01-02", receipt.PurchaseDate); err != nil {
		return fmt.Errorf("invalid date format")
	}

	// Validate purchaseTime
	if _, err := time.Parse("15:04", receipt.PurchaseTime); err != nil {
		return fmt.Errorf("invalid time format")
	}

	// Validate total
	totalPattern := regexp.MustCompile(`^\d+\.\d{2}$`)
	if !totalPattern.MatchString(receipt.Total) {
		return fmt.Errorf("invalid total format")
	}

	// Validate items
	if len(receipt.Items) < 1 {
		return fmt.Errorf("at least one item is required")
	}

	// Validate each item
	for _, item := range receipt.Items {
		if !totalPattern.MatchString(item.Price) {
			return fmt.Errorf("invalid item price format")
		}
		descPattern := regexp.MustCompile(`^[\w\s\-]+$`)
		if !descPattern.MatchString(item.ShortDescription) {
			return fmt.Errorf("invalid item description format")
		}
	}

	return nil
}

func (s *ReceiptService) CalculatePoints(receipt models.Receipt) int {
	points := 0

	// Rule 1: Points for retailer name
	alphanumeric := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(alphanumeric.FindAllString(receipt.Retailer, -1))
	fmt.Print(points)

	// Rule 2: Round dollar amount
	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == float64(int(total)) {
		points += 50
	}
	fmt.Print(points)

	// Rule 3: Multiple of 0.25
	if math.Mod(total*100, 25) == 0 {
		points += 25
	}
	fmt.Print(points)

	// Rule 4: Every two items
	points += (len(receipt.Items) / 2) * 5
	fmt.Print(points)

	// Rule 5: Description length multiple of 3
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}
	fmt.Print(points)

	// Rule 7: Odd day
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 == 1 {
		points += 6
	}
	fmt.Print(points)

	// Rule 8: Time between 2:00 and 4:00 PM
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	hour := purchaseTime.Hour()
	if hour >= 14 && hour < 16 {
		points += 10
	}
	fmt.Print(points)

	return points
}
