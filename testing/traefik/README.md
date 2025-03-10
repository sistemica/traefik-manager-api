# Traefik Manager Configuration Demo

This demo showcases how to serve dynamically generated configurations to Traefik v3 using the Traefik Manager. It demonstrates how to:
- Generate Traefik configurations on the fly
- Serve configurations via HTTP provider
- Route traffic to a demo service (httpbin)

## Project Structure

```
.
├── docker-compose.yml  # Orchestrates Traefik, and httpbin
└── traefik.yml        # Static Traefik configuration
```

## Components

- **Traefik Manager**: A Go HTTP server that generates and serves Traefik's dynamic configuration - needs to be started seperately
- **Traefik**: The reverse proxy/load balancer that reads config from our server
- **Httpbin**: A demo service that our configuration routes to

## Quick Start

1. Start the stack:
```bash
docker compose up -d
```

2. Access the services:
- Traefik Dashboard: http://localhost:8080/dashboard/
- Httpbin service: http://localhost/httpbin/get

## How It Works

1. **Config Generation**: The Go server (Traefik Manager) generates a dynamic configuration that includes:
   - A router for the httpbin service
   - Service backend configuration
   - Middleware to strip path prefixes

2. **Configuration Polling**: Traefik polls the config server every 5 seconds (configurable in `traefik.yml`) to get updated configurations

3. **Network Setup**: All services run in a Docker network called `traefik-net`, allowing them to communicate with each other

## Configuration Details

### Traefik Configuration (`traefik.yml`)
```yaml
providers:
  http:
    endpoint: "http://config-server:9000/api/config"
    pollInterval: "5s"
```

### Dynamic Configuration (served by Go server)
The config server generates and serves:
- Router rules
- Service definitions
- Middleware configurations

### Docker Setup
- All services run in a dedicated network
- Config server is built from the Go source
- Traefik exposes ports 80 (HTTP) and 8080 (Dashboard)


## Notes

- The dashboard is enabled in insecure mode for demo purposes - don't use this in production
- The config server provides a basic example - add error handling and validation for production use

## Requirements

- Docker
- Docker Compose