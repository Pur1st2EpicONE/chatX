// Package logger provides a structured logging interface and implementation
// for the application, including support for different severity levels
// and integration with goose migrations.
package logger

import (
	"chatX/internal/config"
	"chatX/internal/logger/slog"
	"io"
	"log"
	"os"

	"github.com/pressly/goose/v3"
)

// Logger defines the interface for structured logging with different severity levels.
type Logger interface {
	LogFatal(msg string, err error, args ...any) // LogFatal logs a fatal message with an error and optional key-value arguments.
	LogError(string, error, ...any)              // LogError logs an error message with an error and optional key-value arguments.
	LogWarn(msg string, args ...any)             // LogWarn logs a warning message with an error and optional key-value arguments.
	LogInfo(msg string, args ...any)             // LogInfo logs an informational message with optional key-value arguments.
	Debug(msg string, args ...any)               // Debug logs a debug message with optional key-value arguments.
}

// NewLogger creates a new structured Logger instance based on the provided configuration.
//
// It also sets the logger for goose migrations to discard output to prevent clutter.
func NewLogger(config config.Logger) (Logger, *os.File) {
	goose.SetLogger(log.New(io.Discard, "", 0))
	return slog.NewLogger(config)
}
