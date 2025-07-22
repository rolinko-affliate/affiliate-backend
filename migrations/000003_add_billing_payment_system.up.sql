-- #############################################################################
-- ## Organization-Level Billing and Payment System Migration
-- ## This migration adds comprehensive billing and payment functionality
-- ## 
-- ## Includes:
-- ## - Billing accounts (extends organizations with billing info)
-- ## - Payment methods (Stripe payment method storage)
-- ## - Transactions (all billing/payment events)
-- ## - Invoices (postpaid billing)
-- ## - Usage records (daily usage calculation)
-- #############################################################################

-- #############################################################################
-- ## Billing Accounts (Organization-Level Billing Configuration)
-- #############################################################################

-- billing_accounts: Organization-level billing configuration and balance
CREATE TABLE public.billing_accounts (
    billing_account_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    
    -- Stripe Integration
    stripe_customer_id VARCHAR(255) UNIQUE, -- Stripe Customer ID
    stripe_account_id VARCHAR(255), -- Stripe Connected Account ID (for payouts)
    
    -- Billing Configuration
    billing_mode VARCHAR(20) NOT NULL DEFAULT 'prepaid' CHECK (billing_mode IN ('prepaid', 'postpaid')),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD', -- ISO 4217 currency code
    
    -- Balance and Credit (for prepaid accounts)
    balance DECIMAL(15,4) NOT NULL DEFAULT 0.00, -- Current balance
    credit_limit DECIMAL(15,4) DEFAULT 0.00, -- Credit limit for postpaid accounts
    
    -- Payment Configuration
    default_payment_method_id VARCHAR(255), -- Stripe Payment Method ID
    auto_recharge_enabled BOOLEAN DEFAULT FALSE,
    auto_recharge_threshold DECIMAL(15,4) DEFAULT 0.00,
    auto_recharge_amount DECIMAL(15,4) DEFAULT 0.00,
    
    -- Invoice Configuration (for postpaid accounts)
    invoice_day_of_month INTEGER DEFAULT 1 CHECK (invoice_day_of_month >= 1 AND invoice_day_of_month <= 31),
    payment_terms_days INTEGER DEFAULT 30, -- Net payment terms
    
    -- Status and Metadata
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'closed')),
    billing_email VARCHAR(255),
    billing_address JSONB, -- Structured billing address
    tax_info JSONB, -- Tax ID, VAT number, etc.
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Ensure one billing account per organization
    UNIQUE (organization_id)
);

CREATE TRIGGER set_billing_accounts_timestamp
BEFORE UPDATE ON public.billing_accounts
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_billing_accounts_organization_id ON public.billing_accounts(organization_id);
CREATE INDEX idx_billing_accounts_stripe_customer_id ON public.billing_accounts(stripe_customer_id) WHERE stripe_customer_id IS NOT NULL;
CREATE INDEX idx_billing_accounts_billing_mode ON public.billing_accounts(billing_mode);
CREATE INDEX idx_billing_accounts_status ON public.billing_accounts(status);

-- #############################################################################
-- ## Payment Methods (Stripe Payment Method Storage)
-- #############################################################################

-- payment_methods: Store Stripe payment method metadata for organizations
CREATE TABLE public.payment_methods (
    payment_method_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    billing_account_id BIGINT NOT NULL REFERENCES public.billing_accounts(billing_account_id) ON DELETE CASCADE,
    
    -- Stripe Integration
    stripe_payment_method_id VARCHAR(255) NOT NULL UNIQUE, -- Stripe Payment Method ID
    
    -- Payment Method Details
    type VARCHAR(50) NOT NULL, -- 'card', 'bank_account', 'sepa_debit', etc.
    brand VARCHAR(50), -- 'visa', 'mastercard', etc. (for cards)
    last4 VARCHAR(4), -- Last 4 digits
    exp_month INTEGER, -- Expiration month (for cards)
    exp_year INTEGER, -- Expiration year (for cards)
    
    -- Bank Account Details (for ACH, SEPA, etc.)
    bank_name VARCHAR(255),
    account_holder_type VARCHAR(50), -- 'individual', 'company'
    
    -- Status and Configuration
    is_default BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'expired', 'failed')),
    
    -- Metadata
    nickname VARCHAR(255), -- User-friendly name
    metadata JSONB, -- Additional Stripe metadata
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_payment_methods_timestamp
BEFORE UPDATE ON public.payment_methods
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_payment_methods_organization_id ON public.payment_methods(organization_id);
CREATE INDEX idx_payment_methods_billing_account_id ON public.payment_methods(billing_account_id);
CREATE INDEX idx_payment_methods_stripe_pm_id ON public.payment_methods(stripe_payment_method_id);
CREATE INDEX idx_payment_methods_is_default ON public.payment_methods(organization_id, is_default) WHERE is_default = TRUE;

-- #############################################################################
-- ## Transactions (All Billing/Payment Events)
-- #############################################################################

-- transactions: Log all billing/payment events for audit and history
CREATE TABLE public.transactions (
    transaction_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    billing_account_id BIGINT NOT NULL REFERENCES public.billing_accounts(billing_account_id) ON DELETE CASCADE,
    
    -- Transaction Details
    type VARCHAR(50) NOT NULL CHECK (type IN (
        'recharge', 'debit', 'credit', 'refund', 'chargeback', 
        'invoice_payment', 'usage_charge', 'affiliate_payout', 
        'platform_fee', 'adjustment', 'transfer'
    )),
    amount DECIMAL(15,4) NOT NULL, -- Positive for credits, negative for debits
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Balance Tracking
    balance_before DECIMAL(15,4) NOT NULL,
    balance_after DECIMAL(15,4) NOT NULL,
    
    -- References
    reference_type VARCHAR(50), -- 'stripe_payment_intent', 'stripe_invoice', 'usage_record', etc.
    reference_id VARCHAR(255), -- External reference ID
    related_transaction_id BIGINT REFERENCES public.transactions(transaction_id), -- For refunds, adjustments
    
    -- Stripe Integration
    stripe_payment_intent_id VARCHAR(255),
    stripe_invoice_id VARCHAR(255),
    stripe_charge_id VARCHAR(255),
    
    -- Metadata
    description TEXT,
    metadata JSONB, -- Additional transaction data
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'completed' CHECK (status IN (
        'pending', 'completed', 'failed', 'cancelled', 'refunded'
    )),
    
    -- Timestamps
    processed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_transactions_timestamp
BEFORE UPDATE ON public.transactions
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_transactions_organization_id ON public.transactions(organization_id);
CREATE INDEX idx_transactions_billing_account_id ON public.transactions(billing_account_id);
CREATE INDEX idx_transactions_type ON public.transactions(type);
CREATE INDEX idx_transactions_status ON public.transactions(status);
CREATE INDEX idx_transactions_processed_at ON public.transactions(processed_at);
CREATE INDEX idx_transactions_stripe_payment_intent ON public.transactions(stripe_payment_intent_id) WHERE stripe_payment_intent_id IS NOT NULL;
CREATE INDEX idx_transactions_stripe_invoice ON public.transactions(stripe_invoice_id) WHERE stripe_invoice_id IS NOT NULL;
CREATE INDEX idx_transactions_reference ON public.transactions(reference_type, reference_id) WHERE reference_type IS NOT NULL AND reference_id IS NOT NULL;

-- #############################################################################
-- ## Invoices (Postpaid Billing)
-- #############################################################################

-- invoices: Store invoice information for postpaid billing
CREATE TABLE public.invoices (
    invoice_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    billing_account_id BIGINT NOT NULL REFERENCES public.billing_accounts(billing_account_id) ON DELETE CASCADE,
    
    -- Invoice Details
    invoice_number VARCHAR(100) NOT NULL UNIQUE, -- Human-readable invoice number
    
    -- Stripe Integration
    stripe_invoice_id VARCHAR(255) UNIQUE, -- Stripe Invoice ID
    
    -- Financial Details
    subtotal DECIMAL(15,4) NOT NULL DEFAULT 0.00,
    tax_amount DECIMAL(15,4) NOT NULL DEFAULT 0.00,
    total_amount DECIMAL(15,4) NOT NULL DEFAULT 0.00,
    amount_paid DECIMAL(15,4) NOT NULL DEFAULT 0.00,
    amount_due DECIMAL(15,4) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Billing Period
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    
    -- Dates
    invoice_date DATE NOT NULL,
    due_date DATE NOT NULL,
    paid_at TIMESTAMPTZ,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN (
        'draft', 'open', 'paid', 'void', 'uncollectible', 'overdue'
    )),
    
    -- Invoice Content
    line_items JSONB, -- Array of invoice line items
    notes TEXT,
    
    -- Metadata
    metadata JSONB,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_invoices_timestamp
BEFORE UPDATE ON public.invoices
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_invoices_organization_id ON public.invoices(organization_id);
CREATE INDEX idx_invoices_billing_account_id ON public.invoices(billing_account_id);
CREATE INDEX idx_invoices_stripe_invoice_id ON public.invoices(stripe_invoice_id) WHERE stripe_invoice_id IS NOT NULL;
CREATE INDEX idx_invoices_status ON public.invoices(status);
CREATE INDEX idx_invoices_due_date ON public.invoices(due_date);
CREATE INDEX idx_invoices_period ON public.invoices(period_start, period_end);
CREATE INDEX idx_invoices_invoice_number ON public.invoices(invoice_number);

-- #############################################################################
-- ## Usage Records (Daily Usage Calculation)
-- #############################################################################

-- usage_records: Store daily usage/spend calculations per organization
CREATE TABLE public.usage_records (
    usage_record_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    billing_account_id BIGINT NOT NULL REFERENCES public.billing_accounts(billing_account_id) ON DELETE CASCADE,
    
    -- Usage Period
    usage_date DATE NOT NULL,
    
    -- Usage Metrics
    clicks INTEGER DEFAULT 0,
    conversions INTEGER DEFAULT 0,
    impressions INTEGER DEFAULT 0,
    
    -- Financial Metrics
    advertiser_spend DECIMAL(15,4) DEFAULT 0.00, -- Amount spent by advertiser
    affiliate_payout DECIMAL(15,4) DEFAULT 0.00, -- Amount paid to affiliates
    platform_revenue DECIMAL(15,4) DEFAULT 0.00, -- Platform commission/fee
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Processing Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN (
        'pending', 'calculated', 'billed', 'paid', 'failed'
    )),
    
    -- Allocation Details
    allocated_at TIMESTAMPTZ, -- When funds were allocated
    billed_at TIMESTAMPTZ, -- When usage was billed
    
    -- Metadata
    campaign_breakdown JSONB, -- Per-campaign usage breakdown
    affiliate_breakdown JSONB, -- Per-affiliate payout breakdown
    metadata JSONB,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Ensure one record per organization per day
    UNIQUE (organization_id, usage_date)
);

CREATE TRIGGER set_usage_records_timestamp
BEFORE UPDATE ON public.usage_records
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_usage_records_organization_id ON public.usage_records(organization_id);
CREATE INDEX idx_usage_records_billing_account_id ON public.usage_records(billing_account_id);
CREATE INDEX idx_usage_records_usage_date ON public.usage_records(usage_date);
CREATE INDEX idx_usage_records_status ON public.usage_records(status);
CREATE INDEX idx_usage_records_allocated_at ON public.usage_records(allocated_at) WHERE allocated_at IS NOT NULL;

-- #############################################################################
-- ## Webhook Events (Stripe Webhook Processing)
-- #############################################################################

-- webhook_events: Store and track Stripe webhook events
CREATE TABLE public.webhook_events (
    webhook_event_id BIGSERIAL PRIMARY KEY,
    
    -- Stripe Event Details
    stripe_event_id VARCHAR(255) NOT NULL UNIQUE, -- Stripe Event ID
    event_type VARCHAR(100) NOT NULL, -- e.g., 'payment_intent.succeeded'
    
    -- Processing Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN (
        'pending', 'processed', 'failed', 'ignored'
    )),
    
    -- Event Data
    event_data JSONB NOT NULL, -- Full Stripe event payload
    
    -- Processing Details
    processed_at TIMESTAMPTZ,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    
    -- Related Records
    organization_id BIGINT REFERENCES public.organizations(organization_id),
    transaction_id BIGINT REFERENCES public.transactions(transaction_id),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER set_webhook_events_timestamp
BEFORE UPDATE ON public.webhook_events
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE INDEX idx_webhook_events_stripe_event_id ON public.webhook_events(stripe_event_id);
CREATE INDEX idx_webhook_events_event_type ON public.webhook_events(event_type);
CREATE INDEX idx_webhook_events_status ON public.webhook_events(status);
CREATE INDEX idx_webhook_events_organization_id ON public.webhook_events(organization_id) WHERE organization_id IS NOT NULL;
CREATE INDEX idx_webhook_events_created_at ON public.webhook_events(created_at);

-- #############################################################################
-- ## Functions and Triggers for Balance Management
-- #############################################################################

-- Function to update billing account balance
CREATE OR REPLACE FUNCTION update_billing_account_balance()
RETURNS TRIGGER AS $$
BEGIN
    -- Update the billing account balance when a transaction is inserted
    IF TG_OP = 'INSERT' AND NEW.status = 'completed' THEN
        UPDATE public.billing_accounts 
        SET balance = NEW.balance_after,
            updated_at = NOW()
        WHERE billing_account_id = NEW.billing_account_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update billing account balance
CREATE TRIGGER trigger_update_billing_account_balance
AFTER INSERT ON public.transactions
FOR EACH ROW
EXECUTE FUNCTION update_billing_account_balance();

-- Function to ensure only one default payment method per organization
CREATE OR REPLACE FUNCTION ensure_single_default_payment_method()
RETURNS TRIGGER AS $$
BEGIN
    -- If setting a payment method as default, unset all others for the same organization
    IF NEW.is_default = TRUE THEN
        UPDATE public.payment_methods 
        SET is_default = FALSE,
            updated_at = NOW()
        WHERE organization_id = NEW.organization_id 
        AND payment_method_id != NEW.payment_method_id 
        AND is_default = TRUE;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to ensure only one default payment method per organization
CREATE TRIGGER trigger_ensure_single_default_payment_method
BEFORE INSERT OR UPDATE ON public.payment_methods
FOR EACH ROW
EXECUTE FUNCTION ensure_single_default_payment_method();

-- #############################################################################
-- ## Initial Data and Configuration
-- #############################################################################

-- Create billing accounts for existing organizations
INSERT INTO public.billing_accounts (organization_id, billing_mode, currency)
SELECT organization_id, 'prepaid', 'USD'
FROM public.organizations
WHERE NOT EXISTS (
    SELECT 1 FROM public.billing_accounts 
    WHERE billing_accounts.organization_id = organizations.organization_id
);

-- #############################################################################
-- ## Comments and Documentation
-- #############################################################################

COMMENT ON TABLE public.billing_accounts IS 'Organization-level billing configuration and balance tracking';
COMMENT ON TABLE public.payment_methods IS 'Stripe payment method storage and metadata';
COMMENT ON TABLE public.transactions IS 'All billing and payment transaction history';
COMMENT ON TABLE public.invoices IS 'Invoice management for postpaid billing';
COMMENT ON TABLE public.usage_records IS 'Daily usage calculation and billing records';
COMMENT ON TABLE public.webhook_events IS 'Stripe webhook event processing and tracking';

COMMENT ON COLUMN public.billing_accounts.billing_mode IS 'prepaid: pay-as-you-go with balance, postpaid: monthly invoicing';
COMMENT ON COLUMN public.transactions.amount IS 'Positive for credits/income, negative for debits/expenses';
COMMENT ON COLUMN public.usage_records.advertiser_spend IS 'Total amount spent by advertiser for the day';
COMMENT ON COLUMN public.usage_records.affiliate_payout IS 'Total amount paid to affiliates for the day';
COMMENT ON COLUMN public.usage_records.platform_revenue IS 'Platform commission/fee for the day';

-- Analyze new tables for query optimization
ANALYZE public.billing_accounts;
ANALYZE public.payment_methods;
ANALYZE public.transactions;
ANALYZE public.invoices;
ANALYZE public.usage_records;
ANALYZE public.webhook_events;
