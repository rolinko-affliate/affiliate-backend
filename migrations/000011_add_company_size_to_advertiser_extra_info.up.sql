-- Add company_size field to advertiser_extra_info table
ALTER TABLE advertiser_extra_info 
ADD COLUMN company_size VARCHAR(50);

-- Add check constraint for valid company sizes
ALTER TABLE advertiser_extra_info 
ADD CONSTRAINT advertiser_extra_info_company_size_check 
CHECK (company_size IS NULL OR company_size IN (
    'startup',           -- 1-10 employees
    'small',            -- 11-50 employees  
    'medium',           -- 51-200 employees
    'large',            -- 201-1000 employees
    'enterprise'        -- 1000+ employees
));

-- Add index for company_size filtering
CREATE INDEX idx_advertiser_extra_info_company_size 
ON advertiser_extra_info(company_size) 
WHERE company_size IS NOT NULL;

-- Add comment for documentation
COMMENT ON COLUMN advertiser_extra_info.company_size IS 'Company size category: startup (1-10), small (11-50), medium (51-200), large (201-1000), enterprise (1000+)';