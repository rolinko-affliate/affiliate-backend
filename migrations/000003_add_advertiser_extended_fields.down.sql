-- Remove extended fields from advertisers table

ALTER TABLE public.advertisers 
DROP COLUMN IF EXISTS internal_notes,
DROP COLUMN IF EXISTS default_currency_id,
DROP COLUMN IF EXISTS platform_name,
DROP COLUMN IF EXISTS platform_url,
DROP COLUMN IF EXISTS platform_username,
DROP COLUMN IF EXISTS accounting_contact_email,
DROP COLUMN IF EXISTS offer_id_macro,
DROP COLUMN IF EXISTS affiliate_id_macro,
DROP COLUMN IF EXISTS attribution_method,
DROP COLUMN IF EXISTS email_attribution_method,
DROP COLUMN IF EXISTS attribution_priority,
DROP COLUMN IF EXISTS reporting_timezone_id;