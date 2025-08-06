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
	viper.SetDefault("MockMode", false)

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

	log.Printf("Loaded configuration: %+v", AppConfig)

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
