-- Remove company_size field from advertiser_extra_info table
DROP INDEX IF EXISTS idx_advertiser_extra_info_company_size;
ALTER TABLE advertiser_extra_info DROP COLUMN IF EXISTS company_size;