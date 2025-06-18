-- #############################################################################
-- ## Complete Affiliate Platform Database Schema Initialization
-- ## This migration creates the complete database schema for the affiliate platform
-- ## Combines all previous migrations into a single comprehensive schema
-- ## 
-- ## Includes:
-- ## - Core platform tables (organizations, roles, profiles)
-- ## - Domain entities (advertisers, affiliates, campaigns)
-- ## - Provider mapping tables for external integrations
-- ## - Analytics tables for advertiser and publisher data
-- ## - Tracking links and provider mappings
-- ## - Extended fields and updated payout models
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

-- advertisers: Core advertiser information with extended fields
CREATE TABLE public.advertisers (
    advertiser_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255),
    billing_details JSONB, -- Store address, tax ID, etc.
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'inactive', 'rejected')),
    
    -- Extended fields from migration 000003
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
-- ## Campaign Tables with Updated Payout Model
-- #############################################################################

-- campaigns: Platform-level campaigns with updated payout/revenue model
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
    
    -- Updated payout and revenue configuration (from migration 000004)
    billing_model VARCHAR(20) DEFAULT 'click' CHECK (billing_model IN ('click', 'conversion')),
    payout_structure VARCHAR(20) DEFAULT 'fixed' CHECK (payout_structure IN ('fixed', 'percentage')),
    payout_amount DECIMAL(10,2) DEFAULT 1.00,
    revenue_structure VARCHAR(20) DEFAULT 'fixed' CHECK (revenue_structure IN ('fixed', 'percentage')),
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

-- campaign_provider_mappings: Maps platform campaigns to provider offers
CREATE TABLE public.campaign_provider_mappings (
    mapping_id BIGSERIAL PRIMARY KEY,
    campaign_id BIGINT NOT NULL REFERENCES public.campaigns(campaign_id) ON DELETE CASCADE,
    provider_type VARCHAR(50) NOT NULL DEFAULT 'everflow' CHECK (provider_type IN ('everflow')),
    provider_offer_id VARCHAR(255), -- Provider's Offer ID
    provider_config JSONB, -- All provider-specific offer configuration
    is_active_on_provider BOOLEAN DEFAULT FALSE,
    last_synced_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (campaign_id, provider_type)
);

CREATE TRIGGER set_campaign_prov_map_timestamp
BEFORE UPDATE ON public.campaign_provider_mappings
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_campaign_prov_map_campaign_id ON public.campaign_provider_mappings(campaign_id);
CREATE INDEX idx_campaign_prov_map_provider_offer_id ON public.campaign_provider_mappings(provider_offer_id);

-- #############################################################################
-- ## Tracking Links Tables
-- #############################################################################

-- tracking_links: Core tracking link information (provider-agnostic)
CREATE TABLE public.tracking_links (
    tracking_link_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    campaign_id BIGINT NOT NULL REFERENCES public.campaigns(campaign_id) ON DELETE CASCADE,
    affiliate_id BIGINT NOT NULL REFERENCES public.affiliates(affiliate_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paused', 'archived')),
    
    -- Core tracking link fields (provider-agnostic)
    tracking_url TEXT,
    source_id VARCHAR(255),
    sub1 VARCHAR(255),
    sub2 VARCHAR(255),
    sub3 VARCHAR(255),
    sub4 VARCHAR(255),
    sub5 VARCHAR(255),
    
    -- Link configuration
    is_encrypt_parameters BOOLEAN DEFAULT FALSE,
    is_redirect_link BOOLEAN DEFAULT FALSE,
    
    -- General purpose fields
    internal_notes TEXT,
    tags JSONB, -- Array of tags/labels for categorizing tracking links
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Ensure unique tracking link per campaign-affiliate combination with same parameters
    UNIQUE (campaign_id, affiliate_id, source_id, sub1, sub2, sub3, sub4, sub5)
);

-- Add trigger for automatic updated_at timestamp
CREATE TRIGGER set_tracking_links_timestamp
BEFORE UPDATE ON public.tracking_links
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Core indexes for performance
CREATE INDEX idx_tracking_links_organization_id ON public.tracking_links(organization_id);
CREATE INDEX idx_tracking_links_campaign_id ON public.tracking_links(campaign_id);
CREATE INDEX idx_tracking_links_affiliate_id ON public.tracking_links(affiliate_id);
CREATE INDEX idx_tracking_links_status ON public.tracking_links(status);
CREATE INDEX idx_tracking_links_created_at ON public.tracking_links(created_at);

-- Indexes for tracking parameters (for efficient filtering)
CREATE INDEX idx_tracking_links_source_id ON public.tracking_links(source_id) WHERE source_id IS NOT NULL;
CREATE INDEX idx_tracking_links_sub1 ON public.tracking_links(sub1) WHERE sub1 IS NOT NULL;
CREATE INDEX idx_tracking_links_sub2 ON public.tracking_links(sub2) WHERE sub2 IS NOT NULL;
CREATE INDEX idx_tracking_links_sub3 ON public.tracking_links(sub3) WHERE sub3 IS NOT NULL;

-- GIN index for JSONB tags field
CREATE INDEX idx_tracking_links_tags_gin ON public.tracking_links USING GIN(tags) WHERE tags IS NOT NULL;

-- tracking_link_provider_mappings: Maps platform tracking links to external providers
CREATE TABLE public.tracking_link_provider_mappings (
    mapping_id BIGSERIAL PRIMARY KEY,
    tracking_link_id BIGINT NOT NULL REFERENCES public.tracking_links(tracking_link_id) ON DELETE CASCADE,
    provider_type VARCHAR(50) NOT NULL DEFAULT 'everflow' CHECK (provider_type IN ('everflow')),
    provider_tracking_link_id VARCHAR(255), -- Provider's tracking link identifier (if any)
    
    -- Provider-specific data stored as JSONB
    provider_data JSONB, -- Store all provider-specific fields (Everflow fields, etc.)
    
    -- Synchronization metadata
    sync_status VARCHAR(50) CHECK (sync_status IS NULL OR sync_status IN ('pending', 'synced', 'failed', 'out_of_sync')),
    last_sync_at TIMESTAMPTZ,
    sync_error TEXT,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Ensure unique mapping per tracking link and provider
    UNIQUE (tracking_link_id, provider_type)
);

-- Add trigger for automatic updated_at timestamp
CREATE TRIGGER set_tracking_link_prov_map_timestamp
BEFORE UPDATE ON public.tracking_link_provider_mappings
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Indexes for tracking link provider mappings
CREATE INDEX idx_tl_prov_map_tracking_link_id ON public.tracking_link_provider_mappings(tracking_link_id);
CREATE INDEX idx_tl_prov_map_provider_type ON public.tracking_link_provider_mappings(provider_type);
CREATE INDEX idx_tl_prov_map_sync_status ON public.tracking_link_provider_mappings(sync_status);
CREATE INDEX idx_tl_prov_map_provider_data_gin ON public.tracking_link_provider_mappings USING GIN (provider_data);

-- #############################################################################
-- ## Analytics Tables
-- #############################################################################

-- analytics_advertisers: Store advertiser analytics data
CREATE TABLE public.analytics_advertisers (
    id BIGSERIAL PRIMARY KEY,
    domain VARCHAR(255) NOT NULL UNIQUE,
    
    -- Metadata
    description TEXT,
    favicon_image_url TEXT,
    screenshot_image_url TEXT,
    
    -- Affiliate Networks
    affiliate_networks JSONB, -- Array of network names
    
    -- Contact Information
    contact_emails JSONB, -- Array of contact email objects
    
    -- Keywords
    keywords JSONB, -- Array of keyword objects with scores
    
    -- Verticals
    verticals JSONB, -- Array of vertical objects with names, ranks, scores
    
    -- Partner Information
    partner_information JSONB, -- Complex partner data structure
    
    -- Related Advertisers
    related_advertisers JSONB, -- Array of related advertiser domains
    
    -- Social Media
    social_media JSONB, -- Social media links and availability
    
    -- Backlinks
    backlinks JSONB, -- Backlink information
    
    -- Additional analytics data
    additional_data JSONB, -- For any extra fields from the API
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- analytics_publishers: Store publisher analytics data  
CREATE TABLE public.analytics_publishers (
    id BIGSERIAL PRIMARY KEY,
    domain VARCHAR(255) NOT NULL UNIQUE,
    
    -- Metadata
    description TEXT,
    favicon_image_url TEXT,
    screenshot_image_url TEXT,
    
    -- Affiliate Networks
    affiliate_networks JSONB, -- Array of network names
    
    -- Country Rankings
    country_rankings JSONB, -- Array of country ranking objects
    
    -- Keywords
    keywords JSONB, -- Array of keyword objects with scores
    
    -- Verticals
    verticals JSONB, -- Array of vertical names
    verticals_v2 JSONB, -- Array of vertical objects with names, ranks, scores
    
    -- Partner Information
    partner_information JSONB, -- Complex partner data structure
    partners JSONB, -- Array of partner domains
    
    -- Related Publishers
    related_publishers JSONB, -- Array of related publisher domains
    
    -- Social Media
    social_media JSONB, -- Social media links and availability
    
    -- Live URLs
    live_urls JSONB, -- Array of live URLs
    
    -- Flags and Scores
    known BOOLEAN DEFAULT FALSE,
    relevance DECIMAL(5,2) DEFAULT 0,
    traffic_score DECIMAL(10,2) DEFAULT 0,
    
    -- Promotype
    promotype VARCHAR(50),
    
    -- Additional analytics data
    additional_data JSONB, -- For any extra fields from the API
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create triggers for updated_at
CREATE TRIGGER set_analytics_advertisers_timestamp
BEFORE UPDATE ON public.analytics_advertisers
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_analytics_publishers_timestamp
BEFORE UPDATE ON public.analytics_publishers
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Create indexes for performance
CREATE INDEX idx_analytics_advertisers_domain ON public.analytics_advertisers(domain);
CREATE INDEX idx_analytics_advertisers_domain_text ON public.analytics_advertisers USING gin(to_tsvector('english', domain));

CREATE INDEX idx_analytics_publishers_domain ON public.analytics_publishers(domain);
CREATE INDEX idx_analytics_publishers_domain_text ON public.analytics_publishers USING gin(to_tsvector('english', domain));
CREATE INDEX idx_analytics_publishers_known ON public.analytics_publishers(known);
CREATE INDEX idx_analytics_publishers_relevance ON public.analytics_publishers(relevance);
CREATE INDEX idx_analytics_publishers_traffic_score ON public.analytics_publishers(traffic_score);

-- GIN indexes for JSONB fields to enable efficient querying
CREATE INDEX idx_analytics_advertisers_affiliate_networks_gin ON public.analytics_advertisers USING GIN(affiliate_networks);
CREATE INDEX idx_analytics_advertisers_keywords_gin ON public.analytics_advertisers USING GIN(keywords);
CREATE INDEX idx_analytics_advertisers_verticals_gin ON public.analytics_advertisers USING GIN(verticals);

CREATE INDEX idx_analytics_publishers_affiliate_networks_gin ON public.analytics_publishers USING GIN(affiliate_networks);
CREATE INDEX idx_analytics_publishers_keywords_gin ON public.analytics_publishers USING GIN(keywords);
CREATE INDEX idx_analytics_publishers_verticals_gin ON public.analytics_publishers USING GIN(verticals);
CREATE INDEX idx_analytics_publishers_verticals_v2_gin ON public.analytics_publishers USING GIN(verticals_v2);

-- #############################################################################
-- ## Add Comments for Documentation
-- #############################################################################

-- Core tables comments
COMMENT ON TABLE public.organizations IS 'Organizations represent distinct tenants or teams using the platform';
COMMENT ON TABLE public.roles IS 'User roles within the platform (Admin, AdvertiserManager, AffiliateManager, etc.)';
COMMENT ON TABLE public.profiles IS 'User profiles linked to Supabase Auth with organization and role associations';

-- Domain entity comments
COMMENT ON TABLE public.advertisers IS 'Core advertiser information with extended fields for platform management';
COMMENT ON TABLE public.affiliates IS 'Core affiliate information with general purpose fields following clean architecture';
COMMENT ON TABLE public.campaigns IS 'Platform campaigns with updated payout/revenue model';

-- Provider mapping comments
COMMENT ON TABLE public.advertiser_provider_mappings IS 'Maps platform advertisers to external providers (Everflow, etc.)';
COMMENT ON TABLE public.affiliate_provider_mappings IS 'Maps platform affiliates to external providers with provider-specific data';
COMMENT ON TABLE public.campaign_provider_mappings IS 'Maps platform campaigns to provider offers';

-- Tracking links comments
COMMENT ON TABLE public.tracking_links IS 'Core tracking links table following clean architecture - provider-agnostic fields only';
COMMENT ON TABLE public.tracking_link_provider_mappings IS 'Provider-specific tracking link mappings and data';

-- Analytics comments
COMMENT ON TABLE public.analytics_advertisers IS 'Analytics data for advertisers including metadata, networks, keywords, and partner information';
COMMENT ON TABLE public.analytics_publishers IS 'Analytics data for publishers including metadata, networks, keywords, and partner information';

-- Column comments for key fields
COMMENT ON COLUMN public.advertisers.internal_notes IS 'Internal notes about the advertiser for team reference';
COMMENT ON COLUMN public.advertisers.default_currency_id IS 'Default currency code for advertiser transactions (e.g., USD, EUR)';
COMMENT ON COLUMN public.advertisers.platform_name IS 'Name of the advertising platform';
COMMENT ON COLUMN public.advertisers.platform_url IS 'URL of the advertising platform';
COMMENT ON COLUMN public.advertisers.platform_username IS 'Username for the advertising platform';
COMMENT ON COLUMN public.advertisers.accounting_contact_email IS 'Email for accounting/billing contact';
COMMENT ON COLUMN public.advertisers.offer_id_macro IS 'Macro for offer ID tracking';
COMMENT ON COLUMN public.advertisers.affiliate_id_macro IS 'Macro for affiliate ID tracking';
COMMENT ON COLUMN public.advertisers.attribution_method IS 'Method used for attribution tracking';
COMMENT ON COLUMN public.advertisers.email_attribution_method IS 'Method used for email attribution tracking';
COMMENT ON COLUMN public.advertisers.attribution_priority IS 'Priority level for attribution';
COMMENT ON COLUMN public.advertisers.reporting_timezone_id IS 'Timezone ID for reporting purposes';

COMMENT ON COLUMN public.affiliates.internal_notes IS 'Internal notes about the affiliate for team reference';
COMMENT ON COLUMN public.affiliates.default_currency_id IS 'Default currency code for affiliate transactions (e.g., USD, EUR)';
COMMENT ON COLUMN public.affiliates.contact_address IS 'Contact address information stored as JSONB';
COMMENT ON COLUMN public.affiliates.billing_info IS 'Billing information and preferences stored as JSONB';
COMMENT ON COLUMN public.affiliates.labels IS 'Array of labels/tags for categorizing affiliates stored as JSONB';
COMMENT ON COLUMN public.affiliates.invoice_amount_threshold IS 'Minimum amount threshold for automatic invoice generation';
COMMENT ON COLUMN public.affiliates.default_payment_terms IS 'Default payment terms in days';

COMMENT ON COLUMN public.campaigns.billing_model IS 'How we charge advertisers: click (per click) or conversion (per conversion)';
COMMENT ON COLUMN public.campaigns.payout_structure IS 'How we pay affiliates: fixed (fixed amount) or percentage (percentage of revenue)';
COMMENT ON COLUMN public.campaigns.payout_amount IS 'Amount to pay affiliates - currency amount for fixed, percentage value for percentage';
COMMENT ON COLUMN public.campaigns.revenue_structure IS 'How we calculate our revenue: fixed (fixed amount) or percentage (percentage of advertiser payment)';
COMMENT ON COLUMN public.campaigns.revenue_amount IS 'Our revenue amount - currency amount for fixed, percentage value for percentage';

COMMENT ON COLUMN public.tracking_links.tracking_url IS 'Generated tracking URL from the provider';
COMMENT ON COLUMN public.tracking_links.source_id IS 'Source identifier for tracking attribution';
COMMENT ON COLUMN public.tracking_links.sub1 IS 'Sub ID 1 for custom tracking parameters';
COMMENT ON COLUMN public.tracking_links.sub2 IS 'Sub ID 2 for custom tracking parameters';
COMMENT ON COLUMN public.tracking_links.sub3 IS 'Sub ID 3 for custom tracking parameters';
COMMENT ON COLUMN public.tracking_links.sub4 IS 'Sub ID 4 for custom tracking parameters';
COMMENT ON COLUMN public.tracking_links.sub5 IS 'Sub ID 5 for custom tracking parameters';
COMMENT ON COLUMN public.tracking_links.is_encrypt_parameters IS 'Whether to encrypt query string parameters';
COMMENT ON COLUMN public.tracking_links.is_redirect_link IS 'Whether to use redirect link for direct linking offers';
COMMENT ON COLUMN public.tracking_links.tags IS 'Array of tags/labels for categorizing tracking links stored as JSONB';

COMMENT ON COLUMN public.tracking_link_provider_mappings.provider_data IS 'Provider-specific data stored as JSONB (e.g., Everflow network IDs, domain IDs, etc.)';
COMMENT ON COLUMN public.tracking_link_provider_mappings.sync_status IS 'Synchronization status with the provider';
COMMENT ON COLUMN public.tracking_link_provider_mappings.sync_error IS 'Error message if synchronization failed';

COMMENT ON COLUMN public.analytics_advertisers.domain IS 'Primary domain of the advertiser';
COMMENT ON COLUMN public.analytics_advertisers.affiliate_networks IS 'Array of affiliate network names the advertiser works with';
COMMENT ON COLUMN public.analytics_advertisers.keywords IS 'Array of keyword objects with scores for SEO/content analysis';
COMMENT ON COLUMN public.analytics_advertisers.verticals IS 'Array of vertical/industry categories with rankings and scores';

COMMENT ON COLUMN public.analytics_publishers.domain IS 'Primary domain of the publisher';
COMMENT ON COLUMN public.analytics_publishers.affiliate_networks IS 'Array of affiliate network names the publisher works with';
COMMENT ON COLUMN public.analytics_publishers.keywords IS 'Array of keyword objects with scores for content analysis';
COMMENT ON COLUMN public.analytics_publishers.verticals IS 'Array of vertical/industry category names';
COMMENT ON COLUMN public.analytics_publishers.verticals_v2 IS 'Enhanced vertical data with names, ranks, and scores';
COMMENT ON COLUMN public.analytics_publishers.known IS 'Flag indicating if this is a known/verified publisher';
COMMENT ON COLUMN public.analytics_publishers.relevance IS 'Relevance score for the publisher';
COMMENT ON COLUMN public.analytics_publishers.traffic_score IS 'Traffic score indicating publisher reach';

-- #############################################################################
-- ## Initial Seed Data
-- #############################################################################

-- Seed initial roles
INSERT INTO public.roles (role_id, name, description) VALUES
  (100000, 'User', 'Default user role with limited access'),
  (1000, 'AdvertiserManager', 'Manages advertisers and their campaigns within their organization'),
  (1001, 'AffiliateManager', 'Manages affiliates and approves applications within their organization'),
  (1, 'Admin', 'Platform Administrator with full access');

-- Create platform owner organization
INSERT INTO public.organizations (name, type) VALUES 
  ('rolinko', 'platform_owner');

-- Create admin user profile
INSERT INTO public.profiles (
    id, 
    organization_id, 
    role_id, 
    email, 
    first_name, 
    last_name
) VALUES (
    '43bad314-bdd3-49d3-9f85-5be4d019c2ae',
    (SELECT organization_id FROM public.organizations WHERE name = 'rolinko' LIMIT 1),
    1,
    'admin@rolinko.com',
    'Platform',
    'Administrator'
);

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
ANALYZE public.campaign_provider_mappings;
ANALYZE public.tracking_links;
ANALYZE public.tracking_link_provider_mappings;
ANALYZE public.analytics_advertisers;
ANALYZE public.analytics_publishers;

-- Add schema tracking comment
COMMENT ON SCHEMA public IS 'Complete Affiliate Platform Schema - All migrations consolidated into single comprehensive schema';

-- Schema initialization complete