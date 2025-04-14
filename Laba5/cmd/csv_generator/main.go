package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/csv_generator/main.go [products|orders] [count]")
	}

	dataType := os.Args[1]
	count := 10000 // default count

	if len(os.Args) >= 3 {
		var err error
		count, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid count: %v", err)
		}
	}

	rand.Seed(time.Now().UnixNano())

	switch dataType {
	case "products":
		generateProducts(count)
	case "orders":
		generateOrders(count)
	default:
		log.Fatalf("Unknown data type: %s", dataType)
	}
}

func generateProducts(count int) {
	file, err := os.Create("sample_products.csv")
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"name", "unit_price", "category"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Error writing header: %v", err)
	}

	categories := []string{"Electronics", "Clothing", "Food", "Books", "Toys", "Home", "Garden", "Tools", "Sports", "Beauty"}

	// Generate product records
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("Product %d", i+1)
		price := fmt.Sprintf("%.2f", 5.0+rand.Float64()*995.0) // Prices between $5 and $1000
		category := categories[rand.Intn(len(categories))]

		record := []string{name, price, category}
		if err := writer.Write(record); err != nil {
			log.Fatalf("Error writing record: %v", err)
		}
	}

	fmt.Printf("Generated %d product records in sample_products.csv\n", count)
}

func generateOrders(count int) {
	file, err := os.Create("sample_orders.csv")
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"customer_name", "product_id", "quantity"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Error writing header: %v", err)
	}

	// Generate order records
	for i := 0; i < count; i++ {
		customerName := fmt.Sprintf("Customer %d", i%1000+1) // Reuse customer names
		productID := strconv.Itoa(rand.Intn(100) + 1)        // Assume product IDs from 1-100
		quantity := strconv.Itoa(rand.Intn(10) + 1)          // Quantities from 1-10

		record := []string{customerName, productID, quantity}
		if err := writer.Write(record); err != nil {
			log.Fatalf("Error writing record: %v", err)
		}
	}

	fmt.Printf("Generated %d order records in sample_orders.csv\n", count)
}
