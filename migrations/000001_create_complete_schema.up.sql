-- #############################################################################
-- ## Complete Database Schema for Affiliate Platform
-- ## Includes all tables with Everflow integration support
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
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_organizations_timestamp
BEFORE UPDATE ON public.organizations
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

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
-- ## Advertiser Tables with Everflow Support
-- #############################################################################

-- advertisers: Stores advertiser information with Everflow integration fields
CREATE TABLE public.advertisers (
    advertiser_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255),
    billing_details JSONB, -- Store address, tax ID, etc.
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'inactive', 'rejected')),
    
    -- Everflow-specific fields
    internal_notes TEXT,
    default_currency_id VARCHAR(10),
    platform_name VARCHAR(255),
    platform_url VARCHAR(500),
    platform_username VARCHAR(255),
    accounting_contact_email VARCHAR(255),
    offer_id_macro VARCHAR(255),
    affiliate_id_macro VARCHAR(255),
    attribution_method VARCHAR(100),
    email_attribution_method VARCHAR(100),
    attribution_priority VARCHAR(100),
    reporting_timezone_id INTEGER,
    is_expose_publisher_reporting BOOLEAN,
    
    -- Everflow synchronization metadata
    everflow_sync_status VARCHAR(50) CHECK (everflow_sync_status IS NULL OR everflow_sync_status IN ('pending', 'synced', 'failed', 'out_of_sync')),
    last_everflow_sync_at TIMESTAMPTZ,
    everflow_sync_error TEXT,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_advertisers_timestamp
BEFORE UPDATE ON public.advertisers
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_advertisers_organization_id ON public.advertisers(organization_id);
CREATE INDEX idx_advertisers_sync_status ON public.advertisers(everflow_sync_status);
CREATE INDEX idx_advertisers_last_sync ON public.advertisers(last_everflow_sync_at);

-- #############################################################################
-- ## Affiliate Tables
-- #############################################################################

-- affiliates: Stores affiliate information, linked to an `organization`.
CREATE TABLE public.affiliates (
    affiliate_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255),
    payment_details JSONB, -- Store payment method (PayPal, Wise, Bank) and details
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'rejected', 'inactive')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_affiliates_timestamp
BEFORE UPDATE ON public.affiliates
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_affiliates_organization_id ON public.affiliates(organization_id);

-- #############################################################################
-- ## Provider Mapping Tables
-- #############################################################################

-- advertiser_provider_mappings: Maps platform `advertisers` to their corresponding IDs and API credentials on external provider systems like Everflow.
CREATE TABLE public.advertiser_provider_mappings (
    mapping_id BIGSERIAL PRIMARY KEY,
    advertiser_id BIGINT NOT NULL REFERENCES public.advertisers(advertiser_id) ON DELETE CASCADE,
    provider_type VARCHAR(50) NOT NULL DEFAULT 'everflow' CHECK (provider_type IN ('everflow')), -- MVP focuses on Everflow
    provider_advertiser_id VARCHAR(255), -- e.g., Everflow's network_advertiser_id
    api_credentials JSONB, -- Store ENCRYPTED API Key, Secret, Token for Everflow
    provider_config JSONB, -- Store other provider-specific config e.g. Everflow network_id
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (advertiser_id, provider_type)
);

CREATE TRIGGER set_adv_prov_map_timestamp
BEFORE UPDATE ON public.advertiser_provider_mappings
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_adv_prov_map_advertiser ON public.advertiser_provider_mappings(advertiser_id);

-- affiliate_provider_mappings: Maps platform `affiliates` to their corresponding IDs on external provider systems like Everflow.
CREATE TABLE public.affiliate_provider_mappings (
    mapping_id BIGSERIAL PRIMARY KEY,
    affiliate_id BIGINT NOT NULL REFERENCES public.affiliates(affiliate_id) ON DELETE CASCADE,
    provider_type VARCHAR(50) NOT NULL DEFAULT 'everflow' CHECK (provider_type IN ('everflow')), -- MVP focuses on Everflow
    provider_affiliate_id VARCHAR(255), -- e.g., Everflow's network_affiliate_id
    provider_config JSONB, -- Store other provider-specific config
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (affiliate_id, provider_type)
);

CREATE TRIGGER set_aff_prov_map_timestamp
BEFORE UPDATE ON public.affiliate_provider_mappings
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_aff_prov_map_affiliate ON public.affiliate_provider_mappings(affiliate_id);

-- #############################################################################
-- ## Campaign Tables with Complete Everflow Offer Support
-- #############################################################################

-- campaigns: Platform-level campaigns with complete Everflow offer fields
CREATE TABLE public.campaigns (
    campaign_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    advertiser_id BIGINT NOT NULL REFERENCES public.advertisers(advertiser_id) ON DELETE CASCADE, -- The advertiser for whom this campaign is run
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date DATE,
    end_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'paused', 'archived')),
    
    -- Core Everflow offer fields
    destination_url TEXT,
    thumbnail_url TEXT,
    preview_url TEXT,
    visibility VARCHAR(50) DEFAULT 'public' CHECK (visibility IN ('public', 'require_approval', 'private')),
    currency_id VARCHAR(10) DEFAULT 'USD',
    conversion_method VARCHAR(50) DEFAULT 'server_postback',
    session_definition VARCHAR(50) DEFAULT 'cookie' CHECK (session_definition IN ('cookie', 'ip', 'ip_user_agent', 'google_ad_id', 'idfa')),
    session_duration INTEGER DEFAULT 24, -- in hours
    internal_notes TEXT,
    terms_and_conditions TEXT,
    is_force_terms_and_conditions BOOLEAN DEFAULT FALSE,
    
    -- Additional Everflow offer fields
    caps_timezone_id INTEGER,
    project_id VARCHAR(255),
    date_live_until DATE,
    html_description TEXT,
    is_using_explicit_terms_and_conditions BOOLEAN,
    is_whitelist_check_enabled BOOLEAN,
    is_view_through_enabled BOOLEAN,
    server_side_url TEXT,
    view_through_destination_url TEXT,
    is_description_plain_text BOOLEAN,
    is_use_direct_linking BOOLEAN,
    app_identifier VARCHAR(255),
    
    -- Caps and limits
    is_caps_enabled BOOLEAN DEFAULT FALSE,
    daily_conversion_cap INTEGER,
    weekly_conversion_cap INTEGER,
    monthly_conversion_cap INTEGER,
    global_conversion_cap INTEGER,
    daily_click_cap INTEGER,
    weekly_click_cap INTEGER,
    monthly_click_cap INTEGER,
    global_click_cap INTEGER,
    
    -- Everflow tracking fields
    encoded_value VARCHAR(255),
    today_clicks INTEGER,
    today_revenue VARCHAR(50),
    time_created INTEGER, -- Unix timestamp
    time_saved INTEGER,   -- Unix timestamp
    
    -- Payout and revenue configuration
    payout_type VARCHAR(20) DEFAULT 'cpa' CHECK (payout_type IN ('cpa', 'cpc', 'cpm', 'cps', 'cpa_cps', 'prv')),
    payout_amount DECIMAL(10,2) DEFAULT 1.00,
    revenue_type VARCHAR(20) DEFAULT 'rpa' CHECK (revenue_type IN ('rpa', 'rpc', 'rpm', 'rps', 'rpa_rps', 'prv')),
    revenue_amount DECIMAL(10,2) DEFAULT 2.00,
    
    -- Additional configuration stored as JSON
    offer_config JSONB, -- Store additional Everflow-specific configuration
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_campaigns_timestamp
BEFORE UPDATE ON public.campaigns
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_campaigns_organization_id ON public.campaigns(organization_id);
CREATE INDEX idx_campaigns_advertiser_id ON public.campaigns(advertiser_id);
CREATE INDEX idx_campaigns_visibility ON public.campaigns(visibility);
CREATE INDEX idx_campaigns_currency_id ON public.campaigns(currency_id);
CREATE INDEX idx_campaigns_status_visibility ON public.campaigns(status, visibility);
CREATE INDEX idx_campaigns_offer_config_gin ON public.campaigns USING GIN (offer_config);

-- campaign_provider_offers: Maps platform `campaigns` to specific offers on provider systems (e.g., an Everflow Offer).
CREATE TABLE public.campaign_provider_offers (
    campaign_provider_offer_id BIGSERIAL PRIMARY KEY, -- Updated column name for consistency
    campaign_id BIGINT NOT NULL REFERENCES public.campaigns(campaign_id) ON DELETE CASCADE,
    provider_type VARCHAR(50) NOT NULL DEFAULT 'everflow' CHECK (provider_type IN ('everflow')),
    provider_offer_ref VARCHAR(255), -- Provider's Offer ID (e.g., Everflow's network_offer_id)
    provider_offer_config JSONB, -- Stores detailed Offer configuration for the provider
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
-- Optional GIN index for querying provider_offer_config if needed in the future
-- CREATE INDEX idx_cp_offers_config_gin ON public.campaign_provider_offers USING GIN (provider_offer_config);

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
INSERT INTO public.organizations (name) VALUES ('rolinko');