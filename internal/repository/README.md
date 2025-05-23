# Repository Module

This module implements the data access layer of the application, providing interfaces and implementations for database operations. It follows the repository pattern to abstract database interactions from the business logic.

## Key Components

### Database Connection

The module provides functions for initializing and managing database connections:
- `InitDB`: Initializes the database connection pool
- `CloseDB`: Closes the database connection pool
- `InitDBConnection`: Creates a new database connection for testing or migrations

### Repository Interfaces

Each entity has a corresponding repository interface that defines its data access operations:

```go
// Example: OrganizationRepository interface
type OrganizationRepository interface {
    CreateOrganization(ctx context.Context, org *domain.Organization) error
    GetOrganizationByID(ctx context.Context, id int64) (*domain.Organization, error)
    UpdateOrganization(ctx context.Context, org *domain.Organization) error
    ListOrganizations(ctx context.Context, limit, offset int) ([]*domain.Organization, error)
    DeleteOrganization(ctx context.Context, id int64) error
}
```

### PostgreSQL Implementations

Each repository interface has a PostgreSQL implementation using the pgx library:

```go
// Example: pgxOrganizationRepository implementation
type pgxOrganizationRepository struct {
    db *pgxpool.Pool
}

func NewPgxOrganizationRepository(db *pgxpool.Pool) OrganizationRepository {
    return &pgxOrganizationRepository{db: db}
}

// Implementation of interface methods...
```

## Repository Types

### ProfileRepository

Manages user profiles in the database:
- Creating and retrieving profiles
- Updating profile information
- Deleting profiles
- Retrieving roles

### OrganizationRepository

Manages organizations in the database:
- Creating and retrieving organizations
- Updating organization information
- Listing organizations with pagination
- Deleting organizations

### AdvertiserRepository

Manages advertisers and their provider mappings:
- Creating and retrieving advertisers
- Updating advertiser information
- Listing advertisers by organization
- Managing advertiser provider mappings

### AffiliateRepository

Manages affiliates and their provider mappings:
- Creating and retrieving affiliates
- Updating affiliate information
- Listing affiliates by organization
- Managing affiliate provider mappings

### CampaignRepository

Manages campaigns and their provider offers:
- Creating and retrieving campaigns
- Updating campaign information
- Listing campaigns by organization or advertiser
- Managing campaign provider offers

## Error Handling

The repositories provide consistent error handling:
- Domain-specific errors (e.g., `domain.ErrNotFound`)
- Detailed error messages with wrapped errors
- SQL error translation to domain errors

## Transaction Support

Some operations support transactions for maintaining data consistency:
- Begin/commit/rollback transaction handling
- Context-based transaction management
- Error handling with automatic rollback

## Database Schema

The repositories work with the database schema defined in the migrations:
- Tables with foreign key relationships
- Indexes for performance optimization
- Timestamps for auditing
- JSON/JSONB fields for flexible data storage