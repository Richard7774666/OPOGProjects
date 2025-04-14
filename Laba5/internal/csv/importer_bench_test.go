package csv

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"os"
	"testing"
)

func BenchmarkParseProductRecord(b *testing.B) {
	record := []string{"Test Product", "19.99", "Test Category"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parseProductRecord(record)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCsvReading(b *testing.B) {
	data := setupProductTestData(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		csvReader := csv.NewReader(reader)
		header, _ := csvReader.Read() // Skip header

		for {
			_, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkMemoryUsage(b *testing.B) {
	repo := &MockProductRepository{}
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	importer := NewProductImporter(repo, logger)

	data := setupProductTestData(50000) // Large dataset to test memory usage

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		importer.Import(reader)
	}
}
