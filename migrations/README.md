# Database Migrations

This directory contains database migration files for the affiliate platform.

## Migration Files

- `000001_create_complete_schema.sql` - Creates the complete database schema with all tables and Everflow integration support

## Schema Overview

The migration creates the following tables:

### Core Platform Tables
- `organizations` - Multi-tenant organization support
- `roles` - User role definitions
- `profiles` - User profiles linked to Supabase Auth

### Advertiser & Affiliate Tables
- `advertisers` - Advertiser entities with complete Everflow integration fields
- `affiliates` - Affiliate entities
- `advertiser_provider_mappings` - Maps advertisers to external providers (Everflow)
- `affiliate_provider_mappings` - Maps affiliates to external providers

### Campaign Tables
- `campaigns` - Campaign entities with complete Everflow offer field support (25+ fields)
- `campaign_provider_offers` - Maps campaigns to provider offers

## Features

- **Complete Everflow Integration**: All necessary fields for Everflow offer creation and management
- **Multi-tenant Support**: Organization-based data isolation
- **Provider Abstraction**: Designed to support multiple affiliate networks (currently Everflow)
- **Comprehensive Offer Support**: Caps, payouts, revenue, tracking, and configuration options
- **Audit Trail**: Created/updated timestamps with automatic triggers

## Campaign Table Fields

The campaigns table includes comprehensive Everflow offer support with the following field categories:

### Core Campaign Fields
- `campaign_id`, `organization_id`, `advertiser_id`, `name`, `description`, `status`
- `start_date`, `end_date`, `created_at`, `updated_at`

### Everflow Offer Fields
- **URLs**: `destination_url`, `thumbnail_url`, `preview_url`, `server_side_url`, `view_through_destination_url`
- **Configuration**: `visibility`, `currency_id`, `conversion_method`, `session_definition`, `session_duration`
- **Content**: `internal_notes`, `terms_and_conditions`, `html_description`, `app_identifier`
- **Settings**: `is_force_terms_and_conditions`, `is_using_explicit_terms_and_conditions`, `is_whitelist_check_enabled`, `is_view_through_enabled`, `is_description_plain_text`, `is_use_direct_linking`
- **Timing**: `caps_timezone_id`, `project_id`, `date_live_until`

### Caps and Limits
- **Conversion Caps**: `daily_conversion_cap`, `weekly_conversion_cap`, `monthly_conversion_cap`, `global_conversion_cap`
- **Click Caps**: `daily_click_cap`, `weekly_click_cap`, `monthly_click_cap`, `global_click_cap`
- **Cap Control**: `is_caps_enabled`

### Tracking and Revenue
- **Tracking**: `encoded_value`, `today_clicks`, `today_revenue`, `time_created`, `time_saved`
- **Payout**: `payout_type`, `payout_amount`
- **Revenue**: `revenue_type`, `revenue_amount`
- **Configuration**: `offer_config` (JSONB for additional settings)

## Running Migrations

Use your preferred migration tool (e.g., golang-migrate, Flyway, etc.) to apply the migration.

Example with golang-migrate:
```bash
migrate -path ./migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" up
```

## Database Recreation

Since this is a single comprehensive migration, you can easily recreate the database from scratch:

```bash
# Drop and recreate database
migrate -path ./migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" down
migrate -path ./migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" up
```

## Schema Features

### Automatic Timestamps
All tables have `created_at` and `updated_at` timestamps that are automatically managed using triggers.

### Foreign Key Relationships
Tables are properly linked with foreign key relationships and appropriate cascade behaviors.

### Indexes
Performance indexes are created on frequently queried columns including:
- Organization and advertiser relationships
- Campaign status and visibility
- Provider mappings
- JSONB configuration fields (GIN indexes)

### Data Validation
- CHECK constraints on enum fields (status, visibility, payout types, etc.)
- NOT NULL constraints on required fields
- UNIQUE constraints where appropriate

### Seed Data
The migration includes initial seed data for:
- Default user roles (Admin, AdvertiserManager, AffiliateManager, User)
- Default organization ('rolinko')

## Managing Migrations

Migrations can be managed using standard migration tools or custom migration commands in your application.

## Creating New Migrations

When adding new features, create additional migration files following the naming convention:
- `000002_feature_name.up.sql`
- `000002_feature_name.down.sql`