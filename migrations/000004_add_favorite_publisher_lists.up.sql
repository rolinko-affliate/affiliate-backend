
-- #############################################################################
-- ## Favorite Publisher Lists Migration
-- ## This migration adds support for organizations to manage favorite publisher lists
-- ## 
-- ## Features:
-- ## - Organizations can create multiple favorite publisher lists
-- ## - Each list has a name and description
-- ## - Publishers (by domain) can be added/removed from lists
-- ## - Proper foreign key constraints and indexes for performance
-- #############################################################################

-- favorite_publisher_lists: Stores the favorite publisher lists for organizations
CREATE TABLE public.favorite_publisher_lists (
    list_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Ensure unique list names per organization
    CONSTRAINT unique_list_name_per_org UNIQUE (organization_id, name)
);

-- Add trigger for automatic timestamp updates
CREATE TRIGGER set_favorite_publisher_lists_timestamp
BEFORE UPDATE ON public.favorite_publisher_lists
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- favorite_publisher_list_items: Stores the publishers in each favorite list
CREATE TABLE public.favorite_publisher_list_items (
    item_id BIGSERIAL PRIMARY KEY,
    list_id BIGINT NOT NULL REFERENCES public.favorite_publisher_lists(list_id) ON DELETE CASCADE,
    publisher_domain VARCHAR(255) NOT NULL,
    notes TEXT, -- Optional notes about why this publisher is in the list
    status VARCHAR(20) DEFAULT 'added' NOT NULL CHECK (status IN ('added', 'contacted', 'accepted')),
    added_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Ensure a publisher domain can only be added once per list
    CONSTRAINT unique_publisher_per_list UNIQUE (list_id, publisher_domain)
);

-- Indexes for performance
CREATE INDEX idx_favorite_publisher_lists_organization_id ON public.favorite_publisher_lists(organization_id);
CREATE INDEX idx_favorite_publisher_lists_name ON public.favorite_publisher_lists(name);
CREATE INDEX idx_favorite_publisher_list_items_list_id ON public.favorite_publisher_list_items(list_id);
CREATE INDEX idx_favorite_publisher_list_items_domain ON public.favorite_publisher_list_items(publisher_domain);
CREATE INDEX idx_favorite_publisher_list_items_status ON public.favorite_publisher_list_items(status);
CREATE INDEX idx_favorite_publisher_list_items_added_at ON public.favorite_publisher_list_items(added_at);

-- Add comments for documentation
COMMENT ON TABLE public.favorite_publisher_lists IS 'Stores favorite publisher lists created by organizations';
COMMENT ON TABLE public.favorite_publisher_list_items IS 'Stores individual publisher domains within favorite lists';

COMMENT ON COLUMN public.favorite_publisher_lists.organization_id IS 'Reference to the organization that owns this list';
COMMENT ON COLUMN public.favorite_publisher_lists.name IS 'Display name of the favorite list';
COMMENT ON COLUMN public.favorite_publisher_lists.description IS 'Optional description of the list purpose';

COMMENT ON COLUMN public.favorite_publisher_list_items.list_id IS 'Reference to the favorite list this item belongs to';
COMMENT ON COLUMN public.favorite_publisher_list_items.publisher_domain IS 'Domain name of the publisher (e.g., example.com)';
COMMENT ON COLUMN public.favorite_publisher_list_items.notes IS 'Optional notes about this publisher in the context of this list';
COMMENT ON COLUMN public.favorite_publisher_list_items.status IS 'Status of publisher interaction: added -> contacted -> accepted';
COMMENT ON COLUMN public.favorite_publisher_list_items.added_at IS 'Timestamp when the publisher was added to the list';