package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// InventoryAdapter defines the interface for fetching inventory data
// This allows easy switching between different data sources (file, API, database, etc.)
type InventoryAdapter interface {
	GetStockLevel(productID, warehouse string) (int, error)
	LoadInventory() error
}

// FileInventoryAdapter implements InventoryAdapter using a JSON file as data source
type FileInventoryAdapter struct {
	filePath  string
	inventory []InventoryItem
}

// NewFileInventoryAdapter creates a new file-based inventory adapter
func NewFileInventoryAdapter(filePath string) *FileInventoryAdapter {
	return &FileInventoryAdapter{
		filePath:  filePath,
		inventory: []InventoryItem{},
	}
}

// LoadInventory loads inventory data from the JSON file
func (f *FileInventoryAdapter) LoadInventory() error {
	// Read the JSON file
	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return fmt.Errorf("failed to read inventory file: %w", err)
	}

	// Parse JSON into inventory items
	err = json.Unmarshal(data, &f.inventory)
	if err != nil {
		return fmt.Errorf("failed to parse inventory JSON: %w", err)
	}

	return nil
}

// GetStockLevel retrieves the stock level for a product at a specific warehouse
func (f *FileInventoryAdapter) GetStockLevel(productID, warehouse string) (int, error) {
	for _, item := range f.inventory {
		if item.ProductID == productID && item.Warehouse == warehouse {
			return item.StockLevel, nil
		}
	}
	return 0, fmt.Errorf("product %s not found in warehouse %s", productID, warehouse)
}

// APIInventoryAdapter is a placeholder for future API-based inventory
// This demonstrates how easy it is to swap implementations with the adapter pattern
type APIInventoryAdapter struct {
	apiURL string
	// Add fields for API client, authentication, etc.
}

// NewAPIInventoryAdapter creates a new API-based inventory adapter
func NewAPIInventoryAdapter(apiURL string) *APIInventoryAdapter {
	return &APIInventoryAdapter{
		apiURL: apiURL,
	}
}

// LoadInventory would initialize the API connection
func (a *APIInventoryAdapter) LoadInventory() error {
	// TODO: Initialize API client, authenticate, etc.
	return fmt.Errorf("API adapter not yet implemented")
}

// GetStockLevel would fetch stock from an external API
func (a *APIInventoryAdapter) GetStockLevel(productID, warehouse string) (int, error) {
	// TODO: Make API call to fetch stock level
	// Example: GET /api/inventory?product={productID}&warehouse={warehouse}
	return 0, fmt.Errorf("API adapter not yet implemented")
}
