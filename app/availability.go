package main

import (
	"fmt"
	"time"
)

// AvailabilityService handles the business logic for checking product availability
type AvailabilityService struct {
	inventoryAdapter InventoryAdapter
}

// NewAvailabilityService creates a new availability service with the given inventory adapter
func NewAvailabilityService(adapter InventoryAdapter) *AvailabilityService {
	return &AvailabilityService{
		inventoryAdapter: adapter,
	}
}

// isWeekend checks if today is Saturday or Sunday
func isWeekend() bool {
	today := time.Now().Weekday()
	return today == time.Saturday || today == time.Sunday
}

// CheckAvailability implements the business logic for checking product availability
// Business Rules:
// 1. 10% reserve buffer is always kept from total stock
// 2. Weekend orders require 2x the normal quantity in stock
// 3. Returns availability status with detailed reason
func (s *AvailabilityService) CheckAvailability(req Request) Response {
	response := Response{
		Warehouse: req.WarehouseLocation,
	}

	// Get stock level from the inventory adapter
	stockLevel, err := s.inventoryAdapter.GetStockLevel(req.ProductID, req.WarehouseLocation)
	if err != nil {
		response.Available = false
		response.AvailableQuantity = 0
		response.Reason = "Product not found in specified warehouse"
		return response
	}

	// Calculate available stock after applying 10% reserve buffer
	reserveBuffer := float64(stockLevel) * 0.10
	availableStock := stockLevel - int(reserveBuffer)
	response.AvailableQuantity = availableStock

	// Check if stock is zero
	if stockLevel == 0 {
		response.Available = false
		response.Reason = "Product is out of stock"
		return response
	}

	// Determine required quantity based on weekend logic
	requiredQuantity := req.Quantity
	if isWeekend() {
		requiredQuantity = req.Quantity * 2
	}

	// Check if we have enough available stock
	if availableStock >= requiredQuantity {
		response.Available = true
		if isWeekend() {
			response.Reason = fmt.Sprintf("Sufficient stock available (weekend: requires %d units in stock for %d order)", requiredQuantity, req.Quantity)
		} else {
			response.Reason = "Sufficient stock available"
		}
	} else {
		response.Available = false
		if isWeekend() {
			response.Reason = fmt.Sprintf("Insufficient stock (weekend: requires %d units, only %d available after reserve)", requiredQuantity, availableStock)
		} else {
			response.Reason = fmt.Sprintf("Insufficient stock (requires %d units, only %d available after reserve)", requiredQuantity, availableStock)
		}
	}

	return response
}
