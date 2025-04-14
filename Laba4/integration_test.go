// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"order-management-api/config"
	"order-management-api/models"
	"order-management-api/repository"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupIntegrationTest(t *testing.T) *http.ServeMux {
	// Load test environment variables
	err := godotenv.Load(".env.test")
	if err != nil {
		t.Log("Warning: .env.test file not found, using default test environment")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	assert.NoError(t, err)

	// Initialize database connection
	db, err := repository.InitDB(cfg.DatabaseURL)
	assert.NoError(t, err)

	// Clean up database for testing
	_, err = db.Exec("TRUNCATE products, orders, order_items RESTART IDENTITY CASCADE")
	assert.NoError(t, err)

	// Set up router with actual dependencies
	router := setupRouter(db)
	
	return router
}

func TestIntegrationCreateAndGetProduct(t *testing.T) {
	router := setupIntegrationTest(t)
	
	// Create a product
	createReq := models.CreateProductRequest{
		Name:      "Integration Test Product",
		UnitPrice: 15.99,
		Category:  "Test",
	}
	
	createBody, _ := json.Marshal(createReq)
	
	// Make create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/products", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var createdProduct models.Product
	json.Unmarshal(w.Body.Bytes(), &createdProduct)
	
	// Test getting the product
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/"+string(createdProduct.ID), nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var retrievedProduct models.Product
	json.Unmarshal(w.Body.Bytes(), &retrievedProduct)
	
	assert.Equal(t, createdProduct.ID, retrievedProduct.ID)
	assert.Equal(t, createReq.Name, retrievedProduct.Name)
	assert.Equal(t, createReq.UnitPrice, retrievedProduct.UnitPrice)
	assert.Equal(t, createReq.Category, retrievedProduct.Category)
}

// Helper function to set up router with real dependencies
func setupRouter(db *sqlx.DB) *http.ServeMux {
	// Initialize repository
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Initialize service
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Setup router
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

	return mux
}
