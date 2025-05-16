# Affiliate Platform Backend

This is the backend service for the Affiliate Platform MVP. It manages all application-specific data and business logic, while leveraging Supabase solely for user authentication (JWT issuance and validation).

## Project Structure

```
affiliate-backend/
├── cmd/
│   └── api/                // Main application entry point
│       └── main.go
├── internal/               // Private application and library code
│   ├── api/                // API handlers, middleware, routing
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── router.go
│   ├── auth/               // Authentication (JWT validation) & Authorization (RBAC)
│   ├── config/             // Configuration loading
│   ├── domain/             // Core business logic entities/models (structs)
│   ├── platform/           // External platform integrations (e.g., everflow client)
│   ├── repository/         // Data access layer (database interactions)
│   └── service/            // Business logic services
├── migrations/             // Database migration files
├── go.mod
├── go.sum
└── .env                    // Environment variables (DO NOT COMMIT if it contains secrets)
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database

### Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and update the values
3. Run database migrations:
   ```
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   migrate -path migrations -database "${DATABASE_URL}" up
   ```
4. Build and run the application:
   ```
   go build -o affiliate-backend ./cmd/api
   ./affiliate-backend
   ```

## API Endpoints

### Public Endpoints

- `POST /api/v1/public/webhooks/supabase/new-user`: Webhook for Supabase auth to create a profile when a new user signs up

### Authenticated Endpoints

- `GET /api/v1/users/me`: Get the current user's profile

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