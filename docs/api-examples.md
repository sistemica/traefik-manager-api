# Traefik Manager API Examples

This document provides practical examples for testing the Traefik Manager API with curl commands. 

## Table of Contents

1. [Setup](#setup)
2. [Health Check](#health-check)
3. [Traefik Provider Endpoint](#traefik-provider-endpoint)
4. [Middleware Operations](#middleware-operations)
5. [Service Operations](#service-operations)
6. [Router Operations](#router-operations)
7. [Common Use Cases](#common-use-cases)

## Setup

Before running these examples, make sure both Traefik and Traefik Manager are running. You can use the docker-compose setup provided in the main project.

```bash
# Set up environment variables for easier testing
export API_URL="http://localhost:9000/api/v1"
export AUTH_TOKEN="your-auth-token"  # If you have auth enabled
```

## Health Check

```bash
# Check if the API is healthy
curl -s "${API_URL}/health" | jq
```

Expected response:
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": "3h2m15s"
}
```

## Traefik Provider Endpoint

Returns the current config which will be consumed by Traefik:

```bash
curl -s "localhost:9000/traefik/provider" | jq
```

## Middleware Operations

### List All Middlewares
```bash
curl -s "${API_URL}/middlewares" | jq
```

### Get Specific Middleware
```bash
curl -s "${API_URL}/middlewares/redirect-https" | jq
```

### Create Middleware
```bash
curl -X POST "${API_URL}/middlewares" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "add-prefix",
    "type": "addPrefix",
    "config": {
      "prefix": "/api"
    }
  }'
```

### Update Middleware
```bash
curl -X PUT "${API_URL}/middlewares/add-prefix" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "addPrefix",
    "config": {
      "prefix": "/v1/api"
    }
  }'
```

### Delete Middleware
```bash
curl -X DELETE "${API_URL}/middlewares/add-prefix"
```

## Service Operations

### List All Services
```bash
curl -s "${API_URL}/services" | jq
```

### Get Specific Service
```bash
curl -s "${API_URL}/services/blog-service" | jq
```

### Create Simple Service
```bash
curl -X POST "${API_URL}/services" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "api-service",
    "url": "http://api-backend:8080"
  }'
```

### Create Load Balancer Service
```bash
curl -X POST "${API_URL}/services" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "balanced-service",
    "loadBalancer": {
      "servers": [
        {"url": "http://server1:8080"},
        {"url": "http://server2:8080"},
        {"url": "http://server3:8080"}
      ],
      "healthCheck": {
        "path": "/health",
        "interval": "10s",
        "timeout": "3s"
      }
    }
  }'
```

### Create Weighted Service
```bash
curl -X POST "${API_URL}/services" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "weighted-service",
    "weighted": {
      "services": [
        {
          "name": {"id": "api-service"},
          "weight": 8
        },
        {
          "name": {"id": "backup-service"},
          "weight": 2
        }
      ]
    }
  }'
```

### Create Mirroring Service
```bash
curl -X POST "${API_URL}/services" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "mirror-service",
    "mirroring": {
      "service": {"id": "main-service"},
      "mirrors": [
        {
          "name": {"id": "analytics-service"},
          "percent": 10
        }
      ]
    }
  }'
```

### Create Failover Service
```bash
curl -X POST "${API_URL}/services" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "failover-service",
    "failover": {
      "service": {"id": "primary-service"},
      "fallback": {"id": "backup-service"}
    }
  }'
```

### Update Service
```bash
curl -X PUT "${API_URL}/services/api-service" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "http://api-backend:9000"
  }'
```

### Delete Service
```bash
curl -X DELETE "${API_URL}/services/api-service"
```

## Router Operations

### List All Routers
```bash
curl -s "${API_URL}/routers" | jq
```

### Get Specific Router
```bash
curl -s "${API_URL}/routers/blog-router" | jq
```

### Create Basic Router
```bash
curl -X POST "${API_URL}/routers" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "api-router",
    "rule": "Host(`api.example.com`)",
    "service": "api-service",
    "entryPoints": ["web", "websecure"],
    "middlewares": ["add-prefix", "redirect-https"]
  }'
```

### Create Router with TLS Configuration
```bash
curl -X POST "${API_URL}/routers" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "secure-router",
    "rule": "Host(`secure.example.com`)",
    "service": "api-service",
    "entryPoints": ["websecure"],
    "tls": {
      "certResolver": "default",
      "domains": [
        {
          "main": "secure.example.com",
          "sans": ["www.secure.example.com"]
        }
      ]
    }
  }'
```

### Create Router with Path-Based Routing
```bash
curl -X POST "${API_URL}/routers" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "path-router",
    "rule": "Host(`example.com`) && PathPrefix(`/api`)",
    "service": "api-service",
    "entryPoints": ["websecure"],
    "priority": 100
  }'
```

### Create Router with Observability Settings
```bash
curl -X POST "${API_URL}/routers" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "monitored-router",
    "rule": "Host(`monitoring.example.com`)",
    "service": "api-service",
    "entryPoints": ["websecure"],
    "observability": {
      "accessLogs": true,
      "tracing": true,
      "metrics": true
    }
  }'
```

### Update Router
```bash
curl -X PUT "${API_URL}/routers/api-router" \
  -H "Content-Type: application/json" \
  -d '{
    "rule": "Host(`api.example.org`)",
    "service": "api-service",
    "entryPoints": ["websecure"],
    "middlewares": ["redirect-https"]
  }'
```

### Delete Router
```bash
curl -X DELETE "${API_URL}/routers/api-router"
```

## Common Use Cases

### Setting up a Basic Website with HTTPS Redirect

1. Create HTTPS redirect middleware:
```bash
curl -X POST "${API_URL}/middlewares" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "https-redirect",
    "type": "redirectScheme",
    "config": {
      "scheme": "https",
      "permanent": true
    }
  }'
```

2. Create backend service:
```bash
curl -X POST "${API_URL}/services" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "website-service",
    "url": "http://website-backend:80"
  }'
```

3. Create HTTP router with redirect:
```bash
curl -X POST "${API_URL}/routers" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "website-http",
    "rule": "Host(`mywebsite.com`)",
    "service": "website-service",
    "entryPoints": ["web"],
    "middlewares": ["https-redirect"]
  }'
```

4. Create HTTPS router:
```bash
curl -X POST "${API_URL}/routers" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "website-https",
    "rule": "Host(`mywebsite.com`)",
    "service": "website-service",
    "entryPoints": ["websecure"],
    "tls": {
      "certResolver": "default"
    }
  }'
```

### Setting up API with Rate Limiting

1. Create rate limit middleware:
```bash
curl -X POST "${API_URL}/middlewares" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "api-rate-limit",
    "type": "rateLimit",
    "config": {
      "average": 100,
      "burst": 50
    }
  }'
```

2. Create API service:
```bash
curl -X POST "${API_URL}/services" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "public-api",
    "url": "http://api-backend:3000"
  }'
```

3. Create API router with rate limiting:
```bash
curl -X POST "${API_URL}/routers" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "api-public",
    "rule": "Host(`api.example.com`)",
    "service": "public-api",
    "entryPoints": ["websecure"],
    "middlewares": ["api-rate-limit"],
    "tls": {}
  }'
```

### Setting up Basic Authentication for Admin Panel

1. Create basic auth middleware (password must be hashed using htpasswd):
```bash
curl -X POST "${API_URL}/middlewares" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "admin-auth",
    "type": "basicAuth",
    "config": {
      "users": ["admin:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/"],
      "realm": "Admin Area"
    }
  }'
```

2. Create admin service:
```bash
curl -X POST "${API_URL}/services" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "admin-service",
    "url": "http://admin-panel:8080"
  }'
```

3. Create admin router with authentication:
```bash
curl -X POST "${API_URL}/routers" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "admin-router",
    "rule": "Host(`admin.example.com`)",
    "service": "admin-service",
    "entryPoints": ["websecure"],
    "middlewares": ["admin-auth"],
    "tls": {}
  }'
```

### Clean Up All Resources

Delete routers:
```bash
curl -X DELETE "${API_URL}/routers/website-http"
curl -X DELETE "${API_URL}/routers/website-https"
curl -X DELETE "${API_URL}/routers/api-public"
curl -X DELETE "${API_URL}/routers/admin-router"
```

Delete services:
```bash
curl -X DELETE "${API_URL}/services/website-service"
curl -X DELETE "${API_URL}/services/public-api"
curl -X DELETE "${API_URL}/services/admin-service"
```

Delete middlewares:
```bash
curl -X DELETE "${API_URL}/middlewares/https-redirect"
curl -X DELETE "${API_URL}/middlewares/api-rate-limit"
curl -X DELETE "${API_URL}/middlewares/admin-auth"
```