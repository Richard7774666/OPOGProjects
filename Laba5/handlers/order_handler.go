package handlers

import (
	"encoding/json"
	"net/http"
	"order-management-api/models"
	"order-management-api/service"
	"strconv"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	service *service.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// Create handles the creation of a new order
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateOrder(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetByID handles retrieving an order by ID
func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrderByID(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// GetAll handles retrieving all orders
func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// AddItem handles adding an item to an order
func (h *OrderHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	orderID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var req models.AddItemRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.AddItemToOrder(orderID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return the updated order
	order, err := h.service.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, "Failed to retrieve updated order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// GetHighestTotal handles retrieving the order with the highest total
func (h *OrderHandler) GetHighestTotal(w http.ResponseWriter, r *http.Request) {
	order, err := h.service.GetOrderWithHighestTotal()
	if err != nil {
		if err.Error() == "no orders found" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "No orders found",
			})
			return
		}
		http.Error(w, "Failed to retrieve order with highest total", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
