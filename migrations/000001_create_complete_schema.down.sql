-- #############################################################################
-- ## Drop Complete Database Schema
-- ## Reverses all table creation and setup
-- #############################################################################

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS public.campaign_provider_offers CASCADE;
DROP TABLE IF EXISTS public.campaigns CASCADE;
DROP TABLE IF EXISTS public.affiliate_provider_mappings CASCADE;
DROP TABLE IF EXISTS public.advertiser_provider_mappings CASCADE;
DROP TABLE IF EXISTS public.affiliates CASCADE;
DROP TABLE IF EXISTS public.advertisers CASCADE;
DROP TABLE IF EXISTS public.profiles CASCADE;
DROP TABLE IF EXISTS public.roles CASCADE;
DROP TABLE IF EXISTS public.organizations CASCADE;

-- Drop helper functions
DROP FUNCTION IF EXISTS trigger_set_timestamp() CASCADE;