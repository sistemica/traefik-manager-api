# Traefik Manager

<p align="center">
  <img src="traefik-manager.jpg" alt="Traefik Manager" width="400"/>
</p>

<p align="center">
  <a href="https://github.com/sistemica/traefik-manager/actions/workflows/test.yml"><img src="https://github.com/sistemica/traefik-manager/actions/workflows/test.yml/badge.svg" alt="Tests"></a>
  <a href="https://github.com/sistemica/traefik-manager/releases/latest"><img src="https://img.shields.io/github/v/release/sistemica/traefik-manager" alt="Release"></a>
  <a href="https://github.com/sistemica/traefik-manager/actions/workflows/docker-publish.yml"><img src="https://github.com/sistemica/traefik-manager/actions/workflows/docker-publish.yml/badge.svg" alt="Docker Build"></a>
  <a href="https://goreportcard.com/report/github.com/sistemica/traefik-manager"><img src="https://goreportcard.com/badge/github.com/sistemica/traefik-manager" alt="Go Report Card"></a>
  <a href="https://github.com/sistemica/traefik-manager/blob/main/LICENSE"><img src="https://img.shields.io/github/license/sistemica/traefik-manager" alt="License"></a>
  <a href="https://hub.docker.com/r/sistemica/traefik-manager"><img src="https://img.shields.io/docker/pulls/sistemica/traefik-manager" alt="Docker Pulls"></a>
  <a href="https://buymeacoffee.com/hannes"><img src="https://img.shields.io/badge/support-buymeacoffee-yellow" alt="Buy Me A Coffee"></a>
</p>

## Introduction

Traefik Manager is a powerful RESTful API that simplifies the management of [Traefik](https://traefik.io/) 3.3's dynamic configuration. It provides a clean, intuitive interface for managing routes, services, and middleware without directly editing configuration files.

Whether you're running a simple home lab or a complex microservices architecture, Traefik Manager helps you:

- **Simplify Management**: Configure your Traefik instance through a simple HTTP API
- **Centralize Configuration**: Store your routing rules, services, and middleware in one place
- **Automate Deployments**: Integrate with your CI/CD pipelines for automated configuration
- **Ensure Consistency**: Validate your configuration before it reaches Traefik

## Features

- HTTP API for managing Traefik 3.3's dynamic configuration
- File-based persistence with in-memory caching
- Dynamic configuration provider endpoint for Traefik
- Flexible authentication options for both API and provider endpoints
- Health check endpoint
- Support for all Traefik middleware types
- Support for various service types (load balancer, weighted, mirroring, failover)
- Multi-architecture support (amd64, arm64)

## Documentation

### Package-Level Documentation

Each package in the `internal` directory contains its own README with detailed information about its purpose and usage:

- [`internal/api`](internal/api/README.md): HTTP API implementation details
  - Learn about route registration, request handling, and API structure

- [`internal/config`](internal/config/README.md): Configuration management
  - Understand how environment variables are parsed and validated
  - See all available configuration options

- [`internal/logger`](internal/logger/README.md): Structured logging
  - Explore logging configurations and best practices
  - Learn about log levels, formats, and output destinations

- [`internal/middleware`](internal/middleware/README.md): HTTP middleware components
  - Discover authentication, logging, and recovery middleware
  - Understand middleware order and configuration

- [`internal/models`](internal/models/README.md): Data models for Traefik resources
  - Detailed explanation of router, service, and middleware models
  - Understand the structure of Traefik configurations

- [`internal/store`](internal/store/README.md): Data persistence layer
  - Learn about the file-based storage implementation
  - Understand resource management and dependency tracking

### API Documentation

For comprehensive API usage examples, refer to our [API Examples](docs/api-examples.md) document. This guide provides:

- Curl commands for all API endpoints
- Examples of creating and managing:
  - Middlewares
  - Services
  - Routers
- Common use cases and configuration patterns

### Additional Resources

- **Architecture Overview**: [`ARCHITECTURE.md`](ARCHITECTURE.md)
  - Deep dive into the system design
  - Understand the interaction between components

- **Testing Documentation**: [`test/README.md`](test/README.md)
  - Learn about our testing strategy
  - Understand unit and integration test approaches

### External Documentation

- [Traefik Documentation](https://doc.traefik.io/traefik/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

## Support the Project

If you find Traefik Manager useful, please consider supporting its development:

<p align="center">
  <a href="https://buymeacoffee.com/hannes">
    <img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee">
  </a>
</p>

Your support helps maintain this project and fund upcoming features including:
- A user-friendly web UI for Traefik Manager
- A Traefik-based Ngrok alternative for exposing local services



## Project Structure

```
.
├── cmd
│   └── server               # Main entry point for the application
│       ├── main.go          # Application startup
│       └── README.md
├── data                     # Default data storage location
│   └── traefik-manager.json # Persisted configuration
├── docs                     # Documentation
│   └── api-examples.md      # API usage examples
├── internal                 # Internal packages
│   ├── api                  # API implementation
│   │   ├── handlers         # Request handlers for each resource type
│   │   ├── routes           # Route definitions
│   │   └── server           # HTTP server setup
│   ├── config               # Configuration loading and validation
│   ├── logger               # Structured logging
│   ├── middleware           # HTTP middleware (auth, logging, recovery)
│   ├── models               # Data models for Traefik resources
│   ├── store                # Data persistence
│   └── traefik              # Traefik-specific models and mapping
├── scripts                  # Utility scripts
├── testing                  # Testing configurations
│   └── traefik              # Traefik test setup
├── .github/workflows        # GitHub Actions workflows
├── Dockerfile               # Multi-stage Docker build
├── docker-compose.yml       # Docker Compose setup with Traefik and httpbin
├── go.mod                   # Go module definition
├── go.sum                   # Go module checksums
├── Makefile                 # Build and development commands
└── README.md                # This file
```

## API Endpoints

### Health Check

- `GET /api/v1/health` - Get service health status

### Traefik Configuration Provider

- `GET /traefik/provider` - Dynamic configuration provider endpoint for Traefik

### Routers

- `GET /api/v1/routers` - List all routers
- `GET /api/v1/routers/{id}` - Get a specific router
- `POST /api/v1/routers` - Create a new router
- `PUT /api/v1/routers/{id}` - Update an existing router
- `DELETE /api/v1/routers/{id}` - Delete a router

### Services

- `GET /api/v1/services` - List all services
- `GET /api/v1/services/{id}` - Get a specific service
- `POST /api/v1/services` - Create a new service
- `PUT /api/v1/services/{id}` - Update an existing service
- `DELETE /api/v1/services/{id}` - Delete a service

### Middlewares

- `GET /api/v1/middlewares` - List all middlewares
- `GET /api/v1/middlewares/{id}` - Get a specific middleware
- `POST /api/v1/middlewares` - Create a new middleware
- `PUT /api/v1/middlewares/{id}` - Update an existing middleware
- `DELETE /api/v1/middlewares/{id}` - Delete a middleware

## Authentication

Traefik Manager provides flexible authentication options for both the API endpoints and the Traefik provider endpoint.

### Global Authentication

You can enable API key-based authentication for all endpoints with these environment variables:

```bash
AUTH_ENABLED=true
AUTH_HEADER_NAME=X-API-Key  # Default header name, can be customized
AUTH_KEY=your-secure-api-key
```

When global authentication is enabled, all API endpoints will require the specified API key in the request headers, except for:
- The `/health` endpoint (always public)
- Any paths specifically excluded in the configuration

### Provider Endpoint Authentication

The Traefik provider endpoint (`/traefik/provider`) can have its own authentication settings, which can be different from the global authentication:

```bash
PROVIDER_AUTH_ENABLED=true
PROVIDER_AUTH_HEADER_NAME=X-API-Key  # Default header name, can be customized
PROVIDER_AUTH_KEY=your-provider-key  # Can be different from the global key
```

#### Authentication Scenarios

1. **No Authentication**: If `AUTH_ENABLED=false` and `PROVIDER_AUTH_ENABLED=false`, all endpoints are publicly accessible.

2. **Global Authentication Only**: If `AUTH_ENABLED=true` and `PROVIDER_AUTH_ENABLED=false`, all endpoints require the global API key except for `/health` and the provider endpoint.

3. **Provider Authentication Only**: If `AUTH_ENABLED=false` and `PROVIDER_AUTH_ENABLED=true`, only the provider endpoint requires authentication.

4. **Separate Authentication**: If both are enabled with different keys, the API uses the global key while the provider endpoint uses its specific key.

#### Example Usage

For Traefik to authenticate with the provider endpoint:

```yaml
# traefik.yml
providers:
  http:
    endpoint: "http://traefik-manager:9000/traefik/provider"
    headers:
      X-API-Key: "your-provider-key"
```

For API clients:

```bash
curl -H "X-API-Key: your-secure-api-key" http://localhost:9000/api/v1/routers
```

#### Security Recommendations

1. Use strong, randomly generated API keys
2. Use HTTPS in production environments
3. Regularly rotate API keys
4. Use different keys for the provider and API if possible

## Prerequisites

- Go 1.23.5 or later
- Traefik v3.3
- Docker and Docker Compose (for container deployment)

## Configuration

Traefik Manager is configured using environment variables:

### Server Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_HOST` | Host address to bind to | `""` (all interfaces) |
| `SERVER_PORT` | Port to listen on | `9000` |
| `SERVER_BASE_PATH` | Base path for all API endpoints | `/api/v1` |
| `SERVER_READ_TIMEOUT` | HTTP read timeout | `15s` |
| `SERVER_WRITE_TIMEOUT` | HTTP write timeout | `15s` |

### Storage Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `STORAGE_FILE_PATH` | Path to the storage file | System temporary file |
| `STORAGE_SAVE_INTERVAL` | Debounce interval for saving changes | `5s` |

### Provider Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `PROVIDER_PATH` | Path for the Traefik provider endpoint | `/traefik/provider` |
| `PROVIDER_AUTH_ENABLED` | Enable API key authentication for provider | `false` |
| `PROVIDER_AUTH_HEADER_NAME` | API key header name for provider | `X-API-Key` |
| `PROVIDER_AUTH_KEY` | API key value for provider | `""` (required if PROVIDER_AUTH_ENABLED is true) |

### Logger Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `LOG_LEVEL` | Log level (debug, info, warn, error) | `info` |
| `LOG_FORMAT` | Log format (json, text) | `json` |
| `LOG_FILE_PATH` | Path to log file, empty for stdout | `""` |
| `LOG_USE_COLORS` | Use colors in console output | `true` |

### CORS Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `CORS_ALLOWED_ORIGINS` | Comma-separated list of allowed origins | `*` |
| `CORS_ALLOWED_METHODS` | Comma-separated list of allowed methods | `GET,POST,PUT,DELETE,OPTIONS` |
| `CORS_ALLOWED_HEADERS` | Comma-separated list of allowed headers | `Content-Type,Authorization` |
| `CORS_ALLOW_CREDENTIALS` | Allow credentials | `false` |
| `CORS_MAX_AGE` | Max age in seconds | `300` |

### Auth Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `AUTH_ENABLED` | Enable API key authentication | `false` |
| `AUTH_HEADER_NAME` | API key header name | `X-API-Key` |
| `AUTH_KEY` | API key value | `""` (required if AUTH_ENABLED is true) |

## Getting Started

### Local Development Setup

1. Clone the repository:

```bash
git clone https://github.com/sistemica/traefik-manager.git
cd traefik-manager
```

2. Build and run the application:

```bash
# Build the application
make build

# Run the application
make run
```

### Docker Setup

Build and run using Docker:

```bash
# Build the Docker image
make docker-build

# Run with Docker
make docker-run
```

### Docker Compose Setup

Run the full stack with Traefik and httpbin:

```bash
# Start all services
make docker-compose-up

# View logs
make docker-logs

# Stop all services
make docker-compose-down
```

## Development

### Available Make Commands

```bash
# Show all available commands
make help
```

Build commands:
- `make build` - Build the application
- `make clean` - Remove generated files
- `make all` - Clean, build, and test

Development commands:
- `make run` - Run the application locally
- `make test` - Run tests
- `make fmt` - Format code
- `make lint` - Lint code

Docker commands:
- `make docker-build` - Build Docker image
- `make docker-run` - Run Docker container
- `make docker-compose-up` - Start all services with Docker Compose
- `make docker-compose-down` - Stop all services
- `make docker-logs` - Show logs from Docker Compose services

## Integration with Traefik

### Traefik Configuration

To use Traefik Manager as a provider for Traefik, configure Traefik as follows:

```yaml
providers:
  http:
    endpoint: "http://traefik-manager:9000/traefik/provider"
    pollInterval: "5s"
    # Add authentication if enabled
    headers:
      X-API-Key: "your-provider-key"
```

### Workflow

1. Traefik Manager serves dynamic configuration via the provider endpoint
2. Traefik polls this endpoint at regular intervals (defined by `pollInterval`)
3. When changes are made via the API, they are automatically picked up by Traefik on the next poll

## Contributing

Contributions are welcome! Here's how you can contribute to Traefik Manager:

### Reporting Issues

If you find a bug or have a feature request:
1. Check if the issue already exists in the [Issues](https://github.com/sistemica/traefik-manager/issues)
2. If not, create a new issue with a descriptive title
3. Include steps to reproduce (for bugs) or a detailed description (for features)

### Pull Requests

We welcome code contributions:
1. Fork the repository
2. Create a new branch for your feature/fix (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run the tests (`make test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to your branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Guidelines

- Follow Go best practices and code style
- Write unit tests for your code
- Update documentation for new features
- Keep PR scope focused - one feature or bug fix per PR
- Ensure all tests pass before submitting PRs

### Setup Development Environment

```bash
# Clone repository
git clone https://github.com/sistemica/traefik-manager.git
cd traefik-manager

# Install dependencies
go mod download

# Run tests
make test

# Run local development server
make run
```

## Troubleshooting

### Common Issues

- **API returns 404**: Check that your routes are properly registered and the base path is correct
- **Authentication failures**: Verify that the correct API key is being sent in the specified header
- **Traefik doesn't apply configurations**: Check that the provider URL is correct and that Traefik can reach Traefik Manager

### Debugging

Enable debug logging:

```bash
export LOG_LEVEL=debug
make run
```

## API Examples

For detailed API examples, see [API Examples](docs/api-examples.md).

## Internal Package Details

### cmd/server

The main entry point for the application. It initializes all components and starts the HTTP server.

### internal/api

Contains the HTTP API implementation:

- **handlers**: Request handlers for each resource type (routers, services, middlewares)
- **routes**: Route definitions and registration
- **server**: HTTP server setup and middleware configuration

### internal/config

Manages application configuration loading from environment variables with support for `.env` files.

### internal/logger

Provides structured logging with support for different formats and log levels.

### internal/middleware

HTTP middleware components:

- **auth.go**: API key authentication
- **logging.go**: Request logging
- **recovery.go**: Panic recovery

### internal/models

Data models for Traefik resources:

- **models.go**: Core models for routers, services, etc.
- **middleware_configs.go**: Models for all middleware configurations

### internal/store

Data persistence implementation:

- **file.go**: File-based storage implementation
- **store.go**: Storage interface definition
- **errors.go**: Error types and handling

### internal/traefik

Traefik-specific models and mapping between internal models and Traefik's expected format.

## License

[MIT License](LICENSE)