-- #############################################################################
-- ## Rollback Favorite Publisher Lists Migration
-- ## This migration removes the favorite publisher lists functionality
-- #############################################################################

-- Drop tables in reverse order due to foreign key constraints
DROP TABLE IF EXISTS public.favorite_publisher_list_items;
DROP TABLE IF EXISTS public.favorite_publisher_lists;