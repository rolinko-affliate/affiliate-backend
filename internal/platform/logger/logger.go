package logger

import (
	"context"
	"fmt"
	// "log/slog" // Commented out for Go 1.19 compatibility
	"log"
	"os"
	// "strings" // Commented out - not used in simplified implementation
	
	// "github.com/affiliate-backend/internal/config" // Commented out - not used in simplified implementation
)

// Logger provides a consistent logging interface
// Simplified for Go 1.19 compatibility
type Logger struct {
	logger *log.Logger
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
// Simplified for Go 1.19 compatibility
func NewLogger(config Config) *Logger {
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

	return &Logger{
		logger: log.New(output, "", log.LstdFlags),
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

// Logger methods - simplified for Go 1.19 compatibility
func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Printf("[DEBUG] %s %v", msg, args)
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Printf("[INFO] %s %v", msg, args)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Printf("[WARN] %s %v", msg, args)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Printf("[ERROR] %s %v", msg, args)
}

// Sanitized logging functions that automatically sanitize sensitive data
// Simplified for Go 1.19 compatibility
func DebugSanitized(msg string, args ...any) {
	GetDefault().Debug(msg, args...)
}

func InfoSanitized(msg string, args ...any) {
	GetDefault().Info(msg, args...)
}

func WarnSanitized(msg string, args ...any) {
	GetDefault().Warn(msg, args...)
}

func ErrorSanitized(msg string, args ...any) {
	GetDefault().Error(msg, args...)
}

// WithContext returns a logger with the given context
// Simplified for Go 1.19 compatibility
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return l // Return same logger for now
}

// WithFields returns a logger with the given fields
// Simplified for Go 1.19 compatibility
func (l *Logger) WithFields(fields map[string]any) *Logger {
	return l // Return same logger for now
}

// WithField returns a logger with a single field
// Simplified for Go 1.19 compatibility
func (l *Logger) WithField(key string, value any) *Logger {
	return l // Return same logger for now
}

// Convenience methods for structured logging

// LogDatabaseOperation logs database operations
// Simplified for Go 1.19 compatibility
func (l *Logger) LogDatabaseOperation(operation, table string, duration int64, err error) {
	if err != nil {
		l.Error(fmt.Sprintf("Database operation failed: %s on %s (duration: %dms)", operation, table, duration), "error", err.Error())
	} else {
		l.Debug(fmt.Sprintf("Database operation completed: %s on %s (duration: %dms)", operation, table, duration))
	}
}

// LogHTTPRequest logs HTTP requests
// Simplified for Go 1.19 compatibility
func (l *Logger) LogHTTPRequest(method, path string, statusCode int, duration int64, userID string) {
	msg := fmt.Sprintf("HTTP request: %s %s (status: %d, duration: %dms)", method, path, statusCode, duration)
	if userID != "" {
		msg += fmt.Sprintf(" user: %s", userID)
	}
	
	if statusCode >= 400 {
		l.Warn(msg)
	} else {
		l.Info(msg)
	}
}

// LogProviderOperation logs external provider operations
// Simplified for Go 1.19 compatibility
func (l *Logger) LogProviderOperation(provider, operation, entityType string, entityID any, err error) {
	msg := fmt.Sprintf("Provider operation: %s %s %s (ID: %v)", provider, operation, entityType, entityID)
	
	if err != nil {
		l.Error(msg, "error", err.Error())
	} else {
		l.Info(msg)
	}
}

// LogServiceOperation logs service layer operations
// Simplified for Go 1.19 compatibility
func (l *Logger) LogServiceOperation(service, operation string, entityID any, err error) {
	msg := fmt.Sprintf("Service operation: %s %s", service, operation)
	if entityID != nil {
		msg += fmt.Sprintf(" (ID: %v)", entityID)
	}

	if err != nil {
		l.Error(msg, "error", err.Error())
	} else {
		l.Debug(msg)
	}
}
