package models

// Product - struct representing a product
type Product struct {
	ID        int     `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	UnitPrice float64 `json:"unit_price" db:"unit_price"`
	Category  string  `json:"category" db:"category"`
}

// OrderItem - struct representing an item in an order
type OrderItem struct {
	ID        int     `json:"id" db:"id"`
	OrderID   int     `json:"order_id" db:"order_id"`
	ProductID int     `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity" db:"quantity"`
	UnitPrice float64 `json:"unit_price" db:"unit_price"`
}

// Order - struct representing an order
type Order struct {
	ID           int         `json:"id" db:"id"`
	CustomerName string      `json:"customer_name" db:"customer_name"`
	CreatedAt    string      `json:"created_at" db:"created_at"`
	Items        []OrderItem `json:"items,omitempty"`
	TotalCost    float64     `json:"total_cost,omitempty" db:"total_cost"`
}

// type Order struct {
// 	ID           int         `json:"id" db:"id"`
// 	CustomerName string      `json:"customer_name" db:"customer_name"`
// 	CreatedAt    string      `json:"created_at" db:"created_at"`
// 	Items        []OrderItem `json:"items,omitempty"`
// 	TotalCost    float64     `json:"total_cost,omitempty"`
// }

// CreateProductRequest - request to create a new product
type CreateProductRequest struct {
	Name      string  `json:"name"`
	UnitPrice float64 `json:"unit_price"`
	Category  string  `json:"category"`
}

// CreateOrderRequest - request to create a new order
type CreateOrderRequest struct {
	CustomerName string `json:"customer_name"`
}

// AddItemRequest - request to add an item to an order
type AddItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
