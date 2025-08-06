-- Add 'agency' as a valid organization type
-- This migration updates the CHECK constraint on the organizations table
-- to include 'agency' as a valid organization type

-- Drop the existing constraint
ALTER TABLE public.organizations 
DROP CONSTRAINT IF EXISTS organizations_type_check;

-- Add the new constraint that includes 'agency'
ALTER TABLE public.organizations 
ADD CONSTRAINT organizations_type_check 
CHECK (type IN ('advertiser', 'affiliate', 'platform_owner', 'agency'));