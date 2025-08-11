package config

import (
	"log/slog"
	"os"
	"strings"
)

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

// LogFormat represents the logging format
type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  LogLevel  `json:"level"`
	Format LogFormat `json:"format"`
}

// GetLoggingConfig returns logging configuration from environment variables
func GetLoggingConfig() LoggingConfig {
	config := LoggingConfig{
		Level:  LogLevelInfo, // Default level
		Format: LogFormatText, // Default format
	}

	// Parse log level from environment
	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		switch strings.ToUpper(levelStr) {
		case "DEBUG":
			config.Level = LogLevelDebug
		case "INFO":
			config.Level = LogLevelInfo
		case "WARN", "WARNING":
			config.Level = LogLevelWarn
		case "ERROR":
			config.Level = LogLevelError
		}
	}

	// Parse log format from environment
	if formatStr := os.Getenv("LOG_FORMAT"); formatStr != "" {
		switch strings.ToLower(formatStr) {
		case "json":
			config.Format = LogFormatJSON
		case "text":
			config.Format = LogFormatText
		}
	}

	return config
}

// ToSlogLevel converts LogLevel to slog.Level
func (l LogLevel) ToSlogLevel() slog.Level {
	switch l {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// ConfigureLogger sets up the global logger with the specified configuration
func ConfigureLogger(config LoggingConfig) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: config.Level.ToSlogLevel(),
	}

	switch config.Format {
	case LogFormatJSON:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case LogFormatText:
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// SanitizeForLogging removes or masks sensitive information from log values
func SanitizeForLogging(key string, value interface{}) interface{} {
	keyLower := strings.ToLower(key)
	
	// List of sensitive field patterns
	sensitivePatterns := []string{
		"password", "passwd", "pwd",
		"token", "key", "secret",
		"auth", "credential", "cred",
		"api_key", "apikey",
		"email", "mail", // Email addresses can be sensitive
		"phone", "mobile",
		"ssn", "social",
		"credit", "card", "payment",
	}

	for _, pattern := range sensitivePatterns {
		if strings.Contains(keyLower, pattern) {
			if str, ok := value.(string); ok && str != "" {
				// Mask the value, showing only first and last characters for longer strings
				if len(str) <= 4 {
					return "***"
				}
				return str[:1] + "***" + str[len(str)-1:]
			}
			return "***"
		}
	}

	return value
}

// LogWithSanitization logs with automatic sanitization of sensitive fields
func LogWithSanitization(level slog.Level, msg string, args ...interface{}) {
	if len(args)%2 != 0 {
		// If odd number of args, just log normally
		slog.Log(nil, level, msg, args...)
		return
	}

	// Sanitize key-value pairs
	sanitizedArgs := make([]interface{}, len(args))
	for i := 0; i < len(args); i += 2 {
		key := args[i]
		value := args[i+1]
		
		sanitizedArgs[i] = key
		if keyStr, ok := key.(string); ok {
			sanitizedArgs[i+1] = SanitizeForLogging(keyStr, value)
		} else {
			sanitizedArgs[i+1] = value
		}
	}

	slog.Log(nil, level, msg, sanitizedArgs...)
}

// Convenience functions for sanitized logging
func DebugSanitized(msg string, args ...interface{}) {
	LogWithSanitization(slog.LevelDebug, msg, args...)
}

func InfoSanitized(msg string, args ...interface{}) {
	LogWithSanitization(slog.LevelInfo, msg, args...)
}

func WarnSanitized(msg string, args ...interface{}) {
	LogWithSanitization(slog.LevelWarn, msg, args...)
}

func ErrorSanitized(msg string, args ...interface{}) {
	LogWithSanitization(slog.LevelError, msg, args...)
}