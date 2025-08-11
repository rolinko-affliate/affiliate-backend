#!/bin/bash

# Common utilities for E2E tests
# This file contains shared functions used across all test scenarios

# Configuration
API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
TEST_TIMEOUT="${TEST_TIMEOUT:-30}"
VERBOSE_OUTPUT="${VERBOSE_OUTPUT:-false}"
SKIP_CLEANUP="${SKIP_CLEANUP:-false}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# Test results array
declare -a TEST_RESULTS=()

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    if [[ "$VERBOSE_OUTPUT" == "true" ]]; then
        echo -e "${PURPLE}[DEBUG]${NC} $1"
    fi
}

# Test execution functions
start_test() {
    local test_name="$1"
    local description="$2"
    
    echo -e "\n${CYAN}=== Starting Test: $test_name ===${NC}"
    echo -e "${CYAN}Description: $description${NC}"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    CURRENT_TEST_NAME="$test_name"
    CURRENT_TEST_START_TIME=$(date +%s)
}

end_test() {
    local status="$1"
    local message="$2"
    
    local end_time=$(date +%s)
    local duration=$((end_time - CURRENT_TEST_START_TIME))
    
    case "$status" in
        "PASS")
            PASSED_TESTS=$((PASSED_TESTS + 1))
            log_success "Test '$CURRENT_TEST_NAME' PASSED in ${duration}s: $message"
            TEST_RESULTS+=("PASS|$CURRENT_TEST_NAME|$duration|$message")
            ;;
        "FAIL")
            FAILED_TESTS=$((FAILED_TESTS + 1))
            log_error "Test '$CURRENT_TEST_NAME' FAILED in ${duration}s: $message"
            TEST_RESULTS+=("FAIL|$CURRENT_TEST_NAME|$duration|$message")
            ;;
        "SKIP")
            SKIPPED_TESTS=$((SKIPPED_TESTS + 1))
            log_warning "Test '$CURRENT_TEST_NAME' SKIPPED in ${duration}s: $message"
            TEST_RESULTS+=("SKIP|$CURRENT_TEST_NAME|$duration|$message")
            ;;
    esac
}

# HTTP request functions
make_request() {
    local method="$1"
    local endpoint="$2"
    local data="$3"
    local headers="$4"
    local expected_status="$5"
    
    local url="${API_BASE_URL}${endpoint}"
    local curl_cmd="curl -s -w '%{http_code}|%{time_total}' --max-time $TEST_TIMEOUT"
    
    # Add method
    curl_cmd="$curl_cmd -X $method"
    
    # Add headers
    if [[ -n "$headers" ]]; then
        while IFS= read -r header; do
            curl_cmd="$curl_cmd -H '$header'"
        done <<< "$headers"
    fi
    
    # Add data for POST/PUT requests
    if [[ -n "$data" && ("$method" == "POST" || "$method" == "PUT" || "$method" == "PATCH") ]]; then
        curl_cmd="$curl_cmd -d '$data'"
    fi
    
    # Add URL
    curl_cmd="$curl_cmd '$url'"
    
    log_debug "Executing: $curl_cmd"
    
    # Execute request and capture response
    local response
    response=$(eval "$curl_cmd")
    local exit_code=$?
    
    if [[ $exit_code -ne 0 ]]; then
        log_error "Request failed with exit code $exit_code"
        return 1
    fi
    
    # Parse response (format: body|status_code|time_total)
    # The curl -w format appends status_code|time_total to the body
    # Extract time_total (last part after |)
    local time_total="${response##*|}"
    # Remove time_total from response
    local response_without_time="${response%|*}"
    
    # Now we need to extract the 3-digit status code from the end
    # Status codes are always 3 digits, so extract last 3 characters
    local status_code="${response_without_time: -3}"
    # Body is everything except the last 3 characters
    local body="${response_without_time%???}"
    
    log_debug "Response Status: $status_code"
    log_debug "Response Time: ${time_total}s"
    log_debug "Response Body: $body"
    
    # Store response data globally for test access
    RESPONSE_BODY="$body"
    RESPONSE_STATUS="$status_code"
    RESPONSE_TIME="$time_total"
    
    # Validate expected status if provided
    if [[ -n "$expected_status" && "$status_code" != "$expected_status" ]]; then
        log_error "Expected status $expected_status but got $status_code"
        return 1
    fi
    
    return 0
}

# JSON parsing functions
parse_json() {
    local json="$1"
    local key="$2"
    
    # Simple JSON parsing using grep and sed (for basic cases)
    # For complex JSON parsing, consider using jq if available
    if command -v jq >/dev/null 2>&1; then
        # Handle null, boolean, and string values properly
        local value=$(echo "$json" | jq -r ".$key")
        if [[ "$value" == "null" ]]; then
            echo ""
        else
            echo "$value"
        fi
    else
        # Fallback parsing for simple cases
        echo "$json" | grep -o "\"$key\":[^,}]*" | sed "s/\"$key\"://g" | sed 's/[",]//g' | xargs
    fi
}

# JWT token generation for testing
generate_test_jwt() {
    local user_id="$1"
    local role="$2"
    local email="$3"
    
    # This is a simplified JWT generation for testing
    # In real scenarios, use proper JWT libraries
    local header='{"alg":"HS256","typ":"JWT"}'
    local payload="{\"sub\":\"$user_id\",\"email\":\"$email\",\"role\":\"$role\",\"exp\":$(($(date +%s) + 3600))}"
    
    # Base64 encode (simplified - in real tests, use proper JWT encoding)
    local header_b64=$(echo -n "$header" | base64 -w 0)
    local payload_b64=$(echo -n "$payload" | base64 -w 0)
    
    # For testing purposes, return a mock token format
    echo "${header_b64}.${payload_b64}.mock_signature"
}

# Validation functions
validate_response_field() {
    local field_name="$1"
    local expected_value="$2"
    local actual_value="$3"
    
    if [[ "$actual_value" == "$expected_value" ]]; then
        log_debug "✓ Field '$field_name' matches expected value: $expected_value"
        return 0
    else
        log_error "✗ Field '$field_name' mismatch. Expected: '$expected_value', Got: '$actual_value'"
        return 1
    fi
}

validate_response_contains() {
    local field_name="$1"
    local expected_substring="$2"
    local actual_value="$3"
    
    if [[ "$actual_value" == *"$expected_substring"* ]]; then
        log_debug "✓ Field '$field_name' contains expected substring: $expected_substring"
        return 0
    else
        log_error "✗ Field '$field_name' does not contain expected substring. Expected: '$expected_substring', Got: '$actual_value'"
        return 1
    fi
}

validate_response_not_empty() {
    local field_name="$1"
    local actual_value="$2"
    
    if [[ -n "$actual_value" && "$actual_value" != "null" ]]; then
        log_debug "✓ Field '$field_name' is not empty: $actual_value"
        return 0
    else
        log_error "✗ Field '$field_name' is empty or null"
        return 1
    fi
}

# Test data management
create_test_organization() {
    local org_type="$1"
    local org_name="$2"
    local auth_token="$3"
    
    local data="{\"name\":\"$org_name\",\"type\":\"$org_type\"}"
    local headers="Content-Type: application/json"
    
    if [[ -n "$auth_token" ]]; then
        headers="$headers
Authorization: Bearer $auth_token"
    fi
    
    make_request "POST" "/api/v1/organizations" "$data" "$headers" "201"
    
    if [[ $? -eq 0 ]]; then
        local org_id=$(parse_json "$RESPONSE_BODY" "organization_id")
        echo "$org_id"
        return 0
    else
        return 1
    fi
}

create_test_profile() {
    local user_id="$1"
    local email="$2"
    local first_name="$3"
    local last_name="$4"
    local role_id="$5"
    local auth_token="$6"
    
    local data="{\"id\":\"$user_id\",\"email\":\"$email\",\"first_name\":\"$first_name\",\"last_name\":\"$last_name\",\"role_id\":$role_id}"
    local headers="Content-Type: application/json
Authorization: Bearer $auth_token"
    
    make_request "POST" "/api/v1/profiles" "$data" "$headers" "201"
    
    return $?
}

# Cleanup functions
cleanup_test_data() {
    if [[ "$SKIP_CLEANUP" == "true" ]]; then
        log_info "Skipping cleanup as requested"
        return 0
    fi
    
    log_info "Cleaning up test data..."
    
    # Add cleanup logic here
    # This would typically involve deleting test organizations, profiles, etc.
    # Implementation depends on the specific test data created
    
    log_info "Cleanup completed"
}

# Report generation
generate_test_report() {
    local report_file="$1"
    
    log_info "Generating test report: $report_file"
    
    cat > "$report_file" << EOF
# Test Execution Report

**Generated:** $(date)
**Total Tests:** $TOTAL_TESTS
**Passed:** $PASSED_TESTS
**Failed:** $FAILED_TESTS
**Skipped:** $SKIPPED_TESTS
**Success Rate:** $(( PASSED_TESTS * 100 / TOTAL_TESTS ))%

## Test Results

EOF
    
    for result in "${TEST_RESULTS[@]}"; do
        IFS='|' read -r status name duration message <<< "$result"
        echo "- **$name** ($duration s): $status - $message" >> "$report_file"
    done
    
    log_success "Test report generated: $report_file"
}

# Summary function
print_test_summary() {
    echo -e "\n${CYAN}=== Test Execution Summary ===${NC}"
    echo -e "Total Tests: $TOTAL_TESTS"
    echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
    echo -e "${RED}Failed: $FAILED_TESTS${NC}"
    echo -e "${YELLOW}Skipped: $SKIPPED_TESTS${NC}"
    
    if [[ $TOTAL_TESTS -gt 0 ]]; then
        local success_rate=$(( PASSED_TESTS * 100 / TOTAL_TESTS ))
        echo -e "Success Rate: ${success_rate}%"
        
        if [[ $FAILED_TESTS -eq 0 ]]; then
            echo -e "\n${GREEN}🎉 All tests passed!${NC}"
        else
            echo -e "\n${RED}❌ Some tests failed. Check the logs above.${NC}"
        fi
    fi
}

# Trap for cleanup on exit
trap cleanup_test_data EXIT