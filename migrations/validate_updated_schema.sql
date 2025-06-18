-- #############################################################################
-- ## Affiliate Platform Database Schema Validation Script
-- ## Updated for single organization (rolinko) and admin user setup
-- #############################################################################

-- Check that all expected tables exist
SELECT 
    'Table Validation' as check_type,
    schemaname,
    tablename,
    '✅ EXISTS' as status
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY tablename;

-- Check that the trigger function exists
SELECT 
    'Function Validation' as check_type,
    'public' as schemaname,
    proname as function_name,
    '✅ EXISTS' as status
FROM pg_proc 
WHERE proname = 'trigger_set_timestamp';

-- Check indexes
SELECT 
    'Index Validation' as check_type,
    schemaname,
    indexname,
    '✅ EXISTS' as status
FROM pg_indexes 
WHERE schemaname = 'public'
ORDER BY indexname;

-- Check foreign key relationships
SELECT 
    'Foreign Key Validation' as check_type,
    tc.table_name,
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name,
    '✅ LINKED' as status
FROM 
    information_schema.table_constraints AS tc 
    JOIN information_schema.key_column_usage AS kcu
      ON tc.constraint_name = kcu.constraint_name
      AND tc.table_schema = kcu.table_schema
    JOIN information_schema.constraint_column_usage AS ccu
      ON ccu.constraint_name = tc.constraint_name
      AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY' 
  AND tc.table_schema = 'public'
ORDER BY tc.table_name, kcu.column_name;

-- Check that seed data was inserted
SELECT 
    'Seed Data Validation' as check_type,
    'roles' as table_name,
    COUNT(*) as record_count,
    CASE 
        WHEN COUNT(*) >= 4 THEN '✅ SEEDED'
        ELSE '❌ MISSING'
    END as status
FROM public.roles
UNION ALL
SELECT 
    'Seed Data Validation' as check_type,
    'organizations (rolinko only)' as table_name,
    COUNT(*) as record_count,
    CASE 
        WHEN COUNT(*) = 1 AND MAX(name) = 'rolinko' AND MAX(type) = 'platform_owner' THEN '✅ SEEDED'
        ELSE '❌ INCORRECT'
    END as status
FROM public.organizations
UNION ALL
SELECT 
    'Seed Data Validation' as check_type,
    'profiles (admin user)' as table_name,
    COUNT(*) as record_count,
    CASE 
        WHEN COUNT(*) = 1 AND MAX(email) = 'admin@rolinko.com' THEN '✅ SEEDED'
        ELSE '❌ INCORRECT'
    END as status
FROM public.profiles
WHERE id = '4cbe2452-88aa-4429-9145-b527d9eebfbf';

-- Detailed verification of the admin user setup
SELECT 
    'Admin User Details' as check_type,
    p.id as user_id,
    p.email,
    p.first_name,
    p.last_name,
    o.name as organization_name,
    o.type as organization_type,
    r.name as role_name,
    '✅ VERIFIED' as status
FROM public.profiles p
JOIN public.organizations o ON p.organization_id = o.organization_id
JOIN public.roles r ON p.role_id = r.role_id
WHERE p.id = '4cbe2452-88aa-4429-9145-b527d9eebfbf';

-- Summary counts
SELECT 
    'Summary' as check_type,
    'Total Tables' as item,
    COUNT(*) as count,
    '📊 COUNT' as status
FROM pg_tables 
WHERE schemaname = 'public'
UNION ALL
SELECT 
    'Summary' as check_type,
    'Total Indexes' as item,
    COUNT(*) as count,
    '📊 COUNT' as status
FROM pg_indexes 
WHERE schemaname = 'public'
UNION ALL
SELECT 
    'Summary' as check_type,
    'Total Organizations' as item,
    COUNT(*) as count,
    '📊 COUNT' as status
FROM public.organizations
UNION ALL
SELECT 
    'Summary' as check_type,
    'Total Profiles' as item,
    COUNT(*) as count,
    '📊 COUNT' as status
FROM public.profiles;