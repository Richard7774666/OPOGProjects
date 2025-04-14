package csv

import (
	"bytes"
	"log"
	"order-management-api/models"
	"os"
	"testing"
)

// MockProductRepository for testing
type MockProductRepository struct {
	products []models.Product
}

func (m *MockProductRepository) Create(product models.Product) (models.Product, error) {
	product.ID = len(m.products) + 1
	m.products = append(m.products, product)
	return product, nil
}

func (m *MockProductRepository) GetByID(id int) (models.Product, error) {
	return models.Product{}, nil
}

func (m *MockProductRepository) GetAll() ([]models.Product, error) {
	return m.products, nil
}

// MockOrderRepository for testing
type MockOrderRepository struct{}

func (m *MockOrderRepository) Create(order models.Order) (models.Order, error) {
	order.ID = 1
	return order, nil
}

func (m *MockOrderRepository) GetByID(id int) (models.Order, error) {
	return models.Order{}, nil
}

func (m *MockOrderRepository) GetAll() ([]models.Order, error) {
	return []models.Order{}, nil
}

func (m *MockOrderRepository) AddItem(orderID, productID, quantity int, unitPrice float64) error {
	return nil
}

func (m *MockOrderRepository) GetOrderWithHighestTotal() (models.Order, error) {
	return models.Order{}, nil
}

// Setup helpers
func setupProductTestData(count int) []byte {
	var buf bytes.Buffer
	buf.WriteString("name,unit_price,category\n")
	for i := 0; i < count; i++ {
		buf.WriteString("Test Product,19.99,Test Category\n")
	}
	return buf.Bytes()
}

// Benchmarks
func BenchmarkProductImport(b *testing.B) {
	repo := &MockProductRepository{}
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	importer := NewProductImporter(repo, logger)

	data := setupProductTestData(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		importer.Import(reader)
	}
}

func BenchmarkProductImportParallel_2Workers(b *testing.B) {
	repo := &MockProductRepository{}
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	importer := NewProductImporter(repo, logger)

	data := setupProductTestData(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		importer.ImportParallel(reader, 2)
	}
}

func BenchmarkProductImportParallel_4Workers(b *testing.B) {
	repo := &MockProductRepository{}
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	importer := NewProductImporter(repo, logger)

	data := setupProductTestData(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		importer.ImportParallel(reader, 4)
	}
}

func BenchmarkProductImportParallel_8Workers(b *testing.B) {
	repo := &MockProductRepository{}
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	importer := NewProductImporter(repo, logger)

	data := setupProductTestData(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		importer.ImportParallel(reader, 8)
	}
}
