# Test Execution Report

**Generated:** Mon Aug  4 02:46:00 UTC 2025
**Total Tests:** 17
**Passed:** 12
**Failed:** 0
**Skipped:** 5
**Success Rate:** 70%

## Test Results

- **create_affiliate_org** (0 s): PASS - Affiliate organization created successfully: ID 37
- **create_affiliate** (0 s): PASS - Affiliate created successfully: Test Affiliate 1 (ID: 7)
- **get_affiliate** (0 s): PASS - Affiliate details retrieved successfully: Test Affiliate 1
- **update_affiliate** (0 s): PASS - Affiliate updated successfully
- **list_affiliates_by_org** (0 s): PASS - Affiliates listed successfully for organization
- **search_affiliates** (0 s): PASS - Affiliate search completed successfully
- **get_visible_campaigns** (0 s): PASS - Visible campaigns retrieved successfully
- **create_association_request** (1 s): PASS - Association request created successfully: ID 9
- **list_associations** (0 s): PASS - Associations listed successfully
- **use_association_invitation** (0 s): PASS - Association invitation used successfully
- **generate_affiliate_tracking_links** (0 s): SKIP - Database schema mismatch: provider_campaign_id column missing (known issue)
- **affiliate_analytics** (0 s): SKIP - Requires external provider data sync (Everflow integration)
- **create_messaging_conversation** (0 s): SKIP - Requires organization membership setup for messaging context
- **add_message_to_conversation** (0 s): SKIP - Depends on conversation creation (organization membership required)
- **get_messaging_conversations** (0 s): SKIP - Requires organization membership setup for messaging context
- **error_handling_invalid_affiliate_data** (0 s): PASS - Invalid affiliate data properly rejected with 400
- **advertiser_manager_permission** (0 s): PASS - Advertiser manager properly denied affiliate creation with 403
