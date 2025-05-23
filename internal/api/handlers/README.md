# API Handlers

This module contains HTTP request handlers for all resources in the application. Handlers are responsible for:

1. Parsing and validating incoming HTTP requests
2. Checking user permissions for the requested operation
3. Calling the appropriate service methods to perform business logic
4. Formatting and returning HTTP responses

## Key Components

### ProfileHandler

Manages user profiles and authentication-related operations:
- Creating and updating user profiles
- Retrieving the current user's profile
- Handling Supabase webhook for new user creation

### OrganizationHandler

Manages organizations with strict permission checks:
- Only Admin users can create organizations
- Users can only view/update/delete organizations they belong to
- Admin users have access to all organizations

### AdvertiserHandler

Manages advertisers with organization-based access control:
- Creating, retrieving, updating, and deleting advertisers
- Managing advertiser provider mappings (e.g., Everflow integration)
- Users can only access advertisers within their organization

### AffiliateHandler

Manages affiliates with organization-based access control:
- Creating, retrieving, updating, and deleting affiliates
- Managing affiliate provider mappings
- Users can only access affiliates within their organization

### CampaignHandler

Manages advertising campaigns:
- Creating, retrieving, updating, and deleting campaigns
- Managing campaign provider offers
- Listing campaigns by organization or advertiser

## Permission Checking

All handlers implement organization-based permission checking:
1. Admin users have full access to all resources
2. Non-admin users can only access resources within their organization
3. Permission checks are performed before any data access or modification

## Request/Response Flow

1. Request validation using Gin's binding mechanism
2. Permission checking based on user role and organization
3. Service method invocation to perform business logic
4. Response formatting with appropriate HTTP status codes
5. Error handling with descriptive error messages

## Error Handling

Handlers provide consistent error responses with appropriate HTTP status codes:
- 400 Bad Request: Invalid input data
- 401 Unauthorized: Missing or invalid authentication
- 403 Forbidden: Insufficient permissions
- 404 Not Found: Resource not found
- 500 Internal Server Error: Unexpected errors

## Swagger Documentation

All handlers include Swagger annotations for automatic API documentation generation.