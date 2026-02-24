package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Initialize the inventory adapter (using file-based adapter)
	// This can be easily swapped with APIInventoryAdapter in the future
	inventoryAdapter := NewFileInventoryAdapter("inventory.json")

	// Load inventory data from JSON file
	err := inventoryAdapter.LoadInventory()
	if err != nil {
		log.Fatalf("Failed to load inventory: %v", err)
	}

	// Initialize the availability service with the inventory adapter
	availabilityService := NewAvailabilityService(inventoryAdapter)

	// Initialize the HTTP handler with the availability service
	handler := NewAvailabilityHandler(availabilityService)

	// Register the endpoints
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/docs", http.StatusMovedPermanently)
			return
		}
		http.NotFound(w, r)
	})
	http.HandleFunc("/api/check-availability", handler.HandleCheckAvailability)
	http.HandleFunc("/docs", HandleSwaggerUI)
	http.HandleFunc("/openapi.json", HandleOpenAPI)

	// Start the server
	port := ":8080"
	fmt.Printf("Starting server on http://localhost%s\n", port)
	fmt.Println("Endpoints:")
	fmt.Println("  - GET  / (redirects to /docs)")
	fmt.Println("  - POST /api/check-availability")
	fmt.Println("  - GET  /docs (Swagger UI Documentation)")
	fmt.Println("  - GET  /openapi.json (OpenAPI Specification)")
	fmt.Printf("Current day: %s (Weekend: %v)\n", time.Now().Weekday(), isWeekend())
	fmt.Println("Inventory loaded from: inventory.json")
	fmt.Println()

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
