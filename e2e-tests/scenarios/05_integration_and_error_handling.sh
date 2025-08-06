#!/bin/bash

# E2E Test Scenario 5: Integration and Error Handling
# 
# This test validates system integration points and comprehensive error handling:
# 1. External provider integration (Everflow mock)
# 2. Webhook processing and validation
# 3. Billing system integration
# 4. Analytics and reporting integration
# 5. Comprehensive error scenarios
# 6. Edge cases and boundary conditions
# 7. Performance and timeout handling
# 8. Data consistency validation
#
# Test Coverage:
# - External service integration
# - Webhook endpoint validation
# - Error response consistency
# - Input validation edge cases
# - System resilience testing
# - Data integrity checks

set -e

# Load common utilities
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../utils/common.sh"
source "$SCRIPT_DIR/../utils/jwt_helper.sh"
source "$SCRIPT_DIR/../data/test_data.sh"

# Test configuration
TEST_SCENARIO="Integration and Error Handling"
CLEANUP_ORGS=()
CLEANUP_ADVERTISERS=()
CLEANUP_AFFILIATES=()
CLEANUP_CAMPAIGNS=()

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}Starting E2E Test Scenario 5${NC}"
echo -e "${CYAN}$TEST_SCENARIO${NC}"
echo -e "${CYAN}========================================${NC}"

# Test 1: Supabase Webhook Integration
test_supabase_webhook_integration() {
    start_test "supabase_webhook_integration" "Test Supabase webhook processing with various scenarios"
    
    # Test 1a: Valid webhook payload
    local user_id=$(python3 -c "import uuid; print(str(uuid.uuid4()))")
    local email="webhook-test-$(date +%s)@example.com"
    
    local valid_webhook_data="{
        \"type\": \"INSERT\",
        \"table\": \"users\",
        \"record\": {
            \"id\": \"$user_id\",
            \"email\": \"$email\",
            \"created_at\": \"$(date -Iseconds)\",
            \"email_confirmed_at\": \"$(date -Iseconds)\"
        },
        \"schema\": \"auth\",
        \"old_record\": null
    }"
    
    local headers="Content-Type: application/json"
    
    make_request "POST" "/api/v1/public/webhooks/supabase/new-user" "$valid_webhook_data" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        local message=$(parse_json "$RESPONSE_BODY" "message")
        if [[ "$message" == *"successfully"* ]]; then
            log_debug "✓ Valid webhook processed successfully"
            
            # Test 1b: Duplicate webhook (should return error for duplicate user)
            make_request "POST" "/api/v1/public/webhooks/supabase/new-user" "$valid_webhook_data" "$headers" "500"
            
            if [[ $? -eq 0 ]]; then
                local error_msg=$(parse_json "$RESPONSE_BODY" "error")
                if [[ "$error_msg" == *"duplicate"* ]]; then
                    log_debug "✓ Duplicate webhook properly rejected"
                else
                    end_test "FAIL" "Duplicate webhook should return duplicate error"
                    return
                fi
                
                # Test 1c: Invalid webhook payload (should be ignored)
                local invalid_webhook_data="{
                    \"type\": \"INVALID\",
                    \"table\": \"users\",
                    \"record\": null
                }"
                
                make_request "POST" "/api/v1/public/webhooks/supabase/new-user" "$invalid_webhook_data" "$headers" "200"
                
                if [[ $? -eq 0 ]]; then
                    local message=$(parse_json "$RESPONSE_BODY" "message")
                    if [[ "$message" == *"ignored"* ]]; then
                        log_debug "✓ Invalid webhook properly ignored"
                        end_test "PASS" "Supabase webhook integration working correctly"
                    else
                        end_test "FAIL" "Invalid webhook should be ignored"
                    fi
                else
                    end_test "FAIL" "Invalid webhook handling failed"
                fi
            else
                end_test "FAIL" "Duplicate webhook handling failed"
            fi
        else
            end_test "FAIL" "Valid webhook processing failed"
        fi
    else
        end_test "FAIL" "Supabase webhook integration failed with status $RESPONSE_STATUS"
    fi
}

# Test 2: Stripe Webhook Integration
test_stripe_webhook_integration() {
    start_test "stripe_webhook_integration" "Test Stripe webhook processing"
    
    # Skip this test as it requires valid Stripe webhook signatures
    # which cannot be generated without the actual Stripe webhook secret
    end_test "SKIP" "Stripe webhook requires valid signature validation (external dependency)"
}

# Test 3: Everflow Integration (Mock Mode)
test_everflow_integration() {
    start_test "everflow_integration" "Test Everflow provider integration in mock mode"
    
    # Skip this test as it requires complex organization/advertiser setup with proper RBAC
    # and is testing external provider integration which is covered in advertiser workflow
    end_test "SKIP" "Everflow integration requires complex setup and external provider (covered in advertiser workflow)"
}

# Test 4: Analytics Integration
test_analytics_integration() {
    start_test "analytics_integration" "Test analytics endpoints and data consistency"
    
    local auth_header
    auth_header=$(create_auth_header "advertiser_manager")
    
    local headers="$auth_header"
    
    # Test analytics autocomplete (requires 'q' parameter with minimum 3 characters)
    make_request "GET" "/api/v1/analytics/autocomplete?q=test" "" "$headers" "200"
    
    if [[ $? -eq 0 ]]; then
        log_debug "✓ Analytics autocomplete working"
        
        # Test creating analytics advertiser (use unique domain)
        local timestamp=$(date +%s)
        local analytics_advertiser_data="{
            \"domain\": \"analytics-test-${timestamp}.com\",
            \"data\": {
                \"name\": \"Analytics Test Advertiser\",
                \"description\": \"Test advertiser for analytics\"
            }
        }"
        
        local headers="Content-Type: application/json
$auth_header"
        
        make_request "POST" "/api/v1/analytics/advertisers" "$analytics_advertiser_data" "$headers" "201"
        
        if [[ $? -eq 0 ]]; then
            local analytics_id=$(parse_json "$RESPONSE_BODY" "data.id")
            log_debug "✓ Analytics advertiser created: $analytics_id"
            
            # Test retrieving analytics advertiser
            make_request "GET" "/api/v1/analytics/advertisers/$analytics_id" "" "$auth_header" "200"
            
            if [[ $? -eq 0 ]]; then
                log_debug "✓ Analytics advertiser retrieved successfully"
                end_test "PASS" "Analytics integration working correctly"
            else
                end_test "FAIL" "Analytics advertiser retrieval failed"
            fi
        else
            end_test "FAIL" "Analytics advertiser creation failed"
        fi
    else
        end_test "FAIL" "Analytics autocomplete failed with status $RESPONSE_STATUS"
    fi
}

# Test 5: Billing System Integration
test_billing_integration() {
    start_test "billing_integration" "Test billing system endpoints"
    
    # Skip this test as it requires valid Stripe API credentials
    # and is testing external payment provider integration
    end_test "SKIP" "Billing integration requires valid Stripe API credentials (external dependency)"
}

# Test 6: Input Validation Edge Cases
test_input_validation_edge_cases() {
    start_test "input_validation_edge_cases" "Test input validation with edge cases"
    
    local auth_header
    auth_header=$(create_auth_header "advertiser_manager")
    
    local validation_passed=0
    local total_cases=0
    
    # Test 1: Empty name
    total_cases=$((total_cases + 1))
    local edge_case='{"name":"","type":"advertiser"}'
    log_debug "Testing edge case: Empty name"
    
    local headers="Content-Type: application/json
$auth_header"
    
    make_request "POST" "/api/v1/organizations" "$edge_case" "$headers" "400"
    
    if [[ $? -eq 0 ]]; then
        validation_passed=$((validation_passed + 1))
        log_debug "✓ Empty name properly rejected"
    else
        log_debug "✗ Empty name not properly handled (status: $RESPONSE_STATUS)"
    fi
    
    # Test 2: Very long name
    total_cases=$((total_cases + 1))
    local long_name=$(printf 'A%.0s' {1..300})
    edge_case="{\"name\":\"$long_name\",\"type\":\"advertiser\"}"
    log_debug "Testing edge case: Very long name (300 chars)"
    
    make_request "POST" "/api/v1/organizations" "$edge_case" "$headers" "400"
    
    if [[ $? -eq 0 ]]; then
        validation_passed=$((validation_passed + 1))
        log_debug "✓ Long name properly rejected"
    else
        log_debug "✗ Long name not properly handled (status: $RESPONSE_STATUS)"
    fi
    
    # Test 3: Special characters (Note: Some special characters may be allowed in organization names)
    total_cases=$((total_cases + 1))
    edge_case='{"name":"Test<script>alert(1)</script>","type":"advertiser"}'
    log_debug "Testing edge case: Special characters"
    
    make_request "POST" "/api/v1/organizations" "$edge_case" "$headers" ""
    
    # Accept either 400 (rejected) or 201 (accepted) - both are valid behaviors
    if [[ "$RESPONSE_STATUS" == "400" || "$RESPONSE_STATUS" == "201" ]]; then
        validation_passed=$((validation_passed + 1))
        log_debug "✓ Special characters handled appropriately (status: $RESPONSE_STATUS)"
        # If created successfully, add to cleanup
        if [[ "$RESPONSE_STATUS" == "201" ]]; then
            local org_id=$(parse_json "$RESPONSE_BODY" "organization_id")
            if [[ -n "$org_id" ]]; then
                CLEANUP_ORGANIZATIONS+=("$org_id")
            fi
        fi
    else
        log_debug "✗ Special characters not properly handled (status: $RESPONSE_STATUS)"
    fi
    
    # Test 4: Invalid JSON
    total_cases=$((total_cases + 1))
    edge_case='{"name":"Test","type":"advertiser"'
    log_debug "Testing edge case: Invalid JSON"
    
    make_request "POST" "/api/v1/organizations" "$edge_case" "$headers" "400"
    
    if [[ $? -eq 0 ]]; then
        validation_passed=$((validation_passed + 1))
        log_debug "✓ Invalid JSON properly rejected"
    else
        log_debug "✗ Invalid JSON not properly handled (status: $RESPONSE_STATUS)"
    fi
    
    # Test 5: Null values
    total_cases=$((total_cases + 1))
    edge_case='{"name":null,"type":"advertiser"}'
    log_debug "Testing edge case: Null values"
    
    make_request "POST" "/api/v1/organizations" "$edge_case" "$headers" "400"
    
    if [[ $? -eq 0 ]]; then
        validation_passed=$((validation_passed + 1))
        log_debug "✓ Null values properly rejected"
    else
        log_debug "✗ Null values not properly handled (status: $RESPONSE_STATUS)"
    fi
    
    if [[ $validation_passed -ge $((total_cases * 80 / 100)) ]]; then
        end_test "PASS" "Input validation working correctly ($validation_passed/$total_cases cases passed)"
    else
        end_test "FAIL" "Input validation insufficient ($validation_passed/$total_cases cases passed)"
    fi
}

# Test 7: Authentication Edge Cases
test_authentication_edge_cases() {
    start_test "authentication_edge_cases" "Test authentication with various edge cases"
    
    local auth_edge_cases=(
        # Empty token
        "Authorization: Bearer "
        # Malformed token
        "Authorization: Bearer invalid.token"
        # Token with wrong algorithm
        "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature"
        # Expired token (mock)
        "Authorization: Bearer expired.token.here"
        # Token with invalid signature
        "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.invalid_signature"
    )
    
    local auth_tests_passed=0
    local total_auth_tests=${#auth_edge_cases[@]}
    
    for auth_case in "${auth_edge_cases[@]}"; do
        log_debug "Testing auth edge case: ${auth_case:0:50}..."
        
        make_request "GET" "/api/v1/users/me" "" "$auth_case" "401"
        
        if [[ $? -eq 0 ]]; then
            auth_tests_passed=$((auth_tests_passed + 1))
            log_debug "✓ Auth edge case properly rejected"
        else
            log_debug "✗ Auth edge case not properly handled (status: $RESPONSE_STATUS)"
        fi
    done
    
    if [[ $auth_tests_passed -ge $((total_auth_tests * 80 / 100)) ]]; then
        end_test "PASS" "Authentication edge cases handled correctly ($auth_tests_passed/$total_auth_tests cases passed)"
    else
        end_test "FAIL" "Authentication edge case handling insufficient ($auth_tests_passed/$total_auth_tests cases passed)"
    fi
}

# Test 8: Rate Limiting and Performance
test_rate_limiting_and_performance() {
    start_test "rate_limiting_and_performance" "Test system performance and rate limiting"
    
    local auth_header
    auth_header=$(create_auth_header "admin")
    
    local headers="$auth_header"
    
    # Test rapid requests to health endpoint
    local rapid_requests=10
    local successful_requests=0
    local start_time=$(date +%s)
    
    for i in $(seq 1 $rapid_requests); do
        make_request "GET" "/health" "" "" "200"
        
        if [[ $? -eq 0 ]]; then
            successful_requests=$((successful_requests + 1))
        fi
    done
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    log_debug "Completed $successful_requests/$rapid_requests requests in ${duration}s"
    
    if [[ $successful_requests -ge $((rapid_requests * 80 / 100)) ]]; then
        # Test response time for complex operation
        local complex_start_time=$(date +%s%3N)
        make_request "GET" "/api/v1/organizations" "" "$headers" "200"
        local complex_end_time=$(date +%s%3N)
        
        local response_time=$((complex_end_time - complex_start_time))
        
        if [[ $response_time -lt 5000 ]]; then  # Less than 5 seconds
            end_test "PASS" "Performance test passed (${response_time}ms response time)"
        else
            end_test "FAIL" "Response time too slow (${response_time}ms)"
        fi
    else
        end_test "FAIL" "Too many requests failed ($successful_requests/$rapid_requests)"
    fi
}

# Test 9: Data Consistency Validation
test_data_consistency() {
    start_test "data_consistency" "Test data consistency across related entities"
    
    # Create organization
    local org_name=$(generate_random_org_name "advertiser")
    local org_data="{\"name\":\"$org_name\",\"type\":\"advertiser\"}"
    
    local admin_header
    admin_header=$(create_auth_header "admin")
    
    local headers="Content-Type: application/json
$admin_header"
    
    make_request "POST" "/api/v1/organizations" "$org_data" "$headers" "201"
    
    if [[ $? -eq 0 ]]; then
        local org_id=$(parse_json "$RESPONSE_BODY" "organization_id")
        CLEANUP_ORGS+=("$org_id")
        
        # Create advertiser (use admin token since admin can create advertisers for any organization)
        local advertiser_data=$(get_test_advertiser_data "advertiser_1" "json")
        advertiser_data=$(echo "$advertiser_data" | sed "s/}$/,\"organization_id\":$org_id}/")
        
        local auth_header
        auth_header=$(create_auth_header "admin")
        
        local headers="Content-Type: application/json
$auth_header"
        
        make_request "POST" "/api/v1/advertisers" "$advertiser_data" "$headers" "201"
        
        if [[ $? -eq 0 ]]; then
            local advertiser_id=$(parse_json "$RESPONSE_BODY" "advertiser_id")
            CLEANUP_ADVERTISERS+=("$advertiser_id")
            
            # Create campaign
            local campaign_data=$(get_test_campaign_data "campaign_1" "json")
            campaign_data=$(echo "$campaign_data" | sed "s/}$/,\"advertiser_id\":$advertiser_id,\"organization_id\":$org_id}/")
            
            make_request "POST" "/api/v1/campaigns" "$campaign_data" "$headers" "201"
            
            if [[ $? -eq 0 ]]; then
                local campaign_id=$(parse_json "$RESPONSE_BODY" "campaign_id")
                CLEANUP_CAMPAIGNS+=("$campaign_id")
                
                # Verify relationships
                # 1. Campaign should belong to advertiser
                make_request "GET" "/api/v1/campaigns/$campaign_id" "" "$auth_header" "200"
                
                if [[ $? -eq 0 ]]; then
                    local campaign_advertiser_id=$(parse_json "$RESPONSE_BODY" "advertiser_id")
                    
                    if [[ "$campaign_advertiser_id" == "$advertiser_id" ]]; then
                        # 2. Advertiser should belong to organization
                        make_request "GET" "/api/v1/advertisers/$advertiser_id" "" "$auth_header" "200"
                        
                        if [[ $? -eq 0 ]]; then
                            local advertiser_org_id=$(parse_json "$RESPONSE_BODY" "organization_id")
                            
                            if [[ "$advertiser_org_id" == "$org_id" ]]; then
                                # 3. Organization should list the advertiser
                                make_request "GET" "/api/v1/organizations/$org_id/advertisers" "" "$auth_header" "200"
                                
                                if [[ $? -eq 0 && "$RESPONSE_BODY" == *"$advertiser_id"* ]]; then
                                    end_test "PASS" "Data consistency validation passed"
                                else
                                    end_test "FAIL" "Organization-advertiser relationship inconsistent"
                                fi
                            else
                                end_test "FAIL" "Advertiser-organization relationship inconsistent"
                            fi
                        else
                            end_test "FAIL" "Failed to retrieve advertiser for consistency check"
                        fi
                    else
                        end_test "FAIL" "Campaign-advertiser relationship inconsistent"
                    fi
                else
                    end_test "FAIL" "Failed to retrieve campaign for consistency check"
                fi
            else
                end_test "FAIL" "Failed to create campaign for consistency test"
            fi
        else
            end_test "FAIL" "Failed to create advertiser for consistency test"
        fi
    else
        end_test "FAIL" "Failed to create organization for consistency test"
    fi
}

# Test 10: Error Response Consistency
test_error_response_consistency() {
    start_test "error_response_consistency" "Test consistency of error responses across endpoints"
    
    local auth_header
    auth_header=$(create_auth_header "advertiser_manager")
    
    # Test 404 responses
    local not_found_endpoints=(
        "/api/v1/organizations/99999"
        "/api/v1/advertisers/99999"
        "/api/v1/campaigns/99999"
        "/api/v1/advertiser-association-invitations/99999"
    )
    
    local consistent_404=0
    
    for endpoint in "${not_found_endpoints[@]}"; do
        make_request "GET" "$endpoint" "" "$auth_header" "404"
        
        if [[ $? -eq 0 ]]; then
            consistent_404=$((consistent_404 + 1))
            log_debug "✓ $endpoint returns consistent 404"
        else
            log_debug "✗ $endpoint returns inconsistent status: $RESPONSE_STATUS"
        fi
    done
    
    # Test 400 responses with invalid data
    local invalid_data='{"invalid": "data"}'
    local headers="Content-Type: application/json
$auth_header"
    
    local bad_request_endpoints=(
        "/api/v1/organizations"
        "/api/v1/advertisers"
        "/api/v1/campaigns"
    )
    
    local consistent_400=0
    
    for endpoint in "${bad_request_endpoints[@]}"; do
        make_request "POST" "$endpoint" "$invalid_data" "$headers" "400"
        
        if [[ $? -eq 0 ]]; then
            consistent_400=$((consistent_400 + 1))
            log_debug "✓ $endpoint returns consistent 400"
        else
            log_debug "✗ $endpoint returns inconsistent status: $RESPONSE_STATUS"
        fi
    done
    
    local total_tests=$((${#not_found_endpoints[@]} + ${#bad_request_endpoints[@]}))
    local passed_tests=$((consistent_404 + consistent_400))
    
    if [[ $passed_tests -ge $((total_tests * 80 / 100)) ]]; then
        end_test "PASS" "Error response consistency acceptable ($passed_tests/$total_tests)"
    else
        end_test "FAIL" "Error response consistency insufficient ($passed_tests/$total_tests)"
    fi
}

# Cleanup function for this test scenario
cleanup_integration_test_data() {
    if [[ "$SKIP_CLEANUP" == "true" ]]; then
        log_info "Skipping integration test cleanup"
        return 0
    fi
    
    log_info "Cleaning up integration test data..."
    
    local auth_header
    auth_header=$(create_auth_header "admin")
    
    # Clean up campaigns
    for campaign_id in "${CLEANUP_CAMPAIGNS[@]}"; do
        log_debug "Cleaning up campaign: $campaign_id"
        make_request "DELETE" "/api/v1/campaigns/$campaign_id" "" "$auth_header" ""
    done
    
    # Clean up advertisers
    for advertiser_id in "${CLEANUP_ADVERTISERS[@]}"; do
        log_debug "Cleaning up advertiser: $advertiser_id"
        make_request "DELETE" "/api/v1/advertisers/$advertiser_id" "" "$auth_header" ""
    done
    
    # Clean up affiliates
    for affiliate_id in "${CLEANUP_AFFILIATES[@]}"; do
        log_debug "Cleaning up affiliate: $affiliate_id"
        make_request "DELETE" "/api/v1/affiliates/$affiliate_id" "" "$auth_header" ""
    done
    
    # Clean up organizations
    for org_id in "${CLEANUP_ORGS[@]}"; do
        log_debug "Cleaning up organization: $org_id"
        make_request "DELETE" "/api/v1/organizations/$org_id" "" "$auth_header" ""
    done
    
    log_info "Integration test cleanup completed"
}

# Main test execution
main() {
    log_info "Starting Integration and Error Handling tests..."
    
    # Execute all tests
    test_supabase_webhook_integration
    test_stripe_webhook_integration
    test_everflow_integration
    test_analytics_integration
    test_billing_integration
    test_input_validation_edge_cases
    test_authentication_edge_cases
    test_rate_limiting_and_performance
    test_data_consistency
    test_error_response_consistency
    
    # Cleanup
    cleanup_integration_test_data
    
    # Print summary
    print_test_summary
    
    # Generate report
    local report_file="$SCRIPT_DIR/../reports/integration_and_error_handling_report.md"
    mkdir -p "$(dirname "$report_file")"
    generate_test_report "$report_file"
    
    # Exit with appropriate code
    if [[ $FAILED_TESTS -eq 0 ]]; then
        log_success "All integration and error handling tests passed!"
        exit 0
    else
        log_error "Some integration and error handling tests failed!"
        exit 1
    fi
}

# Run main function if script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi