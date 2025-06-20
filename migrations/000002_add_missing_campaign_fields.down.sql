-- Rollback migration: Remove the added campaign fields

ALTER TABLE public.campaigns 
DROP COLUMN IF EXISTS internal_notes,
DROP COLUMN IF EXISTS conversion_method,
DROP COLUMN IF EXISTS session_definition,
DROP COLUMN IF EXISTS session_duration,
DROP COLUMN IF EXISTS terms_and_conditions,
DROP COLUMN IF EXISTS is_caps_enabled,
DROP COLUMN IF EXISTS daily_conversion_cap,
DROP COLUMN IF EXISTS weekly_conversion_cap,
DROP COLUMN IF EXISTS monthly_conversion_cap,
DROP COLUMN IF EXISTS global_conversion_cap,
DROP COLUMN IF EXISTS daily_click_cap,
DROP COLUMN IF EXISTS weekly_click_cap,
DROP COLUMN IF EXISTS monthly_click_cap,
DROP COLUMN IF EXISTS global_click_cap;