package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"order-management-api/models"
	"order-management-api/repository"
	"strconv"
	"sync"
	"time"
)

// ProductImporter handles importing products from CSV
type ProductImporter struct {
	repo   *repository.ProductRepository
	logger *log.Logger
}

// NewProductImporter creates a new product importer
func NewProductImporter(repo *repository.ProductRepository, logger *log.Logger) *ProductImporter {
	return &ProductImporter{
		repo:   repo,
		logger: logger,
	}
}

// Import imports products from a CSV file
func (i *ProductImporter) Import(reader io.Reader) (int, error) {
	csvReader := csv.NewReader(reader)

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate header
	expectedHeader := []string{"name", "unit_price", "category"}
	if !validateHeader(header, expectedHeader) {
		return 0, fmt.Errorf("invalid CSV header, expected: %v, got: %v", expectedHeader, header)
	}

	i.logger.Println("Starting product import")
	startTime := time.Now()

	// Process rows
	count := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			i.logger.Printf("Error reading CSV row: %v", err)
			continue
		}

		product, err := parseProductRecord(record)
		if err != nil {
			i.logger.Printf("Error parsing product: %v", err)
			continue
		}

		_, err = i.repo.Create(product)
		if err != nil {
			i.logger.Printf("Error saving product to database: %v", err)
			continue
		}

		count++
		if count%1000 == 0 {
			i.logger.Printf("Imported %d products so far...", count)
		}
	}

	duration := time.Since(startTime)
	i.logger.Printf("Import completed: %d products imported in %v", count, duration)

	return count, nil
}

// ImportParallel imports products from CSV in parallel for better performance
func (i *ProductImporter) ImportParallel(reader io.Reader, workers int) (int, error) {
	csvReader := csv.NewReader(reader)

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate header
	expectedHeader := []string{"name", "unit_price", "category"}
	if !validateHeader(header, expectedHeader) {
		return 0, fmt.Errorf("invalid CSV header, expected: %v, got: %v", expectedHeader, header)
	}

	i.logger.Println("Starting parallel product import")
	startTime := time.Now()

	// Channel for work distribution
	records := make(chan []string)
	errors := make(chan error, workers)
	var wg sync.WaitGroup

	// Counter for successful imports
	var counter int
	var counterMutex sync.Mutex

	// Start worker goroutines
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range records {
				product, err := parseProductRecord(record)
				if err != nil {
					errors <- fmt.Errorf("parse error: %w", err)
					continue
				}

				_, err = i.repo.Create(product)
				if err != nil {
					errors <- fmt.Errorf("database error: %w", err)
					continue
				}

				counterMutex.Lock()
				counter++
				counterMutex.Unlock()
			}
		}()
	}

	// Producer goroutine - read CSV and send to workers
	go func() {
		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errors <- fmt.Errorf("read error: %w", err)
				continue
			}
			records <- record
		}
		close(records)
	}()

	// Error reporting goroutine
	go func() {
		errorCount := 0
		for err := range errors {
			errorCount++
			if errorCount <= 100 { // Limit logged errors
				i.logger.Printf("Import error: %v", err)
			}
		}
	}()

	// Wait for all workers to complete
	wg.Wait()
	close(errors)

	duration := time.Since(startTime)
	i.logger.Printf("Parallel import completed: %d products imported in %v using %d workers",
		counter, duration, workers)

	return counter, nil
}

// OrderImporter handles importing orders from CSV
type OrderImporter struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
	logger      *log.Logger
}

// NewOrderImporter creates a new order importer
func NewOrderImporter(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository, logger *log.Logger) *OrderImporter {
	return &OrderImporter{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		logger:      logger,
	}
}

// Import imports orders from a CSV file
func (i *OrderImporter) Import(reader io.Reader) (int, error) {
	csvReader := csv.NewReader(reader)

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate header
	expectedHeader := []string{"customer_name", "product_id", "quantity"}
	if !validateHeader(header, expectedHeader) {
		return 0, fmt.Errorf("invalid CSV header, expected: %v, got: %v", expectedHeader, header)
	}

	i.logger.Println("Starting order import")
	startTime := time.Now()

	count := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			i.logger.Printf("Error reading CSV row: %v", err)
			continue
		}

		// Parse record fields
		customerName := record[0]
		productID, err := strconv.Atoi(record[1])
		if err != nil {
			i.logger.Printf("Invalid product ID: %v", err)
			continue
		}

		quantity, err := strconv.Atoi(record[2])
		if err != nil {
			i.logger.Printf("Invalid quantity: %v", err)
			continue
		}

		// Create order
		order := models.Order{
			CustomerName: customerName,
		}

		createdOrder, err := i.orderRepo.Create(order)
		if err != nil {
			i.logger.Printf("Error creating order: %v", err)
			continue
		}

		// Get product for price
		product, err := i.productRepo.GetByID(productID)
		if err != nil {
			i.logger.Printf("Error retrieving product: %v", err)
			continue
		}

		// Add item to order
		err = i.orderRepo.AddItem(createdOrder.ID, productID, quantity, product.UnitPrice)
		if err != nil {
			i.logger.Printf("Error adding item to order: %v", err)
			continue
		}

		count++
		if count%100 == 0 {
			i.logger.Printf("Imported %d orders so far...", count)
		}
	}

	duration := time.Since(startTime)
	i.logger.Printf("Import completed: %d orders imported in %v", count, duration)

	return count, nil
}

// Helper functions

// validateHeader checks if the CSV header matches the expected format
func validateHeader(header, expectedHeader []string) bool {
	if len(header) != len(expectedHeader) {
		return false
	}

	for i, field := range expectedHeader {
		if header[i] != field {
			return false
		}
	}

	return true
}

// parseProductRecord converts a CSV record to a Product model
func parseProductRecord(record []string) (models.Product, error) {
	if len(record) != 3 {
		return models.Product{}, fmt.Errorf("invalid record format, expected 3 fields")
	}

	price, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return models.Product{}, fmt.Errorf("invalid price format: %w", err)
	}

	return models.Product{
		Name:      record[0],
		UnitPrice: price,
		Category:  record[2],
	}, nil
}
