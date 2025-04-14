package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"order-management-api/handlers"
	"order-management-api/models"
	"order-management-api/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the database
type MockDB struct {
	mock.Mock
}

// MockOrderRepository is a mock implementation of OrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(order models.Order) (models.Order, error) {
	args := m.Called(order)
	return args.Get(0).(models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByID(id int) (models.Order, error) {
	args := m.Called(id)
	return args.Get(0).(models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetAll() ([]models.Order, error) {
	args := m.Called()
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockOrderRepository) AddItem(orderID int, productID int, quantity int, unitPrice float64) error {
	args := m.Called(orderID, productID, quantity, unitPrice)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrderWithHighestTotal() (models.Order, error) {
	args := m.Called()
	return args.Get(0).(models.Order), args.Error(1)
}

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product models.Product) (models.Product, error) {
	args := m.Called(product)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductRepository) GetByID(id int) (models.Product, error) {
	args := m.Called(id)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductRepository) GetAll() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

// setupMockServer sets up a test HTTP server with mocked dependencies
func setupMockServer() (*httptest.Server, *MockProductRepository, *MockOrderRepository) {
	mockProductRepo := new(MockProductRepository)
	mockOrderRepo := new(MockOrderRepository)

	productService := service.NewProductService(mockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	mux := http.NewServeMux()

	// Product endpoints
	mux.HandleFunc("POST /api/products", productHandler.Create)
	mux.HandleFunc("GET /api/products", productHandler.GetAll)
	mux.HandleFunc("GET /api/products/{id}", productHandler.GetByID)

	// Order endpoints
	mux.HandleFunc("POST /api/orders", orderHandler.Create)
	mux.HandleFunc("GET /api/orders", orderHandler.GetAll)
	mux.HandleFunc("GET /api/orders/highest", orderHandler.GetHighestTotal)
	mux.HandleFunc("GET /api/orders/{id}", orderHandler.GetByID)
	mux.HandleFunc("POST /api/orders/{id}/items", orderHandler.AddItem)

	server := httptest.NewServer(mux)

	return server, mockProductRepo, mockOrderRepo
}

func TestCreateProduct(t *testing.T) {
	// Setup
	server, mockProductRepo, _ := setupMockServer()
	defer server.Close()

	// Define expected behavior
	expectedProduct := models.Product{
		ID:        1,
		Name:      "Test Product",
		UnitPrice: 10.99,
		Category:  "Test Category",
	}

	mockProductRepo.On("Create", mock.MatchedBy(func(p models.Product) bool {
		return p.Name == "Test Product" && p.UnitPrice == 10.99 && p.Category == "Test Category"
	})).Return(expectedProduct, nil)

	// Create request body
	reqBody := models.CreateProductRequest{
		Name:      "Test Product",
		UnitPrice: 10.99,
		Category:  "Test Category",
	}

	jsonBody, _ := json.Marshal(reqBody)

	// Make request
	resp, err := http.Post(server.URL+"/api/products", "application/json", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Parse response
	var product models.Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, expectedProduct.ID, product.ID)
	assert.Equal(t, expectedProduct.Name, product.Name)
	assert.Equal(t, expectedProduct.UnitPrice, product.UnitPrice)
	assert.Equal(t, expectedProduct.Category, product.Category)

	mockProductRepo.AssertExpectations(t)
}

func TestGetHighestTotal(t *testing.T) {
	// Setup
	server, _, mockOrderRepo := setupMockServer()
	defer server.Close()

	// Define expected behavior
	expectedOrder := models.Order{
		ID:           1,
		CustomerName: "Test Customer",
		CreatedAt:    "2025-04-09T14:36:48Z",
		Items: []models.OrderItem{
			{
				ID:        1,
				OrderID:   1,
				ProductID: 1,
				Quantity:  2,
				UnitPrice: 10.99,
			},
		},
		TotalCost: 21.98,
	}

	mockOrderRepo.On("GetOrderWithHighestTotal").Return(expectedOrder, nil)

	// Make request
	resp, err := http.Get(server.URL + "/api/orders/highest")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var order models.Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, expectedOrder.ID, order.ID)
	assert.Equal(t, expectedOrder.CustomerName, order.CustomerName)
	assert.Equal(t, expectedOrder.TotalCost, order.TotalCost)

	mockOrderRepo.AssertExpectations(t)
}
