# Configuration Module

This module handles application configuration loading and management. It provides a centralized way to access configuration values from environment variables and default settings.

## Key Components

### Config Struct

The `Config` struct defines the application configuration parameters:

```go
type Config struct {
    Port              string
    DatabaseURL       string
    SupabaseJWTSecret string
    EncryptionKey     string
    Environment       string
}
```

### AppConfig

The `AppConfig` variable is a global instance of the `Config` struct that holds the loaded configuration:

```go
var AppConfig Config
```

### LoadConfig

The `LoadConfig` function loads configuration values from environment variables:

```go
func LoadConfig() {
    // Load environment variables from .env file
    godotenv.Load()
    
    // Set configuration values from environment variables
    AppConfig.Port = getEnv("PORT", "8080")
    AppConfig.DatabaseURL = getEnv("DATABASE_URL", "")
    AppConfig.SupabaseJWTSecret = getEnv("SUPABASE_JWT_SECRET", "")
    AppConfig.EncryptionKey = getEnv("ENCRYPTION_KEY", "")
    AppConfig.Environment = getEnv("ENVIRONMENT", "development")
}
```

## Environment Variables

The module handles the following environment variables:

- `PORT`: The port on which the server listens (default: "8080")
- `DATABASE_URL`: PostgreSQL connection string
- `SUPABASE_JWT_SECRET`: Secret key for validating Supabase JWT tokens
- `ENCRYPTION_KEY`: Key for encrypting sensitive data
- `ENVIRONMENT`: Application environment (development/production)

## Helper Functions

The module includes helper functions for working with environment variables:

- `getEnv`: Gets an environment variable with a default value
- `getEnvBool`: Gets a boolean environment variable with a default value
- `getEnvInt`: Gets an integer environment variable with a default value

## Usage

The configuration module is used throughout the application to access configuration values:

```go
// Initialize database connection
db, err := repository.InitDBConnection(config.AppConfig.DatabaseURL)

// Check environment for CORS configuration
if config.AppConfig.Environment == "development" {
    // Enable CORS for development
}

// Use JWT secret for token validation
token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
    return []byte(config.AppConfig.SupabaseJWTSecret), nil
})
```

## Environment-Specific Behavior

The configuration module supports different behavior based on the environment:

- **Development**: Enables CORS, detailed error messages, etc.
- **Production**: Disables CORS, hides detailed error messages, etc.

## Security Considerations

- Sensitive configuration values (JWT secret, encryption key) are never logged
- Environment variables can be loaded from a .env file (not committed to version control)
- Default values are provided for non-critical configuration