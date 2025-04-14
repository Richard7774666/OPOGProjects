package main

import (
	"flag"
	"fmt"
	"log"
	"order-management-api/config"
	"order-management-api/internal/csv"
	"order-management-api/repository"
	"os"
)

func main() {
	// Define command line flags
	fileFlag := flag.String("file", "", "Path to CSV file to import")
	typeFlag := flag.String("type", "", "Type of data to import (products or orders)")
	workersFlag := flag.Int("workers", 1, "Number of worker goroutines for parallel import")
	flag.Parse()

	if *fileFlag == "" || *typeFlag == "" {
		fmt.Println("Usage: go run cmd/import/main.go -file <path> -type <products|orders> [-workers <num>]")
		os.Exit(1)
	}

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
	logger := log.New(os.Stdout, "[CSV Import] ", log.LstdFlags)

	// Open CSV file
	file, err := os.Open(*fileFlag)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Create repositories
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Perform import based on type
	switch *typeFlag {
	case "products":
		importer := csv.NewProductImporter(productRepo, logger)
		var count int
		if *workersFlag > 1 {
			count, err = importer.ImportParallel(file, *workersFlag)
		} else {
			count, err = importer.Import(file)
		}
		if err != nil {
			log.Fatalf("Import failed: %v", err)
		}
		fmt.Printf("Successfully imported %d products\n", count)

	case "orders":
		importer := csv.NewOrderImporter(orderRepo, productRepo, logger)
		count, err := importer.Import(file)
		if err != nil {
			log.Fatalf("Import failed: %v", err)
		}
		fmt.Printf("Successfully imported %d orders\n", count)

	default:
		log.Fatalf("Unknown data type: %s", *typeFlag)
	}
}
