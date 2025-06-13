-- #############################################################################
-- ## Analytics Service Database Schema
-- ## Creates tables for storing advertiser and publisher analytics data
-- #############################################################################

-- analytics_advertisers: Store advertiser analytics data
CREATE TABLE public.analytics_advertisers (
    id BIGSERIAL PRIMARY KEY,
    domain VARCHAR(255) NOT NULL UNIQUE,
    
    -- Metadata
    description TEXT,
    favicon_image_url TEXT,
    screenshot_image_url TEXT,
    
    -- Affiliate Networks
    affiliate_networks JSONB, -- Array of network names
    
    -- Contact Information
    contact_emails JSONB, -- Array of contact email objects
    
    -- Keywords
    keywords JSONB, -- Array of keyword objects with scores
    
    -- Verticals
    verticals JSONB, -- Array of vertical objects with names, ranks, scores
    
    -- Partner Information
    partner_information JSONB, -- Complex partner data structure
    
    -- Related Advertisers
    related_advertisers JSONB, -- Array of related advertiser domains
    
    -- Social Media
    social_media JSONB, -- Social media links and availability
    
    -- Backlinks
    backlinks JSONB, -- Backlink information
    
    -- Additional analytics data
    additional_data JSONB, -- For any extra fields from the API
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- analytics_publishers: Store publisher analytics data  
CREATE TABLE public.analytics_publishers (
    id BIGSERIAL PRIMARY KEY,
    domain VARCHAR(255) NOT NULL UNIQUE,
    
    -- Metadata
    description TEXT,
    favicon_image_url TEXT,
    screenshot_image_url TEXT,
    
    -- Affiliate Networks
    affiliate_networks JSONB, -- Array of network names
    
    -- Country Rankings
    country_rankings JSONB, -- Array of country ranking objects
    
    -- Keywords
    keywords JSONB, -- Array of keyword objects with scores
    
    -- Verticals
    verticals JSONB, -- Array of vertical names
    verticals_v2 JSONB, -- Array of vertical objects with names, ranks, scores
    
    -- Partner Information
    partner_information JSONB, -- Complex partner data structure
    partners JSONB, -- Array of partner domains
    
    -- Related Publishers
    related_publishers JSONB, -- Array of related publisher domains
    
    -- Social Media
    social_media JSONB, -- Social media links and availability
    
    -- Live URLs
    live_urls JSONB, -- Array of live URLs
    
    -- Flags and Scores
    known BOOLEAN DEFAULT FALSE,
    relevance DECIMAL(5,2) DEFAULT 0,
    traffic_score DECIMAL(10,2) DEFAULT 0,
    
    -- Promotype
    promotype VARCHAR(50),
    
    -- Additional analytics data
    additional_data JSONB, -- For any extra fields from the API
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create triggers for updated_at
CREATE TRIGGER set_analytics_advertisers_timestamp
BEFORE UPDATE ON public.analytics_advertisers
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_analytics_publishers_timestamp
BEFORE UPDATE ON public.analytics_publishers
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Create indexes for performance
CREATE INDEX idx_analytics_advertisers_domain ON public.analytics_advertisers(domain);
CREATE INDEX idx_analytics_advertisers_domain_text ON public.analytics_advertisers USING gin(to_tsvector('english', domain));

CREATE INDEX idx_analytics_publishers_domain ON public.analytics_publishers(domain);
CREATE INDEX idx_analytics_publishers_domain_text ON public.analytics_publishers USING gin(to_tsvector('english', domain));
CREATE INDEX idx_analytics_publishers_known ON public.analytics_publishers(known);
CREATE INDEX idx_analytics_publishers_relevance ON public.analytics_publishers(relevance);
CREATE INDEX idx_analytics_publishers_traffic_score ON public.analytics_publishers(traffic_score);

-- GIN indexes for JSONB fields to enable efficient querying
CREATE INDEX idx_analytics_advertisers_affiliate_networks_gin ON public.analytics_advertisers USING GIN(affiliate_networks);
CREATE INDEX idx_analytics_advertisers_keywords_gin ON public.analytics_advertisers USING GIN(keywords);
CREATE INDEX idx_analytics_advertisers_verticals_gin ON public.analytics_advertisers USING GIN(verticals);

CREATE INDEX idx_analytics_publishers_affiliate_networks_gin ON public.analytics_publishers USING GIN(affiliate_networks);
CREATE INDEX idx_analytics_publishers_keywords_gin ON public.analytics_publishers USING GIN(keywords);
CREATE INDEX idx_analytics_publishers_verticals_gin ON public.analytics_publishers USING GIN(verticals);
CREATE INDEX idx_analytics_publishers_verticals_v2_gin ON public.analytics_publishers USING GIN(verticals_v2);

-- Add comments for documentation
COMMENT ON TABLE public.analytics_advertisers IS 'Analytics data for advertisers including metadata, networks, keywords, and partner information';
COMMENT ON TABLE public.analytics_publishers IS 'Analytics data for publishers including metadata, networks, keywords, and partner information';

COMMENT ON COLUMN public.analytics_advertisers.domain IS 'Primary domain of the advertiser';
COMMENT ON COLUMN public.analytics_advertisers.affiliate_networks IS 'Array of affiliate network names the advertiser works with';
COMMENT ON COLUMN public.analytics_advertisers.keywords IS 'Array of keyword objects with scores for SEO/content analysis';
COMMENT ON COLUMN public.analytics_advertisers.verticals IS 'Array of vertical/industry categories with rankings and scores';

COMMENT ON COLUMN public.analytics_publishers.domain IS 'Primary domain of the publisher';
COMMENT ON COLUMN public.analytics_publishers.affiliate_networks IS 'Array of affiliate network names the publisher works with';
COMMENT ON COLUMN public.analytics_publishers.keywords IS 'Array of keyword objects with scores for content analysis';
COMMENT ON COLUMN public.analytics_publishers.verticals IS 'Array of vertical/industry category names';
COMMENT ON COLUMN public.analytics_publishers.verticals_v2 IS 'Enhanced vertical data with names, ranks, and scores';
COMMENT ON COLUMN public.analytics_publishers.known IS 'Flag indicating if this is a known/verified publisher';
COMMENT ON COLUMN public.analytics_publishers.relevance IS 'Relevance score for the publisher';
COMMENT ON COLUMN public.analytics_publishers.traffic_score IS 'Traffic score indicating publisher reach';