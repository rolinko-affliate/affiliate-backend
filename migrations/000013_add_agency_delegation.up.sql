-- #############################################################################
-- ## Agency Delegation Migration
-- ## This migration adds support for managing delegation relationships between
-- ## agency and advertiser organizations, including granular permission control
-- ## and expiration management.
-- #############################################################################

-- agency_delegations: Manages delegation relationships between agency and advertiser organizations
CREATE TABLE public.agency_delegations (
    delegation_id BIGSERIAL PRIMARY KEY,
    agency_org_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    advertiser_org_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' 
        CHECK (status IN ('pending', 'active', 'suspended', 'revoked')),
    
    -- Permissions granted to the agency (JSONB array of permission strings)
    permissions JSONB NOT NULL DEFAULT '[]'::jsonb,
    
    -- Delegation metadata
    delegated_by_user_id UUID, -- References profiles.id (auth.uid()) - advertiser user who created delegation
    accepted_by_user_id UUID,  -- References profiles.id (auth.uid()) - agency user who accepted delegation
    message TEXT, -- Optional message with delegation
    
    -- Expiration settings
    expires_at TIMESTAMPTZ, -- Optional expiration date
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    accepted_at TIMESTAMPTZ,
    
    -- Ensure unique delegation per agency-advertiser pair
    UNIQUE (agency_org_id, advertiser_org_id),
    
    -- Ensure agency and advertiser are different organizations
    CHECK (agency_org_id != advertiser_org_id),
    
    -- Ensure expiration date is in the future when set
    CHECK (expires_at IS NULL OR expires_at > created_at)
);

-- Add trigger for automatic updated_at timestamp
CREATE TRIGGER set_agency_delegations_timestamp
BEFORE UPDATE ON public.agency_delegations
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Indexes for performance
CREATE INDEX idx_agency_delegations_agency_org_id ON public.agency_delegations(agency_org_id);
CREATE INDEX idx_agency_delegations_advertiser_org_id ON public.agency_delegations(advertiser_org_id);
CREATE INDEX idx_agency_delegations_status ON public.agency_delegations(status);
CREATE INDEX idx_agency_delegations_created_at ON public.agency_delegations(created_at);
CREATE INDEX idx_agency_delegations_expires_at ON public.agency_delegations(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX idx_agency_delegations_delegated_by ON public.agency_delegations(delegated_by_user_id) WHERE delegated_by_user_id IS NOT NULL;
CREATE INDEX idx_agency_delegations_accepted_by ON public.agency_delegations(accepted_by_user_id) WHERE accepted_by_user_id IS NOT NULL;

-- GIN index for JSONB permissions field to enable efficient querying
CREATE INDEX idx_agency_delegations_permissions_gin 
ON public.agency_delegations USING GIN(permissions);

-- Composite indexes for common query patterns
CREATE INDEX idx_agency_delegations_agency_status ON public.agency_delegations(agency_org_id, status);
CREATE INDEX idx_agency_delegations_advertiser_status ON public.agency_delegations(advertiser_org_id, status);
CREATE INDEX idx_agency_delegations_active_not_expired ON public.agency_delegations(agency_org_id, advertiser_org_id, status) 
WHERE status = 'active';

-- Add foreign key constraints for user references (optional, as users might be deleted)
-- These are not enforced with CASCADE to allow user deletion without affecting delegations
ALTER TABLE public.agency_delegations 
ADD CONSTRAINT fk_agency_delegations_delegated_by_user 
FOREIGN KEY (delegated_by_user_id) REFERENCES public.profiles(id) ON DELETE SET NULL;

ALTER TABLE public.agency_delegations 
ADD CONSTRAINT fk_agency_delegations_accepted_by_user 
FOREIGN KEY (accepted_by_user_id) REFERENCES public.profiles(id) ON DELETE SET NULL;

-- Add organization type validation function
CREATE OR REPLACE FUNCTION validate_agency_delegation_org_types()
RETURNS TRIGGER AS $$
BEGIN
    -- Check that agency_org_id is actually an agency
    IF NOT EXISTS (
        SELECT 1 FROM public.organizations 
        WHERE organization_id = NEW.agency_org_id 
        AND type = 'agency'
    ) THEN
        RAISE EXCEPTION 'Agency organization ID % is not of type agency', NEW.agency_org_id;
    END IF;
    
    -- Check that advertiser_org_id is actually an advertiser
    IF NOT EXISTS (
        SELECT 1 FROM public.organizations 
        WHERE organization_id = NEW.advertiser_org_id 
        AND type = 'advertiser'
    ) THEN
        RAISE EXCEPTION 'Advertiser organization ID % is not of type advertiser', NEW.advertiser_org_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add trigger to validate organization types
CREATE TRIGGER validate_agency_delegation_org_types_trigger
BEFORE INSERT OR UPDATE ON public.agency_delegations
FOR EACH ROW
EXECUTE FUNCTION validate_agency_delegation_org_types();

-- Add comments for documentation
COMMENT ON TABLE public.agency_delegations IS 'Manages delegation relationships between agency and advertiser organizations with granular permission control';
COMMENT ON COLUMN public.agency_delegations.permissions IS 'JSONB array of permission strings that the agency is granted for the advertiser organization';
COMMENT ON COLUMN public.agency_delegations.delegated_by_user_id IS 'UUID of the advertiser user who created the delegation';
COMMENT ON COLUMN public.agency_delegations.accepted_by_user_id IS 'UUID of the agency user who accepted the delegation';
COMMENT ON COLUMN public.agency_delegations.expires_at IS 'Optional expiration timestamp for the delegation';
COMMENT ON COLUMN public.agency_delegations.status IS 'Current status of the delegation: pending, active, suspended, or revoked';