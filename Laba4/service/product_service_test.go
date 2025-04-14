package service

import (
	"errors"
	"order-management-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of the product repository interface
type MockProductRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockProductRepository) Create(product models.Product) (models.Product, error) {
	args := m.Called(product)
	return args.Get(0).(models.Product), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockProductRepository) GetByID(id int) (models.Product, error) {
	args := m.Called(id)
	return args.Get(0).(models.Product), args.Error(1)
}

// GetAll mocks the GetAll method
func (m *MockProductRepository) GetAll() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

func TestCreateProduct_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	req := models.CreateProductRequest{
		Name:      "Test Product",
		UnitPrice: 10.50,
		Category:  "Test Category",
	}

	expectedProduct := models.Product{
		ID:        1,
		Name:      req.Name,
		UnitPrice: req.UnitPrice,
		Category:  req.Category,
	}

	mockRepo.On("Create", mock.MatchedBy(func(p models.Product) bool {
		return p.Name == req.Name && p.UnitPrice == req.UnitPrice && p.Category == req.Category
	})).Return(expectedProduct, nil)

	// Act
	product, err := service.CreateProduct(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_ValidationErrors(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	testCases := []struct {
		name    string
		request models.CreateProductRequest
		errMsg  string
	}{
		{
			name:    "Empty name",
			request: models.CreateProductRequest{Name: "", UnitPrice: 10.5, Category: "Test"},
			errMsg:  "product name is required",
		},
		{
			name:    "Zero price",
			request: models.CreateProductRequest{Name: "Test", UnitPrice: 0, Category: "Test"},
			errMsg:  "product price must be positive",
		},
		{
			name:    "Negative price",
			request: models.CreateProductRequest{Name: "Test", UnitPrice: -5.0, Category: "Test"},
			errMsg:  "product price must be positive",
		},
		{
			name:    "Empty category",
			request: models.CreateProductRequest{Name: "Test", UnitPrice: 10.5, Category: ""},
			errMsg:  "product category is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			product, err := service.CreateProduct(tc.request)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, tc.errMsg, err.Error())
			assert.Equal(t, models.Product{}, product)
		})
	}
}

func TestGetProductByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	expectedProduct := models.Product{
		ID:        1,
		Name:      "Test Product",
		UnitPrice: 10.50,
		Category:  "Test Category",
	}

	mockRepo.On("GetByID", 1).Return(expectedProduct, nil)

	// Act
	product, err := service.GetProductByID(1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_InvalidID(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	// Act
	product, err := service.GetProductByID(0)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid product ID", err.Error())
	assert.Equal(t, models.Product{}, product)
}

func TestGetProductByID_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("GetByID", 1).Return(models.Product{}, errors.New("database error"))

	// Act
	product, err := service.GetProductByID(1)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.Equal(t, models.Product{}, product)
	mockRepo.AssertExpectations(t)
}

func TestGetAllProducts_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	expectedProducts := []models.Product{
		{
			ID:        1,
			Name:      "Test Product 1",
			UnitPrice: 10.50,
			Category:  "Category A",
		},
		{
			ID:        2,
			Name:      "Test Product 2",
			UnitPrice: 15.75,
			Category:  "Category B",
		},
	}

	mockRepo.On("GetAll").Return(expectedProducts, nil)

	// Act
	products, err := service.GetAllProducts()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	assert.Equal(t, 2, len(products))
	mockRepo.AssertExpectations(t)
}

func TestGetAllProducts_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("GetAll").Return([]models.Product{}, errors.New("database error"))

	// Act
	products, err := service.GetAllProducts()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.Empty(t, products)
	mockRepo.AssertExpectations(t)
}
