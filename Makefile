.PHONY: help build test clean proto mock run docker-build docker-run

# Default target
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  proto       - Generate protobuf code"
	@echo "  mock        - Generate mock files"
	@echo "  run         - Run the application"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"

# Build the application
build: proto
	@echo "Building OpenTofu Station..."
	go build -o bin/opentofu-station main.go

# Run tests
test: proto
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f *.pb.go
	rm -f mock/*.go

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative spec.proto

# Generate mock files
mock: proto
	@echo "Generating mock files..."
	go install github.com/matryer/moq@latest
	moq -out=mock/mock.go -pkg=mock . TerraformStationService

# Run the application
run: build
	@echo "Running OpenTofu Station..."
	./bin/opentofu-station

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t opentofu-station .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 opentofu-station

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Generate documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060

# Create sample opentofu directory
sample-tofu:
	@echo "Creating sample OpenTofu directory..."
	mkdir -p tofu
	@echo 'terraform { required_version = ">= 1.0" }' > tofu/main.tf
	@echo 'resource "local_file" "hello" { filename = "${path.module}/hello.txt" content = "Hello, World from OpenTofu!" }' >> tofu/main.tf
	@echo 'output "hello_content" { value = local_file.hello.content }' >> tofu/main.tf

# Test OpenTofu commands
tofu-test:
	@echo "Testing OpenTofu commands..."
	cd tofu && tofu init
	cd tofu && tofu plan
	cd tofu && tofu apply -auto-approve
	cd tofu && tofu output
