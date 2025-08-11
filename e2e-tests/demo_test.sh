#!/bin/bash

# Demo Test Script
# 
# This script demonstrates the E2E test functionality with a quick health check
# and basic API validation to show how the test framework works.

set -e

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Load common utilities
source "$SCRIPT_DIR/utils/common.sh"
source "$SCRIPT_DIR/utils/jwt_helper.sh"

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}E2E Test Framework Demo${NC}"
echo -e "${CYAN}========================================${NC}"

# Demo Test 1: Health Check
start_test "demo_health_check" "Demonstrate basic API health check"

make_request "GET" "/health" "" "" "200"

if [[ $? -eq 0 ]]; then
    status=$(parse_json "$RESPONSE_BODY" "status")
    if [[ "$status" == "UP" ]]; then
        end_test "PASS" "API server is healthy and responding"
    else
        end_test "FAIL" "API server returned unexpected status: $status"
    fi
else
    end_test "FAIL" "Failed to reach API health endpoint"
fi

# Demo Test 2: JWT Token Generation
start_test "demo_jwt_generation" "Demonstrate JWT token generation"

admin_token=$(get_test_user_token "admin")

if [[ $? -eq 0 && -n "$admin_token" ]]; then
    log_debug "Generated admin token: ${admin_token:0:50}..."
    end_test "PASS" "JWT token generated successfully"
else
    end_test "FAIL" "Failed to generate JWT token"
fi

# Demo Test 3: Authenticated Request
start_test "demo_authenticated_request" "Demonstrate authenticated API request"

auth_header=$(create_auth_header "admin")

if [[ $? -eq 0 ]]; then
    make_request "GET" "/api/v1/users/me" "" "$auth_header" "200"
    
    if [[ $? -eq 0 ]]; then
        user_email=$(parse_json "$RESPONSE_BODY" "email")
        if [[ -n "$user_email" ]]; then
            end_test "PASS" "Authenticated request successful, user: $user_email"
        else
            end_test "PASS" "Authenticated request successful (user data may be mock)"
        fi
    else
        end_test "FAIL" "Authenticated request failed with status $RESPONSE_STATUS"
    fi
else
    end_test "FAIL" "Failed to create authentication header"
fi

# Demo Test 4: Error Handling
start_test "demo_error_handling" "Demonstrate error handling for invalid requests"

make_request "GET" "/api/v1/users/me" "" "" "401"

if [[ $? -eq 0 ]]; then
    end_test "PASS" "Unauthorized request properly rejected with 401"
else
    end_test "FAIL" "Error handling test failed - expected 401 but got $RESPONSE_STATUS"
fi

# Print demo summary
echo -e "\n${CYAN}========================================${NC}"
echo -e "${CYAN}Demo Test Summary${NC}"
echo -e "${CYAN}========================================${NC}"

print_test_summary

echo -e "\n${BLUE}Demo completed! This shows how the E2E test framework works:${NC}"
echo -e "${BLUE}1. Health checks validate server availability${NC}"
echo -e "${BLUE}2. JWT tokens are generated for different user roles${NC}"
echo -e "${BLUE}3. Authenticated requests are made with proper headers${NC}"
echo -e "${BLUE}4. Error scenarios are tested and validated${NC}"
echo -e "${BLUE}5. Comprehensive logging and reporting is provided${NC}"

echo -e "\n${GREEN}To run the full test suite:${NC}"
echo -e "${GREEN}  ./run_all_tests.sh${NC}"

echo -e "\n${GREEN}To run individual scenarios:${NC}"
echo -e "${GREEN}  ./run_scenario.sh authentication_flow${NC}"
echo -e "${GREEN}  ./run_scenario.sh advertiser_workflow${NC}"
echo -e "${GREEN}  ./run_scenario.sh affiliate_workflow${NC}"

echo -e "\n${GREEN}For more information, see:${NC}"
echo -e "${GREEN}  ./README.md${NC}"
echo -e "${GREEN}  ./E2E_TEST_DOCUMENTATION.md${NC}"

# Exit with appropriate code
if [[ $FAILED_TESTS -eq 0 ]]; then
    exit 0
else
    exit 1
fi