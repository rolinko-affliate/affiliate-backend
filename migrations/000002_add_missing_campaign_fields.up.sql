-- Add missing fields to campaigns table to match domain model
-- This migration adds fields that are expected by the domain model but missing from the database schema

ALTER TABLE public.campaigns 
ADD COLUMN internal_notes TEXT,
ADD COLUMN conversion_method VARCHAR(50),
ADD COLUMN session_definition VARCHAR(50),
ADD COLUMN session_duration INTEGER, -- in hours
ADD COLUMN terms_and_conditions TEXT,
ADD COLUMN is_caps_enabled BOOLEAN DEFAULT FALSE,
ADD COLUMN daily_conversion_cap INTEGER,
ADD COLUMN weekly_conversion_cap INTEGER,
ADD COLUMN monthly_conversion_cap INTEGER,
ADD COLUMN global_conversion_cap INTEGER,
ADD COLUMN daily_click_cap INTEGER,
ADD COLUMN weekly_click_cap INTEGER,
ADD COLUMN monthly_click_cap INTEGER,
ADD COLUMN global_click_cap INTEGER;

-- Add check constraints for the new fields
ALTER TABLE public.campaigns 
ADD CONSTRAINT check_conversion_method 
CHECK (conversion_method IS NULL OR conversion_method IN ('server_postback', 'pixel', 'hybrid'));

ALTER TABLE public.campaigns 
ADD CONSTRAINT check_session_definition 
CHECK (session_definition IS NULL OR session_definition IN ('cookie', 'ip', 'fingerprint'));

-- Add comments for documentation
COMMENT ON COLUMN public.campaigns.internal_notes IS 'Internal notes for campaign management';
COMMENT ON COLUMN public.campaigns.conversion_method IS 'Method used for tracking conversions: server_postback, pixel, or hybrid';
COMMENT ON COLUMN public.campaigns.session_definition IS 'How sessions are defined: cookie, ip, or fingerprint';
COMMENT ON COLUMN public.campaigns.session_duration IS 'Session duration in hours';
COMMENT ON COLUMN public.campaigns.terms_and_conditions IS 'Campaign-specific terms and conditions';
COMMENT ON COLUMN public.campaigns.is_caps_enabled IS 'Whether caps are enabled for this campaign';
COMMENT ON COLUMN public.campaigns.daily_conversion_cap IS 'Daily conversion limit';
COMMENT ON COLUMN public.campaigns.weekly_conversion_cap IS 'Weekly conversion limit';
COMMENT ON COLUMN public.campaigns.monthly_conversion_cap IS 'Monthly conversion limit';
COMMENT ON COLUMN public.campaigns.global_conversion_cap IS 'Global conversion limit';
COMMENT ON COLUMN public.campaigns.daily_click_cap IS 'Daily click limit';
COMMENT ON COLUMN public.campaigns.weekly_click_cap IS 'Weekly click limit';
COMMENT ON COLUMN public.campaigns.monthly_click_cap IS 'Monthly click limit';
COMMENT ON COLUMN public.campaigns.global_click_cap IS 'Global click limit';