# OpenTofu Station

A Go-based API toolkit for managing OpenTofu provisioning operations with enterprise-grade features including database persistence, comprehensive logging, and robust error handling.

## Features

- **OpenTofu Command Execution**: Execute any OpenTofu command with proper validation and error handling
- **Database Persistence**: Store operation history, plans, and state information in PostgreSQL or SQLite
- **Comprehensive API**: Full gRPC/Protobuf interface for all OpenTofu operations
- **Configuration Management**: Flexible configuration with environment variables and config files
- **Docker Support**: Ready-to-use Docker containers with docker-compose
- **Testing Framework**: Comprehensive test suite with mocking support
- **Error Handling**: Structured error handling with error codes and detailed messages
- **Working Directory Management**: Safe and validated working directory handling
- **Timeout Management**: Configurable timeouts for long-running operations

## Prerequisites

- Go 1.21 or later
- OpenTofu CLI installed and accessible
- PostgreSQL (optional, SQLite is supported for development)
- Protocol Buffers compiler (protoc)

## Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/ForestMars/TerraformStation.git
cd TerraformStation
```

2. Install dependencies:
```bash
go mod download
```

3. Generate protobuf code:
```bash
make proto
```

4. Build the application:
```bash
make build
```

### Using Docker

1. Build and run with docker-compose:
```bash
docker-compose up --build
```

2. Or build the Docker image manually:
```bash
docker build -t opentofu-station .
docker run -p 8080:8080 opentofu-station
```

## Quick Start

### Basic Usage

1. Create a working directory for your OpenTofu files:
```bash
make sample-tofu
```

2. Test OpenTofu commands manually:
```bash
make tofu-test
```

3. Run the application:
```bash
make run
```

4. The application will start and execute a test OpenTofu version command.

### Configuration

The application can be configured using:

- **Environment variables**: Set database and OpenTofu configuration
- **Command line flags**: Override specific settings
- **Configuration file**: YAML-based configuration

Example configuration:
```yaml
opentofu_path: "tofu"
working_directory: "./tofu"
timeout: "30m"
database:
  driver: "sqlite"
  database: "opentofu_station.db"
```

### Command Line Options

```bash
./opentofu-station \
  --opentofu /usr/local/bin/tofu \
  --workdir ./my-opentofu-project \
  --port 9090 \
  --host 0.0.0.0 \
  --db-driver sqlite
```

## Architecture

### Core Components

- **API Layer**: Protobuf-based service definitions
- **Service Layer**: Business logic implementation
- **Database Layer**: GORM-based data persistence
- **Executor Layer**: OpenTofu command execution
- **Configuration Layer**: Flexible configuration management

### Project Structure

```
TerraformStation/
├── api.go              # Service interface definitions
├── config.go           # Configuration structures
├── errors.go           # Error handling and types
├── utils.go            # Utility functions
├── models.go           # Database models
├── database.go         # Database management
├── main.go             # Application entry point
├── factory/            # Dependency injection
├── internal/           # Implementation details
├── mock/               # Generated mock files
├── config/             # Configuration files
├── tofu/               # OpenTofu configuration files
├── spec.proto          # Protobuf definitions
├── Makefile            # Build and development tasks
├── Dockerfile          # Container definition
└── docker-compose.yml  # Development environment
```

## Testing

Run the test suite:

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/...
```

### Test Dependencies

The project uses:
- **testify**: Assertion and mocking library
- **SQLite**: In-memory database for testing
- **GORM**: Database ORM for testing

## Development

### Available Make Targets

```bash
make help          # Show available targets
make build         # Build the application
make test          # Run tests
make clean         # Clean build artifacts
make proto         # Generate protobuf code
make mock          # Generate mock files
make run           # Run the application
make fmt           # Format code
make lint          # Lint code
make deps          # Install dependencies
make sample-tofu   # Create sample OpenTofu directory
make tofu-test     # Test OpenTofu commands manually
```

### Adding New Features

1. **Update Protobuf**: Modify `spec.proto` for new API endpoints
2. **Regenerate Code**: Run `make proto` to update generated code
3. **Implement Service**: Add implementation in `internal/impl.go`
4. **Add Tests**: Create tests in the appropriate package
5. **Update Models**: Add database models if needed

### Code Style

- Follow Go standard formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Include error handling for all operations
- Write tests for new functionality

## Database Schema

The application uses the following database tables:

- **terraform_operations**: Stores all OpenTofu command executions
- **terraform_plans**: Stores plan results and metadata
- **terraform_applies**: Stores apply results and resource counts
- **terraform_states**: Stores state information and metadata

## Security Considerations

- Working directory validation prevents directory traversal attacks
- Command execution with proper timeout limits
- Database connection security (SSL, authentication)
- Input validation for all user-provided data

## Error Handling

The application provides structured error handling with:

- **Error Codes**: Standardized error codes for different failure types
- **Detailed Messages**: Human-readable error descriptions
- **Context Information**: Additional details for debugging
- **Error Recovery**: Graceful handling of common failure scenarios

## Monitoring and Observability

- **Health Checks**: Built-in health check endpoints
- **Metrics**: Prometheus-compatible metrics (planned)
- **Logging**: Structured logging with configurable levels
- **Tracing**: Request tracing for debugging (planned)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:

- Create an issue on GitHub
- Check the documentation
- Review the test examples
- Examine the configuration options

## Roadmap

- [ ] HTTP REST API endpoints
- [ ] Web UI dashboard
- [ ] Advanced OpenTofu state management
- [ ] Multi-cloud provider support
- [ ] Backup and restore functionality
- [ ] Advanced monitoring and alerting
- [ ] Plugin system for custom operations
- [ ] CI/CD integration examples

