-- #############################################################################
-- ## Rollback Advertiser Association Invitation Migration
-- ## This migration removes the advertiser association invitation system
-- #############################################################################

-- Drop the invitation usage log table first (due to foreign key dependencies)
DROP TABLE IF EXISTS public.invitation_usage_log;

-- Drop the advertiser association invitations table
DROP TABLE IF EXISTS public.advertiser_association_invitations;