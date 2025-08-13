# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git protobuf

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate protobuf code
RUN protoc --go_out=. --go_opt=paths=source_relative spec.proto

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o opentofu-station main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates opentofu

# Create non-root user
RUN addgroup -g 1001 -S opentofu && \
    adduser -u 1001 -S opentofu -G opentofu

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/opentofu-station .

# Create opentofu working directory
RUN mkdir -p /app/tofu && \
    chown -R opentofu:opentofu /app

# Switch to non-root user
USER opentofu

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./opentofu-station"]
