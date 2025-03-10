# Configuration Package

This package provides configuration management for the Traefik Manager application. It loads configuration from environment variables with support for `.env` files.

## Usage

```go
import "github.com/sistemica/traefik-manager/internal/config"

func main() {
    // Load configuration from .env file
    cfg, err := config.LoadConfig(".env")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Use configuration
    serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
    fmt.Printf("Starting server on %s\n", serverAddr)
    // ...
}
```

## Environment Variables

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
| `STORAGE_FILE_PATH` | Path to the storage file | `./data/traefik-manager.json` |
| `STORAGE_SAVE_INTERVAL` | Debounce interval for saving changes | `5s` |

### Traefik Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `TRAEFIK_PROVIDER_PATH` | Path for the Traefik provider endpoint | `/api/provider/traefik` |
| `TRAEFIK_AUTH_ENABLED` | Enable API key authentication | `false` |
| `TRAEFIK_AUTH_HEADER_NAME` | API key header name | `X-API-Key` |
| `TRAEFIK_AUTH_KEY` | API key value | `""` (required if AUTH_ENABLED is true) |

> **Future Development (Optional)**: Future versions may include `TRAEFIK_API_URL` for active monitoring of Traefik instances.

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

## .env File Example

```env
# Server configuration
SERVER_PORT=9000
SERVER_BASE_PATH=/api/v1

# Storage configuration
STORAGE_FILE_PATH=./data/traefik-manager.json

# TRAEFIK_API_URL=http://traefik:8080  # For future monitoring features

# Logger configuration
LOG_LEVEL=debug

# CORS configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://admin.example.com

# Auth configuration
AUTH_ENABLED=true
AUTH_KEY=your-secret-api-key
```