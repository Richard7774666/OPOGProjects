package repository

import (
	"order-management-api/models"

	"github.com/jmoiron/sqlx"
)

// ProductRepository handles database operations for products
type ProductRepository struct {
	db *sqlx.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create adds a new product to the database
func (r *ProductRepository) Create(product models.Product) (models.Product, error) {
	query := `
		INSERT INTO products (name, unit_price, category)
		VALUES ($1, $2, $3)
		RETURNING id, name, unit_price, category
	`

	var createdProduct models.Product
	err := r.db.QueryRowx(
		query,
		product.Name,
		product.UnitPrice,
		product.Category,
	).StructScan(&createdProduct)

	return createdProduct, err
}

// GetByID retrieves a product by its ID
func (r *ProductRepository) GetByID(id int) (models.Product, error) {
	var product models.Product
	query := `SELECT id, name, unit_price, category FROM products WHERE id = $1`
	err := r.db.Get(&product, query, id)
	return product, err
}

// GetAll retrieves all products
func (r *ProductRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	query := `SELECT id, name, unit_price, category FROM products`
	err := r.db.Select(&products, query)
	return products, err
}
