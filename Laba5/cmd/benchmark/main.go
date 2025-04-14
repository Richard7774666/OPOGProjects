package main

import (
	"fmt"
	"log"
	"order-management-api/config"
	"order-management-api/internal/csv"
	"order-management-api/repository"
	"os"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := repository.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Set up logger
	logger := log.New(os.Stdout, "[Benchmark] ", log.LstdFlags)

	// Open test CSV file
	file, err := os.Open("sample_products.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Clear the products table
	_, err = db.Exec("TRUNCATE products CASCADE")
	if err != nil {
		log.Fatalf("Failed to clear products table: %v", err)
	}

	// Test standard sequential import
	productRepo := repository.NewProductRepository(db)
	importer := csv.NewProductImporter(productRepo, logger)

	startTime := time.Now()
	count, err := importer.Import(file)
	duration := time.Since(startTime)
	if err != nil {
		log.Fatalf("Sequential import failed: %v", err)
	}

	fmt.Printf("Sequential import: %d products in %v\n", count, duration)

	// Reset file position
	file.Seek(0, 0)

	// Clear the products table
	_, err = db.Exec("TRUNCATE products CASCADE")
	if err != nil {
		log.Fatalf("Failed to clear products table: %v", err)
	}

	// Test parallel import
	startTime = time.Now()
	count, err = importer.ImportParallel(file, 4)
	duration = time.Since(startTime)
	if err != nil {
		log.Fatalf("Parallel import failed: %v", err)
	}

	fmt.Printf("Parallel import (4 workers): %d products in %v\n", count, duration)

	// Reset file position
	file.Seek(0, 0)

	// Clear the products table
	_, err = db.Exec("TRUNCATE products CASCADE")
	if err != nil {
		log.Fatalf("Failed to clear products table: %v", err)
	}

	// Test bulk import
	bulkImporter := csv.NewBulkProductImporter(db, logger)
	bulkImporter.SetBatchSize(500)

	startTime = time.Now()
	count, err = bulkImporter.Import(file)
	duration = time.Since(startTime)
	if err != nil {
		log.Fatalf("Bulk import failed: %v", err)
	}

	fmt.Printf("Bulk import (batch size 500): %d products in %v\n", count, duration)
}
