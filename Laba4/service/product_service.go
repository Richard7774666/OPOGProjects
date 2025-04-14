package service

import (
	"errors"
	"order-management-api/models"
	"order-management-api/repository"
)

// ProductService handles business logic for products
type ProductService struct {
	repo *repository.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(req models.CreateProductRequest) (models.Product, error) {
	// Validate request
	if req.Name == "" {
		return models.Product{}, errors.New("product name is required")
	}
	if req.UnitPrice <= 0 {
		return models.Product{}, errors.New("product price must be positive")
	}
	if req.Category == "" {
		return models.Product{}, errors.New("product category is required")
	}

	product := models.Product{
		Name:      req.Name,
		UnitPrice: req.UnitPrice,
		Category:  req.Category,
	}

	return s.repo.Create(product)
}

// GetProductByID retrieves a product by its ID
func (s *ProductService) GetProductByID(id int) (models.Product, error) {
	if id <= 0 {
		return models.Product{}, errors.New("invalid product ID")
	}
	return s.repo.GetByID(id)
}

// GetAllProducts retrieves all products
func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.repo.GetAll()
}
