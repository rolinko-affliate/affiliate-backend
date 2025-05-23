# Service Module

This module implements the business logic layer of the application. Services act as intermediaries between the API handlers and repositories, encapsulating business rules, validation, and orchestration of data operations.

## Key Components

### Service Interfaces

Each entity has a corresponding service interface that defines its business operations:

```go
// Example: OrganizationService interface
type OrganizationService interface {
    CreateOrganization(ctx context.Context, name string) (*domain.Organization, error)
    GetOrganizationByID(ctx context.Context, id int64) (*domain.Organization, error)
    UpdateOrganization(ctx context.Context, org *domain.Organization) error
    ListOrganizations(ctx context.Context, page, pageSize int) ([]*domain.Organization, error)
    DeleteOrganization(ctx context.Context, id int64) error
}
```

### Service Implementations

Each service interface has a concrete implementation that uses repositories for data access:

```go
// Example: organizationService implementation
type organizationService struct {
    orgRepo repository.OrganizationRepository
}

func NewOrganizationService(orgRepo repository.OrganizationRepository) OrganizationService {
    return &organizationService{orgRepo: orgRepo}
}

// Implementation of interface methods...
```

## Service Types

### ProfileService

Manages user profile business logic:
- Creating new user profiles
- Retrieving and updating profiles
- Handling profile deletion
- Managing user roles

### OrganizationService

Manages organization business logic:
- Creating and validating organizations
- Retrieving organization details
- Updating organization information
- Listing organizations with pagination

### AdvertiserService

Manages advertiser business logic:
- Creating and validating advertisers
- Retrieving advertiser details
- Updating advertiser information
- Managing advertiser provider mappings
- Integration with external provider systems (e.g., Everflow)

### AffiliateService

Manages affiliate business logic:
- Creating and validating affiliates
- Retrieving affiliate details
- Updating affiliate information
- Managing affiliate provider mappings

### CampaignService

Manages campaign business logic:
- Creating and validating campaigns
- Retrieving campaign details
- Updating campaign information
- Managing campaign provider offers
- Integration with external provider systems

## Business Rules

Services implement various business rules:
- Input validation and sanitization
- Entity relationship validation
- Status management and transitions
- Permission and access control logic
- Integration with external systems

## Error Handling

Services provide consistent error handling:
- Domain-specific errors with clear messages
- Validation errors for invalid input
- Not found errors for missing entities
- Permission errors for unauthorized access

## Integration with External Systems

Some services integrate with external systems:
- Everflow API for affiliate network integration
- Encryption service for sensitive data
- Other third-party services as needed

## Pagination

Services handle pagination for list operations:
- Page and page size parameters
- Default values for pagination
- Conversion between page/pageSize and limit/offset