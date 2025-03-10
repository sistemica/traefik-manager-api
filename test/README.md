# Traefik Manager Tests

This directory contains tests for the Traefik Manager application. The tests are organized into unit tests (located in the package directories) and integration tests (in the `integration` subdirectory).

## Test Structure

```
test/
├── integration/       # Integration tests that verify the entire system
│   └── integration_test.go    # Full API integration tests
└── README.md          # This file
```

Unit tests are located alongside the code they test in their respective packages:
- `internal/store/store_test.go` - Tests for the data store
- `internal/config/config_test.go` - Tests for configuration loading
- `internal/middleware/auth_middleware_test.go` - Tests for authentication middleware
- `internal/api/handlers/middleware_handler_test.go` - Tests for API handlers

## Running Tests

### Run All Tests

```bash
# Run all tests
go test ./...
```

### Run Unit Tests Only

```bash
# Run unit tests for a specific package
go test ./internal/store
go test ./internal/config
go test ./internal/middleware
go test ./internal/api/handlers
```

### Run Integration Tests

```bash
# Run integration tests only
go test ./test/integration

# Run with verbose output
go test -v ./test/integration
```

### Run Tests with Coverage

```bash
# Run tests with coverage report
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Integration Tests

The integration tests start a real HTTP server with a temporary file store and run tests against it. These tests verify that:

1. The API endpoints work correctly
2. Data is properly persisted
3. Dependency checking between resources works
4. The Traefik provider endpoint returns the correct configuration

Integration tests are more resource-intensive and slower than unit tests, but they provide valuable verification that the system works as a whole.

### CI/CD Configuration

In CI/CD environments, you can control whether integration tests run using environment variables:

```bash
# Skip integration tests in CI
export CI=true

# Force running integration tests in CI
export CI=true
export RUN_INTEGRATION_TESTS=true
```

## Writing New Tests

### Unit Tests

Add new unit tests in the same package as the code they test with the `_test.go` suffix.

```go
// Example test in internal/store/store_test.go
func TestNewFeature(t *testing.T) {
    // Test setup
    store, err := NewFileStore(tmpPath)
    if err != nil {
        t.Fatalf("Failed to create store: %v", err)
    }
    
    // Test assertions
    result, err := store.SomeOperation()
    if err != nil {
        t.Fatalf("Operation failed: %v", err)
    }
    if result != expectedResult {
        t.Errorf("Expected %v, got %v", expectedResult, result)
    }
}
```

### Integration Tests

Add new integration tests to `test/integration/integration_test.go` or create additional test files in that directory.

```go
// Example integration test
t.Run("New Feature E2E", func(t *testing.T) {
    // Make API requests to test the feature
    resp, err := client.Get(fmt.Sprintf("%s/new-feature", baseURL))
    if err != nil {
        t.Fatalf("Request failed: %v", err)
    }
    defer resp.Body.Close()
    
    // Assert expected results
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
    }
    
    // Verify the response body
    var result SomeType
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }
    
    if result.SomeField != expectedValue {
        t.Errorf("Expected %v, got %v", expectedValue, result.SomeField)
    }
})
```

## Test Best Practices

1. **Use table-driven tests** for testing multiple inputs with the same logic
2. **Clean up after tests** to avoid interference between tests
3. **Use descriptive test names** to make failures easier to understand
4. **Test error cases** in addition to the happy path
5. **Use subtests** (`t.Run()`) to organize related tests
6. **Avoid global state** that could cause tests to interfere with each other