-- Add missing roles for agency management and platform ownership
-- This migration adds roles that are referenced in the codebase but missing from the database

-- Add AgencyManager role (referenced in router but not in database)
INSERT INTO public.roles (role_id, name, description) VALUES
  (1002, 'AgencyManager', 'Manages agency operations and delegation relationships')
ON CONFLICT (name) DO NOTHING;

-- Add PlatformOwner role for platform owner organization users
INSERT INTO public.roles (role_id, name, description) VALUES
  (2, 'PlatformOwner', 'Platform owner with administrative privileges over all organizations and delegations')
ON CONFLICT (name) DO NOTHING;

-- Update role descriptions for clarity
UPDATE public.roles SET description = 'Platform Administrator with full system access' WHERE name = 'Admin';
UPDATE public.roles SET description = 'Platform Owner with administrative privileges over all organizations and delegations' WHERE name = 'PlatformOwner';
UPDATE public.roles SET description = 'Manages agency operations, delegation relationships, and client accounts' WHERE name = 'AgencyManager';