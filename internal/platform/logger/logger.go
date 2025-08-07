package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// Logger wraps slog.Logger to provide a consistent logging interface
type Logger struct {
	*slog.Logger
}

// LogLevel represents the logging level
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

var (
	// Default logger instance
	defaultLogger *Logger
)

// Config holds logger configuration
type Config struct {
	Level     LogLevel `mapstructure:"LOG_LEVEL"`
	Format    string   `mapstructure:"LOG_FORMAT"`     // "json" or "text"
	Output    string   `mapstructure:"LOG_OUTPUT"`     // "stdout", "stderr", or file path
	AddSource bool     `mapstructure:"LOG_ADD_SOURCE"` // Add source file and line number
}

// DefaultConfig returns default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:     LevelInfo,
		Format:    "text",
		Output:    "stdout",
		AddSource: false,
	}
}

// NewLogger creates a new logger with the given configuration
func NewLogger(config Config) *Logger {
	var level slog.Level
	switch strings.ToUpper(string(config.Level)) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var output *os.File
	switch config.Output {
	case "stderr":
		output = os.Stderr
	case "stdout", "":
		output = os.Stdout
	default:
		// For file output, we would open the file here
		// For now, default to stdout
		output = os.Stdout
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: config.AddSource,
	}

	switch config.Format {
	case "json":
		handler = slog.NewJSONHandler(output, opts)
	default:
		handler = slog.NewTextHandler(output, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// InitDefault initializes the default logger
func InitDefault(config Config) {
	defaultLogger = NewLogger(config)
}

// GetDefault returns the default logger instance
func GetDefault() *Logger {
	if defaultLogger == nil {
		defaultLogger = NewLogger(DefaultConfig())
	}
	return defaultLogger
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	GetDefault().Debug(msg, args...)
}

// Info logs an info message
func Info(msg string, args ...any) {
	GetDefault().Info(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	GetDefault().Warn(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...any) {
	GetDefault().Error(msg, args...)
}

// Fatal logs an error message and exits the program
func Fatal(msg string, args ...any) {
	GetDefault().Error(msg, args...)
	os.Exit(1)
}

// WithContext returns a logger with the given context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		Logger: l.Logger.With(),
	}
}

// WithFields returns a logger with the given fields
func (l *Logger) WithFields(fields map[string]any) *Logger {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		Logger: l.Logger.With(args...),
	}
}

// WithField returns a logger with a single field
func (l *Logger) WithField(key string, value any) *Logger {
	return &Logger{
		Logger: l.Logger.With(key, value),
	}
}

// Convenience methods for structured logging

// LogDatabaseOperation logs database operations
func (l *Logger) LogDatabaseOperation(operation, table string, duration int64, err error) {
	fields := map[string]any{
		"operation": operation,
		"table":     table,
		"duration":  duration,
	}

	if err != nil {
		fields["error"] = err.Error()
		l.WithFields(fields).Error("Database operation failed")
	} else {
		l.WithFields(fields).Debug("Database operation completed")
	}
}

// LogHTTPRequest logs HTTP requests
func (l *Logger) LogHTTPRequest(method, path string, statusCode int, duration int64, userID string) {
	fields := map[string]any{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"duration":    duration,
	}

	if userID != "" {
		fields["user_id"] = userID
	}

	if statusCode >= 400 {
		l.WithFields(fields).Warn("HTTP request completed with error")
	} else {
		l.WithFields(fields).Info("HTTP request completed")
	}
}

// LogProviderOperation logs external provider operations
func (l *Logger) LogProviderOperation(provider, operation, entityType string, entityID any, err error) {
	fields := map[string]any{
		"provider":    provider,
		"operation":   operation,
		"entity_type": entityType,
		"entity_id":   entityID,
	}

	if err != nil {
		fields["error"] = err.Error()
		l.WithFields(fields).Error("Provider operation failed")
	} else {
		l.WithFields(fields).Info("Provider operation completed")
	}
}

// LogServiceOperation logs service layer operations
func (l *Logger) LogServiceOperation(service, operation string, entityID any, err error) {
	fields := map[string]any{
		"service":   service,
		"operation": operation,
	}

	if entityID != nil {
		fields["entity_id"] = entityID
	}

	if err != nil {
		fields["error"] = err.Error()
		l.WithFields(fields).Error("Service operation failed")
	} else {
		l.WithFields(fields).Debug("Service operation completed")
	}
}
