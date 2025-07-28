-- Remove 'agency' as a valid organization type
-- This migration reverts the CHECK constraint on the organizations table
-- to exclude 'agency' as a valid organization type

-- First, ensure no organizations have type 'agency' before removing the constraint
-- This will fail if there are any 'agency' type organizations, which is intentional
-- to prevent data loss
UPDATE public.organizations 
SET type = 'platform_owner' 
WHERE type = 'agency';

-- Drop the existing constraint
ALTER TABLE public.organizations 
DROP CONSTRAINT IF EXISTS organizations_type_check;

-- Add back the original constraint without 'agency'
ALTER TABLE public.organizations 
ADD CONSTRAINT organizations_type_check 
CHECK (type IN ('advertiser', 'affiliate', 'platform_owner'));