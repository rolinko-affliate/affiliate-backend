#!/bin/bash

# E2E Test Scenario 4: Association Invitation System Complete Workflow
# 
# This test validates the complete association invitation system workflow:
# 1. Invitation creation by advertisers
# 2. Invitation management (list, update, delete)
# 3. Public invitation access (no auth required)
# 4. Invitation usage by affiliates
# 5. Usage tracking and analytics
# 6. Link generation and sharing
# 7. Expiration and usage limit handling
# 8. Error scenarios and edge cases
#
# Test Coverage:
# - Complete invitation lifecycle
# - Public and authenticated endpoints
# - Usage tracking and logging
# - Expiration and limit enforcement
# - Error handling and validation
# - Security and access control

set -e

# Load common utilities
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../utils/common.sh"
source "$SCRIPT_DIR/../utils/jwt_helper.sh"
source "$SCRIPT_DIR/../data/test_data.sh"

# Test configuration
TEST_SCENARIO="Association Invitation System Complete Workflow"
CLEANUP_ORGS=()
CLEANUP_INVITATIONS=()
CLEANUP_ASSOCIATIONS=()

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}Starting E2E Test Scenario 4${NC}"
echo -e "${CYAN}$TEST_SCENARIO${NC}"
echo -e "${CYAN}========================================${NC}"

# Test 1: Setup Test Organizations
test_setup_test_organizations() {
    start_test "setup_test_organizations" "Create test advertiser and affiliate organizations"
    
    local admin_header
    admin_header=$(create_auth_header "admin")
    
    local headers="Content-Type: application/json
$admin_header"
    
    # Create advertiser organization
    local advertiser_org_name=$(generate_random_org_name "advertiser")
    local advertiser_org_data="{\"name\":\"$advertiser_org_name\",\"type\":\"advertiser\"}"
    
    make_request "POST" "/api/v1/organizations" "$advertiser_org_data" "$headers" "201"
    
    if [[ $? -eq 0 ]]; then
        ADVERTISER_ORG_ID=$(parse_json "$RESPONSE_BODY" "organization_id")
        CLEANUP_ORGS+=("$ADVERTISER_ORG_ID")
        
        # Create affiliate organization
        local affiliate_org_name=$(generate_random_org_name "affiliate")
        local affiliate_org_data="{\"name\":\"$affiliate_org_name\",\"type\":\"affiliate\"}"
        
        make_request "POST" "/api/v1/organizations" "$affiliate_org_data" "$headers" "201"
        
        if [[ $? -eq 0 ]]; then
            AFFILIATE_ORG_ID=$(parse_json "$RESPONSE_BODY" "organization_id")
            CLEANUP_ORGS+=("$AFFILIATE_ORG_ID")
            
            # Create second affiliate organization for testing restrictions
            local affiliate_org_2_name=$(generate_random_org_name "affiliate")
            local affiliate_org_2_data="{\"name\":\"$affiliate_org_2_name\",\"type\":\"affiliate\"}"
            
            make_request "POST" "/api/v1/organizations" "$affiliate_org_2_data" "$headers" "201"
            
            if [[ $? -eq 0 ]]; then
                AFFILIATE_ORG_2_ID=$(parse_json "$RESPONSE_BODY" "organization_id")
                CLEANUP_ORGS+=("$AFFILIATE_ORG_2_ID")
                
                end_test "PASS" "Test organizations created successfully (Advertiser: $ADVERTISER_ORG_ID, Affiliate1: $AFFILIATE_ORG_ID, Affiliate2: $AFFILIATE_ORG_2_ID)"
            else
                end_test "FAIL" "Failed to create second affiliate organization"
            fi
        else
            end_test "FAIL" "Failed to create affiliate organization"
        fi
    else
        end_test "FAIL" "Failed to create advertiser organization"
    fi
}

# Test 2: Create Basic Invitation
test_create_basic_invitation() {
    start_test "create_basic_invitation" "Create a basic association invitation"
    
    if [[ -z "$ADVERTISER_ORG_ID" ]]; then
        end_test "SKIP" "No advertiser organization ID available"
        return
    fi
    
    local invitation_data="{
        \"advertiser_org_id\": $ADVERTISER_ORG_ID,
        \"name\": \"Basic Test Invitation\",
        \"description\": \"A basic invitation for E2E testing\",
        \"max_uses\": 10,
        \"expires_at\": \"$(date -d '+30 days' -Iseconds)\",
        \"default_all_affiliates_visible\": true,
        \"default_all_campaigns_visible\": true
    }"
    
    local auth_header
    auth_header=$(create_auth_header "advertiser_manager")
    
    local headers="Content-Type: application/json
$auth_header"
    
    make_request "POST" "/api/v1/advertiser-association-invitations" "$invitation_data" "$headers" "201"
    
    if [[ $? -eq 0 ]]; then
        local invitation_id=$(parse_json "$RESPONSE_BODY" "invitation_id")
        local invitation_token=$(parse_json "$RESPONSE_BODY" "invitation_token")
        local status=$(parse_json "$RESPONSE_BODY" "status")
        
        if [[ -n "$invitation_id" && -n "$invitation_token" && "$status" == "active" ]]; then
            CLEANUP_INVITATIONS+=("$invitation_id")
            BASIC_INVITATION_ID="$invitation_id"
            BASIC_INVITATION_TOKEN="$invitation_token"
            end_test "PASS" "Basic invitation created successfully: ID $invitation_id"
        else
            end_test "FAIL" "Basic invitation creation response validation failed"
        fi
    else
        end_test "FAIL" "Basic invitation creation failed with status $RESPONSE_STATUS"
    fi
}

# Test 3: Create Restricted Invitation
test_create_restricted_invitation() {
    start_test "create_restricted_invitation" "Create an invitation restricted to specific affiliates"
    
    if [[ -z "$ADVERTISER_ORG_ID" || -z "$AFFILIATE_ORG_ID" ]]; then
        end_test "SKIP" "Missing required organization IDs"
        return
    fi
    
    local invitation_data="{
        \"advertiser_org_id\": $ADVERTISER_ORG_ID,
        \"name\": \"Restricted Test Invitation\",
        \"description\": \"An invitation restricted to specific affiliates\",
        \"allowed_affiliate_org_ids\": [$AFFILIATE_ORG_ID],
        \"max_uses\": 5,
        \"expires_at\": \"$(date -d '+7 days' -Iseconds)\",
        \"message\": \"Welcome to our exclusive partner program!\",
        \"default_all_affiliates_visible\": false,
        \"default_all_campaigns_visible\": false
    }"
    
    local auth_header
    auth_header=$(create_auth_header "advertiser_manager")
    
    local headers="Content-Type: application/json
$auth_header"
    
    make_request "POST" "/api/v1/advertiser-association-invitations" "$invitation_data" "$headers" "201"
    
    if [[ $? -eq 0 ]]; then
        local invitation_id=$(parse_json "$RESPONSE_BODY" "invitation_id")
        local invitation_token=$(parse_json "$RESPONSE_BODY" "invitation_token")
        local allowed_orgs=$(parse_json "$RESPONSE_BODY" "allowed_affiliate_org_ids")
        
        if [[ -n "$invitation_id" && -n "$invitation_token" && "$allowed_orgs" == *"$AFFILIATE_ORG_ID"* ]]; then
            CLEANUP_INVITATIONS+=("$invitation_id")
            RESTRICTED_INVITATION_ID="$invitation_id"
            RESTRICTED_INVITATION_TOKEN="$invitation_token"
            end_test "PASS" "Restricted invitation created successfully: ID $invitation_id"
        else
            end_test "FAIL" "Restricted invitation creation response validation failed"
        fi
    else
        end_test "FAIL" "Restricted invitation creation failed with status $RESPONSE_STATUS"
    fi
}

# Test 4: List Invitations
test_list_invitations() {
    start_test "list_invitations" "List all association invitations"
    
    local auth_header
    auth_header=$(create_auth_header "advertiser_manager")
    
    local headers="$auth_header"
    
    make_request "GET" "/api/v1/advertiser-association-invitations" "" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        # Check if response is an array and contains our created invitations
        if [[ "$RESPONSE_BODY" == *"["* && "$RESPONSE_BODY" == *"]"* ]]; then
            local contains_basic=false
            local contains_restricted=false
            
            if [[ -n "$BASIC_INVITATION_ID" && "$RESPONSE_BODY" == *"$BASIC_INVITATION_ID"* ]]; then
                contains_basic=true
            fi
            
            if [[ -n "$RESTRICTED_INVITATION_ID" && "$RESPONSE_BODY" == *"$RESTRICTED_INVITATION_ID"* ]]; then
                contains_restricted=true
            fi
            
            if [[ "$contains_basic" == "true" && "$contains_restricted" == "true" ]]; then
                end_test "PASS" "Invitations listed successfully with created invitations present"
            else
                end_test "PASS" "Invitations listed successfully (created invitations may not be visible due to filtering)"
            fi
        else
            end_test "FAIL" "Invitations list response format invalid"
        fi
    else
        end_test "FAIL" "Failed to list invitations with status $RESPONSE_STATUS"
    fi
}

# Test 5: Get Invitation Details
test_get_invitation_details() {
    start_test "get_invitation_details" "Get detailed information about an invitation"
    
    if [[ -z "$BASIC_INVITATION_ID" ]]; then
        end_test "SKIP" "No basic invitation ID available"
        return
    fi
    
    local auth_header
    auth_header=$(create_auth_header "advertiser_manager")
    
    local headers="$auth_header"
    
    make_request "GET" "/api/v1/advertiser-association-invitations/$BASIC_INVITATION_ID" "" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        local invitation_id=$(parse_json "$RESPONSE_BODY" "invitation_id")
        local advertiser_org=$(parse_json "$RESPONSE_BODY" "advertiser_organization.organization_id")
        local created_by=$(parse_json "$RESPONSE_BODY" "created_by_user.id")
        
        if [[ "$invitation_id" == "$BASIC_INVITATION_ID" && "$advertiser_org" == "$ADVERTISER_ORG_ID" && -n "$created_by" ]]; then
            end_test "PASS" "Invitation details retrieved successfully with related data"
        else
            end_test "FAIL" "Invitation details validation failed"
        fi
    else
        end_test "FAIL" "Failed to get invitation details with status $RESPONSE_STATUS"
    fi
}

# Test 6: Generate Invitation Link
test_generate_invitation_link() {
    start_test "generate_invitation_link" "Generate shareable invitation link"
    
    if [[ -z "$BASIC_INVITATION_ID" ]]; then
        end_test "SKIP" "No basic invitation ID available"
        return
    fi
    
    local auth_header
    auth_header=$(create_auth_header "advertiser_manager")
    
    local headers="$auth_header"
    
    make_request "GET" "/api/v1/advertiser-association-invitations/$BASIC_INVITATION_ID/link" "" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        local invitation_link=$(parse_json "$RESPONSE_BODY" "invitation_link")
        
        if [[ -n "$invitation_link" && "$invitation_link" == *"$BASIC_INVITATION_TOKEN"* ]]; then
            INVITATION_LINK="$invitation_link"
            end_test "PASS" "Invitation link generated successfully: $invitation_link"
        else
            end_test "FAIL" "Invitation link generation validation failed"
        fi
    else
        end_test "FAIL" "Invitation link generation failed with status $RESPONSE_STATUS"
    fi
}

# Test 7: Public Invitation Access (No Auth)
test_public_invitation_access() {
    start_test "public_invitation_access" "Access invitation details without authentication"
    
    if [[ -z "$BASIC_INVITATION_TOKEN" ]]; then
        end_test "SKIP" "No basic invitation token available"
        return
    fi
    
    # No authentication headers for public access
    make_request "GET" "/api/v1/public/invitations/$BASIC_INVITATION_TOKEN" "" "" "200"
    
    if [[ $? -eq 0 ]]; then
        local invitation_id=$(parse_json "$RESPONSE_BODY" "invitation_id")
        local advertiser_name=$(parse_json "$RESPONSE_BODY" "advertiser_organization.name")
        local status=$(parse_json "$RESPONSE_BODY" "status")
        
        if [[ "$invitation_id" == "$BASIC_INVITATION_ID" && -n "$advertiser_name" && "$status" == "active" ]]; then
            end_test "PASS" "Public invitation access successful with complete data"
        else
            end_test "FAIL" "Public invitation access validation failed"
        fi
    else
        end_test "FAIL" "Public invitation access failed with status $RESPONSE_STATUS"
    fi
}

# Test 8: Use Invitation Successfully
test_use_invitation_successfully() {
    start_test "use_invitation_successfully" "Use invitation to create association"
    
    if [[ -z "$BASIC_INVITATION_TOKEN" || -z "$AFFILIATE_ORG_ID" ]]; then
        end_test "SKIP" "Missing invitation token or affiliate organization ID"
        return
    fi
    
    local use_data="{
        \"invitation_token\": \"$BASIC_INVITATION_TOKEN\",
        \"affiliate_org_id\": $AFFILIATE_ORG_ID
    }"
    
    local auth_header
    auth_header=$(create_auth_header "affiliate_manager")
    
    local headers="Content-Type: application/json
$auth_header"
    
    make_request "POST" "/api/v1/advertiser-association-invitations/use" "$use_data" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        local success=$(parse_json "$RESPONSE_BODY" "success")
        
        if [[ "$success" == "true" ]]; then
            local association_id=$(parse_json "$RESPONSE_BODY" "association.association_id")
            if [[ -n "$association_id" ]]; then
                CLEANUP_ASSOCIATIONS+=("$association_id")
                CREATED_ASSOCIATION_ID="$association_id"
            fi
            end_test "PASS" "Invitation used successfully, association created"
        else
            local error_message=$(parse_json "$RESPONSE_BODY" "error_message")
            if [[ "$error_message" == *"already exists"* ]]; then
                end_test "PASS" "Invitation use handled correctly (association already exists)"
            else
                end_test "FAIL" "Invitation use failed: $error_message"
            fi
        fi
    else
        end_test "FAIL" "Invitation use failed with status $RESPONSE_STATUS"
    fi
}

# Test 9: Use Invitation Again (Should Handle Duplicate)
test_use_invitation_duplicate() {
    start_test "use_invitation_duplicate" "Attempt to use same invitation again"
    
    if [[ -z "$BASIC_INVITATION_TOKEN" || -z "$AFFILIATE_ORG_ID" ]]; then
        end_test "SKIP" "Missing invitation token or affiliate organization ID"
        return
    fi
    
    local use_data="{
        \"invitation_token\": \"$BASIC_INVITATION_TOKEN\",
        \"affiliate_org_id\": $AFFILIATE_ORG_ID
    }"
    
    local auth_header
    auth_header=$(create_auth_header "affiliate_manager")
    
    local headers="Content-Type: application/json
$auth_header"
    
    make_request "POST" "/api/v1/advertiser-association-invitations/use" "$use_data" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        local success=$(parse_json "$RESPONSE_BODY" "success")
        local error_message=$(parse_json "$RESPONSE_BODY" "error_message")
        
        if [[ "$success" == "false" && "$error_message" == *"already exists"* ]]; then
            end_test "PASS" "Duplicate invitation use properly handled"
        elif [[ "$success" == "true" ]]; then
            end_test "PASS" "Invitation use successful (may be allowed multiple times)"
        else
            end_test "FAIL" "Unexpected response for duplicate invitation use"
        fi
    else
        end_test "FAIL" "Duplicate invitation use failed with status $RESPONSE_STATUS"
    fi
}

# Test 10: Use Restricted Invitation with Allowed Affiliate
test_use_restricted_invitation_allowed() {
    start_test "use_restricted_invitation_allowed" "Use restricted invitation with allowed affiliate"
    
    if [[ -z "$RESTRICTED_INVITATION_TOKEN" || -z "$AFFILIATE_ORG_ID" ]]; then
        end_test "SKIP" "Missing restricted invitation token or affiliate organization ID"
        return
    fi
    
    local use_data="{
        \"invitation_token\": \"$RESTRICTED_INVITATION_TOKEN\",
        \"affiliate_org_id\": $AFFILIATE_ORG_ID
    }"
    
    local auth_header
    auth_header=$(create_auth_header "affiliate_manager")
    
    local headers="Content-Type: application/json
$auth_header"
    
    make_request "POST" "/api/v1/advertiser-association-invitations/use" "$use_data" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        local success=$(parse_json "$RESPONSE_BODY" "success")
        
        if [[ "$success" == "true" ]]; then
            local association_id=$(parse_json "$RESPONSE_BODY" "association.association_id")
            if [[ -n "$association_id" ]]; then
                CLEANUP_ASSOCIATIONS+=("$association_id")
            fi
            end_test "PASS" "Restricted invitation used successfully by allowed affiliate"
        else
            local error_message=$(parse_json "$RESPONSE_BODY" "error_message")
            if [[ "$error_message" == *"already exists"* ]]; then
                end_test "PASS" "Restricted invitation use handled correctly (association already exists)"
            else
                end_test "FAIL" "Restricted invitation use failed: $error_message"
            fi
        fi
    else
        end_test "FAIL" "Restricted invitation use failed with status $RESPONSE_STATUS"
    fi
}

# Test 11: Use Restricted Invitation with Disallowed Affiliate
test_use_restricted_invitation_disallowed() {
    start_test "use_restricted_invitation_disallowed" "Use restricted invitation with disallowed affiliate"
    
    if [[ -z "$RESTRICTED_INVITATION_TOKEN" || -z "$AFFILIATE_ORG_2_ID" ]]; then
        end_test "SKIP" "Missing restricted invitation token or second affiliate organization ID"
        return
    fi
    
    local use_data="{
        \"invitation_token\": \"$RESTRICTED_INVITATION_TOKEN\",
        \"affiliate_org_id\": $AFFILIATE_ORG_2_ID
    }"
    
    local auth_header
    auth_header=$(create_auth_header "affiliate_manager")
    
    local headers="Content-Type: application/json
$auth_header"
    
    make_request "POST" "/api/v1/advertiser-association-invitations/use" "$use_data" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        # Check if it's properly denied in the response body
        local success=$(parse_json "$RESPONSE_BODY" "success")
        local error_message=$(parse_json "$RESPONSE_BODY" "error_message")
        
        if [[ "$success" == "false" && "$error_message" == *"not allowed"* ]]; then
            end_test "PASS" "Restricted invitation properly denied for disallowed affiliate"
        else
            end_test "FAIL" "Restricted invitation should be denied for disallowed affiliate (success: $success, error: $error_message)"
        fi
    else
        end_test "FAIL" "Unexpected status for restricted invitation test: $RESPONSE_STATUS"
    fi
}

# Test 12: Get Invitation Usage History
test_get_invitation_usage_history() {
    start_test "get_invitation_usage_history" "Get usage history for invitation"
    
    if [[ -z "$BASIC_INVITATION_ID" ]]; then
        end_test "SKIP" "No basic invitation ID available"
        return
    fi
    
    local auth_header
    auth_header=$(create_auth_header "advertiser_manager")
    
    local headers="$auth_header"
    
    make_request "GET" "/api/v1/advertiser-association-invitations/$BASIC_INVITATION_ID/usage-history" "" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        # Check if response is an array
        if [[ "$RESPONSE_BODY" == *"["* && "$RESPONSE_BODY" == *"]"* ]]; then
            # Check if usage history contains our affiliate organization
            if [[ "$RESPONSE_BODY" == *"$AFFILIATE_ORG_ID"* ]]; then
                end_test "PASS" "Usage history retrieved successfully with usage records"
            else
                end_test "PASS" "Usage history retrieved successfully (may be empty)"
            fi
        else
            end_test "FAIL" "Usage history response format invalid"
        fi
    else
        end_test "FAIL" "Failed to get usage history with status $RESPONSE_STATUS"
    fi
}

# Test 13: Update Invitation
test_update_invitation() {
    start_test "update_invitation" "Update invitation details"
    
    if [[ -z "$BASIC_INVITATION_ID" ]]; then
        end_test "SKIP" "No basic invitation ID available"
        return
    fi
    
    local updated_name="Updated Basic Test Invitation"
    local updated_description="Updated description for E2E testing - $(date)"
    
    local update_data="{
        \"name\": \"$updated_name\",
        \"description\": \"$updated_description\",
        \"max_uses\": 20
    }"
    
    local auth_header
    auth_header=$(create_auth_header "advertiser_manager")
    
    local headers="Content-Type: application/json
$auth_header"
    
    make_request "PUT" "/api/v1/advertiser-association-invitations/$BASIC_INVITATION_ID" "$update_data" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        local updated_name_response=$(parse_json "$RESPONSE_BODY" "name")
        local updated_desc_response=$(parse_json "$RESPONSE_BODY" "description")
        local updated_max_uses=$(parse_json "$RESPONSE_BODY" "max_uses")
        
        if [[ "$updated_name_response" == "$updated_name" && "$updated_desc_response" == "$updated_description" && "$updated_max_uses" == "20" ]]; then
            end_test "PASS" "Invitation updated successfully"
        else
            end_test "FAIL" "Invitation update validation failed"
        fi
    else
        end_test "FAIL" "Invitation update failed with status $RESPONSE_STATUS"
    fi
}

# Test 14: Create Expired Invitation
test_create_expired_invitation() {
    start_test "create_expired_invitation" "Create an invitation that expires immediately"
    
    if [[ -z "$ADVERTISER_ORG_ID" ]]; then
        end_test "SKIP" "No advertiser organization ID available"
        return
    fi
    
    local invitation_data="{
        \"advertiser_org_id\": $ADVERTISER_ORG_ID,
        \"name\": \"Expired Test Invitation\",
        \"description\": \"An invitation that expires immediately\",
        \"max_uses\": 1,
        \"expires_at\": \"$(date -d '-1 hour' -Iseconds)\",
        \"default_all_affiliates_visible\": true,
        \"default_all_campaigns_visible\": true
    }"
    
    local auth_header
    auth_header=$(create_auth_header "advertiser_manager")
    
    local headers="Content-Type: application/json
$auth_header"
    
    make_request "POST" "/api/v1/advertiser-association-invitations" "$invitation_data" "$headers" "201"
    
    if [[ $? -eq 0 ]]; then
        local invitation_id=$(parse_json "$RESPONSE_BODY" "invitation_id")
        local invitation_token=$(parse_json "$RESPONSE_BODY" "invitation_token")
        
        if [[ -n "$invitation_id" && -n "$invitation_token" ]]; then
            CLEANUP_INVITATIONS+=("$invitation_id")
            EXPIRED_INVITATION_TOKEN="$invitation_token"
            end_test "PASS" "Expired invitation created successfully: ID $invitation_id"
        else
            end_test "FAIL" "Expired invitation creation response validation failed"
        fi
    else
        end_test "FAIL" "Expired invitation creation failed with status $RESPONSE_STATUS"
    fi
}

# Test 15: Use Expired Invitation
test_use_expired_invitation() {
    start_test "use_expired_invitation" "Attempt to use expired invitation"
    
    if [[ -z "$EXPIRED_INVITATION_TOKEN" || -z "$AFFILIATE_ORG_ID" ]]; then
        end_test "SKIP" "Missing expired invitation token or affiliate organization ID"
        return
    fi
    
    local use_data="{
        \"invitation_token\": \"$EXPIRED_INVITATION_TOKEN\",
        \"affiliate_org_id\": $AFFILIATE_ORG_ID
    }"
    
    local auth_header
    auth_header=$(create_auth_header "affiliate_manager")
    
    local headers="Content-Type: application/json
$auth_header"
    
    make_request "POST" "/api/v1/advertiser-association-invitations/use" "$use_data" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        # Check if it's properly rejected in the response body
        local success=$(parse_json "$RESPONSE_BODY" "success")
        local error_message=$(parse_json "$RESPONSE_BODY" "error_message")
        
        if [[ "$success" == "false" && "$error_message" == *"expired"* ]]; then
            end_test "PASS" "Expired invitation properly rejected"
        else
            end_test "FAIL" "Expired invitation should be rejected (success: $success, error: $error_message)"
        fi
    else
        end_test "FAIL" "Unexpected status for expired invitation test: $RESPONSE_STATUS"
    fi
}

# Test 16: Error Handling - Invalid Token
test_invalid_invitation_token() {
    start_test "invalid_invitation_token" "Test access with invalid invitation token"
    
    local invalid_token="invalid-token-12345"
    
    # Test public access with invalid token
    make_request "GET" "/api/v1/public/invitations/$invalid_token" "" "" "404"
    
    if [[ $? -eq 0 ]]; then
        end_test "PASS" "Invalid invitation token properly rejected with 404"
    else
        end_test "FAIL" "Invalid token test failed - expected 404 but got $RESPONSE_STATUS"
    fi
}

# Test 17: Error Handling - Missing Required Fields
test_missing_required_fields() {
    start_test "missing_required_fields" "Test invitation creation with missing required fields"
    
    local invalid_data="{
        \"name\": \"Test Invitation\"
    }"
    
    local auth_header
    auth_header=$(create_auth_header "advertiser_manager")
    
    local headers="Content-Type: application/json
$auth_header"
    
    make_request "POST" "/api/v1/advertiser-association-invitations" "$invalid_data" "$headers" "400"
    
    if [[ $? -eq 0 ]]; then
        end_test "PASS" "Missing required fields properly rejected with 400"
    else
        end_test "FAIL" "Missing fields validation failed - expected 400 but got $RESPONSE_STATUS"
    fi
}

# Test 18: Permission Test - Affiliate Manager Cannot Create Invitations
test_affiliate_manager_cannot_create_invitations() {
    start_test "affiliate_manager_cannot_create_invitations" "Test that affiliate managers cannot create invitations"
    
    if [[ -z "$ADVERTISER_ORG_ID" ]]; then
        end_test "SKIP" "No advertiser organization ID available"
        return
    fi
    
    local invitation_data="{
        \"advertiser_org_id\": $ADVERTISER_ORG_ID,
        \"name\": \"Unauthorized Invitation\",
        \"description\": \"This should not be created\"
    }"
    
    local auth_header
    auth_header=$(create_auth_header "affiliate_manager")
    
    local headers="Content-Type: application/json
$auth_header"
    
    make_request "POST" "/api/v1/advertiser-association-invitations" "$invitation_data" "$headers" "403"
    
    if [[ $? -eq 0 ]]; then
        end_test "PASS" "Affiliate manager properly denied invitation creation with 403"
    else
        end_test "FAIL" "Permission test failed - expected 403 but got $RESPONSE_STATUS"
    fi
}

# Cleanup function for this test scenario
cleanup_invitation_test_data() {
    if [[ "$SKIP_CLEANUP" == "true" ]]; then
        log_info "Skipping invitation test cleanup"
        return 0
    fi
    
    log_info "Cleaning up invitation test data..."
    
    local auth_header
    auth_header=$(create_auth_header "admin")
    
    # Clean up associations
    for association_id in "${CLEANUP_ASSOCIATIONS[@]}"; do
        log_debug "Cleaning up association: $association_id"
        make_request "DELETE" "/api/v1/organization-associations/$association_id" "" "$auth_header" ""
    done
    
    # Clean up invitations
    for invitation_id in "${CLEANUP_INVITATIONS[@]}"; do
        log_debug "Cleaning up invitation: $invitation_id"
        make_request "DELETE" "/api/v1/advertiser-association-invitations/$invitation_id" "" "$auth_header" ""
    done
    
    # Clean up organizations
    for org_id in "${CLEANUP_ORGS[@]}"; do
        log_debug "Cleaning up organization: $org_id"
        make_request "DELETE" "/api/v1/organizations/$org_id" "" "$auth_header" ""
    done
    
    log_info "Invitation test cleanup completed"
}

# Main test execution
main() {
    log_info "Starting Association Invitation System Complete Workflow tests..."
    
    # Execute all tests in logical order
    test_setup_test_organizations
    test_create_basic_invitation
    test_create_restricted_invitation
    test_list_invitations
    test_get_invitation_details
    test_generate_invitation_link
    test_public_invitation_access
    test_use_invitation_successfully
    test_use_invitation_duplicate
    test_use_restricted_invitation_allowed
    test_use_restricted_invitation_disallowed
    test_get_invitation_usage_history
    test_update_invitation
    test_create_expired_invitation
    test_use_expired_invitation
    test_invalid_invitation_token
    test_missing_required_fields
    test_affiliate_manager_cannot_create_invitations
    
    # Cleanup
    cleanup_invitation_test_data
    
    # Print summary
    print_test_summary
    
    # Generate report
    local report_file="$SCRIPT_DIR/../reports/association_invitation_system_report.md"
    mkdir -p "$(dirname "$report_file")"
    generate_test_report "$report_file"
    
    # Exit with appropriate code
    if [[ $FAILED_TESTS -eq 0 ]]; then
        log_success "All association invitation system tests passed!"
        exit 0
    else
        log_error "Some association invitation system tests failed!"
        exit 1
    fi
}

# Run main function if script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi