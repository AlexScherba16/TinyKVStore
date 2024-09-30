package zap

import (
	"TinyKVStore/internal/config"
	"TinyKVStore/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	defaultLogLevel = zapcore.DebugLevel
)

// zapLoggerCfg private zap logger config, use it after validation config.LoggerConfig.
type zapLoggerCfg struct {
	level zapcore.Level
}

// validateConfig ensures that the correct settings will be used for the logger configuration.
func validateConfig(loggerCfg *config.LoggerConfig) zapLoggerCfg {
	level, err := zapcore.ParseLevel(loggerCfg.Level)
	if err != nil {
		level = defaultLogLevel
	}

	return zapLoggerCfg{
		level: level,
	}
}

// zapLogger private zap.Logger wrapper logger abstraction,
// concrete implementation of ILogger interface.
type zapLogger struct {
	logger *zap.Logger
}

// NewLogger creates a new zapLogger instance with color output configuration and log level settings.
// Returns the ILogger interface for decoupling from the concrete logging implementation,
// or error if something went wrong.
func NewLogger(loggerCfg *config.LoggerConfig) (logger.ILogger, error) {
	validatedZapCfg := validateConfig(loggerCfg)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), validatedZapCfg.level)
	return &zapLogger{
		logger: zap.New(consoleCore),
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

// Debug logs error messages at the Debug level.
func (log *zapLogger) Debug(msg string) {
	log.logger.Debug(msg)
}
