package utils

import (
	"math/rand"
	"time"

	"github.com/KapDmitry/WB_L0/internal/order"
)

// generateRandomString генерирует случайную строку заданной длины
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	//rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// generateRandomItem генерирует структуру Item с рандомными значениями
func generateRandomItem() order.Item {
	return order.Item{
		ChrtID:      rand.Int63(),
		TrackNumber: generateRandomString(8),
		Price:       rand.Intn(100),
		RID:         generateRandomString(8),
		Name:        generateRandomString(6),
		Sale:        rand.Intn(50),
		Size:        generateRandomString(3),
		TotalPrice:  rand.Intn(100),
		NMID:        rand.Int63(),
		Brand:       generateRandomString(5),
		Status:      rand.Intn(200),
	}
}

// generateRandomOrder генерирует структуру Order с рандомными значениями
func GenerateRandomOrder() order.Order {
	numItems := rand.Intn(5) + 1 // Генерируем от 1 до 5 Item'ов
	var items []order.Item

	for i := 0; i < numItems; i++ {
		items = append(items, generateRandomItem())
	}

	return order.Order{
		OrderUID:    generateRandomString(10),
		TrackNumber: generateRandomString(8),
		Entry:       generateRandomString(3),
		OrderDelivery: order.Delivery{
			Name:    generateRandomString(8),
			Phone:   generateRandomString(10),
			ZIP:     generateRandomString(7),
			City:    generateRandomString(6),
			Address: generateRandomString(10),
			Region:  generateRandomString(6),
			Email:   generateRandomString(10) + "@example.com",
		},
		OrderPayment: order.Payment{
			Transaction:  generateRandomString(10),
			RequestID:    generateRandomString(8),
			Currency:     "USD",
			Provider:     generateRandomString(6),
			Amount:       rand.Intn(1000),
			PaymentDT:    int(time.Now().Unix()),
			Bank:         generateRandomString(5),
			DeliveryCost: rand.Intn(500),
			GoodsTotal:   rand.Intn(500),
			CustomFee:    rand.Intn(100),
		},
		Items:             items,
		Locale:            "en",
		InternalSignature: generateRandomString(12),
		CustomerID:        generateRandomString(6),
		DeliveryService:   generateRandomString(6),
		Shardkey:          generateRandomString(1),
		SMID:              rand.Int63(),
		DateCreated:       time.Now().Format(time.RFC3339),
		OOFShard:          generateRandomString(1),
	}
}
