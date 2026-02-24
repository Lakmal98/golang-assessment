# Testing Guide

## Running Tests

### Using Go

```bash
cd app
go test -v
```

### Test Coverage

To run tests with coverage:

```bash
cd app
go test -v -cover
```

### Generate Coverage Report

```bash
cd app
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Suite

The test suite includes comprehensive edge cases:

### Test Cases Covered

1. **TestCheckAvailability_SufficientStock**
   - Verifies availability when stock is sufficient
   - Tests reserve buffer calculation (10%)
   - Expected: Available = true

2. **TestCheckAvailability_OutOfStock**
   - Tests handling of zero stock items
   - Expected: Available = false with "out of stock" message

3. **TestCheckAvailability_InsufficientStock**
   - Tests when available stock (after reserve) is less than requested
   - Stock: 25, Available: 23, Requesting: 24
   - Verifies proper insufficient stock messaging

4. **TestCheckAvailability_ExactlyAtThreshold**
   - Edge case: Requesting exactly the available quantity after reserve
   - Tests boundary conditions

5. **TestCheckAvailability_JustAboveThreshold**
   - Edge case: Stock level just above minimum threshold
   - Verifies rounding behavior of reserve buffer

6. **TestCheckAvailability_VeryLowStock**
   - Tests items with only 1 unit in stock
   - Verifies reserve buffer calculation on low quantities

7. **TestCheckAvailability_VeryLowStockInsufficientQuantity**
   - Tests rejection when requesting more than very low stock

8. **TestCheckAvailability_HighStock**
   - Tests large inventory handling
   - Verifies correct reserve buffer on large numbers

9. **TestCheckAvailability_ProductNotFound**
   - Tests non-existent product handling
   - Expected: Proper error message

10. **TestCheckAvailability_WrongWarehouse**
    - Tests product lookup in wrong warehouse
    - Verifies warehouse-specific inventory

11. **TestCheckAvailability_LargeInventory**
    - Tests handling of large stock quantities
    - Verifies scaling of reserve buffer logic

12. **TestCheckAvailability_JustOverAvailable**
    - Edge case: Requesting 1 more than available
    - Tests precise boundary detection

13. **TestCheckAvailability_ReserveBufferCalculation**
    - Parametric test verifying 10% reserve across multiple products
    - Validates calculation accuracy for various stock levels

## Mock Data for Testing

The test suite uses a `MockInventoryAdapter` that implements the `InventoryAdapter` interface:

```go
type MockInventoryAdapter struct {
    inventory map[string]map[string]int // productID -> warehouse -> stock
}
```

### Test Inventory Data

| Product ID | Warehouse | Stock Level | Available After Reserve |
|------------|-----------|-------------|------------------------|
| PROD-123 | DE-Berlin | 100 | 90 |
| PROD-123 | US-NewYork | 50 | 45 |
| PROD-456 | DE-Berlin | 25 | 23 |
| PROD-456 | US-NewYork | 75 | 67 |
| PROD-789 | DE-Berlin | 0 | 0 |
| PROD-789 | US-NewYork | 15 | 14 |
| PROD-101 | DE-Berlin | 10 | 9 |
| PROD-101 | US-NewYork | 11 | 10 |
| PROD-202 | DE-Berlin | 1 | 1 |
| PROD-202 | UK-London | 500 | 450 |
| PROD-303 | UK-London | 20 | 18 |
| PROD-404 | DE-Berlin | 200 | 180 |
| PROD-505 | US-NewYork | 5 | 5 |

## Edge Cases Tested

### Reserve Buffer Edge Cases
- Stock = 1 (reserve = 0, available = 1)
- Stock = 5 (reserve = 0, available = 5)
- Stock = 10 (reserve = 1, available = 9)
- Stock = 100 (reserve = 10, available = 90)
- Stock = 500 (reserve = 50, available = 450)

### Availability Edge Cases
- Exactly at threshold (request = available)
- Just below threshold (request < available)
- Just above threshold (request > available)
- Out of stock (stock = 0)
- Very low stock (stock = 1)

### Product/Warehouse Edge Cases
- Product exists in warehouse A but not B
- Product doesn't exist at all
- Warehouse name mismatch

## Manual Testing of Weekend Logic

**Note:** Weekend logic is NOT covered by automated tests. Manual testing is required.

The weekend logic doubles the required quantity on Saturday/Sunday. To manually test:

### Method 1: Test on Actual Weekend
1. Wait until Saturday or Sunday
2. Start the server: `cd app && go run .`
3. Make a request for 5 units of PROD-123 in DE-Berlin:
```bash
curl -X POST http://localhost:8080/api/check-availability \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "PROD-123",
    "quantity": 5,
    "warehouse_location": "DE-Berlin"
  }'
```
4. Verify the response reason mentions "weekend: requires 10 units in stock for 5 order"

### Method 2: Temporarily Modify Code for Testing
To test without waiting for the weekend:

1. Open `app/availability.go`
2. Temporarily modify the `isWeekend()` function:
```go
// For testing: force weekend mode
func isWeekend() bool {
    return true  // Always return true for testing
}
```
3. Restart server and test as above
4. **Important:** Revert changes after testing

### Expected Behavior

**Weekday Request:**
- Request: 5 units
- Required stock (after 10% reserve): 5 units
- For PROD-123 in DE-Berlin: Available (90 available, 5 needed)

**Weekend Request:**
- Request: 5 units  
- Required stock (after 10% reserve): **10 units** (2x)
- For PROD-123 in DE-Berlin: Available (90 available, 10 needed)
- Response includes: "weekend: requires 10 units in stock for 5 order"

**Weekend Stress Test:**
- Request: 50 units of PROD-123
- Weekend required: 100 units
- Available: 90 units (100 - 10%)
- Expected: Available = false, "Insufficient stock (weekend: requires 100 units, only 90 available after reserve)"

## Continuous Testing

For development with continuous testing:

```bash
# Install a file watcher like entr
cd app
ls *.go | entr -c go test -v
```

This will automatically run tests when any Go file changes.
