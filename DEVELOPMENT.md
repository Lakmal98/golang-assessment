# Development Notes

## AI & Human Collaboration

### Development Process

**AI-Generated Components (GitHub Copilot):**
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
- 🔧 Requested documentation cleanup (removed unimplemented features)
- 🔧 Fixed test failures and validation logic

**Human Verification:**
- ✓ Reviewed architecture and design patterns
- ✓ Validated business logic correctness
- ✓ Verified Docker configuration
- ✓ Confirmed test coverage meets requirements
- ✓ Ensured documentation accuracy

### Code Evolution Timeline

1. **Initial Implementation** - Single-file monolithic API
2. **Refactoring** - Multi-file structure with layers
3. **Data Abstraction** - Adapter pattern for flexible storage
4. **Testing** - Comprehensive test suite with mocking
5. **Containerization** - Docker with development watch mode
6. **Documentation** - OpenAPI/Swagger integration
7. **Organization** - Structured documentation split
8. **Quality Assurance** - Test fixes and documentation cleanup

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

**Tests failing:**
```bash
# Run tests with verbose output
cd app
go test -v

# Check for import errors
go mod tidy
```

## Development Workflow

### With Docker Watch Mode (Recommended)

```bash
# Start with auto-reload
docker compose watch

# Edit files in app/
# - .go files: Auto-rebuild container
# - inventory.json: Auto-sync without rebuild
```

### Local Development

```bash
# Run server
cd app
go run .

# Run tests in watch mode (requires entr)
ls *.go | entr -c go test -v
```

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run tests with coverage
go test -v -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```
