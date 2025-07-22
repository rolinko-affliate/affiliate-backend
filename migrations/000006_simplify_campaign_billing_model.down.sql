-- Rollback simplify campaign billing model
-- Restore the original complex billing model structure

-- Add back the old complex billing fields
ALTER TABLE public.campaigns 
ADD COLUMN billing_model VARCHAR(20) DEFAULT 'click' CHECK (billing_model IN ('click', 'conversion')),
ADD COLUMN payout_structure VARCHAR(20) DEFAULT 'fixed' CHECK (payout_structure IN ('fixed', 'percentage')),
ADD COLUMN payout_amount DECIMAL(10,2) DEFAULT 1.00,
ADD COLUMN revenue_structure VARCHAR(20) DEFAULT 'fixed' CHECK (revenue_structure IN ('fixed', 'percentage')),
ADD COLUMN revenue_amount DECIMAL(10,2) DEFAULT 2.00;

-- Migrate data back from simplified to complex structure
-- NOTE: This is a best-effort reverse migration. Some data combinations may not be perfectly restored.

-- Map fixed_revenue → revenue_structure = 'fixed', revenue_amount = fixed_revenue
UPDATE public.campaigns 
SET revenue_structure = 'fixed', revenue_amount = fixed_revenue 
WHERE fixed_revenue IS NOT NULL;

-- Map fixed_click_amount → billing_model = 'click', payout_structure = 'fixed', payout_amount = fixed_click_amount
UPDATE public.campaigns 
SET billing_model = 'click', payout_structure = 'fixed', payout_amount = fixed_click_amount 
WHERE fixed_click_amount IS NOT NULL;

-- Map fixed_conversion_amount → billing_model = 'conversion', payout_structure = 'fixed', payout_amount = fixed_conversion_amount
UPDATE public.campaigns 
SET billing_model = 'conversion', payout_structure = 'fixed', payout_amount = fixed_conversion_amount 
WHERE fixed_conversion_amount IS NOT NULL;

-- Map percentage_conversion_amount → billing_model = 'conversion', payout_structure = 'percentage', payout_amount = percentage_conversion_amount
UPDATE public.campaigns 
SET billing_model = 'conversion', payout_structure = 'percentage', payout_amount = percentage_conversion_amount 
WHERE percentage_conversion_amount IS NOT NULL;

-- Drop the simplified billing fields
ALTER TABLE public.campaigns 
DROP COLUMN fixed_revenue,
DROP COLUMN fixed_click_amount,
DROP COLUMN fixed_conversion_amount,
DROP COLUMN percentage_conversion_amount;

-- Restore comments for the old fields
COMMENT ON COLUMN public.campaigns.billing_model IS 'How we charge advertisers: click (per click) or conversion (per conversion)';
COMMENT ON COLUMN public.campaigns.payout_structure IS 'How we pay affiliates: fixed (fixed amount) or percentage (percentage of revenue)';
COMMENT ON COLUMN public.campaigns.payout_amount IS 'Amount to pay affiliates - currency amount for fixed, percentage value for percentage';
COMMENT ON COLUMN public.campaigns.revenue_structure IS 'How we calculate our revenue: fixed (fixed amount) or percentage (percentage of advertiser payment)';
COMMENT ON COLUMN public.campaigns.revenue_amount IS 'Our revenue amount - currency amount for fixed, percentage value for percentage';