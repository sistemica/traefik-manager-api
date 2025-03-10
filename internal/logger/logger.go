package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Levels for logging
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	PanicLevel = "panic"
)

// Formats for logging output
const (
	JSONFormat = "json"
	TextFormat = "text"
)

// Config contains the logger configuration
type Config struct {
	// Level is the minimum log level to output (debug, info, warn, error, etc.)
	Level string
	// Format is the log format (json, text)
	Format string
	// FilePath is the path to the log file, empty for stdout
	FilePath string
	// TimeFormat is the format to use for timestamps
	TimeFormat string
	// UseColor enables colored output in text format
	UseColor bool
}

// Setup initializes the global logger with the given configuration
func Setup(cfg Config) error {
	// Set time format
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = time.RFC3339
	}
	zerolog.TimeFieldFormat = cfg.TimeFormat

	// Set default level
	if cfg.Level == "" {
		cfg.Level = InfoLevel
	}

	// Parse log level
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)

	// Configure output
	var output io.Writer = os.Stdout

	// Use log file if specified
	if cfg.FilePath != "" {
		file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		output = file
	}

	// Configure output format
	if cfg.Format == TextFormat {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: cfg.TimeFormat,
			NoColor:    !cfg.UseColor,
		}
	}

	// Set global logger
	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()

	// Log initialization
	log.Info().
		Str("level", cfg.Level).
		Str("format", cfg.Format).
		Str("output", outputDescription(cfg.FilePath)).
		Msg("Logger initialized")

	return nil
}

// Debug logs a debug message
func Debug() *zerolog.Event {
	return log.Debug()
}

// Info logs an info message
func Info() *zerolog.Event {
	return log.Info()
}

// Warn logs a warning message
func Warn() *zerolog.Event {
	return log.Warn()
}

// Error logs an error message
func Error() *zerolog.Event {
	return log.Error()
}

// Fatal logs a fatal message and exits
func Fatal() *zerolog.Event {
	return log.Fatal()
}

// Panic logs a panic message and panics
func Panic() *zerolog.Event {
	return log.Panic()
}

// With creates a sub-logger with the given component name
func With() zerolog.Context {
	return log.With()
}

// Logger returns the global logger
func Logger() zerolog.Logger {
	return log.Logger
}

// parseLevel parses a string level into a zerolog.Level
func parseLevel(level string) (zerolog.Level, error) {
	switch strings.ToLower(level) {
	case DebugLevel:
		return zerolog.DebugLevel, nil
	case InfoLevel:
		return zerolog.InfoLevel, nil
	case WarnLevel:
		return zerolog.WarnLevel, nil
	case ErrorLevel:
		return zerolog.ErrorLevel, nil
	case FatalLevel:
		return zerolog.FatalLevel, nil
	case PanicLevel:
		return zerolog.PanicLevel, nil
	default:
		return zerolog.InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
}

// outputDescription returns a description of the log output
func outputDescription(filePath string) string {
	if filePath == "" {
		return "stdout"
	}
	return filePath
}
