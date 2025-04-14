package main

import (
	"fmt"
	"log"
	"net/http"

	"order-management-api/config"
	"order-management-api/handlers"
	"order-management-api/repository"
	"order-management-api/service"
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

	// Start server
	fmt.Printf("Server starting on port %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
