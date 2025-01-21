package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Service1f/model"
)

// OrderHandler handles the creation of orders
func OrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate customer and product existence
	if !validateCustomer(order.CustomerID) {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}
	if !validateProduct(order.ProductID) {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Create the order
	id, err := models.CreateOrder(order)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create order: %v", err), http.StatusInternalServerError)
		return
	}

	order.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// validateCustomer checks if a customer exists
func validateCustomer(customerID int) bool {
	resp, err := http.Get(fmt.Sprintf("%s%d",models.BaseURLCustomer, customerID))
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

// validateProduct checks if a product exists
func validateProduct(productID int) bool {
	resp, err := http.Get(fmt.Sprintf("%s%d", models.BaseURLProduct,productID))
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
