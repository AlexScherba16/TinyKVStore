package zap

import (
	"TinyKVStore/internal/logger"
	stdout "github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger private zap.Logger wrapper logger abstraction,
// concrete implementation of ILogger interface.
type zapLogger struct {
	logger *zap.Logger
}

// NewLogger creates a new zapLogger instance with color output configuration and log level settings.
// Returns the ILogger interface for decoupling from the concrete logging implementation,
// or error if something went wrong.
func NewLogger() (logger.ILogger, error) {
	colorOutput := zapcore.AddSync(stdout.NewColorableStdout())
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		colorOutput,
		zapcore.InfoLevel,
	)

	return &zapLogger{
		logger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)),
	}, nil
}

// Info logs informational messages at the Info level.
func (log *zapLogger) Info(msg string) {
	log.logger.Info(msg)
}

// Warn logs warning messages at the Warn level.
func (log *zapLogger) Warn(msg string) {
	log.logger.Warn(msg)
}

// Error logs error messages at the Error level.
func (log *zapLogger) Error(msg string) {
	log.logger.Error(msg)
}
