-- #############################################################################
-- ## Complete Affiliate Platform Database Schema Initialization
-- ## This migration creates the complete database schema with refactored affiliate model
-- ## Includes: Core tables, provider mappings, and general purpose affiliate fields
-- #############################################################################

-- #############################################################################
-- ## Helper Functions
-- #############################################################################

-- Function to update updated_at timestamp automatically
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- #############################################################################
-- ## Core Platform Tables
-- #############################################################################

-- organizations: Represents distinct tenants or teams using the platform.
CREATE TABLE public.organizations (
    organization_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'platform_owner' 
        CHECK (type IN ('advertiser', 'affiliate', 'platform_owner')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_organizations_timestamp
BEFORE UPDATE ON public.organizations
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_organizations_type ON public.organizations(type);

-- roles: Defines the different user roles within the platform.
CREATE TABLE public.roles (
    role_id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL, -- e.g., 'Admin', 'AdvertiserManager', 'AffiliateManager', 'Affiliate'
    description TEXT
);

-- profiles: Stores custom user information. The `id` column stores the `auth.uid()` obtained from Supabase Auth.
CREATE TABLE public.profiles (
    id UUID PRIMARY KEY, -- Stores the auth.uid() from Supabase Auth
    organization_id BIGINT REFERENCES public.organizations(organization_id) ON DELETE SET NULL,
    role_id INT REFERENCES public.roles(role_id) ON DELETE RESTRICT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL, -- Store email for easier lookup and consistency
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_profiles_timestamp
BEFORE UPDATE ON public.profiles
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_profiles_organization_id ON public.profiles(organization_id);
CREATE INDEX idx_profiles_role_id ON public.profiles(role_id);
CREATE INDEX idx_profiles_email ON public.profiles(email);

-- #############################################################################
-- ## Core Domain Tables (Clean - No Provider-Specific Fields)
-- #############################################################################

-- advertisers: Core advertiser information only
CREATE TABLE public.advertisers (
    advertiser_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255),
    billing_details JSONB, -- Store address, tax ID, etc.
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'inactive', 'rejected')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_advertisers_timestamp
BEFORE UPDATE ON public.advertisers
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_advertisers_organization_id ON public.advertisers(organization_id);

-- affiliates: Core affiliate information with general purpose fields (refactored model)
CREATE TABLE public.affiliates (
    affiliate_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255),
    payment_details JSONB, -- Store payment method (PayPal, Wise, Bank) and details
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'rejected', 'inactive')),
    
    -- General purpose fields (moved from provider-specific data)
    internal_notes TEXT,
    default_currency_id VARCHAR(10),
    contact_address JSONB,
    billing_info JSONB,
    labels JSONB,
    invoice_amount_threshold DECIMAL(15,2),
    default_payment_terms INTEGER,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_affiliates_timestamp
BEFORE UPDATE ON public.affiliates
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Core affiliate indexes
CREATE INDEX idx_affiliates_organization_id ON public.affiliates(organization_id);

-- General purpose field indexes for performance optimization
CREATE INDEX idx_affiliates_default_currency_id 
ON public.affiliates(default_currency_id) 
WHERE default_currency_id IS NOT NULL;

CREATE INDEX idx_affiliates_invoice_amount_threshold 
ON public.affiliates(invoice_amount_threshold) 
WHERE invoice_amount_threshold IS NOT NULL;

CREATE INDEX idx_affiliates_default_payment_terms 
ON public.affiliates(default_payment_terms) 
WHERE default_payment_terms IS NOT NULL;

-- GIN indexes for JSONB fields to enable efficient querying
CREATE INDEX idx_affiliates_contact_address_gin 
ON public.affiliates USING GIN(contact_address) 
WHERE contact_address IS NOT NULL;

CREATE INDEX idx_affiliates_billing_info_gin 
ON public.affiliates USING GIN(billing_info) 
WHERE billing_info IS NOT NULL;

CREATE INDEX idx_affiliates_labels_gin 
ON public.affiliates USING GIN(labels) 
WHERE labels IS NOT NULL;

-- Add column comments for documentation
COMMENT ON COLUMN public.affiliates.internal_notes IS 'Internal notes about the affiliate for team reference';
COMMENT ON COLUMN public.affiliates.default_currency_id IS 'Default currency code for affiliate transactions (e.g., USD, EUR)';
COMMENT ON COLUMN public.affiliates.contact_address IS 'Contact address information stored as JSONB';
COMMENT ON COLUMN public.affiliates.billing_info IS 'Billing information and preferences stored as JSONB';
COMMENT ON COLUMN public.affiliates.labels IS 'Array of labels/tags for categorizing affiliates stored as JSONB';
COMMENT ON COLUMN public.affiliates.invoice_amount_threshold IS 'Minimum amount threshold for automatic invoice generation';
COMMENT ON COLUMN public.affiliates.default_payment_terms IS 'Default payment terms in days';

-- #############################################################################
-- ## Provider Mapping Tables (All Provider-Specific Data)
-- #############################################################################

-- advertiser_provider_mappings: Maps platform advertisers to external providers
CREATE TABLE public.advertiser_provider_mappings (
    mapping_id BIGSERIAL PRIMARY KEY,
    advertiser_id BIGINT NOT NULL REFERENCES public.advertisers(advertiser_id) ON DELETE CASCADE,
    provider_type VARCHAR(50) NOT NULL DEFAULT 'everflow' CHECK (provider_type IN ('everflow')),
    provider_advertiser_id VARCHAR(255), -- e.g., Everflow's network_advertiser_id
    api_credentials JSONB, -- Store ENCRYPTED API Key, Secret, Token
    provider_config JSONB, -- Store provider-specific config
    provider_data JSONB, -- Store all provider-specific fields (Everflow fields, etc.)
    
    -- Synchronization metadata
    sync_status VARCHAR(50) CHECK (sync_status IS NULL OR sync_status IN ('pending', 'synced', 'failed', 'out_of_sync')),
    last_sync_at TIMESTAMPTZ,
    sync_error TEXT,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (advertiser_id, provider_type)
);

CREATE TRIGGER set_adv_prov_map_timestamp
BEFORE UPDATE ON public.advertiser_provider_mappings
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_adv_prov_map_advertiser ON public.advertiser_provider_mappings(advertiser_id);
CREATE INDEX idx_adv_prov_map_sync_status ON public.advertiser_provider_mappings(sync_status);
CREATE INDEX idx_adv_prov_map_provider_data_gin ON public.advertiser_provider_mappings USING GIN (provider_data);

-- affiliate_provider_mappings: Maps platform affiliates to external providers
-- Note: provider_data now contains only provider-specific fields (not general purpose fields)
CREATE TABLE public.affiliate_provider_mappings (
    mapping_id BIGSERIAL PRIMARY KEY,
    affiliate_id BIGINT NOT NULL REFERENCES public.affiliates(affiliate_id) ON DELETE CASCADE,
    provider_type VARCHAR(50) NOT NULL DEFAULT 'everflow' CHECK (provider_type IN ('everflow')),
    provider_affiliate_id VARCHAR(255), -- e.g., Everflow's network_affiliate_id
    api_credentials JSONB, -- Store ENCRYPTED API Key, Secret, Token (if needed per affiliate)
    provider_config JSONB, -- Store provider-specific config
    provider_data JSONB, -- Store ONLY provider-specific fields (Everflow-specific fields only)
    
    -- Synchronization metadata
    sync_status VARCHAR(50) CHECK (sync_status IS NULL OR sync_status IN ('pending', 'synced', 'failed', 'out_of_sync')),
    last_sync_at TIMESTAMPTZ,
    sync_error TEXT,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (affiliate_id, provider_type)
);

CREATE TRIGGER set_aff_prov_map_timestamp
BEFORE UPDATE ON public.affiliate_provider_mappings
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_aff_prov_map_affiliate ON public.affiliate_provider_mappings(affiliate_id);
CREATE INDEX idx_aff_prov_map_sync_status ON public.affiliate_provider_mappings(sync_status);
CREATE INDEX idx_aff_prov_map_provider_data_gin ON public.affiliate_provider_mappings USING GIN (provider_data);

-- #############################################################################
-- ## Campaign Tables
-- #############################################################################

-- campaigns: Platform-level campaigns with core fields only
CREATE TABLE public.campaigns (
    campaign_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    advertiser_id BIGINT NOT NULL REFERENCES public.advertisers(advertiser_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date DATE,
    end_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'paused', 'archived')),
    
    -- Core campaign fields
    destination_url TEXT,
    thumbnail_url TEXT,
    preview_url TEXT,
    visibility VARCHAR(50) DEFAULT 'public' CHECK (visibility IN ('public', 'require_approval', 'private')),
    currency_id VARCHAR(10) DEFAULT 'USD',
    
    -- Payout and revenue configuration
    payout_type VARCHAR(20) DEFAULT 'cpa' CHECK (payout_type IN ('cpa', 'cpc', 'cpm', 'cps', 'cpa_cps', 'prv')),
    payout_amount DECIMAL(10,2) DEFAULT 1.00,
    revenue_type VARCHAR(20) DEFAULT 'rpa' CHECK (revenue_type IN ('rpa', 'rpc', 'rpm', 'rps', 'rpa_rps', 'prv')),
    revenue_amount DECIMAL(10,2) DEFAULT 2.00,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_campaigns_timestamp
BEFORE UPDATE ON public.campaigns
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_campaigns_organization_id ON public.campaigns(organization_id);
CREATE INDEX idx_campaigns_advertiser_id ON public.campaigns(advertiser_id);
CREATE INDEX idx_campaigns_status_visibility ON public.campaigns(status, visibility);

-- campaign_provider_offers: Maps platform campaigns to provider offers
CREATE TABLE public.campaign_provider_offers (
    campaign_provider_offer_id BIGSERIAL PRIMARY KEY,
    campaign_id BIGINT NOT NULL REFERENCES public.campaigns(campaign_id) ON DELETE CASCADE,
    provider_type VARCHAR(50) NOT NULL DEFAULT 'everflow' CHECK (provider_type IN ('everflow')),
    provider_offer_ref VARCHAR(255), -- Provider's Offer ID
    provider_offer_config JSONB, -- All provider-specific offer configuration
    is_active_on_provider BOOLEAN DEFAULT FALSE,
    last_synced_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_cp_offers_timestamp
BEFORE UPDATE ON public.campaign_provider_offers
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_cp_offers_campaign_id ON public.campaign_provider_offers(campaign_id);
CREATE INDEX idx_cp_offers_provider_ref ON public.campaign_provider_offers(provider_offer_ref);

-- #############################################################################
-- ## No Data Migration Required
-- ## This is a fresh schema initialization - no existing data to migrate
-- #############################################################################
-- #############################################################################
-- ## Initial Seed Data
-- #############################################################################

-- Seed initial roles
INSERT INTO public.roles (role_id, name, description) VALUES
  (100000, 'User', 'Default user role with limited access'),
  (1000, 'AdvertiserManager', 'Manages advertisers and their campaigns within their organization'),
  (1001, 'AffiliateManager', 'Manages affiliates and approves applications within their organization'),
  (1, 'Admin', 'Platform Administrator with full access');

-- Create default organization
INSERT INTO public.organizations (name, type) VALUES ('rolinko', 'platform_owner');

-- #############################################################################
-- ## Final Optimization and Validation
-- #############################################################################

-- Analyze tables for query planner optimization
ANALYZE public.organizations;
ANALYZE public.roles;
ANALYZE public.profiles;
ANALYZE public.advertisers;
ANALYZE public.affiliates;
ANALYZE public.advertiser_provider_mappings;
ANALYZE public.affiliate_provider_mappings;
ANALYZE public.campaigns;
ANALYZE public.campaign_provider_offers;

-- Add schema tracking comment
COMMENT ON SCHEMA public IS 'Affiliate Platform Schema - Initialized with refactored affiliate model';

-- Schema initialization complete