-- #############################################################################
-- ## Rollback: Revert Extra Info Tables to Reference Entity IDs
-- ## This rollback migration reverts the changes to reference advertiser_id/affiliate_id
-- #############################################################################

-- Drop current foreign key constraints and indexes
ALTER TABLE public.advertiser_extra_info DROP CONSTRAINT advertiser_extra_info_organization_id_fkey;
ALTER TABLE public.affiliate_extra_info DROP CONSTRAINT affiliate_extra_info_organization_id_fkey;

DROP INDEX idx_advertiser_extra_info_organization_id;
DROP INDEX idx_affiliate_extra_info_organization_id;

-- Update data back to entity IDs (this assumes 1:1 mapping still exists)
UPDATE public.advertiser_extra_info 
SET organization_id = (
    SELECT a.advertiser_id 
    FROM public.advertisers a 
    WHERE a.organization_id = advertiser_extra_info.organization_id
);

UPDATE public.affiliate_extra_info 
SET organization_id = (
    SELECT af.affiliate_id 
    FROM public.affiliates af 
    WHERE af.organization_id = affiliate_extra_info.organization_id
);

-- Rename columns back
ALTER TABLE public.advertiser_extra_info RENAME COLUMN organization_id TO advertiser_id;
ALTER TABLE public.affiliate_extra_info RENAME COLUMN organization_id TO affiliate_id;

-- Restore original foreign key constraints
ALTER TABLE public.advertiser_extra_info 
ADD CONSTRAINT advertiser_extra_info_advertiser_id_fkey 
FOREIGN KEY (advertiser_id) REFERENCES public.advertisers(advertiser_id) ON DELETE CASCADE;

ALTER TABLE public.affiliate_extra_info 
ADD CONSTRAINT affiliate_extra_info_affiliate_id_fkey 
FOREIGN KEY (affiliate_id) REFERENCES public.affiliates(affiliate_id) ON DELETE CASCADE;

-- Restore original indexes
CREATE INDEX idx_advertiser_extra_info_advertiser_id ON public.advertiser_extra_info(advertiser_id);
CREATE INDEX idx_affiliate_extra_info_affiliate_id ON public.affiliate_extra_info(affiliate_id);

-- Restore original unique constraints
ALTER TABLE public.advertiser_extra_info DROP CONSTRAINT advertiser_extra_info_organization_id_key;
ALTER TABLE public.affiliate_extra_info DROP CONSTRAINT affiliate_extra_info_organization_id_key;

ALTER TABLE public.advertiser_extra_info ADD CONSTRAINT advertiser_extra_info_advertiser_id_key UNIQUE (advertiser_id);
ALTER TABLE public.affiliate_extra_info ADD CONSTRAINT affiliate_extra_info_affiliate_id_key UNIQUE (affiliate_id);

-- Restore original comments
COMMENT ON TABLE public.advertiser_extra_info IS 'Additional information for advertisers including website and platform type';
COMMENT ON COLUMN public.advertiser_extra_info.advertiser_id IS 'Reference to the advertiser this extra info belongs to';

COMMENT ON TABLE public.affiliate_extra_info IS 'Additional information for affiliates including website, type, description, and logo';
COMMENT ON COLUMN public.affiliate_extra_info.affiliate_id IS 'Reference to the affiliate this extra info belongs to';