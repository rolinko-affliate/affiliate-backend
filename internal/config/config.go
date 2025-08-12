package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`

	DatabaseHost     string `mapstructure:"DATABASE_HOST"`
	DatabasePort     string `mapstructure:"DATABASE_PORT"`
	DatabaseName     string `mapstructure:"DATABASE_NAME"`
	DatabaseUser     string `mapstructure:"DATABASE_USER"`
	DatabasePassword string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseSSLMode  string `mapstructure:"DATABASE_SSL_MODE"` // e.g. "disable", "require", "verify-ca", "verify-full"

	SupabaseJWTSecret string `mapstructure:"SUPABASE_JWT_SECRET"`
	// Key for encrypting/decrypting sensitive data like Everflow API keys
	EncryptionKey string `mapstructure:"ENCRYPTION_KEY"` // 32-byte AES key, base64 encoded
	Environment   string `mapstructure:"ENVIRONMENT"`    // "development" or "production"
	DebugMode     bool   `mapstructure:"DEBUG_MODE"`     // Enable debug logging for API requests/responses
	MockMode      bool   `mapstructure:"MOCK_MODE"`      // Enable mock integration service instead of real provider
	
	// Logging configuration
	LogLevel     string `mapstructure:"LOG_LEVEL"`      // "DEBUG", "INFO", "WARN", "ERROR"
	LogFormat    string `mapstructure:"LOG_FORMAT"`     // "json" or "text"
	LogOutput    string `mapstructure:"LOG_OUTPUT"`     // "stdout", "stderr", or file path
	LogAddSource bool   `mapstructure:"LOG_ADD_SOURCE"` // Add source file and line number
	
	// Everflow API configuration
	EverflowAPIKey string `mapstructure:"EVERFLOW_API_KEY"` // Everflow API key for authentication
}

var AppConfig Config

func LoadConfig() {
	viper.AddConfigPath(".")    // Look for config in current directory
	viper.SetConfigName(".env") // Name of config file (without extension)
	viper.SetConfigType("env")  // Type of config file

	// Set defaults (optional)
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENVIRONMENT", "production")
	viper.SetDefault("DEBUG_MODE", false)
	viper.SetDefault("MOCK_MODE", false)

	viper.SetDefault("DATABASE_SSL_MODE", "disable") // Default to no SSL for local dev
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_PORT", "5432")
	viper.SetDefault("DATABASE_NAME", "myapp")
	viper.SetDefault("DATABASE_USER", "postgres")
	viper.SetDefault("DATABASE_PASSWORD", "password") // Default password, should be overridden

	viper.SetDefault("SUPABASE_JWT_SECRET", "") // Default password, should be overridden
	viper.SetDefault("ENCRYPTION_KEY", "")      // Default password, should be overridden
	viper.SetDefault("EVERFLOW_API_KEY", "")    // Default password, should be overridden
	viper.SetDefault("MockMode", false)

	// Logging defaults
	viper.SetDefault("LOG_LEVEL", "INFO")
	viper.SetDefault("LOG_FORMAT", "text")
	viper.SetDefault("LOG_OUTPUT", "stdout")
	viper.SetDefault("LOG_ADD_SOURCE", false)

	viper.AutomaticEnv() // Read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file (.env) not found, relying on environment variables.")
		} else {
			log.Fatalf("Error reading config file: %s", err)
		}
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	// Log configuration without sensitive data
	logSafeConfig()

	// Basic validation
	if AppConfig.DatabaseURL == "" {
		log.Println("DATABASE_URL not set, constructing from individual components")
		// if DATABASE_URL is not set, construct it from individual components likepostgres://postgres:postgres@localhost:5432/affiliate_platform?sslmode=disable
		AppConfig.DatabaseURL = "postgres://" +
			AppConfig.DatabaseUser + ":" +
			AppConfig.DatabasePassword + "@" +
			AppConfig.DatabaseHost + ":" +
			AppConfig.DatabasePort + "/" +
			AppConfig.DatabaseName + "?sslmode=" +
			AppConfig.DatabaseSSLMode
	}
	if AppConfig.SupabaseJWTSecret == "" {
		log.Fatal("SUPABASE_JWT_SECRET must be set")
	}
	if AppConfig.EncryptionKey == "" {
		log.Fatal("ENCRYPTION_KEY must be set for securing provider credentials")
	}
	if AppConfig.EverflowAPIKey == "" {
		log.Println("EVERFLOW_API_KEY not set, some features may be limited")
	}
}

// IsDebugMode returns true if debug mode is enabled (global function)
func IsDebugMode() bool {
	return AppConfig.DebugMode
}

// IsDevelopment returns true if running in development environment (global function)
func IsDevelopment() bool {
	return AppConfig.Environment == "development"
}

// IsDebugMode returns true if debug mode is enabled (method on Config)
func (c *Config) IsDebugMode() bool {
	return c.DebugMode
}

// IsDevelopment returns true if running in development environment (method on Config)
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsMockMode returns true if mock mode is enabled (global function)
func IsMockMode() bool {
	return AppConfig.MockMode
}

// IsMockMode returns true if mock mode is enabled (method on Config)
func (c *Config) IsMockMode() bool {
	return c.MockMode
}

// logSafeConfig logs configuration without sensitive data
func logSafeConfig() {
	safeConfig := struct {
		Port             string `json:"port"`
		DatabaseHost     string `json:"database_host"`
		DatabasePort     string `json:"database_port"`
		DatabaseName     string `json:"database_name"`
		DatabaseUser     string `json:"database_user"`
		DatabaseSSLMode  string `json:"database_ssl_mode"`
		Environment      string `json:"environment"`
		DebugMode        bool   `json:"debug_mode"`
		MockMode         bool   `json:"mock_mode"`
		LogLevel         string `json:"log_level"`
		LogFormat        string `json:"log_format"`
		LogOutput        string `json:"log_output"`
		LogAddSource     bool   `json:"log_add_source"`
		HasJWTSecret     bool   `json:"has_jwt_secret"`
		HasEncryptionKey bool   `json:"has_encryption_key"`
		HasEverflowKey   bool   `json:"has_everflow_key"`
	}{
		Port:             AppConfig.Port,
		DatabaseHost:     AppConfig.DatabaseHost,
		DatabasePort:     AppConfig.DatabasePort,
		DatabaseName:     AppConfig.DatabaseName,
		DatabaseUser:     AppConfig.DatabaseUser,
		DatabaseSSLMode:  AppConfig.DatabaseSSLMode,
		Environment:      AppConfig.Environment,
		DebugMode:        AppConfig.DebugMode,
		MockMode:         AppConfig.MockMode,
		LogLevel:         AppConfig.LogLevel,
		LogFormat:        AppConfig.LogFormat,
		LogOutput:        AppConfig.LogOutput,
		LogAddSource:     AppConfig.LogAddSource,
		HasJWTSecret:     AppConfig.SupabaseJWTSecret != "",
		HasEncryptionKey: AppConfig.EncryptionKey != "",
		HasEverflowKey:   AppConfig.EverflowAPIKey != "",
	}
	
	log.Printf("Loaded configuration: %+v", safeConfig)
}

// GetLoggerConfig returns logger configuration from app config
func GetLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:     AppConfig.LogLevel,
		Format:    AppConfig.LogFormat,
		Output:    AppConfig.LogOutput,
		AddSource: AppConfig.LogAddSource,
	}
}

// LoggerConfig represents logger configuration
type LoggerConfig struct {
	Level     string
	Format    string
	Output    string
	AddSource bool
}
