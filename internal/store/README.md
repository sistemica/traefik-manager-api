# Store Package

The store package provides the data persistence layer for Traefik Manager. It defines a common interface for storing and retrieving Traefik configuration components, with a file-based implementation.

## Store Interface

The `Store` interface defines methods for CRUD operations on all Traefik configuration components:
- Middlewares
- Routers
- Services

It also provides methods for checking dependencies between components and generating the dynamic configuration for Traefik.

## File-based Store Implementation

The `FileStore` implementation provides a thread-safe, file-based persistence mechanism with the following features:

### Thread-safe operations

All methods use a read-write mutex to ensure thread safety when accessing the store data.

### Efficient persistence

- Debounced saving to avoid excessive disk writes
- Structured JSON storage format
- Automatic loading of existing data on startup

### Reference integrity

- Checks for existence of referenced objects (e.g., middlewares referenced by routers)
- Prevents deletion of resources that are in use

### Composite service handling

- Proper creation of related components (router, service, middlewares)
- Consistent naming scheme for related components
- Handles dependencies correctly on updates and deletes

### Dynamic configuration provider

- Generates Traefik-compatible dynamic configuration
- Includes all routers, services, and middlewares

## Usage

Initialize the file store with a path to the storage file:

```go
store, err := store.NewFileStore("/path/to/storage/file.json")
if err != nil {
    log.Fatalf("Failed to create store: %v", err)
}
```

Inject this store into your handlers:

```go
middlewareHandler := handlers.NewMiddlewareHandler(store)
routerHandler := handlers.NewRouterHandler(store)
serviceHandler := handlers.NewServiceHandler(store)
```

The store automatically loads existing data if the file exists, and saves changes asynchronously to minimize performance impact.

## Error Handling

The store package defines several error types:
- `ErrNotFound`: Returned when a requested resource doesn't exist
- `ErrAlreadyExists`: Returned when trying to create a resource that already exists
- `ErrResourceInUse`: Returned when attempting to delete a resource that is referenced by other resources