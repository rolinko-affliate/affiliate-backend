# Test Execution Report

**Generated:** Tue Aug  5 08:59:05 UTC 2025
**Total Tests:** 10
**Passed:** 10
**Failed:** 0
**Skipped:** 0
**Success Rate:** 100%

## Test Results

- **health_check** (0 s): PASS - API server is healthy
- **supabase_webhook** (0 s): PASS - Webhook processed successfully (profile creation may be async)
- **jwt_token_generation** (0 s): PASS - JWT tokens generated successfully for all user types
- **get_user_profile** (0 s): PASS - User profile retrieved successfully: admin@rolinko.com
- **profile_creation** (0 s): PASS - Profile created successfully: profile-test-1754384345@example.com
- **profile_update** (0 s): PASS - Profile updated successfully: Updated
- **unauthorized_access** (0 s): PASS - Unauthorized access properly rejected with 401
- **invalid_token_access** (0 s): PASS - Invalid token properly rejected with 401
- **rbac_enforcement** (0 s): PASS - RBAC properly enforced - regular user denied admin access
- **profile_upsert** (0 s): PASS - Profile upsert functionality working correctly
