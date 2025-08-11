# Test Execution Report

**Generated:** Tue Aug  5 08:28:20 UTC 2025
**Total Tests:** 15
**Passed:** 15
**Failed:** 0
**Skipped:** 0
**Success Rate:** 100%

## Test Results

- **create_advertiser_org** (0 s): PASS - Advertiser organization created successfully: ID 11
- **list_organizations** (0 s): PASS - Organizations listed successfully
- **get_organization_details** (0 s): PASS - Organization details retrieved successfully
- **create_advertiser** (0 s): PASS - Advertiser created successfully: Test Advertiser 1 (ID: 3)
- **get_advertiser** (0 s): PASS - Advertiser details retrieved successfully: Test Advertiser 1
- **update_advertiser** (0 s): PASS - Advertiser updated successfully
- **create_campaign** (0 s): PASS - Campaign created successfully: Summer Sale Campaign (ID: 3)
- **list_campaigns_by_advertiser** (0 s): PASS - Campaigns listed successfully for advertiser
- **create_association_invitation** (1 s): PASS - Association invitation created successfully: ID 4
- **generate_invitation_link** (0 s): PASS - Invitation link generated successfully
- **list_association_invitations** (0 s): PASS - Association invitations listed successfully
- **everflow_sync** (0 s): PASS - Everflow sync completed successfully (mock mode)
- **advertiser_analytics** (0 s): PASS - Advertiser analytics endpoint working (no data for new advertiser)
- **error_handling_invalid_data** (0 s): PASS - Invalid data properly rejected with 400
- **permission_denied** (0 s): PASS - Permission denied properly enforced with 403
