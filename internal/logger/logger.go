package logger

// ILogger general logger interface.
//
// Methods:
//
//   - Info(msg string): logs a message at Info level.
//
//   - Warn(msg string): logs a message at warning level.
//
//   - Error(msg string): logs a message at error level.
//
//   - Debug(msg string): logs a message at debug level.
type ILogger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}
