# Product Availability Check API

A REST API endpoint in Go that checks product availability across warehouses with business rules for reserve buffers and weekend shipping delays.

Note: Used GitHub Copilot for code and documentation generation.

**📚 Additional Documentation:**
- [Architecture & Design](ARCHITECTURE.md) - System architecture, design patterns, and technical details
- [Testing Guide](TESTS.md) - Test suite documentation and how to run tests

---

## Quick Start

### Using Docker (Recommended)

```bash
# Start the application
docker compose up

# Access the API
open http://localhost:8080
```

### Using Go Directly

```bash
cd app
go run .
```

The server will start on **http://localhost:8080** and automatically redirect to the Swagger documentation.

---

## Project Structure

```
golang-assesment/
├── app/                    # Application source code
│   ├── main.go            # Entry point & dependency injection
│   ├── handler.go         # HTTP request handlers
│   ├── availability.go    # Business logic service
│   ├── inventory.go       # Adapter pattern for data access
│   ├── models.go          # Data structures
│   ├── openapi.go         # OpenAPI/Swagger specification
│   ├── inventory.json     # Mock inventory data
│   └── *_test.go          # Unit tests
├── Dockerfile             # Multi-stage Docker build
├── docker-compose.yaml    # Docker Compose with watch mode
├── README.md              # This file (setup & deployment)
├── ARCHITECTURE.md        # System architecture & design
└── TESTS.md              # Testing guide
```

---

## Setup & Running

### Prerequisites

**Option 1: Docker**
- Docker (20.10+)
- Docker Compose V2+

**Option 2: Local Go**
- Go 1.16 or later

### Installation

1. **Clone or navigate to the project:**
   ```bash
   cd /path/to/golang-assessment
   ```

2. **Choose your deployment method:**

#### Method 1: Docker (Production-like)

```bash
# Build and start
docker compose up

# Start in background
docker compose up -d

# View logs
docker compose logs -f

# Stop
docker compose down
```

#### Method 2: Docker with Watch Mode (Development)

Hot-reload on code changes:

```bash
# Start with automatic rebuild on file changes
docker compose watch

# Edit files in app/ - changes auto-apply
# Go file changes trigger rebuild
# inventory.json changes sync without rebuild
```

#### Method 3: Local Go Development

```bash
cd app
go run .
```

---

## API Documentation

Once running, access:

- **Root:** [http://localhost:8080/](http://localhost:8080/) → Redirects to Swagger UI
- **Swagger UI:** [http://localhost:8080/docs](http://localhost:8080/docs) - Interactive API documentation
- **OpenAPI Spec:** [http://localhost:8080/openapi.json](http://localhost:8080/openapi.json) - OpenAPI 3.0 JSON

---

## Deployment

### Docker Deployment

The application uses a **multi-stage Dockerfile** for production deployment:

**Build Image:**
```bash
docker build -t product-availability-api .
```

**Run Container:**
```bash
docker run -p 8080:8080 product-availability-api
```

**Deployment Features:**
- ✅ Multi-stage build (~10MB final image)
- ✅ Alpine Linux base for security
- ✅ Statically compiled Go binary
- ✅ No runtime dependencies
- ✅ Production-ready configuration

### Docker Compose Deployment

**Production deployment:**
```bash
docker compose up -d
```

**Update deployment:**
```bash
docker compose up -d --build
```

**Health check:**
```bash
curl http://localhost:8080/docs
```

---

## API Usage

### Main Endpoint

**POST** `/api/check-availability`

**Request:**
```json
{
  "product_id": "PROD-123",
  "quantity": 5,
  "warehouse_location": "DE-Berlin"
}
```

**Response:**
```json
{
  "available": true,
  "available_quantity": 90,
  "reason": "Sufficient stock available",
  "warehouse": "DE-Berlin"
}
```

### Quick Test

```bash
curl -X POST http://localhost:8080/api/check-availability \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "PROD-123",
    "quantity": 5,
    "warehouse_location": "DE-Berlin"
  }'
```

**More examples available in the [Swagger UI](http://localhost:8080/docs)**

---

## Business Rules

1. **Reserve Buffer (10%)**
   - 10% of total stock is always reserved
   - Only available stock can be ordered
   - Example: Stock=100 → Available=90

2. **Weekend Shipping Delay (2x)**
   - Saturday/Sunday orders require 2x quantity in stock
   - Compensates for shipping delays
   - Example: Weekend order of 5 requires 10 available

3. **Availability Logic**
   - Product available if: `available_stock >= required_quantity`
   - Weekend adjustment applied automatically
   - Detailed reasons provided in response

---

## Mock Data

Inventory stored in `app/inventory.json`:

| Product ID | Warehouse | Stock | Available | Notes |
|------------|-----------|-------|-----------|-------|
| PROD-123 | DE-Berlin | 100 | 90 | Normal |
| PROD-123 | US-NewYork | 50 | 45 | Normal |
| PROD-456 | DE-Berlin | 25 | 22 | Normal |
| PROD-456 | US-NewYork | 75 | 67 | Normal |
| PROD-789 | DE-Berlin | 0 | 0 | ⚠️ Out of stock |
| PROD-789 | US-NewYork | 15 | 13 | Low stock |
| PROD-101 | DE-Berlin | 10 | 9 | Edge case |
| PROD-202 | DE-Berlin | 1 | 1 | ⚠️ Very low |
| PROD-202 | UK-London | 500 | 450 | High stock |
| PROD-404 | DE-Berlin | 200 | 180 | Large inventory |

_Full inventory available in `app/inventory.json`_

---

## AI & Human Collaboration

### Development Process

**AI-Generated Components:**
- ✅ Initial monolithic REST API structure
- ✅ Refactored multi-file architecture with adapter pattern
- ✅ Business logic implementation (reserve buffer, weekend rules)
- ✅ OpenAPI 3.0 specification with Swagger UI integration
- ✅ Comprehensive test suite with edge cases
- ✅ Docker multi-stage build configuration
- ✅ Docker Compose with watch mode setup
- ✅ Complete documentation (README, ARCHITECTURE, TESTS)

**Human-Directed Changes:**
- 🔧 Requested professional multi-file structure
- 🔧 Requested adapter pattern for data access
- 🔧 Requested JSON file storage instead of hardcoded data
- 🔧 Requested Docker setup with watch mode
- 🔧 Requested code organization into app/ directory
- 🔧 Requested Swagger documentation auto-generation
- 🔧 Requested root path redirect to /docs
- 🔧 Requested documentation split into multiple files

**Human Verification:**
- ✓ Reviewed architecture and design patterns
- ✓ Validated business logic correctness
- ✓ Verified Docker configuration
- ✓ Confirmed test coverage meets requirements

### Code Evolution Timeline

1. **Initial Implementation** - Single-file monolithic API
2. **Refactoring** - Multi-file structure with layers
3. **Data Abstraction** - Adapter pattern for flexible storage
4. **Testing** - Comprehensive test suite with mocking
5. **Containerization** - Docker with development watch mode
6. **Documentation** - OpenAPI/Swagger integration
7. **Organization** - Structured documentation (this split)

---

## Assumptions & Design Decisions

### Technical Assumptions

1. **Standard Library Only** - Uses only Go standard library (net/http, encoding/json) per requirements
2. **Port 8080** - Default application port (configurable via environment if needed)
3. **Time Zone** - Weekend detection uses server's local time zone
4. **In-Memory Loading** - Inventory loaded once on startup for performance

### Design Decisions

1. **Reserve Buffer Calculation**
   - Uses integer truncation: `int(stock * 0.10)`
   - Example: Stock=15 → Reserve=1 → Available=14

2. **Data Storage**
   - JSON file chosen for simplicity and requirements
   - Adapter pattern allows future database integration
   - File loaded at startup (not per-request)

3. **Weekend Logic**
   - Checked using Go's `time.Weekday()`
   - Saturday and Sunday trigger 2x requirement
   - Applied before availability comparison

4. **Error Handling**
   - Input validation returns 400 Bad Request
   - Method validation returns 405 Method Not Allowed
   - Product not found returns 200 OK with available=false

5. **Architecture Pattern**
   - Clean layered architecture
   - Dependency injection for testability
   - Adapter pattern for data source flexibility

6. **Docker Strategy**
   - Multi-stage build for minimal image size
   - Alpine Linux for security and size
   - Watch mode for development productivity

7. **Documentation Approach**
   - OpenAPI 3.0 for standard compliance
   - Swagger UI via CDN (no external dependencies)
   - Separate MD files for different audiences

### Case Sensitivity

- **Product IDs:** Case-sensitive (PROD-123 ≠ prod-123)
- **Warehouse Locations:** Case-sensitive (DE-Berlin ≠ de-berlin)
- **JSON Fields:** Must match exactly as specified

### Validation Rules

- **quantity:** Must be > 0 (positive integer)
- **product_id:** Required, non-empty string
- **warehouse_location:** Required, non-empty string

---

## Testing

Run the test suite:

```bash
cd app
go test -v
```

For detailed testing information, see [TESTS.md](TESTS.md).

---

## Troubleshooting

**Port already in use:**
```bash
# Find and kill process on port 8080
lsof -ti:8080 | xargs kill -9

# Or use a different port (modify docker-compose.yaml)
```

**Docker build fails:**
```bash
# Clean Docker cache
docker system prune -a

# Rebuild without cache
docker compose build --no-cache
```

**Inventory not loading:**
- Ensure `app/inventory.json` exists
- Check JSON syntax validity
- Verify file permissions

---

## License

This project was created as a technical assessment task.

---

## Contact & Support

For technical details, see [ARCHITECTURE.md](ARCHITECTURE.md)  
For testing information, see [TESTS.md](TESTS.md)
