# Build stage
FROM golang:latest AS builder

# Set working directory
WORKDIR /build

# Copy go mod files
COPY app/go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY app/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /build/main .
COPY --from=builder /build/inventory.json .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
