# Logger Package

This package provides structured logging for the Traefik Manager application using zerolog. It implements a high-performance, JSON-focused logging system with various output formats and log levels.

## Usage

### Initialize Logger

```go
import (
    "github.com/sistemica/traefik-manager/internal/logger"
)

func main() {
    // Configure and initialize logger
    err := logger.Setup(logger.Config{
        Level:      logger.InfoLevel,
        Format:     logger.JSONFormat,
        FilePath:   "",  // Empty for stdout
        TimeFormat: time.RFC3339,
        UseColor:   true,
    })
    if err != nil {
        panic(err)
    }
    
    // Now you can use the logger
    logger.Info().Str("service", "api").Msg("Application started")
}
```

### Logging Examples

```go
// Simple message
logger.Info().Msg("This is an info message")

// With fields
logger.Debug().
    Str("user", "john").
    Int("id", 123).
    Msg("User logged in")

// Error with stack trace
err := someFunction()
if err != nil {
    logger.Error().
        Err(err).
        Str("operation", "database query").
        Msg("Failed to fetch user data")
}

// Create a sub-logger with component context
httpLogger := logger.With().Str("component", "http").Logger()
httpLogger.Info().Msg("Server started")
```

## Log Levels

The logger supports the following log levels (in order of verbosity):

- `debug` - Debug information, verbose output
- `info` - Normal operational information
- `warn` - Warning conditions, not errors but might require attention
- `error` - Error conditions that should be addressed
- `fatal` - Fatal conditions that cause the application to exit
- `panic` - Panic conditions that cause the application to panic

## Output Formats

The logger supports two output formats:

- `json` - Structured JSON output, ideal for machine processing and log aggregation
- `text` - Human-readable text output, ideal for development environments

## Configuration

The logger can be configured with the following options:

- `Level` - The minimum log level to output (debug, info, warn, error, etc.)
- `Format` - The log format (json, text)
- `FilePath` - The path to the log file, empty for stdout
- `TimeFormat` - The format to use for timestamps
- `UseColor` - Enables colored output in text format

## Best Practices

1. **Use structured logging**: Add context to your logs with fields rather than embedding them in the message
   ```go
   // Good
   logger.Info().Str("user", "john").Int("attempts", 3).Msg("Login successful")
   
   // Avoid
   logger.Info().Msg(fmt.Sprintf("User john logged in successfully after 3 attempts"))
   ```

2. **Use appropriate log levels**: Reserve error logs for actual errors, use debug for detailed information
   ```go
   // Debug for detailed info
   logger.Debug().Str("query", query).Msg("Executing database query")
   
   // Info for normal operations
   logger.Info().Int("users", count).Msg("Users fetched successfully")
   
   // Warn for important but non-error conditions
   logger.Warn().Int("remaining", 2).Msg("Rate limit threshold approaching")
   
   // Error for error conditions
   logger.Error().Err(err).Msg("Database query failed")
   ```

3. **Add context to your logs**: Include relevant information to help diagnose issues
   ```go
   logger.Error().
       Err(err).
       Str("user", userID).
       Str("operation", "password reset").
       Str("requestID", requestID).
       Msg("Operation failed")
   ```