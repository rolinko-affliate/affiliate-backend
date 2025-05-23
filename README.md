# Affiliate Platform Backend

This is the backend service for the Affiliate Platform MVP. It manages all application-specific data and business logic, while leveraging Supabase solely for user authentication (JWT issuance and validation).

## Project Structure

```
affiliate-backend/
├── cmd/
│   ├── api/                // Main application entry point
│   │   └── main.go
│   └── migrate/            // Database migration tool
│       └── main.go
├── internal/               // Private application and library code
│   ├── api/                // API handlers, middleware, routing
│   │   ├── handlers/       // HTTP request handlers for all resources
│   │   ├── middleware/     // Authentication and authorization middleware
│   │   └── router.go       // API route definitions
│   ├── auth/               // Authentication (JWT validation) & Authorization (RBAC)
│   ├── config/             // Configuration loading and environment variables
│   ├── domain/             // Core business logic entities/models (structs)
│   ├── platform/           // External platform integrations
│   │   ├── crypto/         // Encryption/decryption utilities
│   │   └── everflow/       // Everflow API client integration
│   ├── repository/         // Data access layer (database interactions)
│   └── service/            // Business logic services
├── migrations/             // Database migration files
├── go.mod
├── go.sum
└── .env                    // Environment variables (DO NOT COMMIT if it contains secrets)
```

## Security Features

### Authentication
- JWT-based authentication using Supabase Auth
- Token validation and verification
- Secure session management

### Authorization
- Role-Based Access Control (RBAC) system
- Organization-based access control
- Permission checks at the handler level
- Admin users have full access to all resources
- Non-admin users can only access resources within their organization

### Data Protection
- Encryption of sensitive data using AES-256
- Secure storage of API credentials
- Input validation and sanitization

## Key Features

- **Multi-tenancy**: Organizations are isolated from each other
- **Role-based permissions**: Different access levels based on user roles
- **Affiliate Management**: Create and manage affiliate accounts
- **Advertiser Management**: Create and manage advertiser accounts
- **Campaign Management**: Create and manage advertising campaigns
- **Provider Integration**: Integration with external affiliate networks (Everflow)
- **RESTful API**: Clean and consistent API design
- **Database Migrations**: Versioned database schema changes
- **Swagger Documentation**: Auto-generated API documentation

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database

### Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and update the values
   - Set `ENVIRONMENT=development` for local development (enables CORS)
   - Set `ENVIRONMENT=production` for production deployment (disables CORS)
3. Install dependencies:
   ```
   make deps
   ```
   
   If you encounter issues with the migration packages, install them manually:
   ```
   go get -u github.com/golang-migrate/migrate/v4
   go get -u github.com/golang-migrate/migrate/v4/database/postgres
   go get -u github.com/golang-migrate/migrate/v4/source/file
   ```
   
4. Run database migrations:
   ```
   make migrate-up
   ```
5. Build and run the application:
   ```
   make build
   make run
   ```
   
   Or run with auto-migration:
   ```
   make run-with-migrate
   ```

## Database Migrations

The application uses a migration system to manage database schema changes. The migration tool is located in the `cmd/migrate` directory and can be used to apply, rollback, and check the status of migrations.

### Migration Commands

- Apply all pending migrations:
  ```
  make migrate-up
  ```

- Rollback the most recent migration:
  ```
  make migrate-down
  ```

- Rollback all migrations:
  ```
  make migrate-reset
  ```

- Check if migrations are up to date:
  ```
  make migrate-check
  ```

- Show current database version:
  ```
  make migrate-version
  ```

- Show detailed migration status:
  ```
  make migrate-status
  ```

- Create a new migration file:
  ```
  make migrate-create NAME=add_new_table
  ```

### Migration Files

Migration files are stored in the `migrations` directory and follow the naming convention `000001_name.up.sql` and `000001_name.down.sql`. The `up.sql` file contains the SQL statements to apply the migration, and the `down.sql` file contains the SQL statements to rollback the migration.

### Auto-Migration

The API server can automatically apply pending migrations on startup by using the `--auto-migrate` flag:

```
./affiliate-backend --auto-migrate
```

## CORS Configuration

The API server includes CORS (Cross-Origin Resource Sharing) configuration that behaves differently based on the environment:

- In **development** mode (`ENVIRONMENT=development`), CORS is enabled and allows requests from any origin.
- In **production** mode (`ENVIRONMENT=production`), CORS is disabled, restricting cross-origin requests.

This configuration helps secure the API in production while allowing for easier development and testing in local environments.

To set the environment:
- In `.env` file: Set `ENVIRONMENT=development` or `ENVIRONMENT=production`
- In Docker: The `docker-compose.yml` sets `ENVIRONMENT=development` for local development
- In production: The `Dockerfile` sets `ENVIRONMENT=production` by default

## API Endpoints

### Public Endpoints

- `POST /api/v1/public/webhooks/supabase/new-user`: Webhook for Supabase auth to create a profile when a new user signs up

### Authenticated Endpoints

- `GET /api/v1/users/me`: Get the current user's profile

## API Documentation

The API is documented using OpenAPI 3.0 specification. You can generate and view the documentation using the following commands:

### Generate OpenAPI Specification

```bash
# Generate OpenAPI docs in docs/swagger directory
make openapi

# Export OpenAPI spec to JSON file (openapi.json)
make openapi-json

# Export OpenAPI spec to YAML file (openapi.yaml)
make openapi-yaml
```

### View API Documentation

```bash
# Serve Swagger UI documentation on http://localhost:8090
make serve-docs
```

### Integrate Swagger UI with the Application

To add Swagger UI directly to your application:

1. Install the required dependencies:
   ```bash
   make install-swagger-ui
   ```

2. Add the following code to your router setup in `internal/api/router.go`:
   ```go
   import (
       // ... existing imports
       swaggerFiles "github.com/swaggo/files"
       ginSwagger "github.com/swaggo/gin-swagger"
       _ "github.com/affiliate-backend/docs/swagger" // Import generated docs
   )

   // In your SetupRouter function:
   // Swagger documentation
   r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
   ```

3. Rebuild and run your application, then access Swagger UI at `/swagger/index.html`

## Authentication

This service uses Supabase JWT tokens for authentication. Include the JWT token in the `Authorization` header as a Bearer token:

```
Authorization: Bearer <token>
```

## Authorization

The service uses Role-Based Access Control (RBAC) for authorization. The following roles are available:

- `Admin`: System administrator with full access
- `AdvertiserManager`: Manages advertisers and campaigns
- `AffiliateManager`: Manages affiliates
- `Affiliate`: Regular affiliate user

## Environment Variables

- `PORT`: The port to run the server on (default: 8080)
- `DATABASE_URL`: PostgreSQL connection string
- `SUPABASE_JWT_SECRET`: Supabase JWT secret for validating tokens
- `ENCRYPTION_KEY`: 32-byte base64 encoded AES key for encrypting sensitive data

## License

[MIT](LICENSE)