# API Middleware

This module contains middleware components for the Gin web framework that handle cross-cutting concerns such as authentication, authorization, and CORS.

## Key Components

### AuthMiddleware

Handles JWT-based authentication using Supabase:
- Validates the JWT token in the Authorization header
- Extracts user ID from the token and stores it in the request context
- Rejects requests with missing or invalid tokens

```go
// Usage
router.Use(middleware.AuthMiddleware())
```

### RBACMiddleware

Implements Role-Based Access Control (RBAC):
- Retrieves the user's profile using the user ID from the context
- Checks if the user's role is in the list of allowed roles
- Adds the user's role and organization ID to the context for downstream handlers
- Rejects requests from users with insufficient permissions

```go
// Usage
router.Use(middleware.RBACMiddleware(profileService, "Admin", "AdvertiserManager"))
```

### CORSMiddleware

Configures Cross-Origin Resource Sharing (CORS) based on the environment:
- In development mode, allows requests from any origin
- In production mode, disables CORS
- Configures allowed methods, headers, and credentials

```go
// Usage
router.Use(middleware.CORSMiddleware())
```

## Context Keys

The middleware components store important information in the Gin context:
- `UserIDKey`: The user's ID from the JWT token
- `UserEmailKey`: The user's email (if available)
- `UserRoleKey`: The user's role name
- `organizationID`: The user's organization ID

Downstream handlers can access these values to perform permission checks and business logic.

## Security Features

- JWT validation with proper signature verification
- Role-based access control for route protection
- Organization-based access control in handlers
- Environment-aware CORS configuration

## Error Handling

Middleware components provide clear error responses with appropriate HTTP status codes:
- 401 Unauthorized: Missing or invalid authentication
- 403 Forbidden: Insufficient permissions
- 500 Internal Server Error: Unexpected errors during middleware processing