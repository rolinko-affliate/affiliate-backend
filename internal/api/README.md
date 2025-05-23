# API Module

This module is responsible for setting up the HTTP API server, defining routes, and connecting handlers with middleware. It serves as the entry point for all HTTP requests to the application.

## Key Components

### Router

The `router.go` file defines the API routes and connects them with the appropriate handlers and middleware:
- Groups routes by resource type (profiles, organizations, advertisers, etc.)
- Applies authentication middleware to protected routes
- Applies RBAC middleware with appropriate role requirements
- Maps HTTP methods and paths to handler functions

## Route Structure

The API follows a RESTful design with the following route structure:

- `/api/v1/public/*`: Public endpoints that don't require authentication
- `/api/v1/users/me`: Endpoint for retrieving the current user's profile
- `/api/v1/profiles/*`: Profile management endpoints
- `/api/v1/organizations/*`: Organization management endpoints
- `/api/v1/advertisers/*`: Advertiser management endpoints
- `/api/v1/affiliates/*`: Affiliate management endpoints
- `/api/v1/campaigns/*`: Campaign management endpoints

## Middleware Application

The router applies middleware in a specific order:
1. Global middleware (logging, recovery, CORS)
2. Authentication middleware for protected routes
3. RBAC middleware with appropriate role requirements

## RouterOptions

The `RouterOptions` struct provides dependencies to the router:
- `ProfileHandler`: Handles profile-related requests
- `ProfileService`: Used by RBAC middleware to check user roles
- `OrganizationHandler`: Handles organization-related requests
- `AdvertiserHandler`: Handles advertiser-related requests
- `AffiliateHandler`: Handles affiliate-related requests
- `CampaignHandler`: Handles campaign-related requests

## Security Features

- All routes except public ones require authentication
- Role-based access control for different resource types
- Organization-based access control in handlers
- CORS configuration based on environment

## Usage

The router is initialized in the main application entry point:

```go
// Initialize handlers
profileHandler := handlers.NewProfileHandler(profileService)
// ... initialize other handlers

// Setup router
router := api.SetupRouter(api.RouterOptions{
    ProfileHandler: profileHandler,
    ProfileService: profileService,
    // ... other handlers
})

// Start server
server := &http.Server{
    Addr:    ":" + port,
    Handler: router,
}
server.ListenAndServe()
```