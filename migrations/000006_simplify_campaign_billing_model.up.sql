-- Simplify campaign billing model
-- Replace complex billing_model + payout_structure + revenue_structure system
-- with simple fixed amounts and percentage fields

-- Add new simplified billing fields
ALTER TABLE public.campaigns 
ADD COLUMN fixed_revenue DECIMAL(10,2),
ADD COLUMN fixed_click_amount DECIMAL(10,2),
ADD COLUMN fixed_conversion_amount DECIMAL(10,2),
ADD COLUMN percentage_conversion_amount DECIMAL(5,2); -- Percentage with 2 decimal places (e.g., 15.50 for 15.5%)

-- Migrate existing data to new structure
-- NOTE: revenue_structure = 'percentage' cases cannot be directly mapped to the new simplified model
-- and will lose the percentage-based revenue calculation. Consider manual review if such data exists.

-- Show summary of existing data before migration
DO $$
BEGIN
    RAISE NOTICE 'Campaign billing migration summary:';
    RAISE NOTICE 'Total campaigns: %', (SELECT COUNT(*) FROM public.campaigns);
    RAISE NOTICE 'Revenue structure - fixed: %, percentage: %', 
        (SELECT COUNT(*) FROM public.campaigns WHERE revenue_structure = 'fixed'),
        (SELECT COUNT(*) FROM public.campaigns WHERE revenue_structure = 'percentage');
    RAISE NOTICE 'Billing model - click: %, conversion: %',
        (SELECT COUNT(*) FROM public.campaigns WHERE billing_model = 'click'),
        (SELECT COUNT(*) FROM public.campaigns WHERE billing_model = 'conversion');
    RAISE NOTICE 'Payout structure - fixed: %, percentage: %',
        (SELECT COUNT(*) FROM public.campaigns WHERE payout_structure = 'fixed'),
        (SELECT COUNT(*) FROM public.campaigns WHERE payout_structure = 'percentage');
END $$;

-- Map revenue_structure = 'fixed' → fixed_revenue = revenue_amount
UPDATE public.campaigns 
SET fixed_revenue = revenue_amount 
WHERE revenue_structure = 'fixed' AND revenue_amount IS NOT NULL;

-- Map billing_model = 'click' AND payout_structure = 'fixed' → fixed_click_amount = payout_amount
UPDATE public.campaigns 
SET fixed_click_amount = payout_amount 
WHERE billing_model = 'click' AND payout_structure = 'fixed' AND payout_amount IS NOT NULL;

-- Map billing_model = 'conversion' AND payout_structure = 'fixed' → fixed_conversion_amount = payout_amount
UPDATE public.campaigns 
SET fixed_conversion_amount = payout_amount 
WHERE billing_model = 'conversion' AND payout_structure = 'fixed' AND payout_amount IS NOT NULL;

-- Map billing_model = 'conversion' AND payout_structure = 'percentage' → percentage_conversion_amount = payout_amount
UPDATE public.campaigns 
SET percentage_conversion_amount = payout_amount 
WHERE billing_model = 'conversion' AND payout_structure = 'percentage' AND payout_amount IS NOT NULL;

-- Drop the old complex billing fields
ALTER TABLE public.campaigns 
DROP COLUMN billing_model,
DROP COLUMN payout_structure,
DROP COLUMN payout_amount,
DROP COLUMN revenue_structure,
DROP COLUMN revenue_amount;

-- Add check constraints for the new fields
ALTER TABLE public.campaigns 
ADD CONSTRAINT check_fixed_revenue_positive 
CHECK (fixed_revenue IS NULL OR fixed_revenue >= 0);

ALTER TABLE public.campaigns 
ADD CONSTRAINT check_fixed_click_amount_positive 
CHECK (fixed_click_amount IS NULL OR fixed_click_amount >= 0);

ALTER TABLE public.campaigns 
ADD CONSTRAINT check_fixed_conversion_amount_positive 
CHECK (fixed_conversion_amount IS NULL OR fixed_conversion_amount >= 0);

ALTER TABLE public.campaigns 
ADD CONSTRAINT check_percentage_conversion_amount_valid 
CHECK (percentage_conversion_amount IS NULL OR (percentage_conversion_amount >= 0 AND percentage_conversion_amount <= 100));

-- Add comments for documentation
COMMENT ON COLUMN public.campaigns.fixed_revenue IS 'Fixed revenue amount the platform earns per conversion';
COMMENT ON COLUMN public.campaigns.fixed_click_amount IS 'Fixed amount paid to affiliates per click';
COMMENT ON COLUMN public.campaigns.fixed_conversion_amount IS 'Fixed amount paid to affiliates per conversion';
COMMENT ON COLUMN public.campaigns.percentage_conversion_amount IS 'Percentage of revenue paid to affiliates per conversion (0-100)';