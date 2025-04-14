package main

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
)

// Product - struct representing a product
type Product struct {
	ID        int
	Name      string
	UnitPrice float64
	Category  string
}

// Order - struct representing an order
type Order struct {
	ID           int
	CustomerName string
	Items        map[Product]int
}

// NewProduct - function to create a new product
func NewProduct(id int, name string, price float64, category string) Product {
	return Product{ID: id, Name: name, UnitPrice: price, Category: category}
}

// NewOrder - function to create a new order
func NewOrder(id int, customerName string) Order {
	return Order{ID: id, CustomerName: customerName, Items: make(map[Product]int)}
}

// AddItem - method to add a product to the order with a specified quantity
func (o *Order) AddItem(product Product, quantity int) {
	o.Items[product] += quantity
}

// TotalCost - method to calculate the total cost of the order
func (o Order) TotalCost() float64 {
	totalCost := 0.0
	for product, quantity := range o.Items {
		totalCost += product.UnitPrice * float64(quantity)
	}
	return totalCost
}

// FindOrderWithHighestTotal - function to find the order with the highest total amount
func FindOrderWithHighestTotal(orders []Order) Order {
	if len(orders) == 0 {
		return Order{} // Return an empty order if the slice is empty
	}

	highestOrder := orders[0]
	highestTotal := highestOrder.TotalCost()

	for _, currentOrder := range orders[1:] {
		currentTotal := currentOrder.TotalCost()
		if currentTotal > highestTotal {
			highestTotal = currentTotal
			highestOrder = currentOrder
		}
	}
	return highestOrder
}

func GenerateFakeProduct() Product {
	return Product{
		ID:        gofakeit.Number(100, 250),
		Name:      gofakeit.Name(),
		UnitPrice: gofakeit.Price(1, 100),
		Category:  gofakeit.Word(),
	}
}

func GenerateFakeOrder(id int) Order {
	return Order{
		ID:           id,
		CustomerName: gofakeit.Name(),
		Items:        make(map[Product]int), 
	}
}

// PrintOrderDetails - function to print the details of an order including item names, quantities, and prices
func PrintOrderDetails(order Order) {
	fmt.Printf("Order №%d by %s:\n", order.ID, order.CustomerName)
	for product, quantity := range order.Items {
		fmt.Printf("Product: %s, Quantity: %d, Price per unit: %.2f, Total: %.2f\n", product.Name, quantity, product.UnitPrice, product.UnitPrice*float64(quantity))
	}
	fmt.Printf("Total Order Cost: %.2f\n\n", order.TotalCost())
}

func main() {
	gofakeit.Seed(0)

	products := make([]Product, 15)
	for i := 0; i < 15; i++ {
		products[i] = GenerateFakeProduct()
	}

	orders := make([]Order, 7)
	for i := 0; i < 7; i++ {
		order := GenerateFakeOrder(i + 1)
		// Добавляем товары в заказ с случайным количеством
		for j := 0; j < gofakeit.Number(1, 5); j++ {
			product := products[gofakeit.Number(0, len(products)-1)] // случайный товар
			quantity := gofakeit.Number(1, 3)                         // случайное количество
			order.AddItem(product, quantity)
		}
		orders[i] = order
	}

	product1 := NewProduct(1, "Apple", 1.50, "Fruits")
	product2 := NewProduct(2, "Milk", 2.00, "Dairy")
	product3 := NewProduct(3, "Bread", 1.25, "Bakery")

	order1 := NewOrder(101, "John Doe")
	order1.AddItem(product1, 2)
	order1.AddItem(product2, 1)

	order2 := NewOrder(102, "Jane Smith")
	order2.AddItem(product2, 3)
	order2.AddItem(product3, 2)
	order2.AddItem(product1, 1)

	order3 := NewOrder(103, "Peter Jones")
	order3.AddItem(product3, 5)
    // Add custom orders to the orders list
	orders = append(orders, order1, order2, order3)

	for _, order := range orders {
		fmt.Printf("Order №%d by %s: Total cost = %.2f\n", order.ID, order.CustomerName, order.TotalCost())
	}


	highestOrder := FindOrderWithHighestTotal(orders)

	fmt.Println("\nOrder with the highest total:")
	fmt.Printf("ID: %d, Customer: %s, Total: %.2f\n", highestOrder.ID, highestOrder.CustomerName, highestOrder.TotalCost())
	fmt.Println("")
	// Print all orders with their details
	for _, order := range orders {
		PrintOrderDetails(order)
	}
}