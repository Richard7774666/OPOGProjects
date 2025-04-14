package service

import (
	"errors"
	"order-management-api/models"
	"order-management-api/repository"
)

// OrderService handles business logic for orders
type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
}

// NewOrderService creates a new order service
func NewOrderService(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(req models.CreateOrderRequest) (models.Order, error) {
	// Validate request
	if req.CustomerName == "" {
		return models.Order{}, errors.New("customer name is required")
	}

	order := models.Order{
		CustomerName: req.CustomerName,
	}

	return s.orderRepo.Create(order)
}

// GetOrderByID retrieves an order by its ID
func (s *OrderService) GetOrderByID(id int) (models.Order, error) {
	if id <= 0 {
		return models.Order{}, errors.New("invalid order ID")
	}
	return s.orderRepo.GetByID(id)
}

// GetAllOrders retrieves all orders
func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAll()
}

// AddItemToOrder adds an item to an order
func (s *OrderService) AddItemToOrder(orderID int, req models.AddItemRequest) error {
	// Validate request
	if orderID <= 0 {
		return errors.New("invalid order ID")
	}
	if req.ProductID <= 0 {
		return errors.New("invalid product ID")
	}
	if req.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	// Check if order exists
	_, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	// Check if product exists and get its price
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	// Add item to order
	return s.orderRepo.AddItem(orderID, req.ProductID, req.Quantity, product.UnitPrice)
}

// GetOrderWithHighestTotal retrieves the order with the highest total cost
func (s *OrderService) GetOrderWithHighestTotal() (models.Order, error) {
	return s.orderRepo.GetOrderWithHighestTotal()
}
