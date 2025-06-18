-- #############################################################################
-- ## Complete Affiliate Platform Database Schema Rollback
-- ## This migration drops the complete database schema for the affiliate platform
-- ## Reverses the consolidated schema initialization
-- #############################################################################

-- Drop all tables in reverse dependency order

-- Analytics tables (no dependencies)
DROP TABLE IF EXISTS public.analytics_publishers CASCADE;
DROP TABLE IF EXISTS public.analytics_advertisers CASCADE;

-- Provider mapping tables (depend on core domain tables)
DROP TABLE IF EXISTS public.tracking_link_provider_mappings CASCADE;
DROP TABLE IF EXISTS public.campaign_provider_mappings CASCADE;
DROP TABLE IF EXISTS public.affiliate_provider_mappings CASCADE;
DROP TABLE IF EXISTS public.advertiser_provider_mappings CASCADE;

-- Tracking links table (depends on campaigns, affiliates, organizations)
DROP TABLE IF EXISTS public.tracking_links CASCADE;

-- Campaign tables (depend on advertisers, organizations)
DROP TABLE IF EXISTS public.campaigns CASCADE;

-- Core domain tables (depend on organizations)
DROP TABLE IF EXISTS public.affiliates CASCADE;
DROP TABLE IF EXISTS public.advertisers CASCADE;

-- Profile table (depends on organizations, roles)
DROP TABLE IF EXISTS public.profiles CASCADE;

-- Core platform tables
DROP TABLE IF EXISTS public.roles CASCADE;
DROP TABLE IF EXISTS public.organizations CASCADE;

-- Drop helper functions
DROP FUNCTION IF EXISTS trigger_set_timestamp() CASCADE;

-- Remove schema comment
COMMENT ON SCHEMA public IS NULL;

-- Schema rollback complete