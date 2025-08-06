-- #############################################################################
-- ## Agency Delegation Migration Rollback
-- ## This migration removes the agency delegation system
-- #############################################################################

-- Drop the validation trigger and function
DROP TRIGGER IF EXISTS validate_agency_delegation_org_types_trigger ON public.agency_delegations;
DROP FUNCTION IF EXISTS validate_agency_delegation_org_types();

-- Drop the timestamp trigger
DROP TRIGGER IF EXISTS set_agency_delegations_timestamp ON public.agency_delegations;

-- Drop the agency_delegations table (this will also drop all indexes and constraints)
DROP TABLE IF EXISTS public.agency_delegations;