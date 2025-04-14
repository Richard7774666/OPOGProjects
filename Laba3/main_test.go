package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

// Basic tests - success flow

// TestNewProduct verifies that a product is created correctly
func TestNewProduct(t *testing.T) {
	// Create a test product
	p := NewProduct(1, "Test Product", 10.0, "Test Category")

	// Verify product properties
	require.Equal(t, 1, p.ID)
	require.Equal(t, "Test Product", p.Name)
	require.Equal(t, 10.0, p.UnitPrice)
	require.Equal(t, "Test Category", p.Category)
}

// TestNewOrder verifies that an order is created correctly
func TestNewOrder(t *testing.T) {
	// Create a test order
	o := NewOrder(1, "Test Customer")

	// Verify order properties
	require.Equal(t, 1, o.ID)
	require.Equal(t, "Test Customer", o.CustomerName)
	require.NotNil(t, o.Items)
	require.Empty(t, o.Items)
}

// TestAddItem verifies that items can be added to an order
func TestAddItem(t *testing.T) {
	// Create a test order and products
	o := NewOrder(1, "Test Customer")
	p1 := NewProduct(1, "Test Product 1", 10.0, "Test Category")
	p2 := NewProduct(2, "Test Product 2", 20.0, "Test Category")

	// Add items to the order
	o.AddItem(p1, 2)
	o.AddItem(p2, 1)

	// Verify items were added correctly
	require.Equal(t, 2, len(o.Items))
	require.Equal(t, 2, o.Items[p1])
	require.Equal(t, 1, o.Items[p2])

	// Test adding more of an existing item
	o.AddItem(p1, 3)
	require.Equal(t, 5, o.Items[p1]) // Should now be 2+3=5
}

// TestTotalCost verifies that the total cost is calculated correctly
func TestTotalCost(t *testing.T) {
	// Create a test order and products
	o := NewOrder(1, "Test Customer")
	p1 := NewProduct(1, "Test Product 1", 10.0, "Test Category")
	p2 := NewProduct(2, "Test Product 2", 20.0, "Test Category")

	// Calculate total cost when there are no items
	require.Equal(t, 0.0, o.TotalCost())

	// Add items and verify the total cost
	o.AddItem(p1, 2) // 2 * 10.0 = 20.0
	require.Equal(t, 20.0, o.TotalCost())

	o.AddItem(p2, 1) // 20.0 + (1 * 20.0) = 40.0
	require.Equal(t, 40.0, o.TotalCost())
}

// TestFindOrderWithHighestTotal verifies the highest total order is found
func TestFindOrderWithHighestTotal(t *testing.T) {
	// Create test orders
	o1 := NewOrder(1, "Customer 1")
	o1.AddItem(NewProduct(1, "Product 1", 10.0, "Category 1"), 1) // Total: 10.0

	o2 := NewOrder(2, "Customer 2")
	o2.AddItem(NewProduct(2, "Product 2", 20.0, "Category 2"), 2) // Total: 40.0

	o3 := NewOrder(3, "Customer 3")
	o3.AddItem(NewProduct(3, "Product 3", 5.0, "Category 3"), 4) // Total: 20.0

	// Test with empty slice
	emptyResult := FindOrderWithHighestTotal([]Order{})
	require.Equal(t, Order{}, emptyResult)

	// Test with single order
	singleResult := FindOrderWithHighestTotal([]Order{o1})
	require.Equal(t, o1.ID, singleResult.ID)
	require.Equal(t, o1.CustomerName, singleResult.CustomerName)

	// Test with multiple orders
	multiResult := FindOrderWithHighestTotal([]Order{o1, o2, o3})
	require.Equal(t, o2.ID, multiResult.ID)
	require.Equal(t, o2.CustomerName, multiResult.CustomerName)
}

// Edge case tests

// TestTotalCostWithZeroItems verifies that orders with zero-priced items calculate correctly
func TestTotalCostWithZeroItems(t *testing.T) {
	o := NewOrder(1, "Test Customer")
	p1 := NewProduct(1, "Free Product", 0.0, "Free")
	p2 := NewProduct(2, "Paid Product", 10.0, "Paid")

	o.AddItem(p1, 5) // 5 * 0.0 = 0.0
	require.Equal(t, 0.0, o.TotalCost())

	o.AddItem(p2, 1) // 0.0 + (1 * 10.0) = 10.0
	require.Equal(t, 10.0, o.TotalCost())
}

// TestTotalCostWithNegativePrices verifies handling of negative prices
func TestTotalCostWithNegativePrices(t *testing.T) {
	o := NewOrder(1, "Test Customer")
	p1 := NewProduct(1, "Discount Product", -5.0, "Discount")
	p2 := NewProduct(2, "Paid Product", 20.0, "Paid")

	o.AddItem(p1, 2) // 2 * -5.0 = -10.0
	require.Equal(t, -10.0, o.TotalCost())

	o.AddItem(p2, 1) // -10.0 + (1 * 20.0) = 10.0
	require.Equal(t, 10.0, o.TotalCost())
}

// TestFindOrderWithHighestTotalEdgeCases tests edge cases for finding highest total
func TestFindOrderWithHighestTotalEdgeCases(t *testing.T) {
	// Create test orders
	o1 := NewOrder(1, "Customer 1")
	// Empty order with 0 total

	o2 := NewOrder(2, "Customer 2")
	o2.AddItem(NewProduct(2, "Negative Product", -10.0, "Negative"), 1) // Total: -10.0

	o3 := NewOrder(3, "Customer 3")
	o3.AddItem(NewProduct(3, "Free Product", 0.0, "Free"), 5) // Total: 0.0

	// Test case where highest total is 0
	result1 := FindOrderWithHighestTotal([]Order{o1, o2, o3})
	require.Equal(t, o1.ID, result1.ID) // Should pick first order with 0 total

	// Add another order with same total
	o4 := NewOrder(4, "Customer 4")
	result2 := FindOrderWithHighestTotal([]Order{o2, o3, o4})
	require.Equal(t, o3.ID, result2.ID) // Should pick first non-negative total
}

// Table-driven tests

// TestTotalCostTable is a table-driven test for total cost calculation
func TestTotalCostTable(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name          string
		setupOrder    func() Order
		expectedTotal float64
	}{
		{
			name: "Empty order",
			setupOrder: func() Order {
				return NewOrder(1, "Test Customer")
			},
			expectedTotal: 0.0,
		},
		{
			name: "Single item",
			setupOrder: func() Order {
				o := NewOrder(1, "Test Customer")
				o.AddItem(NewProduct(1, "Product", 10.0, "Category"), 2)
				return o
			},
			expectedTotal: 20.0,
		},
		{
			name: "Multiple items",
			setupOrder: func() Order {
				o := NewOrder(1, "Test Customer")
				o.AddItem(NewProduct(1, "Product 1", 10.0, "Category 1"), 1)
				o.AddItem(NewProduct(2, "Product 2", 20.0, "Category 2"), 2)
				return o
			},
			expectedTotal: 50.0,
		},
		{
			name: "Zero-priced items",
			setupOrder: func() Order {
				o := NewOrder(1, "Test Customer")
				o.AddItem(NewProduct(1, "Free Product", 0.0, "Free"), 5)
				return o
			},
			expectedTotal: 0.0,
		},
		{
			name: "Mixed price items",
			setupOrder: func() Order {
				o := NewOrder(1, "Test Customer")
				o.AddItem(NewProduct(1, "Cheap", 5.0, "Budget"), 2)
				o.AddItem(NewProduct(2, "Expensive", 50.0, "Premium"), 1)
				return o
			},
			expectedTotal: 60.0,
		},
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order := tc.setupOrder()
			result := order.TotalCost()
			require.Equal(t, tc.expectedTotal, result)
		})
	}
}

// TestFindOrderWithHighestTotalTable is a table-driven test for finding the highest total order
func TestFindOrderWithHighestTotalTable(t *testing.T) {
	// Define reusable functions to create products and orders
	createProduct := func(id int, price float64) Product {
		return NewProduct(id, "Product "+string(rune(48+id)), price, "Category "+string(rune(48+id)))
	}

	createOrder := func(id int, products map[Product]int) Order {
		o := NewOrder(id, "Customer "+string(rune(48+id)))
		for p, q := range products {
			o.AddItem(p, q)
		}
		return o
	}

	// Define test cases
	testCases := []struct {
		name              string
		orders            []Order
		expectedID        int
		expectedTotalCost float64
	}{
		{
			name:              "Empty orders",
			orders:            []Order{},
			expectedID:        0,
			expectedTotalCost: 0.0,
		},
		{
			name: "Single order",
			orders: []Order{
				createOrder(1, map[Product]int{createProduct(1, 10.0): 1}),
			},
			expectedID:        1,
			expectedTotalCost: 10.0,
		},
		{
			name: "Multiple orders with clear highest",
			orders: []Order{
				createOrder(1, map[Product]int{createProduct(1, 10.0): 1}), // 10.0
				createOrder(2, map[Product]int{createProduct(2, 15.0): 2}), // 30.0
				createOrder(3, map[Product]int{createProduct(3, 5.0): 3}),  // 15.0
			},
			expectedID:        2,
			expectedTotalCost: 30.0,
		},
		{
			name: "Orders with same highest total",
			orders: []Order{
				createOrder(1, map[Product]int{createProduct(1, 20.0): 1}), // 20.0
				createOrder(2, map[Product]int{createProduct(2, 10.0): 2}), // 20.0
				createOrder(3, map[Product]int{createProduct(3, 5.0): 3}),  // 15.0
			},
			expectedID:        1, // First one should be returned
			expectedTotalCost: 20.0,
		},
		{
			name: "Orders with negative totals",
			orders: []Order{
				createOrder(1, map[Product]int{createProduct(1, -10.0): 1}), // -10.0
				createOrder(2, map[Product]int{createProduct(2, -5.0): 1}),  // -5.0
			},
			expectedID:        2, // Less negative should be returned
			expectedTotalCost: -5.0,
		},
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FindOrderWithHighestTotal(tc.orders)
			if len(tc.orders) == 0 {
				require.Equal(t, Order{}, result)
			} else {
				require.Equal(t, tc.expectedID, result.ID)
				require.Equal(t, tc.expectedTotalCost, result.TotalCost())
			}
		})
	}
}

// TestPrintOrderDetails verifies that order details are printed correctly
func TestPrintOrderDetails(t *testing.T) {
	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a test order
	o := NewOrder(101, "Test Customer")
	p1 := NewProduct(1, "Test Product", 10.0, "Test Category")
	o.AddItem(p1, 2)

	// Print order details
	PrintOrderDetails(o)

	// Restore stdout and get the captured output
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify the output contains expected information
	require.Contains(t, output, "Order №101")
	require.Contains(t, output, "Test Customer")
	require.Contains(t, output, "Test Product")
	require.Contains(t, output, "Quantity: 2")
	require.Contains(t, output, "Total Order Cost: 20.00")
}

// Test utility functions

// TestGenerateFakeProduct verifies that fake products are generated correctly
func TestGenerateFakeProduct(t *testing.T) {
	// Set a fixed seed for reproducible tests
	gofakeit.Seed(0)

	product := GenerateFakeProduct()

	// Verify constraints
	require.GreaterOrEqual(t, product.ID, 100)
	require.LessOrEqual(t, product.ID, 250)
	require.NotEmpty(t, product.Name)
	require.GreaterOrEqual(t, product.UnitPrice, 1.0)
	require.LessOrEqual(t, product.UnitPrice, 100.0)
	require.NotEmpty(t, product.Category)
}

// TestGenerateFakeOrder verifies that fake orders are generated correctly
func TestGenerateFakeOrder(t *testing.T) {
	// Set a fixed seed for reproducible tests
	gofakeit.Seed(0)

	order := GenerateFakeOrder(101)

	// Verify constraints
	require.Equal(t, 101, order.ID)
	require.NotEmpty(t, order.CustomerName)
	require.NotNil(t, order.Items)
	require.Empty(t, order.Items) // Items map should be initialized but empty
}

// Test main function indirectly through its components
// Since main() itself is difficult to test directly, we verify
// it would work by testing the functions it calls
func TestMainComponents(t *testing.T) {
	// Set a fixed seed for reproducible tests
	gofakeit.Seed(0)

	// Test generating products and orders as done in main
	products := make([]Product, 3)
	for i := 0; i < 3; i++ {
		products[i] = GenerateFakeProduct()
		require.NotEmpty(t, products[i].Name)
	}

	orders := make([]Order, 2)
	for i := 0; i < 2; i++ {
		orders[i] = GenerateFakeOrder(i + 1)
		// Add some products to the orders
		orders[i].AddItem(products[0], 1)
		require.Equal(t, 1, len(orders[i].Items))
	}

	// Test finding highest total
	highestOrder := FindOrderWithHighestTotal(orders)
	require.NotEqual(t, 0, highestOrder.ID)
}
