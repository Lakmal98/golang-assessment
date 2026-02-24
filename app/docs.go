package main

import (
	"net/http"
)

// HandleDocs serves the API documentation page
func HandleDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Product Availability API - Documentation</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            background: #f5f5f5;
            padding: 20px;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            padding: 40px;
        }
        h1 {
            color: #2c3e50;
            border-bottom: 3px solid #3498db;
            padding-bottom: 10px;
            margin-bottom: 30px;
        }
        h2 {
            color: #34495e;
            margin-top: 30px;
            margin-bottom: 15px;
            padding-bottom: 5px;
            border-bottom: 2px solid #ecf0f1;
        }
        h3 {
            color: #7f8c8d;
            margin-top: 20px;
            margin-bottom: 10px;
        }
        .endpoint {
            background: #ecf0f1;
            padding: 20px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .method {
            display: inline-block;
            background: #27ae60;
            color: white;
            padding: 5px 15px;
            border-radius: 3px;
            font-weight: bold;
            margin-right: 10px;
        }
        .path {
            font-family: 'Courier New', monospace;
            font-size: 18px;
            color: #2c3e50;
        }
        pre {
            background: #2c3e50;
            color: #ecf0f1;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
            margin: 10px 0;
        }
        code {
            font-family: 'Courier New', monospace;
        }
        .info-box {
            background: #e8f4f8;
            border-left: 4px solid #3498db;
            padding: 15px;
            margin: 15px 0;
        }
        .warning-box {
            background: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 15px;
            margin: 15px 0;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 15px 0;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 12px;
            text-align: left;
        }
        th {
            background-color: #3498db;
            color: white;
        }
        tr:nth-child(even) {
            background-color: #f2f2f2;
        }
        .badge {
            display: inline-block;
            padding: 3px 8px;
            border-radius: 3px;
            font-size: 12px;
            font-weight: bold;
        }
        .badge-required {
            background: #e74c3c;
            color: white;
        }
        .badge-integer {
            background: #9b59b6;
            color: white;
        }
        .badge-string {
            background: #16a085;
            color: white;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📦 Product Availability API Documentation</h1>
        
        <div class="info-box">
            <strong>Version:</strong> 1.0.0<br>
            <strong>Base URL:</strong> <code>http://localhost:8080</code><br>
            <strong>Content-Type:</strong> <code>application/json</code>
        </div>

        <h2>Overview</h2>
        <p>This API checks product availability across multiple warehouses with support for reserve buffers and weekend shipping delays.</p>

        <h2>Endpoints</h2>

        <div class="endpoint">
            <div>
                <span class="method">POST</span>
                <span class="path">/api/check-availability</span>
            </div>
            <p style="margin-top: 15px;"><strong>Description:</strong> Check if a product is available at a specific warehouse location</p>
            
            <h3>Request Body</h3>
            <table>
                <thead>
                    <tr>
                        <th>Field</th>
                        <th>Type</th>
                        <th>Required</th>
                        <th>Description</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td><code>product_id</code></td>
                        <td><span class="badge badge-string">string</span></td>
                        <td><span class="badge badge-required">required</span></td>
                        <td>Unique identifier for the product</td>
                    </tr>
                    <tr>
                        <td><code>quantity</code></td>
                        <td><span class="badge badge-integer">integer</span></td>
                        <td><span class="badge badge-required">required</span></td>
                        <td>Requested quantity (must be > 0)</td>
                    </tr>
                    <tr>
                        <td><code>warehouse_location</code></td>
                        <td><span class="badge badge-string">string</span></td>
                        <td><span class="badge badge-required">required</span></td>
                        <td>Warehouse location code</td>
                    </tr>
                </tbody>
            </table>

            <h3>Example Request</h3>
            <pre><code>POST /api/check-availability
Content-Type: application/json

{
  "product_id": "PROD-123",
  "quantity": 5,
  "warehouse_location": "DE-Berlin"
}</code></pre>

            <h3>Response Fields</h3>
            <table>
                <thead>
                    <tr>
                        <th>Field</th>
                        <th>Type</th>
                        <th>Description</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td><code>available</code></td>
                        <td><span class="badge badge-string">boolean</span></td>
                        <td>Whether the product is available in requested quantity</td>
                    </tr>
                    <tr>
                        <td><code>available_quantity</code></td>
                        <td><span class="badge badge-integer">integer</span></td>
                        <td>Total quantity available after applying reserve buffer</td>
                    </tr>
                    <tr>
                        <td><code>reason</code></td>
                        <td><span class="badge badge-string">string</span></td>
                        <td>Detailed reason for the availability status</td>
                    </tr>
                    <tr>
                        <td><code>warehouse</code></td>
                        <td><span class="badge badge-string">string</span></td>
                        <td>Warehouse location that was checked</td>
                    </tr>
                </tbody>
            </table>

            <h3>Success Response (200 OK)</h3>
            <pre><code>{
  "available": true,
  "available_quantity": 90,
  "reason": "Sufficient stock available",
  "warehouse": "DE-Berlin"
}</code></pre>

            <h3>Error Responses</h3>
            <table>
                <thead>
                    <tr>
                        <th>Status Code</th>
                        <th>Description</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td><code>400 Bad Request</code></td>
                        <td>Invalid JSON, missing required fields, or invalid quantity</td>
                    </tr>
                    <tr>
                        <td><code>405 Method Not Allowed</code></td>
                        <td>HTTP method other than POST was used</td>
                    </tr>
                </tbody>
            </table>
        </div>

        <h2>Business Rules</h2>
        
        <div class="warning-box">
            <h3>1. Reserve Buffer (10%)</h3>
            <p>10% of total stock is always reserved and unavailable for orders. The available_quantity reflects this buffer.</p>
            <p><strong>Example:</strong> If stock level is 100, available quantity = 100 - (100 × 0.10) = 90</p>
        </div>

        <div class="warning-box">
            <h3>2. Weekend Shipping Delay</h3>
            <p>Orders placed on Saturday or Sunday require 2x the normal quantity to be in available stock.</p>
            <p><strong>Example:</strong> Ordering 5 units on Saturday requires 10 units in available stock</p>
        </div>

        <div class="warning-box">
            <h3>3. Availability Check</h3>
            <p>Product is marked as available only if the available stock (after reserve buffer) meets or exceeds the required quantity.</p>
        </div>

        <h2>Available Products</h2>
        <table>
            <thead>
                <tr>
                    <th>Product ID</th>
                    <th>Warehouse</th>
                    <th>Stock Level</th>
                    <th>Available After Reserve</th>
                    <th>Notes</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>PROD-123</td>
                    <td>DE-Berlin</td>
                    <td>100</td>
                    <td>90</td>
                    <td>Normal stock</td>
                </tr>
                <tr>
                    <td>PROD-123</td>
                    <td>US-NewYork</td>
                    <td>50</td>
                    <td>45</td>
                    <td>Normal stock</td>
                </tr>
                <tr>
                    <td>PROD-456</td>
                    <td>DE-Berlin</td>
                    <td>25</td>
                    <td>22</td>
                    <td>Normal stock</td>
                </tr>
                <tr>
                    <td>PROD-456</td>
                    <td>US-NewYork</td>
                    <td>75</td>
                    <td>67</td>
                    <td>Normal stock</td>
                </tr>
                <tr>
                    <td>PROD-789</td>
                    <td>DE-Berlin</td>
                    <td>0</td>
                    <td>0</td>
                    <td>⚠️ Out of stock</td>
                </tr>
                <tr>
                    <td>PROD-789</td>
                    <td>US-NewYork</td>
                    <td>15</td>
                    <td>13</td>
                    <td>Low stock</td>
                </tr>
                <tr>
                    <td>PROD-101</td>
                    <td>DE-Berlin</td>
                    <td>10</td>
                    <td>9</td>
                    <td>Edge case: exactly 10</td>
                </tr>
                <tr>
                    <td>PROD-101</td>
                    <td>US-NewYork</td>
                    <td>11</td>
                    <td>10</td>
                    <td>Edge case: just above</td>
                </tr>
                <tr>
                    <td>PROD-202</td>
                    <td>DE-Berlin</td>
                    <td>1</td>
                    <td>1</td>
                    <td>⚠️ Very low stock</td>
                </tr>
                <tr>
                    <td>PROD-202</td>
                    <td>UK-London</td>
                    <td>500</td>
                    <td>450</td>
                    <td>High stock</td>
                </tr>
                <tr>
                    <td>PROD-303</td>
                    <td>UK-London</td>
                    <td>20</td>
                    <td>18</td>
                    <td>Normal stock</td>
                </tr>
                <tr>
                    <td>PROD-404</td>
                    <td>DE-Berlin</td>
                    <td>200</td>
                    <td>180</td>
                    <td>Large inventory</td>
                </tr>
                <tr>
                    <td>PROD-505</td>
                    <td>US-NewYork</td>
                    <td>5</td>
                    <td>4</td>
                    <td>⚠️ Very low stock</td>
                </tr>
            </tbody>
        </table>

        <h2>cURL Examples</h2>
        
        <h3>Check Availability (Sufficient Stock)</h3>
        <pre><code>curl -X POST http://localhost:8080/api/check-availability \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "PROD-123",
    "quantity": 5,
    "warehouse_location": "DE-Berlin"
  }'</code></pre>

        <h3>Check Out of Stock Item</h3>
        <pre><code>curl -X POST http://localhost:8080/api/check-availability \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "PROD-789",
    "quantity": 1,
    "warehouse_location": "DE-Berlin"
  }'</code></pre>

        <h3>Check Low Stock Item</h3>
        <pre><code>curl -X POST http://localhost:8080/api/check-availability \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "PROD-505",
    "quantity": 5,
    "warehouse_location": "US-NewYork"
  }'</code></pre>

        <h3>Product Not In Warehouse</h3>
        <pre><code>curl -X POST http://localhost:8080/api/check-availability \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "PROD-123",
    "quantity": 5,
    "warehouse_location": "UK-London"
  }'</code></pre>

        <h2>Error Examples</h2>

        <h3>Missing Required Field</h3>
        <pre><code>curl -X POST http://localhost:8080/api/check-availability \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 5,
    "warehouse_location": "DE-Berlin"
  }'

# Response: 400 Bad Request
product_id is required</code></pre>

        <h3>Invalid Quantity</h3>
        <pre><code>curl -X POST http://localhost:8080/api/check-availability \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "PROD-123",
    "quantity": 0,
    "warehouse_location": "DE-Berlin"
  }'

# Response: 400 Bad Request
quantity must be greater than 0</code></pre>

        <h3>Wrong HTTP Method</h3>
        <pre><code>curl -X GET http://localhost:8080/api/check-availability

# Response: 405 Method Not Allowed
Method not allowed. Use POST</code></pre>

    </div>
</body>
</html>`

	w.Write([]byte(html))
}
