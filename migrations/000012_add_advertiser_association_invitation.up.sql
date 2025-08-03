-- #############################################################################
-- ## Advertiser Association Invitation Migration
-- ## This migration adds support for advertiser-generated invitation links
-- ## that allow affiliate organizations to create association requests through
-- ## shareable links with optional restrictions on which affiliates can use them.
-- #############################################################################

-- advertiser_association_invitations: Manages invitation links created by advertisers
CREATE TABLE public.advertiser_association_invitations (
    invitation_id BIGSERIAL PRIMARY KEY,
    advertiser_org_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    
    -- Invitation link details
    invitation_token VARCHAR(255) NOT NULL UNIQUE, -- Unique token for the invitation link
    name VARCHAR(255) NOT NULL, -- Human-readable name for the invitation
    description TEXT, -- Optional description of the invitation
    
    -- Access control
    allowed_affiliate_org_ids JSONB, -- Array of affiliate org IDs that can use this invitation (null = unrestricted)
    
    -- Usage limits
    max_uses INTEGER, -- Maximum number of times this invitation can be used (null = unlimited)
    current_uses INTEGER NOT NULL DEFAULT 0, -- Current number of times this invitation has been used
    
    -- Expiration
    expires_at TIMESTAMPTZ, -- When this invitation expires (null = never expires)
    
    -- Status and metadata
    status VARCHAR(20) NOT NULL DEFAULT 'active' 
        CHECK (status IN ('active', 'disabled', 'expired')),
    created_by_user_id UUID NOT NULL REFERENCES public.profiles(id) ON DELETE RESTRICT,
    message TEXT, -- Optional message to display when using the invitation
    
    -- Default visibility settings for associations created through this invitation
    default_all_affiliates_visible BOOLEAN NOT NULL DEFAULT TRUE,
    default_all_campaigns_visible BOOLEAN NOT NULL DEFAULT TRUE,
    default_visible_affiliate_ids JSONB, -- Default affiliate IDs to make visible
    default_visible_campaign_ids JSONB,  -- Default campaign IDs to make visible
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Add trigger for automatic updated_at timestamp
CREATE TRIGGER set_advertiser_association_invitations_timestamp
BEFORE UPDATE ON public.advertiser_association_invitations
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Indexes for performance
CREATE INDEX idx_adv_assoc_invitations_advertiser_org_id ON public.advertiser_association_invitations(advertiser_org_id);
CREATE INDEX idx_adv_assoc_invitations_token ON public.advertiser_association_invitations(invitation_token);
CREATE INDEX idx_adv_assoc_invitations_status ON public.advertiser_association_invitations(status);
CREATE INDEX idx_adv_assoc_invitations_expires_at ON public.advertiser_association_invitations(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX idx_adv_assoc_invitations_created_by ON public.advertiser_association_invitations(created_by_user_id);
CREATE INDEX idx_adv_assoc_invitations_created_at ON public.advertiser_association_invitations(created_at);

-- GIN indexes for JSONB fields to enable efficient querying
CREATE INDEX idx_adv_assoc_invitations_allowed_orgs_gin 
ON public.advertiser_association_invitations USING GIN(allowed_affiliate_org_ids) 
WHERE allowed_affiliate_org_ids IS NOT NULL;

CREATE INDEX idx_adv_assoc_invitations_default_affiliates_gin 
ON public.advertiser_association_invitations USING GIN(default_visible_affiliate_ids) 
WHERE default_visible_affiliate_ids IS NOT NULL;

CREATE INDEX idx_adv_assoc_invitations_default_campaigns_gin 
ON public.advertiser_association_invitations USING GIN(default_visible_campaign_ids) 
WHERE default_visible_campaign_ids IS NOT NULL;

-- Composite indexes for common query patterns
CREATE INDEX idx_adv_assoc_invitations_advertiser_status ON public.advertiser_association_invitations(advertiser_org_id, status);
CREATE INDEX idx_adv_assoc_invitations_status_expires ON public.advertiser_association_invitations(status, expires_at);

-- invitation_usage_log: Track usage of invitations for audit and analytics
CREATE TABLE public.invitation_usage_log (
    usage_id BIGSERIAL PRIMARY KEY,
    invitation_id BIGINT NOT NULL REFERENCES public.advertiser_association_invitations(invitation_id) ON DELETE CASCADE,
    affiliate_org_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    used_by_user_id UUID REFERENCES public.profiles(id) ON DELETE SET NULL,
    association_id BIGINT REFERENCES public.organization_associations(association_id) ON DELETE SET NULL,
    
    -- Usage metadata
    ip_address INET, -- IP address of the user who used the invitation
    user_agent TEXT, -- User agent string
    success BOOLEAN NOT NULL DEFAULT TRUE, -- Whether the invitation usage was successful
    error_message TEXT, -- Error message if usage failed
    
    used_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Indexes for invitation usage log
CREATE INDEX idx_invitation_usage_log_invitation_id ON public.invitation_usage_log(invitation_id);
CREATE INDEX idx_invitation_usage_log_affiliate_org_id ON public.invitation_usage_log(affiliate_org_id);
CREATE INDEX idx_invitation_usage_log_used_by ON public.invitation_usage_log(used_by_user_id) WHERE used_by_user_id IS NOT NULL;
CREATE INDEX idx_invitation_usage_log_association_id ON public.invitation_usage_log(association_id) WHERE association_id IS NOT NULL;
CREATE INDEX idx_invitation_usage_log_used_at ON public.invitation_usage_log(used_at);
CREATE INDEX idx_invitation_usage_log_success ON public.invitation_usage_log(success);

-- Composite indexes for common query patterns
CREATE INDEX idx_invitation_usage_log_invitation_success ON public.invitation_usage_log(invitation_id, success);
CREATE INDEX idx_invitation_usage_log_affiliate_used_at ON public.invitation_usage_log(affiliate_org_id, used_at);

-- Add comments for documentation
COMMENT ON TABLE public.advertiser_association_invitations IS 'Manages invitation links created by advertisers that allow affiliate organizations to create association requests';
COMMENT ON COLUMN public.advertiser_association_invitations.invitation_token IS 'Unique token used in the invitation URL';
COMMENT ON COLUMN public.advertiser_association_invitations.allowed_affiliate_org_ids IS 'JSONB array of affiliate organization IDs that are allowed to use this invitation (null = unrestricted)';
COMMENT ON COLUMN public.advertiser_association_invitations.max_uses IS 'Maximum number of times this invitation can be used (null = unlimited)';
COMMENT ON COLUMN public.advertiser_association_invitations.current_uses IS 'Current number of times this invitation has been used';
COMMENT ON COLUMN public.advertiser_association_invitations.default_visible_affiliate_ids IS 'Default affiliate IDs to make visible when association is created through this invitation';
COMMENT ON COLUMN public.advertiser_association_invitations.default_visible_campaign_ids IS 'Default campaign IDs to make visible when association is created through this invitation';

COMMENT ON TABLE public.invitation_usage_log IS 'Tracks usage of advertiser association invitations for audit and analytics purposes';
COMMENT ON COLUMN public.invitation_usage_log.success IS 'Whether the invitation usage resulted in a successful association creation';