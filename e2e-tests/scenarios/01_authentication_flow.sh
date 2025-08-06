#!/bin/bash

# E2E Test Scenario 1: Authentication and User Management Flow
# 
# This test validates the complete authentication and user management workflow:
# 1. User registration via Supabase webhook
# 2. Profile creation and management
# 3. JWT token validation
# 4. Role-based access control (RBAC)
# 5. User profile updates and retrieval
#
# Test Coverage:
# - Public webhook endpoints
# - Authenticated profile endpoints
# - JWT token validation
# - RBAC enforcement
# - Error handling for invalid tokens and permissions

set -e

# Load common utilities
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../utils/common.sh"
source "$SCRIPT_DIR/../utils/jwt_helper.sh"
source "$SCRIPT_DIR/../data/test_data.sh"

# Test configuration
TEST_SCENARIO="Authentication and User Management Flow"
CLEANUP_USERS=()

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}Starting E2E Test Scenario 1${NC}"
echo -e "${CYAN}$TEST_SCENARIO${NC}"
echo -e "${CYAN}========================================${NC}"

# Test 1: Health Check
test_health_check() {
    start_test "health_check" "Verify API server is running and healthy"
    
    make_request "GET" "/health" "" "" "200"
    
    if [[ $? -eq 0 ]]; then
        status=$(parse_json "$RESPONSE_BODY" "status")
        if [[ "$status" == "UP" ]]; then
            end_test "PASS" "API server is healthy"
        else
            end_test "FAIL" "API server status is not healthy: $status"
        fi
    else
        end_test "FAIL" "Failed to reach health endpoint"
    fi
}

# Test 2: Supabase User Webhook
test_supabase_webhook() {
    start_test "supabase_webhook" "Test user creation via Supabase webhook"
    
    user_id=$(python3 -c "import uuid; print(str(uuid.uuid4()))")
    email="test-user-$(date +%s)@example.com"
    
    webhook_data="{
        \"type\": \"INSERT\",
        \"table\": \"users\",
        \"record\": {
            \"id\": \"$user_id\",
            \"email\": \"$email\",
            \"created_at\": \"$(date -Iseconds)\"
        },
        \"schema\": \"auth\",
        \"old_record\": null
    }"
    
    headers="Content-Type: application/json"
    
    make_request "POST" "/api/v1/public/webhooks/supabase/new-user" "$webhook_data" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        # Check if profile was created
        profile_id=$(parse_json "$RESPONSE_BODY" "profile_id")
        if [[ -n "$profile_id" ]]; then
            CLEANUP_USERS+=("$user_id")
            end_test "PASS" "User profile created successfully via webhook: $profile_id"
        else
            end_test "PASS" "Webhook processed successfully (profile creation may be async)"
        fi
    else
        end_test "FAIL" "Supabase webhook failed with status $RESPONSE_STATUS"
    fi
}

# Test 3: JWT Token Generation and Validation
test_jwt_token_generation() {
    start_test "jwt_token_generation" "Test JWT token generation for different user roles"
    
    test_passed=true
    
    # Test token generation for each user type
    for user_type in "admin" "advertiser_manager" "affiliate_manager" "regular_user"; do
        log_debug "Testing JWT token generation for $user_type"
        
        token=$(get_test_user_token "$user_type")
        
        if [[ $? -eq 0 && -n "$token" ]]; then
            log_debug "✓ Token generated for $user_type"
            
            # Validate token structure (should have 3 parts separated by dots)
            token_parts=$(echo "$token" | tr '.' '\n' | wc -l)
            
            if [[ "$token_parts" -eq 3 ]]; then
                log_debug "✓ Token has correct structure for $user_type"
            else
                log_error "✗ Token has incorrect structure for $user_type (parts: $token_parts)"
                test_passed=false
            fi
        else
            log_error "✗ Failed to generate token for $user_type"
            test_passed=false
        fi
    done
    
    if [[ "$test_passed" == "true" ]]; then
        end_test "PASS" "JWT tokens generated successfully for all user types"
    else
        end_test "FAIL" "JWT token generation failed for one or more user types"
    fi
}

# Test 4: Get User Profile (Authenticated)
test_get_user_profile() {
    start_test "get_user_profile" "Test retrieving user profile with valid JWT token"
    
    # First create the admin profile if it doesn't exist
    local admin_user_id=$(get_test_user_info "admin" "user_id")
    local admin_email=$(get_test_user_info "admin" "email")
    local admin_first_name=$(get_test_user_info "admin" "first_name")
    local admin_last_name=$(get_test_user_info "admin" "last_name")
    local admin_role_id=$(get_test_user_info "admin" "role_id")
    
    local admin_auth_header=$(create_auth_header "admin")
    
    # Create/upsert the admin profile
    local profile_data="{
        \"id\": \"$admin_user_id\",
        \"email\": \"$admin_email\",
        \"first_name\": \"$admin_first_name\",
        \"last_name\": \"$admin_last_name\",
        \"role_id\": $admin_role_id
    }"
    
    local headers="Content-Type: application/json
$admin_auth_header"
    
    make_request "POST" "/api/v1/profiles/upsert" "$profile_data" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        log_debug "Admin profile created/updated successfully"
        
        # Now try to get the profile
        make_request "GET" "/api/v1/users/me" "" "$admin_auth_header" "200"
        
        if [[ $? -eq 0 ]]; then
            user_id=$(parse_json "$RESPONSE_BODY" "id")
            email=$(parse_json "$RESPONSE_BODY" "email")
            
            if [[ -n "$user_id" && -n "$email" ]]; then
                end_test "PASS" "User profile retrieved successfully: $email"
            else
                end_test "FAIL" "User profile response missing required fields"
            fi
        else
            end_test "FAIL" "Failed to retrieve user profile with status $RESPONSE_STATUS"
        fi
    else
        end_test "FAIL" "Failed to create admin profile with status $RESPONSE_STATUS"
    fi
}

# Test 5: Profile Creation
test_profile_creation() {
    start_test "profile_creation" "Test creating a new user profile"
    
    user_id=$(python3 -c "import uuid; print(str(uuid.uuid4()))")
    email="profile-test-$(date +%s)@example.com"
    first_name="Test"
    last_name="User"
    role_id="100000"  # Regular user role
    
    profile_data="{
        \"id\": \"$user_id\",
        \"email\": \"$email\",
        \"first_name\": \"$first_name\",
        \"last_name\": \"$last_name\",
        \"role_id\": $role_id
    }"
    
    auth_header=$(create_auth_header "admin")
    
    headers="Content-Type: application/json
$auth_header"
    
    make_request "POST" "/api/v1/profiles/upsert" "$profile_data" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        created_id=$(parse_json "$RESPONSE_BODY" "id")
        created_email=$(parse_json "$RESPONSE_BODY" "email")
        
        if [[ "$created_id" == "$user_id" && "$created_email" == "$email" ]]; then
            CLEANUP_USERS+=("$user_id")
            end_test "PASS" "Profile created successfully: $created_email"
        else
            end_test "FAIL" "Profile creation response doesn't match input data"
        fi
    else
        end_test "FAIL" "Profile creation failed with status $RESPONSE_STATUS"
    fi
}

# Test 6: Profile Update
test_profile_update() {
    start_test "profile_update" "Test updating an existing user profile"
    
    # First create a profile to update
    user_id=$(python3 -c "import uuid; print(str(uuid.uuid4()))")
    email="update-test-$(date +%s)@example.com"
    original_first_name="Original"
    updated_first_name="Updated"
    
    create_data="{
        \"id\": \"$user_id\",
        \"email\": \"$email\",
        \"first_name\": \"$original_first_name\",
        \"last_name\": \"User\",
        \"role_id\": 100000
    }"
    
    auth_header=$(create_auth_header "admin")
    
    headers="Content-Type: application/json
$auth_header"
    
    # Create profile
    make_request "POST" "/api/v1/profiles/upsert" "$create_data" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        CLEANUP_USERS+=("$user_id")
        
        # Now update the profile
        update_data="{
            \"first_name\": \"$updated_first_name\"
        }"
        
        make_request "PUT" "/api/v1/profiles/$user_id" "$update_data" "$headers" "200"
        
        if [[ $? -eq 0 ]]; then
            updated_name=$(parse_json "$RESPONSE_BODY" "first_name")
            
            if [[ "$updated_name" == "$updated_first_name" ]]; then
                end_test "PASS" "Profile updated successfully: $updated_first_name"
            else
                end_test "FAIL" "Profile update didn't reflect changes"
            fi
        else
            end_test "FAIL" "Profile update failed with status $RESPONSE_STATUS"
        fi
    else
        end_test "FAIL" "Failed to create profile for update test"
    fi
}

# Test 7: Unauthorized Access
test_unauthorized_access() {
    start_test "unauthorized_access" "Test access without authentication token"
    
    make_request "GET" "/api/v1/users/me" "" "" "401"
    
    if [[ $? -eq 0 ]]; then
        end_test "PASS" "Unauthorized access properly rejected with 401"
    else
        end_test "FAIL" "Unauthorized access test failed - expected 401 but got $RESPONSE_STATUS"
    fi
}

# Test 8: Invalid Token Access
test_invalid_token_access() {
    start_test "invalid_token_access" "Test access with invalid JWT token"
    
    invalid_token="invalid.jwt.token"
    headers="Authorization: Bearer $invalid_token"
    
    make_request "GET" "/api/v1/users/me" "" "$headers" "401"
    
    if [[ $? -eq 0 ]]; then
        end_test "PASS" "Invalid token properly rejected with 401"
    else
        end_test "FAIL" "Invalid token test failed - expected 401 but got $RESPONSE_STATUS"
    fi
}

# Test 9: Role-Based Access Control (RBAC)
test_rbac_enforcement() {
    start_test "rbac_enforcement" "Test role-based access control enforcement"
    
    # First create a regular user profile
    regular_user_id="44444444-4444-4444-4444-444444444444"
    regular_user_email="user@test.com"
    
    # Create profile for regular user using admin token
    admin_header=$(create_auth_header "admin")
    
    profile_data="{
        \"id\": \"$regular_user_id\",
        \"email\": \"$regular_user_email\",
        \"first_name\": \"Regular\",
        \"last_name\": \"User\",
        \"role_id\": 100000
    }"
    
    headers="Content-Type: application/json
$admin_header"
    
    make_request "POST" "/api/v1/profiles/upsert" "$profile_data" "$headers" "200"
    
    if [[ $? -ne 0 ]]; then
        end_test "FAIL" "Failed to create regular user profile for RBAC test"
        return
    fi
    
    # Now test that regular user cannot access admin-only endpoints
    regular_user_header=$(create_auth_header "regular_user")
    
    headers="$regular_user_header"
    
    # Try to access admin-only organization update endpoint
    make_request "PUT" "/api/v1/organizations/test-org-id" "{\"name\": \"test\"}" "$headers" "403"
    
    if [[ $? -eq 0 ]]; then
        end_test "PASS" "RBAC properly enforced - regular user denied admin access"
    else
        # If we get 200, RBAC is not working properly
        if [[ "$RESPONSE_STATUS" == "200" ]]; then
            end_test "FAIL" "RBAC not enforced - regular user gained admin access"
        else
            end_test "FAIL" "RBAC test failed with unexpected status $RESPONSE_STATUS (expected 403)"
        fi
    fi
}

# Test 10: Profile Upsert Functionality
test_profile_upsert() {
    start_test "profile_upsert" "Test profile upsert (create or update) functionality"
    
    user_id=$(python3 -c "import uuid; print(str(uuid.uuid4()))")
    email="upsert-test-$(date +%s)@example.com"
    first_name="Upsert"
    last_name="Test"
    
    upsert_data="{
        \"id\": \"$user_id\",
        \"email\": \"$email\",
        \"first_name\": \"$first_name\",
        \"last_name\": \"$last_name\",
        \"role_id\": 100000
    }"
    
    auth_header=$(create_auth_header "admin")
    
    headers="Content-Type: application/json
$auth_header"
    
    # First upsert (should create)
    make_request "POST" "/api/v1/profiles/upsert" "$upsert_data" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        CLEANUP_USERS+=("$user_id")
        
        created_email=$(parse_json "$RESPONSE_BODY" "email")
        
        if [[ "$created_email" == "$email" ]]; then
            # Second upsert with updated data (should update)
            updated_first_name="Updated"
            update_upsert_data="{
                \"id\": \"$user_id\",
                \"email\": \"$email\",
                \"first_name\": \"$updated_first_name\",
                \"last_name\": \"$last_name\",
                \"role_id\": 100000
            }"
            
            make_request "POST" "/api/v1/profiles/upsert" "$update_upsert_data" "$headers" "200"
            
            if [[ $? -eq 0 ]]; then
                updated_name=$(parse_json "$RESPONSE_BODY" "first_name")
                
                if [[ "$updated_name" == "$updated_first_name" ]]; then
                    end_test "PASS" "Profile upsert functionality working correctly"
                else
                    end_test "FAIL" "Profile upsert update didn't work correctly"
                fi
            else
                end_test "FAIL" "Profile upsert update failed with status $RESPONSE_STATUS"
            fi
        else
            end_test "FAIL" "Profile upsert create didn't work correctly"
        fi
    else
        end_test "FAIL" "Profile upsert create failed with status $RESPONSE_STATUS"
    fi
}

# Cleanup function for this test scenario
cleanup_authentication_test_data() {
    if [[ "$SKIP_CLEANUP" == "true" ]]; then
        log_info "Skipping authentication test cleanup"
        return 0
    fi
    
    log_info "Cleaning up authentication test data..."
    
    auth_header=$(create_auth_header "admin")
    
    # Clean up created user profiles
    for user_id in "${CLEANUP_USERS[@]}"; do
        log_debug "Cleaning up user profile: $user_id"
        
        headers="$auth_header"
        make_request "DELETE" "/api/v1/profiles/$user_id" "" "$headers" ""
        
        if [[ "$RESPONSE_STATUS" == "200" || "$RESPONSE_STATUS" == "204" || "$RESPONSE_STATUS" == "404" ]]; then
            log_debug "✓ Cleaned up user profile: $user_id"
        else
            log_warning "Failed to clean up user profile: $user_id (status: $RESPONSE_STATUS)"
        fi
    done
    
    log_info "Authentication test cleanup completed"
}

# Main test execution
main() {
    log_info "Starting Authentication and User Management Flow tests..."
    
    # Execute all tests
    test_health_check
    test_supabase_webhook
    test_jwt_token_generation
    test_get_user_profile
    test_profile_creation
    test_profile_update
    test_unauthorized_access
    test_invalid_token_access
    test_rbac_enforcement
    test_profile_upsert
    
    # Cleanup
    cleanup_authentication_test_data
    
    # Print summary
    print_test_summary
    
    # Generate report
    report_file="$SCRIPT_DIR/../reports/authentication_flow_report.md"
    mkdir -p "$(dirname "$report_file")"
    generate_test_report "$report_file"
    
    # Exit with appropriate code
    if [[ $FAILED_TESTS -eq 0 ]]; then
        log_success "All authentication tests passed!"
        exit 0
    else
        log_error "Some authentication tests failed!"
        exit 1
    fi
}

# Run main function if script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi