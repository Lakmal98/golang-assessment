package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// AvailabilityHandler handles HTTP requests for the availability check endpoint
type AvailabilityHandler struct {
	availabilityService *AvailabilityService
}

// NewAvailabilityHandler creates a new handler with the given availability service
func NewAvailabilityHandler(service *AvailabilityService) *AvailabilityHandler {
	return &AvailabilityHandler{
		availabilityService: service,
	}
}

// HandleCheckAvailability handles POST /api/check-availability requests
func (h *AvailabilityHandler) HandleCheckAvailability(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request
	var req Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validate input fields
	if req.ProductID == "" {
		http.Error(w, "product_id is required", http.StatusBadRequest)
		return
	}
	if req.Quantity <= 0 {
		http.Error(w, "quantity must be greater than 0", http.StatusBadRequest)
		return
	}
	if req.WarehouseLocation == "" {
		http.Error(w, "warehouse_location is required", http.StatusBadRequest)
		return
	}

	// Check availability using the service
	response := h.availabilityService.CheckAvailability(req)

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}