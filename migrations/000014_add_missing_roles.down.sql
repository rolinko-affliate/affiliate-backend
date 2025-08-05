-- Remove the roles added in the up migration
DELETE FROM public.roles WHERE name IN ('AgencyManager', 'PlatformOwner');

-- Revert role description changes
UPDATE public.roles SET description = 'Platform Administrator with full access' WHERE name = 'Admin';