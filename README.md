# Product Availability Check API

REST API in Go that checks product availability across warehouses with reserve buffers and weekend shipping rules.

**Documentation:**
- [ARCHITECTURE.md](ARCHITECTURE.md) - Technical design
- [TESTS.md](TESTS.md) - Testing guide
- [DEVELOPMENT.md](DEVELOPMENT.md) - Development notes

---

## Quick Start

**Using Docker (Recommended):**
```bash
docker compose up
# Access: http://localhost:8080
```

**Using Go:**
```bash
cd app
go run .
```

**Run Tests:**
```bash
cd app
go test -v
```

---

## Building & Deployment

### Docker Production Build

```bash
# Build image (~10MB)
docker build -t product-availability-api .

# Run container
docker run -p 8080:8080 product-availability-api

# Or use Docker Compose
docker compose up -d
```

### Development with Hot-Reload

```bash
docker compose watch
# Auto-rebuilds on .go changes
# Auto-syncs inventory.json
```

---

## API Usage

**Endpoint:** `POST /api/check-availability`

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

**Test:**
```bash
curl -X POST http://localhost:8080/api/check-availability \
  -H "Content-Type: application/json" \
  -d '{"product_id":"PROD-123","quantity":5,"warehouse_location":"DE-Berlin"}'
```

**Interactive Docs:** [http://localhost:8080/docs](http://localhost:8080/docs)

---

## Business Rules

1. **Reserve Buffer:** 10% of stock always reserved (Stock=100 → Available=90)
2. **Weekend Orders:** Saturday/Sunday require 2x quantity in stock
3. **Availability:** Product available if `available_stock >= required_quantity`

---

## Assumptions & Design Decisions

**Key Assumptions:**
- Uses only Go standard library (per requirements)
- Port 8080 as default
- Weekend detection uses server timezone
- JSON file storage (adapter pattern allows future database swap)
- In-memory inventory loaded at startup

**Design Choices:**
- Reserve buffer: `stock - int(stock * 0.10)` (integer truncation)
- Case-sensitive product IDs and warehouse names
- Product not found returns 200 OK with `available: false` (not 404)
- Layered architecture with dependency injection for testability
- Multi-stage Docker build for minimal image size (~10MB)

**Validation:**
- `quantity` must be positive integer
- `product_id` and `warehouse_location` required

---

## Possible Improvements

**Quick Wins:**
- Environment-based configuration (port, file path)
- Structured JSON logging with request IDs
- Map-based inventory lookup (O(1) vs O(n))
- Health check endpoint for monitoring
- Input whitespace trimming

**Production Features:**
- Database integration (PostgreSQL/MySQL)
- Redis caching layer
- JWT authentication
- Rate limiting per API key
- Prometheus metrics
- CORS configuration

**Scaling:**
- Horizontal scaling with load balancer
- Kubernetes deployment
- Multi-warehouse optimization
- Real-time inventory updates (WebSocket)
- Batch availability checks
- Circuit breaker pattern

---

## Project Structure

```
golang-assesment/
├── app/
│   ├── main.go            # Entry point
│   ├── handler.go         # HTTP handlers
│   ├── availability.go    # Business logic
│   ├── inventory.go       # Data adapter
│   ├── models.go          # Structs
│   ├── openapi.go         # API docs
│   ├── inventory.json     # Mock data
│   └── *_test.go          # Tests
├── Dockerfile
├── docker-compose.yaml
├── README.md              # This file
├── ARCHITECTURE.md
├── TESTS.md
└── DEVELOPMENT.md
```

---

**Note:** Built with GitHub Copilot AI assistance.
