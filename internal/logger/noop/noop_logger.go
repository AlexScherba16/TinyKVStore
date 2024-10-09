package noop

import (
	"TinyKVStore/internal/logger"
)

// noopLogger private wrapper logger abstraction,
// concrete implementation of ILogger interface.
type noopLogger struct{}

func NewLogger() (logger.ILogger, error) {
	return &noopLogger{}, nil
}

// Info logs informational messages at the Info level.
func (log *noopLogger) Info(msg string) {}

// Warn logs warning messages at the Warn level.
func (log *noopLogger) Warn(msg string) {}

// Error logs error messages at the Error level.
func (log *noopLogger) Error(msg string) {}

// Debug logs error messages at the Debug level.
func (log *noopLogger) Debug(msg string) {}
