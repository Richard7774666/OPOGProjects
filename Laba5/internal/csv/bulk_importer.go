package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

// BulkProductImporter handles importing products in bulk for better performance
type BulkProductImporter struct {
	db        *sqlx.DB
	logger    *log.Logger
	batchSize int
}

// NewBulkProductImporter creates a new bulk product importer
func NewBulkProductImporter(db *sqlx.DB, logger *log.Logger) *BulkProductImporter {
	return &BulkProductImporter{
		db:        db,
		logger:    logger,
		batchSize: 1000, // Default batch size
	}
}

// Import imports products from a CSV file in bulk
func (i *BulkProductImporter) Import(reader io.Reader) (int, error) {
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

	i.logger.Println("Starting bulk product import")
	startTime := time.Now()

	// Begin transaction
	tx, err := i.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Prepare statement for bulk insert
	stmt, err := tx.Prepare(`
		INSERT INTO products (name, unit_price, category)
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Process rows in batches
	count := 0
	batch := 0

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			i.logger.Printf("Error reading CSV row: %v", err)
			continue
		}

		// Parse product data
		if len(record) != 3 {
			i.logger.Printf("Invalid record format, expected 3 fields")
			continue
		}

		name := record[0]
		price, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			i.logger.Printf("Invalid price format: %v", err)
			continue
		}
		category := record[2]

		// Execute insert
		_, err = stmt.Exec(name, price, category)
		if err != nil {
			i.logger.Printf("Error inserting product: %v", err)
			continue
		}

		count++
		batch++

		// Commit in batches to avoid huge transactions
		if batch >= i.batchSize {
			if err := tx.Commit(); err != nil {
				i.logger.Printf("Error committing batch: %v", err)
				tx.Rollback()
				return count - batch, err
			}

			// Start new transaction
			tx, err = i.db.Beginx()
			if err != nil {
				return count, fmt.Errorf("failed to begin new transaction: %w", err)
			}

			// Prepare new statement
			stmt, err = tx.Prepare(`
				INSERT INTO products (name, unit_price, category)
				VALUES ($1, $2, $3)
			`)
			if err != nil {
				tx.Rollback()
				return count, fmt.Errorf("failed to prepare statement: %w", err)
			}

			i.logger.Printf("Committed batch of %d products (%d total)", batch, count)
			batch = 0
		}
	}

	// Commit final batch
	if batch > 0 {
		if err := tx.Commit(); err != nil {
			i.logger.Printf("Error committing final batch: %v", err)
			tx.Rollback()
			return count - batch, err
		}
		i.logger.Printf("Committed final batch of %d products", batch)
	}

	duration := time.Since(startTime)
	i.logger.Printf("Bulk import completed: %d products imported in %v", count, duration)

	return count, nil
}

// SetBatchSize sets the batch size for bulk imports
func (i *BulkProductImporter) SetBatchSize(size int) {
	if size > 0 {
		i.batchSize = size
	}
}
