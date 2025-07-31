-- #############################################################################
-- ## Add Extra Information Tables for Advertisers and Affiliates
-- ## This migration adds two new tables to store additional information:
-- ## - advertiser_extra_info: website and website type information
-- ## - affiliate_extra_info: website, type, description, and logo information
-- #############################################################################

-- advertiser_extra_info: Additional information for advertisers
CREATE TABLE public.advertiser_extra_info (
    extra_info_id BIGSERIAL PRIMARY KEY,
    advertiser_id BIGINT NOT NULL REFERENCES public.advertisers(advertiser_id) ON DELETE CASCADE,
    website VARCHAR(500),
    website_type VARCHAR(50) CHECK (website_type IS NULL OR website_type IN ('shopify', 'amazon', 'shopline', 'tiktok_shop')),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Ensure one extra info record per advertiser
    UNIQUE (advertiser_id)
);

CREATE TRIGGER set_advertiser_extra_info_timestamp
BEFORE UPDATE ON public.advertiser_extra_info
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_advertiser_extra_info_advertiser_id ON public.advertiser_extra_info(advertiser_id);
CREATE INDEX idx_advertiser_extra_info_website_type ON public.advertiser_extra_info(website_type) WHERE website_type IS NOT NULL;

-- affiliate_extra_info: Additional information for affiliates
CREATE TABLE public.affiliate_extra_info (
    extra_info_id BIGSERIAL PRIMARY KEY,
    affiliate_id BIGINT NOT NULL REFERENCES public.affiliates(affiliate_id) ON DELETE CASCADE,
    website VARCHAR(500),
    affiliate_type VARCHAR(50) CHECK (affiliate_type IS NULL OR affiliate_type IN ('cashback', 'blog', 'incentive', 'content', 'forum', 'sub_affiliate_network')),
    self_description TEXT,
    logo_url VARCHAR(500),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Ensure one extra info record per affiliate
    UNIQUE (affiliate_id)
);

CREATE TRIGGER set_affiliate_extra_info_timestamp
BEFORE UPDATE ON public.affiliate_extra_info
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_affiliate_extra_info_affiliate_id ON public.affiliate_extra_info(affiliate_id);
CREATE INDEX idx_affiliate_extra_info_affiliate_type ON public.affiliate_extra_info(affiliate_type) WHERE affiliate_type IS NOT NULL;

-- Add comments for documentation
COMMENT ON TABLE public.advertiser_extra_info IS 'Additional information for advertisers including website and platform type';
COMMENT ON COLUMN public.advertiser_extra_info.website IS 'Advertiser website URL';
COMMENT ON COLUMN public.advertiser_extra_info.website_type IS 'Type of website platform (shopify, amazon, shopline, tiktok_shop)';

COMMENT ON TABLE public.affiliate_extra_info IS 'Additional information for affiliates including website, type, description, and logo';
COMMENT ON COLUMN public.affiliate_extra_info.website IS 'Affiliate website URL';
COMMENT ON COLUMN public.affiliate_extra_info.affiliate_type IS 'Type of affiliate (cashback, blog, incentive, content, forum, sub_affiliate_network)';
COMMENT ON COLUMN public.affiliate_extra_info.self_description IS 'Affiliate self-description text';
COMMENT ON COLUMN public.affiliate_extra_info.logo_url IS 'URL to affiliate logo image';