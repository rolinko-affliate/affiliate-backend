-- #############################################################################
-- ## Fix Extra Info Tables to Reference Organization ID Instead of Entity IDs
-- ## This migration corrects the design so that extra info tables store
-- ## extended information about organizations directly, not separate entities
-- #############################################################################

-- First, drop existing foreign key constraints and indexes
ALTER TABLE public.advertiser_extra_info DROP CONSTRAINT advertiser_extra_info_advertiser_id_fkey;
ALTER TABLE public.affiliate_extra_info DROP CONSTRAINT affiliate_extra_info_affiliate_id_fkey;

DROP INDEX idx_advertiser_extra_info_advertiser_id;
DROP INDEX idx_affiliate_extra_info_affiliate_id;

-- Rename the columns to reflect their new purpose
ALTER TABLE public.advertiser_extra_info RENAME COLUMN advertiser_id TO organization_id;
ALTER TABLE public.affiliate_extra_info RENAME COLUMN affiliate_id TO organization_id;

-- Update existing data to use organization_id instead of advertiser_id/affiliate_id
UPDATE public.advertiser_extra_info 
SET organization_id = (
    SELECT a.organization_id 
    FROM public.advertisers a 
    WHERE a.advertiser_id = advertiser_extra_info.organization_id
);

UPDATE public.affiliate_extra_info 
SET organization_id = (
    SELECT af.organization_id 
    FROM public.affiliates af 
    WHERE af.affiliate_id = affiliate_extra_info.organization_id
);

-- Add new foreign key constraints to organizations table
ALTER TABLE public.advertiser_extra_info 
ADD CONSTRAINT advertiser_extra_info_organization_id_fkey 
FOREIGN KEY (organization_id) REFERENCES public.organizations(organization_id) ON DELETE CASCADE;

ALTER TABLE public.affiliate_extra_info 
ADD CONSTRAINT affiliate_extra_info_organization_id_fkey 
FOREIGN KEY (organization_id) REFERENCES public.organizations(organization_id) ON DELETE CASCADE;

-- Recreate indexes with new column names
CREATE INDEX idx_advertiser_extra_info_organization_id ON public.advertiser_extra_info(organization_id);
CREATE INDEX idx_affiliate_extra_info_organization_id ON public.affiliate_extra_info(organization_id);

-- Update unique constraints
ALTER TABLE public.advertiser_extra_info DROP CONSTRAINT advertiser_extra_info_advertiser_id_key;
ALTER TABLE public.affiliate_extra_info DROP CONSTRAINT affiliate_extra_info_affiliate_id_key;

ALTER TABLE public.advertiser_extra_info ADD CONSTRAINT advertiser_extra_info_organization_id_key UNIQUE (organization_id);
ALTER TABLE public.affiliate_extra_info ADD CONSTRAINT affiliate_extra_info_organization_id_key UNIQUE (organization_id);

-- Update table comments to reflect the corrected design
COMMENT ON TABLE public.advertiser_extra_info IS 'Extended information for organizations of type advertiser';
COMMENT ON COLUMN public.advertiser_extra_info.organization_id IS 'Reference to the organization this extra info belongs to';

COMMENT ON TABLE public.affiliate_extra_info IS 'Extended information for organizations of type affiliate';
COMMENT ON COLUMN public.affiliate_extra_info.organization_id IS 'Reference to the organization this extra info belongs to';