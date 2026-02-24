package main

// Request represents the incoming availability check request
type Request struct {
	ProductID         string `json:"product_id"`
	Quantity          int    `json:"quantity"`
	WarehouseLocation string `json:"warehouse_location"`
}

// Response represents the availability check response
type Response struct {
	Available         bool   `json:"available"`
	AvailableQuantity int    `json:"available_quantity"`
	Reason            string `json:"reason"`
	Warehouse         string `json:"warehouse"`
}

// InventoryItem represents stock information for a product at a warehouse
type InventoryItem struct {
	ProductID  string `json:"product_id"`
	Warehouse  string `json:"warehouse"`
	StockLevel int    `json:"stock_level"`
}
