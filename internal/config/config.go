package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port              string `mapstructure:"PORT"`
	DatabaseURL       string `mapstructure:"DATABASE_URL"`
	SupabaseJWTSecret string `mapstructure:"SUPABASE_JWT_SECRET"`
	// Key for encrypting/decrypting sensitive data like Everflow API keys
	EncryptionKey string `mapstructure:"ENCRYPTION_KEY"` // 32-byte AES key, base64 encoded
}

var AppConfig Config

func LoadConfig() {
	viper.AddConfigPath(".")      // Look for config in current directory
	viper.SetConfigName(".env")   // Name of config file (without extension)
	viper.SetConfigType("env")    // Type of config file
	viper.AutomaticEnv()          // Read in environment variables that match

	// Set defaults (optional)
	viper.SetDefault("PORT", "8080")

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

	// Basic validation
	if AppConfig.DatabaseURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}
	if AppConfig.SupabaseJWTSecret == "" {
		log.Fatal("SUPABASE_JWT_SECRET must be set")
	}
	if AppConfig.EncryptionKey == "" {
		log.Fatal("ENCRYPTION_KEY must be set for securing provider credentials")
	}
}