# Affiliate Domain Implementation Pattern

## Overview

This document describes the comprehensive domain implementation pattern used in the affiliate-backend Go project, covering the architecture from API layer to external provider integrations. The system follows Clean Architecture principles with clear separation of concerns and dependency inversion.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        API Layer (Gin)                         │
├─────────────────────────────────────────────────────────────────┤
│  Handlers │ Middleware │ Router │ Models │ Error Handling      │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Service Layer                             │
├─────────────────────────────────────────────────────────────────┤
│  Business Logic │ Validation │ Orchestration │ Provider Sync   │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Domain Layer                               │
├─────────────────────────────────────────────────────────────────┤
│  Entities │ Value Objects │ Domain Events │ Business Rules     │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Repository Layer                             │
├─────────────────────────────────────────────────────────────────┤
│  Data Access │ Persistence │ Query Building │ Transactions     │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Provider Layer                               │
├─────────────────────────────────────────────────────────────────┤
│  Integration Service │ Everflow Client │ Mappers │ Mock Service │
└─────────────────────────────────────────────────────────────────┘
```

## Layer-by-Layer Implementation

### 1. API Layer (`/internal/api/`)

#### Router Pattern (`router.go`)
- **Centralized Route Configuration**: Single router setup with grouped routes
- **Dependency Injection**: RouterOptions struct for clean dependency management
- **Middleware Composition**: Layered middleware application (CORS, Auth, RBAC)
- **Role-Based Access Control**: Dynamic RBAC middleware factory

```go
type RouterOptions struct {
    ProfileHandler       *handlers.ProfileHandler
    ProfileService       service.ProfileService
    OrganizationHandler  *handlers.OrganizationHandler
    AdvertiserHandler    *handlers.AdvertiserHandler
    AffiliateHandler     *handlers.AffiliateHandler
    CampaignHandler      *handlers.CampaignHandler
}
```

#### Handler Pattern (`/handlers/`)
- **Resource-Based Handlers**: Separate handlers for each domain entity
- **Constructor Injection**: Dependencies injected via constructors
- **Context Propagation**: Request context passed through all layers
- **Error Standardization**: Consistent error response format

```go
type AffiliateHandler struct {
    affiliateService service.AffiliateService
    profileService   service.ProfileService
}
```

#### Middleware Pattern (`/middleware/`)
- **Authentication Middleware**: JWT validation with Supabase integration
- **RBAC Middleware**: Role-based access control with dynamic role checking
- **CORS Middleware**: Cross-origin resource sharing configuration

### 2. Domain Layer (`/internal/domain/`)

#### Clean Domain Models
The domain models follow Clean Architecture principles with clear separation between core business entities and provider-specific data:

```go
// Core Affiliate Entity
type Affiliate struct {
    AffiliateID             int64     `json:"affiliate_id" db:"affiliate_id"`
    OrganizationID          int64     `json:"organization_id" db:"organization_id"`
    Name                    string    `json:"name" db:"name"`
    ContactEmail            *string   `json:"contact_email,omitempty" db:"contact_email"`
    Status                  string    `json:"status" db:"status"`
    
    // General purpose fields (provider-agnostic)
    InternalNotes           *string   `json:"internal_notes,omitempty" db:"internal_notes"`
    DefaultCurrencyID       *string   `json:"default_currency_id,omitempty" db:"default_currency_id"`
    ContactAddress          *string   `json:"contact_address,omitempty" db:"contact_address"`
    
    CreatedAt               time.Time `json:"created_at" db:"created_at"`
    UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`
}
```

#### Provider Mapping Pattern
Separate entities for provider-specific data to maintain clean domain separation:

```go
// Provider Mapping Entity
type AffiliateProviderMapping struct {
    MappingID           int64     `json:"mapping_id" db:"mapping_id"`
    AffiliateID         int64     `json:"affiliate_id" db:"affiliate_id"`
    ProviderType        string    `json:"provider_type" db:"provider_type"`
    ProviderAffiliateID *string   `json:"provider_affiliate_id,omitempty" db:"provider_affiliate_id"`
    
    // Provider-specific data stored as JSONB
    ProviderData        *string   `json:"provider_data,omitempty" db:"provider_data"`
    
    // Synchronization metadata
    SyncStatus          *string   `json:"sync_status,omitempty" db:"sync_status"`
    LastSyncAt          *time.Time `json:"last_sync_at,omitempty" db:"last_sync_at"`
}
```

#### Provider-Specific Value Objects
```go
// Everflow-specific data structure
type EverflowProviderData struct {
    NetworkAffiliateID              *int32          `json:"network_affiliate_id,omitempty"`
    EnableMediaCostTrackingLinks    *bool           `json:"enable_media_cost_tracking_links,omitempty"`
    ReferrerID                      *int32          `json:"referrer_id,omitempty"`
    NetworkAffiliateTierID          *int32          `json:"network_affiliate_tier_id,omitempty"`
    NetworkEmployeeID               *int32          `json:"network_employee_id,omitempty"`
    Users                           *[]interface{}  `json:"users,omitempty"`
    AdditionalFields                map[string]interface{} `json:"additional_fields,omitempty"`
}
```

### 3. Service Layer (`/internal/service/`)

#### Service Interface Pattern
Clean interfaces defining business operations:

```go
type AffiliateService interface {
    CreateAffiliate(ctx context.Context, affiliate *domain.Affiliate) (*domain.Affiliate, error)
    GetAffiliateByID(ctx context.Context, id int64) (*domain.Affiliate, error)
    UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error
    ListAffiliatesByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Affiliate, error)
    DeleteAffiliate(ctx context.Context, id int64) error
    
    // Provider sync methods
    SyncAffiliateToProvider(ctx context.Context, affiliateID int64) error
    SyncAffiliateFromProvider(ctx context.Context, affiliateID int64) error
    
    // Provider mapping methods
    CreateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) (*domain.AffiliateProviderMapping, error)
    GetAffiliateProviderMapping(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error)
    UpdateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error
    DeleteAffiliateProviderMapping(ctx context.Context, mappingID int64) error
}
```

#### Service Implementation Pattern
```go
type affiliateService struct {
    affiliateRepo           repository.AffiliateRepository
    providerMappingRepo     repository.AffiliateProviderMappingRepository
    orgRepo                 repository.OrganizationRepository
    integrationService      provider.IntegrationService
}
```

#### Business Logic Orchestration
Services orchestrate complex business operations:

1. **Validation**: Business rule validation
2. **Repository Operations**: Data persistence
3. **Provider Integration**: External system synchronization
4. **Error Handling**: Comprehensive error management

### 4. Repository Layer (`/internal/repository/`)

#### Repository Interface Pattern
Data access abstraction with clean interfaces:

```go
type AffiliateRepository interface {
    CreateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error
    GetAffiliateByID(ctx context.Context, affiliateID int64) (*domain.Affiliate, error)
    UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error
    DeleteAffiliate(ctx context.Context, affiliateID int64) error
    GetAffiliatesByOrganization(ctx context.Context, organizationID int64) ([]*domain.Affiliate, error)
    ListAffiliatesByOrganization(ctx context.Context, organizationID int64, limit, offset int) ([]*domain.Affiliate, error)
}
```

#### Database Implementation
- **PostgreSQL with pgx**: High-performance PostgreSQL driver
- **Connection Pooling**: Efficient database connection management
- **Transaction Support**: ACID compliance for complex operations
- **JSONB Support**: Flexible storage for provider-specific data

#### SQL Query Patterns
```sql
-- Clean domain model persistence
INSERT INTO public.affiliates (
    organization_id, name, contact_email, payment_details, status,
    internal_notes, default_currency_id, contact_address, billing_info, labels,
    invoice_amount_threshold, default_payment_terms,
    created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING affiliate_id, created_at, updated_at
```

### 5. Provider Layer (`/internal/platform/`)

#### Integration Service Pattern
Provider-agnostic interface for external integrations:

```go
type IntegrationService interface {
    // Advertisers
    CreateAdvertiser(ctx context.Context, adv domain.Advertiser) (domain.Advertiser, error)
    UpdateAdvertiser(ctx context.Context, adv domain.Advertiser) error
    GetAdvertiser(ctx context.Context, id uuid.UUID) (domain.Advertiser, error)

    // Affiliates
    CreateAffiliate(ctx context.Context, aff domain.Affiliate) (domain.Affiliate, error)
    UpdateAffiliate(ctx context.Context, aff domain.Affiliate) error
    GetAffiliate(ctx context.Context, id uuid.UUID) (domain.Affiliate, error)

    // Campaigns
    CreateCampaign(ctx context.Context, camp domain.Campaign) (domain.Campaign, error)
    UpdateCampaign(ctx context.Context, camp domain.Campaign) error
    GetCampaign(ctx context.Context, id uuid.UUID) (domain.Campaign, error)
}
```

#### Everflow Implementation (`/everflow/`)
Concrete implementation for Everflow provider:

```go
type IntegrationService struct {
    advertiserClient *advertiser.APIClient
    affiliateClient  *affiliate.APIClient
    offerClient      *offer.APIClient
    
    // Repository interfaces for provider mappings
    advertiserRepo AdvertiserRepository
    affiliateRepo  AffiliateRepository
    campaignRepo   CampaignRepository
    
    advertiserProviderMappingRepo AdvertiserProviderMappingRepository
    affiliateProviderMappingRepo  AffiliateProviderMappingRepository
    campaignProviderMappingRepo   CampaignProviderMappingRepository
    
    // Provider mappers
    affiliateProviderMapper *AffiliateProviderMapper
}
```

#### Mock Implementation
Comprehensive mock service for testing and development:

```go
type LoggingMockIntegrationService struct {
    // Implements all IntegrationService methods with logging and simulation
}
```

## Key Design Patterns

### 1. Dependency Inversion
- High-level modules don't depend on low-level modules
- Both depend on abstractions (interfaces)
- Abstractions don't depend on details

### 2. Repository Pattern
- Encapsulates data access logic
- Provides a uniform interface for data operations
- Enables easy testing with mock implementations

### 3. Service Layer Pattern
- Encapsulates business logic
- Orchestrates operations across multiple repositories
- Provides transaction boundaries

### 4. Provider Pattern
- Abstracts external service integrations
- Enables multiple provider implementations
- Supports mock implementations for testing

### 5. Mapper Pattern
- Converts between domain models and external API models
- Handles provider-specific data transformations
- Maintains clean separation between internal and external representations

## Database Schema Design

### Core Tables
- **organizations**: Multi-tenant organization management
- **profiles**: User management with Supabase integration
- **roles**: Role-based access control
- **affiliates**: Clean affiliate domain entities
- **advertisers**: Clean advertiser domain entities
- **campaigns**: Clean campaign domain entities

### Provider Mapping Tables
- **affiliate_provider_mappings**: Affiliate-to-provider relationships
- **advertiser_provider_mappings**: Advertiser-to-provider relationships
- **campaign_provider_mappings**: Campaign-to-provider relationships

### JSONB Usage
Provider-specific data stored as JSONB for flexibility:
- `provider_data`: Everflow-specific fields
- `provider_config`: Provider configuration
- `api_credentials`: Encrypted API credentials

## Authentication & Authorization

### JWT Authentication
- Supabase JWT token validation
- User context propagation
- Secure token parsing and validation

### Role-Based Access Control (RBAC)
- Dynamic role checking
- Organization-based access control
- Resource-level permissions

### Middleware Chain
1. **CORS Middleware**: Cross-origin request handling
2. **Auth Middleware**: JWT validation and user context
3. **RBAC Middleware**: Role-based access control

## Error Handling Strategy

### Layered Error Handling
- **Domain Errors**: Business rule violations
- **Repository Errors**: Data access failures
- **Provider Errors**: External service failures
- **API Errors**: HTTP-specific error responses

### Error Response Format
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code,omitempty"`
    Details string `json:"details,omitempty"`
}
```

## Testing Strategy

### Unit Testing
- Service layer business logic testing
- Repository layer data access testing
- Provider layer integration testing

### Integration Testing
- End-to-end API testing
- Database integration testing
- Provider integration testing

### Mock Testing
- Mock integration service for development
- Mock repositories for unit testing
- Mock providers for isolated testing

## Configuration Management

### Environment-Based Configuration
- Database connection strings
- Provider API credentials
- Feature flags (mock mode)
- JWT secrets

### Configuration Structure
```go
type Config struct {
    DatabaseURL        string `mapstructure:"database_url"`
    Port              string `mapstructure:"port"`
    SupabaseJWTSecret string `mapstructure:"supabase_jwt_secret"`
    MockMode          bool   `mapstructure:"mock_mode"`
}
```

## Deployment Considerations

### Docker Support
- Multi-stage Docker builds
- Environment variable configuration
- Health check endpoints

### Database Migrations
- Version-controlled schema changes
- Up/down migration support
- Migration validation scripts

### Monitoring & Logging
- Structured logging with context
- Request tracing
- Performance metrics
- Error tracking

## Benefits of This Pattern

1. **Maintainability**: Clear separation of concerns
2. **Testability**: Easy mocking and unit testing
3. **Scalability**: Modular architecture supports growth
4. **Flexibility**: Easy to add new providers or features
5. **Security**: Layered security with proper authentication/authorization
6. **Performance**: Efficient database operations and connection pooling
7. **Developer Experience**: Clear patterns and comprehensive documentation

## Future Enhancements

1. **Event Sourcing**: Domain event tracking
2. **CQRS**: Command Query Responsibility Segregation
3. **Microservices**: Service decomposition
4. **GraphQL**: Alternative API layer
5. **Caching**: Redis integration for performance
6. **Message Queues**: Asynchronous processing
7. **API Versioning**: Backward compatibility support

This domain implementation pattern provides a solid foundation for building scalable, maintainable affiliate management systems with clean architecture principles and comprehensive provider integration capabilities.