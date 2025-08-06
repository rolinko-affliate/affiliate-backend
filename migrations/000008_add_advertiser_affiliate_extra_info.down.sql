-- #############################################################################
-- ## Rollback Extra Information Tables for Advertisers and Affiliates
-- ## This migration removes the extra information tables
-- #############################################################################

-- Drop affiliate extra info table
DROP TABLE IF EXISTS public.affiliate_extra_info;

-- Drop advertiser extra info table  
DROP TABLE IF EXISTS public.advertiser_extra_info;