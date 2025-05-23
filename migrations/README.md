# Database Migrations

This directory contains database migration files for the application. Migrations are used to manage database schema changes in a versioned and repeatable way.

## Migration Files

Migration files follow the naming convention:

```
000001_name.up.sql   # SQL to apply the migration
000001_name.down.sql # SQL to rollback the migration
```

The numeric prefix determines the order in which migrations are applied.

## Initial Schema

The initial schema migration (`000001_create_initial_tables.up.sql`) creates the core tables for the application:

- `organizations`: Tenant organizations
- `roles`: User roles for RBAC
- `profiles`: User profiles linked to Supabase Auth
- `advertisers`: Advertiser entities
- `affiliates`: Affiliate entities
- `campaigns`: Advertising campaigns
- `advertiser_provider_mappings`: Links advertisers to external providers
- `affiliate_provider_mappings`: Links affiliates to external providers
- `campaign_provider_offers`: Links campaigns to external provider offers

## Schema Features

The database schema includes several features:

### Automatic Timestamps

All tables have `created_at` and `updated_at` timestamps that are automatically managed:

```sql
-- Function to update updated_at timestamp automatically
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger on each table
CREATE TRIGGER set_organizations_timestamp
BEFORE UPDATE ON public.organizations
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
```

### Foreign Key Relationships

Tables are linked with foreign key relationships:

```sql
-- Example: profiles table with foreign keys
CREATE TABLE public.profiles (
    id UUID PRIMARY KEY,
    organization_id BIGINT REFERENCES public.organizations(organization_id) ON DELETE SET NULL,
    role_id INT REFERENCES public.roles(role_id) ON DELETE RESTRICT NOT NULL,
    -- Other columns...
);
```

### Indexes

Indexes are created for performance optimization:

```sql
-- Example: indexes on profiles table
CREATE INDEX idx_profiles_organization_id ON public.profiles(organization_id);
CREATE INDEX idx_profiles_role_id ON public.profiles(role_id);
CREATE INDEX idx_profiles_email ON public.profiles(email);
```

### JSON/JSONB Fields

Some tables use JSON/JSONB fields for flexible data storage:

```sql
-- Example: JSON fields in advertisers table
CREATE TABLE public.advertisers (
    -- Other columns...
    billing_details JSONB, -- Store address, tax ID, etc.
    -- Other columns...
);
```

### Status Enumerations

Status fields use CHECK constraints to enforce valid values:

```sql
-- Example: status field in advertisers table
CREATE TABLE public.advertisers (
    -- Other columns...
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'inactive', 'rejected')),
    -- Other columns...
);
```

## Seed Data

The initial migration includes seed data for roles and a default organization:

```sql
-- Seed initial roles
INSERT INTO public.roles (role_id, name, description) VALUES
  (100000, 'User', 'Default user role with limited access'),
  (1000, 'AdvertiserManager', 'Manages advertisers and their campaigns within their organization'),
  (1001, 'AffiliateManager', 'Manages affiliates and approves applications within their organization'),
  (1, 'Admin', 'Platform Administrator with full access');

-- Create default organization
INSERT INTO public.organizations (name) VALUES ('rolinko');
```

## Managing Migrations

Migrations are managed using the migration tool in the `cmd/migrate` directory:

```bash
# Apply all pending migrations
go run cmd/migrate/main.go up

# Rollback the most recent migration
go run cmd/migrate/main.go down

# Other migration commands...
```

Or using the Makefile:

```bash
# Apply all pending migrations
make migrate-up

# Rollback the most recent migration
make migrate-down

# Other migration commands...
```

## Creating New Migrations

New migrations can be created using the migration tool:

```bash
# Create a new migration
go run cmd/migrate/main.go create add_new_table
```

Or using the Makefile:

```bash
# Create a new migration
make migrate-create NAME=add_new_table
```

This will create two new files:
- `migrations/000002_add_new_table.up.sql`
- `migrations/000002_add_new_table.down.sql`

Edit these files to add the SQL statements for applying and rolling back the migration.