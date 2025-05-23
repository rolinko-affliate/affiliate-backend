# Domain Module

This module contains the core business entities and models for the application. These structs represent the fundamental data structures and are used throughout the application.

## Key Entities

### Profile

Represents a user profile in the system:
- Links to the Supabase Auth user ID
- Contains user information (email, name)
- Associates users with organizations and roles
- Includes mandatory role name for authorization

```go
type Profile struct {
    ID             uuid.UUID  `json:"id" db:"id"`
    OrganizationID *int64     `json:"organization_id,omitempty" db:"organization_id"`
    RoleID         int        `json:"role_id" db:"role_id"`
    RoleName       string     `json:"role_name" db:"role_name"`
    Email          string     `json:"email" db:"email"`
    FirstName      *string    `json:"first_name,omitempty" db:"first_name"`
    LastName       *string    `json:"last_name,omitempty" db:"last_name"`
    CreatedAt      time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}
```

### Role

Defines user roles for RBAC:
- Admin: Full system access
- AdvertiserManager: Manages advertisers and campaigns
- AffiliateManager: Manages affiliates
- Affiliate: Regular affiliate user

```go
type Role struct {
    RoleID      int     `json:"role_id" db:"role_id"`
    Name        string  `json:"name" db:"name"`
    Description *string `json:"description,omitempty" db:"description"`
}
```

### Organization

Represents a tenant in the multi-tenant system:
- Contains basic organization information
- Serves as a container for advertisers, affiliates, and campaigns
- Used for organization-based access control

```go
type Organization struct {
    OrganizationID int64     `json:"organization_id" db:"organization_id"`
    Name           string    `json:"name" db:"name"`
    CreatedAt      time.Time `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
```

### Advertiser

Represents an advertiser entity:
- Belongs to an organization
- Contains advertiser details and status
- Can be linked to external provider systems

```go
type Advertiser struct {
    AdvertiserID    int64     `json:"advertiser_id" db:"advertiser_id"`
    OrganizationID  int64     `json:"organization_id" db:"organization_id"`
    Name            string    `json:"name" db:"name"`
    ContactEmail    *string   `json:"contact_email,omitempty" db:"contact_email"`
    BillingDetails  *string   `json:"billing_details,omitempty" db:"billing_details"`
    Status          string    `json:"status" db:"status"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
```

### Affiliate

Represents an affiliate entity:
- Belongs to an organization
- Contains affiliate details and status
- Can be linked to external provider systems

```go
type Affiliate struct {
    AffiliateID     int64     `json:"affiliate_id" db:"affiliate_id"`
    OrganizationID  int64     `json:"organization_id" db:"organization_id"`
    Name            string    `json:"name" db:"name"`
    ContactEmail    *string   `json:"contact_email,omitempty" db:"contact_email"`
    PaymentDetails  *string   `json:"payment_details,omitempty" db:"payment_details"`
    Status          string    `json:"status" db:"status"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
```

### Campaign

Represents an advertising campaign:
- Belongs to an organization and advertiser
- Contains campaign details and status
- Can be linked to external provider offers

```go
type Campaign struct {
    CampaignID     int64     `json:"campaign_id" db:"campaign_id"`
    OrganizationID int64     `json:"organization_id" db:"organization_id"`
    AdvertiserID   int64     `json:"advertiser_id" db:"advertiser_id"`
    Name           string    `json:"name" db:"name"`
    Description    *string   `json:"description,omitempty" db:"description"`
    StartDate      *time.Time `json:"start_date,omitempty" db:"start_date"`
    EndDate        *time.Time `json:"end_date,omitempty" db:"end_date"`
    Status         string    `json:"status" db:"status"`
    CreatedAt      time.Time `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
```

### Provider Mappings

Maps internal entities to external provider systems:
- AdvertiserProviderMapping: Links advertisers to external systems
- AffiliateProviderMapping: Links affiliates to external systems
- CampaignProviderOffer: Links campaigns to external offers

## Common Patterns

- All entities include created_at and updated_at timestamps
- Nullable fields are represented as pointers
- JSON tags for API serialization
- Database tags for ORM mapping
- Status fields for entity lifecycle management