# Traefik Manager Architecture Overview

## Dual Interface Design

The Traefik Manager acts as a bridge between two distinct interfaces:

1. **Client-Facing API**
   - REST API for configuration management
   - Simplified model structure for easier use by clients
   - Operations for creating, reading, updating, deleting configurations

2. **Traefik Provider Endpoint**
   - Single `/api/traefik/provider` endpoint that Traefik polls
   - Must strictly adhere to Traefik's expected configuration format
   - Provides translated configurations from our storage


## Data Flow

```
                    Client Interaction                Traefik Interaction
                    ------------------                -------------------

┌─────────────┐    ┌───────────────────┐              ┌────────────┐
│  API Client │───▶│ Traefik Manager   │◀─────────────│   Traefik  │
│             │◀───│ REST API          │              │            │
└─────────────┘    └───────────────────┘              └────────────┘
                           │                                 ▲
                           ▼                                 │
                   ┌───────────────────┐              ┌────────────────┐
                   │     Storage       │───────────────▶  Provider     │
                   │                   │              │  Endpoint      │
                   └───────────────────┘              └────────────────┘
```


## Practical Implementation

### API Models
- **Simplified for usability**: The client-facing API can be simpler and more intuitive than Traefik's format
- **Example**: Offering a `url` field directly on a Service for simple cases, rather than requiring a full LoadBalancer definition
- **Batch operations**: Supporting composite services to create related components in one call

### Provider Endpoint
- The `/api/traefik/provider` endpoint must return exactly what Traefik expects
- No simplifications or deviations from Traefik's schema are possible here
- Must map our internal representation to Traefik's expected format

## Package Structure

to be filled
