package repository

import (
	"database/sql"
	"fmt"
	"order-management-api/models"

	"github.com/jmoiron/sqlx"
)

// OrderRepository handles database operations for orders
type OrderRepository struct {
	db *sqlx.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create adds a new order to the database
func (r *OrderRepository) Create(order models.Order) (models.Order, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return models.Order{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO orders (customer_name)
		VALUES ($1)
		RETURNING id, customer_name, created_at
	`

	var createdOrder models.Order
	err = tx.QueryRowx(query, order.CustomerName).StructScan(&createdOrder)
	if err != nil {
		return models.Order{}, err
	}

	err = tx.Commit()
	if err != nil {
		return models.Order{}, err
	}

	return createdOrder, nil
}

// GetByID retrieves an order by its ID with all its items
func (r *OrderRepository) GetByID(id int) (models.Order, error) {
	order := models.Order{}
	orderQuery := `SELECT id, customer_name, created_at FROM orders WHERE id = $1`
	err := r.db.Get(&order, orderQuery, id)
	if err != nil {
		return models.Order{}, err
	}

	itemsQuery := `
		SELECT id, order_id, product_id, quantity, unit_price
		FROM order_items
		WHERE order_id = $1
	`
	err = r.db.Select(&order.Items, itemsQuery, id)
	if err != nil && err != sql.ErrNoRows {
		return models.Order{}, err
	}

	// Calculate total cost
	var totalCost float64
	for _, item := range order.Items {
		totalCost += item.UnitPrice * float64(item.Quantity)
	}
	order.TotalCost = totalCost

	return order, nil
}

// GetAll retrieves all orders
func (r *OrderRepository) GetAll() ([]models.Order, error) {
	var orders []models.Order
	query := `SELECT id, customer_name, created_at FROM orders`
	err := r.db.Select(&orders, query)
	if err != nil {
		return nil, err
	}

	// For each order, get its items and calculate total cost
	for i := range orders {
		itemsQuery := `
			SELECT id, order_id, product_id, quantity, unit_price
			FROM order_items
			WHERE order_id = $1
		`
		err = r.db.Select(&orders[i].Items, itemsQuery, orders[i].ID)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		// Calculate total cost
		var totalCost float64
		for _, item := range orders[i].Items {
			totalCost += item.UnitPrice * float64(item.Quantity)
		}
		orders[i].TotalCost = totalCost
	}

	return orders, nil
}

// AddItem adds an item to an order
func (r *OrderRepository) AddItem(orderID int, productID int, quantity int, unitPrice float64) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, unit_price)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, orderID, productID, quantity, unitPrice)
	return err
}

// GetOrderWithHighestTotal retrieves the order with the highest total cost
func (r *OrderRepository) GetOrderWithHighestTotal() (models.Order, error) {
	query := `
        SELECT o.id, o.customer_name, o.created_at, COALESCE(SUM(oi.quantity * oi.unit_price), 0) as total_cost
        FROM orders o
        LEFT JOIN order_items oi ON o.id = oi.order_id
        GROUP BY o.id
        ORDER BY total_cost DESC
        LIMIT 1
    `

	var order models.Order
	err := r.db.QueryRowx(query).StructScan(&order)
	if err != nil {
		fmt.Printf("Database error: %v\n", err) // Add explicit error logging
		return models.Order{}, fmt.Errorf("database error: %v", err)
	}

	// Get items for this order
	itemsQuery := `
        SELECT id, order_id, product_id, quantity, unit_price
        FROM order_items
        WHERE order_id = $1
    `
	err = r.db.Select(&order.Items, itemsQuery, order.ID)
	if err != nil && err != sql.ErrNoRows {
		return models.Order{}, err
	}

	return order, nil
}
