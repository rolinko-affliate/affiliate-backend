# Test Execution Report

**Generated:** Tue Aug  5 08:59:08 UTC 2025
**Total Tests:** 18
**Passed:** 18
**Failed:** 0
**Skipped:** 0
**Success Rate:** 100%

## Test Results

- **setup_test_organizations** (0 s): PASS - Test organizations created successfully (Advertiser: 33, Affiliate1: 34, Affiliate2: 35)
- **create_basic_invitation** (0 s): PASS - Basic invitation created successfully: ID 6
- **create_restricted_invitation** (0 s): PASS - Restricted invitation created successfully: ID 7
- **list_invitations** (0 s): PASS - Invitations listed successfully with created invitations present
- **get_invitation_details** (0 s): PASS - Invitation details retrieved successfully with related data
- **generate_invitation_link** (1 s): PASS - Invitation link generated successfully: http://localhost:50220/invitations/4f6e6bcbe457271a5dd1290770256c1b6bcdd9054c247097841cce6478739abb
- **public_invitation_access** (0 s): PASS - Public invitation access successful with complete data
- **use_invitation_successfully** (0 s): PASS - Invitation used successfully, association created
- **use_invitation_duplicate** (0 s): PASS - Duplicate invitation use properly handled
- **use_restricted_invitation_allowed** (0 s): PASS - Restricted invitation use handled correctly (association already exists)
- **use_restricted_invitation_disallowed** (0 s): PASS - Restricted invitation properly denied for disallowed affiliate
- **get_invitation_usage_history** (0 s): PASS - Usage history retrieved successfully with usage records
- **update_invitation** (0 s): PASS - Invitation updated successfully
- **create_expired_invitation** (0 s): PASS - Expired invitation created successfully: ID 8
- **use_expired_invitation** (0 s): PASS - Expired invitation properly rejected
- **invalid_invitation_token** (0 s): PASS - Invalid invitation token properly rejected with 404
- **missing_required_fields** (0 s): PASS - Missing required fields properly rejected with 400
- **affiliate_manager_cannot_create_invitations** (0 s): PASS - Affiliate manager properly denied invitation creation with 403
