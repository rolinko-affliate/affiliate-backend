-- #############################################################################
-- ## Analytics Service Database Schema Rollback
-- ## Drops analytics tables and related indexes
-- #############################################################################

-- Drop indexes first
DROP INDEX IF EXISTS idx_analytics_publishers_verticals_v2_gin;
DROP INDEX IF EXISTS idx_analytics_publishers_verticals_gin;
DROP INDEX IF EXISTS idx_analytics_publishers_keywords_gin;
DROP INDEX IF EXISTS idx_analytics_publishers_affiliate_networks_gin;

DROP INDEX IF EXISTS idx_analytics_advertisers_verticals_gin;
DROP INDEX IF EXISTS idx_analytics_advertisers_keywords_gin;
DROP INDEX IF EXISTS idx_analytics_advertisers_affiliate_networks_gin;

DROP INDEX IF EXISTS idx_analytics_publishers_traffic_score;
DROP INDEX IF EXISTS idx_analytics_publishers_relevance;
DROP INDEX IF EXISTS idx_analytics_publishers_known;
DROP INDEX IF EXISTS idx_analytics_publishers_domain_text;
DROP INDEX IF EXISTS idx_analytics_publishers_domain;

DROP INDEX IF EXISTS idx_analytics_advertisers_domain_text;
DROP INDEX IF EXISTS idx_analytics_advertisers_domain;

-- Drop triggers
DROP TRIGGER IF EXISTS set_analytics_publishers_timestamp ON public.analytics_publishers;
DROP TRIGGER IF EXISTS set_analytics_advertisers_timestamp ON public.analytics_advertisers;

-- Drop tables
DROP TABLE IF EXISTS public.analytics_publishers;
DROP TABLE IF EXISTS public.analytics_advertisers;