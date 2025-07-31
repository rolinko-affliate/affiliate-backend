-- #############################################################################
-- ## Organization Associations Migration
-- ## This migration adds support for managing associations between advertiser
-- ## and affiliate organizations, including invitation/request system and
-- ## visibility controls for affiliates and campaigns.
-- #############################################################################

-- organization_associations: Manages relationships between advertiser and affiliate organizations
CREATE TABLE public.organization_associations (
    association_id BIGSERIAL PRIMARY KEY,
    advertiser_org_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    affiliate_org_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' 
        CHECK (status IN ('pending', 'active', 'suspended', 'rejected')),
    association_type VARCHAR(20) NOT NULL 
        CHECK (association_type IN ('invitation', 'request')),
    
    -- Visibility settings (JSONB arrays of IDs)
    visible_affiliate_ids JSONB, -- Array of affiliate IDs visible to advertiser
    visible_campaign_ids JSONB,  -- Array of campaign IDs visible to affiliate
    
    -- Default visibility flags (when true, all affiliates/campaigns are visible)
    all_affiliates_visible BOOLEAN NOT NULL DEFAULT TRUE,
    all_campaigns_visible BOOLEAN NOT NULL DEFAULT TRUE,
    
    -- Request/invitation metadata
    requested_by_user_id UUID, -- References profiles.id (auth.uid())
    approved_by_user_id UUID,  -- References profiles.id (auth.uid())
    message TEXT, -- Optional message with request/invitation
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    approved_at TIMESTAMPTZ,
    
    -- Ensure unique association per advertiser-affiliate pair
    UNIQUE (advertiser_org_id, affiliate_org_id),
    
    -- Ensure advertiser and affiliate are different organizations
    CHECK (advertiser_org_id != affiliate_org_id)
);

-- Add trigger for automatic updated_at timestamp
CREATE TRIGGER set_organization_associations_timestamp
BEFORE UPDATE ON public.organization_associations
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Indexes for performance
CREATE INDEX idx_org_associations_advertiser_org_id ON public.organization_associations(advertiser_org_id);
CREATE INDEX idx_org_associations_affiliate_org_id ON public.organization_associations(affiliate_org_id);
CREATE INDEX idx_org_associations_status ON public.organization_associations(status);
CREATE INDEX idx_org_associations_type ON public.organization_associations(association_type);
CREATE INDEX idx_org_associations_created_at ON public.organization_associations(created_at);
CREATE INDEX idx_org_associations_requested_by ON public.organization_associations(requested_by_user_id) WHERE requested_by_user_id IS NOT NULL;
CREATE INDEX idx_org_associations_approved_by ON public.organization_associations(approved_by_user_id) WHERE approved_by_user_id IS NOT NULL;

-- GIN indexes for JSONB fields to enable efficient querying
CREATE INDEX idx_org_associations_visible_affiliates_gin 
ON public.organization_associations USING GIN(visible_affiliate_ids) 
WHERE visible_affiliate_ids IS NOT NULL;

CREATE INDEX idx_org_associations_visible_campaigns_gin 
ON public.organization_associations USING GIN(visible_campaign_ids) 
WHERE visible_campaign_ids IS NOT NULL;

-- Composite indexes for common query patterns
CREATE INDEX idx_org_associations_advertiser_status ON public.organization_associations(advertiser_org_id, status);
CREATE INDEX idx_org_associations_affiliate_status ON public.organization_associations(affiliate_org_id, status);
CREATE INDEX idx_org_associations_status_type ON public.organization_associations(status, association_type);

-- Add foreign key constraints for user references (optional, as users might be deleted)
-- These are not enforced with CASCADE to allow user deletion without affecting associations
ALTER TABLE public.organization_associations 
ADD CONSTRAINT fk_org_associations_requested_by_user 
FOREIGN KEY (requested_by_user_id) REFERENCES public.profiles(id) ON DELETE SET NULL;

ALTER TABLE public.organization_associations 
ADD CONSTRAINT fk_org_associations_approved_by_user 
FOREIGN KEY (approved_by_user_id) REFERENCES public.profiles(id) ON DELETE SET NULL;

-- Add constraint to ensure organization types are correct
-- This constraint ensures advertiser_org_id points to an advertiser organization
-- and affiliate_org_id points to an affiliate organization
ALTER TABLE public.organization_associations 
ADD CONSTRAINT check_advertiser_org_type 
CHECK (
    EXISTS (
        SELECT 1 FROM public.organizations 
        WHERE organization_id = advertiser_org_id 
        AND type = 'advertiser'
    )
);

ALTER TABLE public.organization_associations 
ADD CONSTRAINT check_affiliate_org_type 
CHECK (
    EXISTS (
        SELECT 1 FROM public.organizations 
        WHERE organization_id = affiliate_org_id 
        AND type = 'affiliate'
    )
);

-- Add comments for documentation
COMMENT ON TABLE public.organization_associations IS 'Manages associations between advertiser and affiliate organizations, including invitation/request system and visibility controls';
COMMENT ON COLUMN public.organization_associations.association_type IS 'Type of association: invitation (advertiser invites affiliate) or request (affiliate requests to join advertiser)';
COMMENT ON COLUMN public.organization_associations.visible_affiliate_ids IS 'JSONB array of affiliate IDs that are visible to the advertiser organization';
COMMENT ON COLUMN public.organization_associations.visible_campaign_ids IS 'JSONB array of campaign IDs that are visible to the affiliate organization';
COMMENT ON COLUMN public.organization_associations.all_affiliates_visible IS 'When true, all affiliates in the affiliate organization are visible to the advertiser';
COMMENT ON COLUMN public.organization_associations.all_campaigns_visible IS 'When true, all campaigns in the advertiser organization are visible to the affiliate';