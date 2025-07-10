-- #############################################################################
-- ## Rollback Organization-Level Billing and Payment System
-- ## This migration removes all billing and payment system tables and functions
-- #############################################################################

-- Drop triggers first
DROP TRIGGER IF EXISTS trigger_update_billing_account_balance ON public.transactions;
DROP TRIGGER IF EXISTS trigger_ensure_single_default_payment_method ON public.payment_methods;

-- Drop functions
DROP FUNCTION IF EXISTS update_billing_account_balance();
DROP FUNCTION IF EXISTS ensure_single_default_payment_method();

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS public.webhook_events;
DROP TABLE IF EXISTS public.usage_records;
DROP TABLE IF EXISTS public.invoices;
DROP TABLE IF EXISTS public.transactions;
DROP TABLE IF EXISTS public.payment_methods;
DROP TABLE IF EXISTS public.billing_accounts;