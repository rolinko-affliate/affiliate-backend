-- #############################################################################
-- ## Organization Associations Migration Rollback
-- ## This migration removes the organization associations table and related
-- ## indexes, constraints, and triggers.
-- #############################################################################

-- Drop the organization_associations table (this will cascade and remove all related indexes and constraints)
DROP TABLE IF EXISTS public.organization_associations CASCADE;

-- Note: The trigger function trigger_set_timestamp() is shared with other tables,
-- so we don't drop it here. It was created in the initial migration.