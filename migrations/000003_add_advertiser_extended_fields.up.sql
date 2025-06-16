-- Add extended fields to advertisers table to match domain model
-- These fields were missing from the initial schema but are referenced in the domain model

ALTER TABLE public.advertisers 
ADD COLUMN internal_notes TEXT,
ADD COLUMN default_currency_id VARCHAR(10),
ADD COLUMN platform_name VARCHAR(255),
ADD COLUMN platform_url VARCHAR(500),
ADD COLUMN platform_username VARCHAR(255),
ADD COLUMN accounting_contact_email VARCHAR(255),
ADD COLUMN offer_id_macro VARCHAR(255),
ADD COLUMN affiliate_id_macro VARCHAR(255),
ADD COLUMN attribution_method VARCHAR(100),
ADD COLUMN email_attribution_method VARCHAR(100),
ADD COLUMN attribution_priority VARCHAR(100),
ADD COLUMN reporting_timezone_id INTEGER;

-- Add column comments for documentation
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