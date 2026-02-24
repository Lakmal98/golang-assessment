package main

import (
	"encoding/json"
	"net/http"
)

// OpenAPI specification for the API
var openAPISpec = map[string]interface{}{
	"openapi": "3.0.0",
	"info": map[string]interface{}{
		"title":       "Product Availability API",
		"description": "REST API for checking product availability across warehouses with reserve buffer and weekend shipping delay support",
		"version":     "1.0.0",
		"contact": map[string]interface{}{
			"name": "API Support",
		},
	},
	"servers": []map[string]interface{}{
		{
			"url":         "http://localhost:8080",
			"description": "Development server",
		},
	},
	"paths": map[string]interface{}{
		"/api/check-availability": map[string]interface{}{
			"post": map[string]interface{}{
				"summary":     "Check product availability",
				"description": "Check if a product is available at a specific warehouse location. Applies 10% reserve buffer and weekend 2x quantity rules.",
				"operationId": "checkAvailability",
				"tags":        []string{"Availability"},
				"requestBody": map[string]interface{}{
					"required": true,
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"$ref": "#/components/schemas/AvailabilityRequest",
							},
							"examples": map[string]interface{}{
								"normalRequest": map[string]interface{}{
									"summary": "Normal availability check",
									"value": map[string]interface{}{
										"product_id":         "PROD-123",
										"quantity":           5,
										"warehouse_location": "DE-Berlin",
									},
								},
								"lowStock": map[string]interface{}{
									"summary": "Check low stock item",
									"value": map[string]interface{}{
										"product_id":         "PROD-505",
										"quantity":           3,
										"warehouse_location": "US-NewYork",
									},
								},
							},
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Successful availability check",
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/AvailabilityResponse",
								},
								"examples": map[string]interface{}{
									"available": map[string]interface{}{
										"summary": "Product available",
										"value": map[string]interface{}{
											"available":          true,
											"available_quantity": 90,
											"reason":             "Sufficient stock available",
											"warehouse":          "DE-Berlin",
										},
									},
									"outOfStock": map[string]interface{}{
										"summary": "Product out of stock",
										"value": map[string]interface{}{
											"available":          false,
											"available_quantity": 0,
											"reason":             "Product is out of stock",
											"warehouse":          "DE-Berlin",
										},
									},
									"insufficient": map[string]interface{}{
										"summary": "Insufficient stock",
										"value": map[string]interface{}{
											"available":          false,
											"available_quantity": 4,
											"reason":             "Insufficient stock (requires 5 units, only 4 available after reserve)",
											"warehouse":          "US-NewYork",
										},
									},
									"notFound": map[string]interface{}{
										"summary": "Product not in warehouse",
										"value": map[string]interface{}{
											"available":          false,
											"available_quantity": 0,
											"reason":             "Product not found in specified warehouse",
											"warehouse":          "UK-London",
										},
									},
								},
							},
						},
					},
					"400": map[string]interface{}{
						"description": "Bad request - invalid input",
						"content": map[string]interface{}{
							"text/plain": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "string",
								},
								"examples": map[string]interface{}{
									"missingProductId": map[string]interface{}{
										"summary": "Missing product_id",
										"value":   "product_id is required",
									},
									"invalidQuantity": map[string]interface{}{
										"summary": "Invalid quantity",
										"value":   "quantity must be greater than 0",
									},
									"missingWarehouse": map[string]interface{}{
										"summary": "Missing warehouse_location",
										"value":   "warehouse_location is required",
									},
								},
							},
						},
					},
					"405": map[string]interface{}{
						"description": "Method not allowed",
						"content": map[string]interface{}{
							"text/plain": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "string",
								},
								"example": "Method not allowed. Use POST",
							},
						},
					},
				},
			},
		},
	},
	"components": map[string]interface{}{
		"schemas": map[string]interface{}{
			"AvailabilityRequest": map[string]interface{}{
				"type":     "object",
				"required": []string{"product_id", "quantity", "warehouse_location"},
				"properties": map[string]interface{}{
					"product_id": map[string]interface{}{
						"type":        "string",
						"description": "Unique identifier for the product",
						"example":     "PROD-123",
					},
					"quantity": map[string]interface{}{
						"type":        "integer",
						"description": "Requested quantity (must be greater than 0)",
						"minimum":     1,
						"example":     5,
					},
					"warehouse_location": map[string]interface{}{
						"type":        "string",
						"description": "Warehouse location code",
						"example":     "DE-Berlin",
						"enum":        []string{"DE-Berlin", "US-NewYork", "UK-London"},
					},
				},
			},
			"AvailabilityResponse": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"available": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether the product is available in the requested quantity",
						"example":     true,
					},
					"available_quantity": map[string]interface{}{
						"type":        "integer",
						"description": "Total quantity available after applying 10% reserve buffer",
						"example":     90,
					},
					"reason": map[string]interface{}{
						"type":        "string",
						"description": "Detailed reason for the availability status",
						"example":     "Sufficient stock available",
					},
					"warehouse": map[string]interface{}{
						"type":        "string",
						"description": "Warehouse location that was checked",
						"example":     "DE-Berlin",
					},
				},
			},
		},
	},
	"tags": []map[string]interface{}{
		{
			"name":        "Availability",
			"description": "Product availability checking operations",
		},
	},
}

// HandleOpenAPI serves the OpenAPI specification as JSON
func HandleOpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(openAPISpec)
}

// HandleSwaggerUI serves the Swagger UI interface via CDN
func HandleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Product Availability API - Swagger UI</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.10.5/swagger-ui.css" />
    <style>
        body {
            margin: 0;
            padding: 0;
        }
        .topbar {
            display: none;
        }
        .swagger-ui .info {
            margin: 30px 0;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.10.5/swagger-ui-bundle.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.10.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: "/openapi.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                defaultModelsExpandDepth: 1,
                defaultModelExpandDepth: 1,
                docExpansion: "list",
                filter: true,
                tryItOutEnabled: true
            });
        };
    </script>
</body>
</html>`
	
	w.Write([]byte(html))
}
