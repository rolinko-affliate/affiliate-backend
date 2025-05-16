-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS public.campaign_provider_offers;
DROP TABLE IF EXISTS public.campaigns;
DROP TABLE IF EXISTS public.affiliate_provider_mappings;
DROP TABLE IF EXISTS public.affiliates;
DROP TABLE IF EXISTS public.advertiser_provider_mappings;
DROP TABLE IF EXISTS public.advertisers;
DROP TABLE IF EXISTS public.profiles;
DROP TABLE IF EXISTS public.roles;
DROP TABLE IF EXISTS public.organizations;

-- Drop the trigger function
DROP FUNCTION IF EXISTS trigger_set_timestamp();