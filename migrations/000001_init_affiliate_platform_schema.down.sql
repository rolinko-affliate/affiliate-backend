-- #############################################################################
-- ## Complete Affiliate Platform Database Schema Rollback
-- ## This migration completely removes the affiliate platform schema
-- ## WARNING: This will destroy ALL data in the platform
-- #############################################################################

-- Remove schema tracking comment
COMMENT ON SCHEMA public IS NULL;

-- #############################################################################
-- ## Drop Tables (in reverse dependency order)
-- #############################################################################

-- Drop campaign-related tables
DROP TABLE IF EXISTS public.campaign_provider_offers CASCADE;
DROP TABLE IF EXISTS public.campaigns CASCADE;

-- Drop provider mapping tables
DROP TABLE IF EXISTS public.affiliate_provider_mappings CASCADE;
DROP TABLE IF EXISTS public.advertiser_provider_mappings CASCADE;

-- Drop core domain tables
DROP TABLE IF EXISTS public.affiliates CASCADE;
DROP TABLE IF EXISTS public.advertisers CASCADE;

-- Drop platform tables
DROP TABLE IF EXISTS public.profiles CASCADE;
DROP TABLE IF EXISTS public.roles CASCADE;
DROP TABLE IF EXISTS public.organizations CASCADE;

-- #############################################################################
-- ## Drop Helper Functions
-- #############################################################################

DROP FUNCTION IF EXISTS trigger_set_timestamp() CASCADE;

-- #############################################################################
-- ## Final Cleanup
-- #############################################################################

-- Note: This migration completely removes the affiliate platform schema
-- All data will be lost and cannot be recovered without a backup

DO $$
BEGIN
    RAISE NOTICE '=== Affiliate Platform Schema Rollback Complete ===';
    RAISE NOTICE 'All tables, indexes, triggers, and functions have been removed';
    RAISE NOTICE 'All data has been permanently deleted';
    RAISE NOTICE '=== Database Reset to Clean State ===';
END $$;