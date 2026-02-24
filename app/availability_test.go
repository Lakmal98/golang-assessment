package main

import (
	"fmt"
	"testing"
)

// MockInventoryAdapter is a test double for the InventoryAdapter interface
type MockInventoryAdapter struct {
	inventory map[string]map[string]int // productID -> warehouse -> stock
}

// NewMockInventoryAdapter creates a new mock adapter with predefined test data
func NewMockInventoryAdapter() *MockInventoryAdapter {
	return &MockInventoryAdapter{
		inventory: map[string]map[string]int{
			"PROD-123": {
				"DE-Berlin":  100,
				"US-NewYork": 50,
			},
			"PROD-456": {
				"DE-Berlin":  25,
				"US-NewYork": 75,
			},
			"PROD-789": {
				"DE-Berlin":  0,  // Out of stock
				"US-NewYork": 15, // Low stock
			},
			"PROD-101": {
				"DE-Berlin":  10, // Exactly 10 units
				"US-NewYork": 11, // Just above threshold
			},
			"PROD-202": {
				"DE-Berlin": 1,   // Very low stock
				"UK-London": 500, // High stock
			},
			"PROD-303": {
				"UK-London": 20, // Exactly at edge
			},
			"PROD-404": {
				"DE-Berlin": 200, // Large inventory
			},
			"PROD-505": {
				"US-NewYork": 5, // Very low stock
			},
		},
	}
}

func (m *MockInventoryAdapter) LoadInventory() error {
	return nil
}

func (m *MockInventoryAdapter) GetStockLevel(productID, warehouse string) (int, error) {
	if warehouses, ok := m.inventory[productID]; ok {
		if stock, ok := warehouses[warehouse]; ok {
			return stock, nil
		}
	}
	return 0, fmt.Errorf("product %s not found in warehouse %s", productID, warehouse)
}

func TestCheckAvailability_SufficientStock(t *testing.T) {
	adapter := NewMockInventoryAdapter()
	service := NewAvailabilityService(adapter)

	req := Request{
		ProductID:         "PROD-123",
		Quantity:          5,
		WarehouseLocation: "DE-Berlin",
	}

	resp := service.CheckAvailability(req)

	if !resp.Available {
		t.Errorf("Expected available=true, got false. Reason: %s", resp.Reason)
	}
	if resp.AvailableQuantity != 90 {
		t.Errorf("Expected available_quantity=90 (100 - 10%%), got %d", resp.AvailableQuantity)
	}
}

func TestCheckAvailability_OutOfStock(t *testing.T) {
	adapter := NewMockInventoryAdapter()
	service := NewAvailabilityService(adapter)

	req := Request{
		ProductID:         "PROD-789",
		Quantity:          1,
		WarehouseLocation: "DE-Berlin",
	}

	resp := service.CheckAvailability(req)

	if resp.Available {
		t.Error("Expected available=false for out of stock item")
	}
	if resp.Reason != "Product is out of stock" {
		t.Errorf("Expected 'Product is out of stock', got '%s'", resp.Reason)
	}
	if resp.AvailableQuantity != 0 {
		t.Errorf("Expected available_quantity=0, got %d", resp.AvailableQuantity)
	}
}

func TestCheckAvailability_InsufficientStock(t *testing.T) {
	adapter := NewMockInventoryAdapter()
	service := NewAvailabilityService(adapter)

	// Stock: 25, Reserve: 2 (int(2.5)), Available: 23, Required: 24 - should fail
	req := Request{
		ProductID:         "PROD-456",
		Quantity:          24,
		WarehouseLocation: "DE-Berlin",
	}

	resp := service.CheckAvailability(req)

	if resp.Available {
		t.Error("Expected available=false for insufficient stock")
	}
	if resp.AvailableQuantity != 23 {
		t.Errorf("Expected available_quantity=23, got %d", resp.AvailableQuantity)
	}
}

func TestCheckAvailability_ExactlyAtThreshold(t *testing.T) {
	adapter := NewMockInventoryAdapter()
	service := NewAvailabilityService(adapter)

	// PROD-101 has 10 units, after 10% reserve = 9 available
	req := Request{
		ProductID:         "PROD-101",
		Quantity:          9,
		WarehouseLocation: "DE-Berlin",
	}

	resp := service.CheckAvailability(req)

	if !resp.Available {
		t.Errorf("Expected available=true for exactly 9 units, got false. Available: %d", resp.AvailableQuantity)
	}
}

func TestCheckAvailability_JustAboveThreshold(t *testing.T) {
	adapter := NewMockInventoryAdapter()
	service := NewAvailabilityService(adapter)

	// PROD-101 has 11 units, after 10% reserve (1) = 10 available
	req := Request{
		ProductID:         "PROD-101",
		Quantity:          10,
		WarehouseLocation: "US-NewYork",
	}

	resp := service.CheckAvailability(req)

	if !resp.Available {
		t.Errorf("Expected available=true, got false. Available quantity: %d", resp.AvailableQuantity)
	}
}

func TestCheckAvailability_VeryLowStock(t *testing.T) {
	adapter := NewMockInventoryAdapter()
	service := NewAvailabilityService(adapter)

	// PROD-202 has 1 unit, after 10% reserve (0) = 1 available
	req := Request{
		ProductID:         "PROD-202",
		Quantity:          1,
		WarehouseLocation: "DE-Berlin",
	}

	resp := service.CheckAvailability(req)

	if !resp.Available {
		t.Errorf("Expected available=true for 1 unit when stock=1, got false")
	}
}

func TestCheckAvailability_VeryLowStockInsufficientQuantity(t *testing.T) {
	adapter := NewMockInventoryAdapter()
	service := NewAvailabilityService(adapter)

	// PROD-202 has 1 unit, after 10% reserve (0) = 1 available, requesting 2
	req := Request{
		ProductID:         "PROD-202",
		Quantity:          2,
		WarehouseLocation: "DE-Berlin",
	}

	resp := service.CheckAvailability(req)

	if resp.Available {
		t.Error("Expected available=false when requesting 2 units but only 1 available")
	}
}

func TestCheckAvailability_HighStock(t *testing.T) {
	adapter := NewMockInventoryAdapter()
	service := NewAvailabilityService(adapter)

	// PROD-202 has 500 units in UK-London, after 10% reserve (50) = 450 available
	req := Request{
		ProductID:         "PROD-202",
		Quantity:          400,
		WarehouseLocation: "UK-London",
	}

	resp := service.CheckAvailability(req)

	if !resp.Available {
		t.Errorf("Expected available=true for large order with high stock, got false")
	}
	if resp.AvailableQuantity != 450 {
		t.Errorf("Expected available_quantity=450, got %d", resp.AvailableQuantity)
	}
}

func TestCheckAvailability_ProductNotFound(t *testing.T) {
	adapter := NewMockInventoryAdapter()
	service := NewAvailabilityService(adapter)

	req := Request{
		ProductID:         "PROD-999",
		Quantity:          5,
		WarehouseLocation: "DE-Berlin",
	}

	resp := service.CheckAvailability(req)

	if resp.Available {
		t.Error("Expected available=false for non-existent product")
	}
	if resp.Reason != "Product not found in specified warehouse" {
		t.Errorf("Expected 'Product not found' message, got '%s'", resp.Reason)
	}
}

func TestCheckAvailability_WrongWarehouse(t *testing.T) {
	adapter := NewMockInventoryAdapter()
	service := NewAvailabilityService(adapter)

	req := Request{
		ProductID:         "PROD-123",
		Quantity:          5,
		WarehouseLocation: "UK-London", // PROD-123 not in UK-London
	}

	resp := service.CheckAvailability(req)

	if resp.Available {
		t.Error("Expected available=false for product in wrong warehouse")
	}
}

func TestCheckAvailability_LargeInventory(t *testing.T) {
	adapter := NewMockInventoryAdapter()
	service := NewAvailabilityService(adapter)

	// PROD-404 has 200 units, after 10% reserve (20) = 180 available
	req := Request{
		ProductID:         "PROD-404",
		Quantity:          180,
		WarehouseLocation: "DE-Berlin",
	}

	resp := service.CheckAvailability(req)

	if !resp.Available {
		t.Errorf("Expected available=true, got false. Available: %d", resp.AvailableQuantity)
	}
}

func TestCheckAvailability_JustOverAvailable(t *testing.T) {
	adapter := NewMockInventoryAdapter()
	service := NewAvailabilityService(adapter)

	// PROD-404 has 200 units, after 10% reserve (20) = 180 available
	// Requesting 181 should fail
	req := Request{
		ProductID:         "PROD-404",
		Quantity:          181,
		WarehouseLocation: "DE-Berlin",
	}

	resp := service.CheckAvailability(req)

	if resp.Available {
		t.Error("Expected available=false when requesting more than available")
	}
}

func TestCheckAvailability_ReserveBufferCalculation(t *testing.T) {
	adapter := NewMockInventoryAdapter()
	service := NewAvailabilityService(adapter)

	tests := []struct {
		productID string
		warehouse string
		stock     int
		expected  int
	}{
		{"PROD-123", "DE-Berlin", 100, 90}, // 100 - int(10.0) = 90
		{"PROD-123", "US-NewYork", 50, 45}, // 50 - int(5.0) = 45
		{"PROD-456", "DE-Berlin", 25, 23},  // 25 - int(2.5) = 23
		{"PROD-789", "US-NewYork", 15, 14}, // 15 - int(1.5) = 14
		{"PROD-101", "DE-Berlin", 10, 9},   // 10 - int(1.0) = 9
		{"PROD-505", "US-NewYork", 5, 5},   // 5 - int(0.5) = 5
	}

	for _, tt := range tests {
		req := Request{
			ProductID:         tt.productID,
			Quantity:          1,
			WarehouseLocation: tt.warehouse,
		}

		resp := service.CheckAvailability(req)

		if resp.AvailableQuantity != tt.expected {
			t.Errorf("Product %s in %s: expected available_quantity=%d, got %d (stock=%d)",
				tt.productID, tt.warehouse, tt.expected, resp.AvailableQuantity, tt.stock)
		}
	}
}
