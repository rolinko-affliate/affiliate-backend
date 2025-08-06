#!/bin/bash

# E2E Test Suite Runner
# 
# This script runs all E2E test scenarios for the Affiliate Platform API.
# It provides comprehensive testing coverage including:
# - Authentication and user management
# - Advertiser complete workflow
# - Affiliate complete workflow  
# - Association invitation system
# - Integration and error handling
#
# Usage:
#   ./run_all_tests.sh [options]
#
# Options:
#   --scenario <name>     Run specific scenario only
#   --skip-cleanup        Skip test data cleanup
#   --verbose            Enable verbose output
#   --timeout <seconds>   Set request timeout (default: 30)
#   --api-url <url>       Set API base URL (default: http://localhost:8080)
#   --parallel           Run scenarios in parallel (experimental)
#   --report-format <fmt> Report format: html, json, markdown (default: markdown)
#   --help               Show this help message

set -e

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Load common utilities
source "$SCRIPT_DIR/utils/common.sh"
source "$SCRIPT_DIR/utils/jwt_helper.sh"

# Configuration
SCENARIO_TO_RUN=""
PARALLEL_EXECUTION=false
REPORT_FORMAT="markdown"
GENERATE_SUMMARY_REPORT=true

# Test scenarios
declare -a SCENARIOS=(
    "01_authentication_flow.sh|Authentication and User Management Flow"
    "02_advertiser_workflow.sh|Advertiser Complete Workflow"
    "03_affiliate_workflow.sh|Affiliate Complete Workflow"
    "04_association_invitation_system.sh|Association Invitation System Complete Workflow"
    "05_integration_and_error_handling.sh|Integration and Error Handling"
)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Global test results
TOTAL_SCENARIOS=0
PASSED_SCENARIOS=0
FAILED_SCENARIOS=0
SKIPPED_SCENARIOS=0

declare -a SCENARIO_RESULTS=()

# Help function
show_help() {
    cat << EOF
E2E Test Suite Runner for Affiliate Platform API

Usage: $0 [options]

Options:
    --scenario <name>     Run specific scenario only (e.g., authentication_flow)
    --skip-cleanup        Skip test data cleanup after tests
    --verbose            Enable verbose output for debugging
    --timeout <seconds>   Set HTTP request timeout (default: 30)
    --api-url <url>       Set API base URL (default: http://localhost:8080)
    --parallel           Run scenarios in parallel (experimental)
    --report-format <fmt> Report format: html, json, markdown (default: markdown)
    --help               Show this help message

Available Scenarios:
EOF

    for scenario in "${SCENARIOS[@]}"; do
        IFS='|' read -r script_name description <<< "$scenario"
        local scenario_name="${script_name%.*}"
        echo "    ${scenario_name}  ${description}"
    done

    cat << EOF

Examples:
    $0                                    # Run all scenarios
    $0 --scenario authentication_flow    # Run only authentication tests
    $0 --verbose --skip-cleanup          # Run with verbose output, no cleanup
    $0 --api-url http://staging.api.com  # Run against staging environment
    $0 --parallel                        # Run scenarios in parallel

Environment Variables:
    API_BASE_URL         Base URL for API (default: http://localhost:8080)
    TEST_TIMEOUT         Request timeout in seconds (default: 30)
    VERBOSE_OUTPUT       Enable verbose logging (default: false)
    SKIP_CLEANUP         Skip test data cleanup (default: false)
    JWT_SECRET           JWT secret for token generation

EOF
}

# Parse command line arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --scenario)
                SCENARIO_TO_RUN="$2"
                shift 2
                ;;
            --skip-cleanup)
                export SKIP_CLEANUP="true"
                shift
                ;;
            --verbose)
                export VERBOSE_OUTPUT="true"
                shift
                ;;
            --timeout)
                export TEST_TIMEOUT="$2"
                shift 2
                ;;
            --api-url)
                export API_BASE_URL="$2"
                shift 2
                ;;
            --parallel)
                PARALLEL_EXECUTION=true
                shift
                ;;
            --report-format)
                REPORT_FORMAT="$2"
                shift 2
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                echo "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# Pre-flight checks
run_preflight_checks() {
    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN}E2E Test Suite Pre-flight Checks${NC}"
    echo -e "${CYAN}========================================${NC}"
    
    # Check if API server is running
    echo -e "${BLUE}[CHECK]${NC} Testing API server connectivity..."
    
    if curl -s --max-time 10 "$API_BASE_URL/health" > /dev/null; then
        echo -e "${GREEN}✓${NC} API server is accessible at $API_BASE_URL"
    else
        echo -e "${RED}✗${NC} API server is not accessible at $API_BASE_URL"
        echo -e "${RED}Please ensure the server is running before running tests${NC}"
        exit 1
    fi
    
    # Check JWT helper functionality
    echo -e "${BLUE}[CHECK]${NC} Testing JWT token generation..."
    
    if test_jwt_functions > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC} JWT token generation working correctly"
    else
        echo -e "${RED}✗${NC} JWT token generation failed"
        echo -e "${RED}Please check JWT_SECRET configuration${NC}"
        exit 1
    fi
    
    # Check required tools
    echo -e "${BLUE}[CHECK]${NC} Checking required tools..."
    
    local required_tools=("curl" "python3" "date" "grep" "sed")
    local missing_tools=()
    
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            missing_tools+=("$tool")
        fi
    done
    
    if [[ ${#missing_tools[@]} -eq 0 ]]; then
        echo -e "${GREEN}✓${NC} All required tools are available"
    else
        echo -e "${RED}✗${NC} Missing required tools: ${missing_tools[*]}"
        exit 1
    fi
    
    # Check Python dependencies
    echo -e "${BLUE}[CHECK]${NC} Checking Python dependencies..."
    
    if python3 -c "import jwt, requests" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} Python dependencies are available"
    else
        echo -e "${YELLOW}⚠${NC} Python dependencies missing (jwt, requests)"
        echo -e "${YELLOW}Installing dependencies...${NC}"
        
        if pip3 install pyjwt requests > /dev/null 2>&1; then
            echo -e "${GREEN}✓${NC} Python dependencies installed successfully"
        else
            echo -e "${RED}✗${NC} Failed to install Python dependencies"
            exit 1
        fi
    fi
    
    echo -e "${GREEN}All pre-flight checks passed!${NC}\n"
}

# Run single scenario
run_scenario() {
    local script_name="$1"
    local description="$2"
    local scenario_name="${script_name%.*}"
    
    echo -e "\n${PURPLE}========================================${NC}"
    echo -e "${PURPLE}Running Scenario: $scenario_name${NC}"
    echo -e "${PURPLE}Description: $description${NC}"
    echo -e "${PURPLE}========================================${NC}"
    
    local scenario_start_time=$(date +%s)
    local scenario_script="$SCRIPT_DIR/scenarios/$script_name"
    
    if [[ ! -f "$scenario_script" ]]; then
        echo -e "${RED}Error: Scenario script not found: $scenario_script${NC}"
        FAILED_SCENARIOS=$((FAILED_SCENARIOS + 1))
        SCENARIO_RESULTS+=("FAIL|$scenario_name|0|Script not found")
        return 1
    fi
    
    # Make script executable
    chmod +x "$scenario_script"
    
    # Run scenario
    if bash "$scenario_script"; then
        local scenario_end_time=$(date +%s)
        local scenario_duration=$((scenario_end_time - scenario_start_time))
        
        echo -e "${GREEN}✓ Scenario '$scenario_name' completed successfully in ${scenario_duration}s${NC}"
        PASSED_SCENARIOS=$((PASSED_SCENARIOS + 1))
        SCENARIO_RESULTS+=("PASS|$scenario_name|$scenario_duration|All tests passed")
    else
        local scenario_end_time=$(date +%s)
        local scenario_duration=$((scenario_end_time - scenario_start_time))
        
        echo -e "${RED}✗ Scenario '$scenario_name' failed in ${scenario_duration}s${NC}"
        FAILED_SCENARIOS=$((FAILED_SCENARIOS + 1))
        SCENARIO_RESULTS+=("FAIL|$scenario_name|$scenario_duration|Some tests failed")
    fi
}

# Run scenarios in parallel
run_scenarios_parallel() {
    echo -e "${YELLOW}Running scenarios in parallel mode...${NC}"
    
    local pids=()
    local temp_dir="/tmp/e2e_test_results_$$"
    mkdir -p "$temp_dir"
    
    for scenario in "${SCENARIOS[@]}"; do
        IFS='|' read -r script_name description <<< "$scenario"
        local scenario_name="${script_name%.*}"
        
        # Skip if specific scenario requested and this isn't it
        if [[ -n "$SCENARIO_TO_RUN" && "$scenario_name" != "$SCENARIO_TO_RUN" ]]; then
            continue
        fi
        
        echo -e "${BLUE}Starting scenario: $scenario_name${NC}"
        
        # Run scenario in background
        (
            run_scenario "$script_name" "$description" > "$temp_dir/$scenario_name.log" 2>&1
            echo $? > "$temp_dir/$scenario_name.exit_code"
        ) &
        
        pids+=($!)
        TOTAL_SCENARIOS=$((TOTAL_SCENARIOS + 1))
    done
    
    # Wait for all scenarios to complete
    echo -e "${BLUE}Waiting for all scenarios to complete...${NC}"
    
    for pid in "${pids[@]}"; do
        wait "$pid"
    done
    
    # Collect results
    for scenario in "${SCENARIOS[@]}"; do
        IFS='|' read -r script_name description <<< "$scenario"
        local scenario_name="${script_name%.*}"
        
        # Skip if specific scenario requested and this isn't it
        if [[ -n "$SCENARIO_TO_RUN" && "$scenario_name" != "$SCENARIO_TO_RUN" ]]; then
            continue
        fi
        
        local exit_code_file="$temp_dir/$scenario_name.exit_code"
        local log_file="$temp_dir/$scenario_name.log"
        
        if [[ -f "$exit_code_file" ]]; then
            local exit_code=$(cat "$exit_code_file")
            
            if [[ "$exit_code" -eq 0 ]]; then
                PASSED_SCENARIOS=$((PASSED_SCENARIOS + 1))
                SCENARIO_RESULTS+=("PASS|$scenario_name|0|Completed successfully")
                echo -e "${GREEN}✓ $scenario_name completed successfully${NC}"
            else
                FAILED_SCENARIOS=$((FAILED_SCENARIOS + 1))
                SCENARIO_RESULTS+=("FAIL|$scenario_name|0|Failed with exit code $exit_code")
                echo -e "${RED}✗ $scenario_name failed${NC}"
            fi
            
            # Show log if verbose
            if [[ "$VERBOSE_OUTPUT" == "true" && -f "$log_file" ]]; then
                echo -e "${CYAN}--- Log for $scenario_name ---${NC}"
                cat "$log_file"
                echo -e "${CYAN}--- End log for $scenario_name ---${NC}"
            fi
        else
            FAILED_SCENARIOS=$((FAILED_SCENARIOS + 1))
            SCENARIO_RESULTS+=("FAIL|$scenario_name|0|No exit code found")
            echo -e "${RED}✗ $scenario_name - no exit code found${NC}"
        fi
    done
    
    # Cleanup temp directory
    rm -rf "$temp_dir"
}

# Run scenarios sequentially
run_scenarios_sequential() {
    echo -e "${BLUE}Running scenarios sequentially...${NC}"
    
    for scenario in "${SCENARIOS[@]}"; do
        IFS='|' read -r script_name description <<< "$scenario"
        local scenario_name="${script_name%.*}"
        
        # Skip if specific scenario requested and this isn't it
        if [[ -n "$SCENARIO_TO_RUN" && "$scenario_name" != "$SCENARIO_TO_RUN" ]]; then
            echo -e "${YELLOW}Skipping scenario: $scenario_name${NC}"
            SKIPPED_SCENARIOS=$((SKIPPED_SCENARIOS + 1))
            continue
        fi
        
        TOTAL_SCENARIOS=$((TOTAL_SCENARIOS + 1))
        run_scenario "$script_name" "$description"
    done
}

# Generate comprehensive test report
generate_comprehensive_report() {
    local report_dir="$SCRIPT_DIR/reports"
    mkdir -p "$report_dir"
    
    local timestamp=$(date +"%Y%m%d_%H%M%S")
    local report_file="$report_dir/comprehensive_test_report_$timestamp"
    
    case "$REPORT_FORMAT" in
        "html")
            report_file="${report_file}.html"
            generate_html_report "$report_file"
            ;;
        "json")
            report_file="${report_file}.json"
            generate_json_report "$report_file"
            ;;
        *)
            report_file="${report_file}.md"
            generate_markdown_report "$report_file"
            ;;
    esac
    
    echo -e "${GREEN}Comprehensive test report generated: $report_file${NC}"
}

# Generate markdown report
generate_markdown_report() {
    local report_file="$1"
    
    cat > "$report_file" << EOF
# E2E Test Suite Comprehensive Report

**Generated:** $(date)
**API Base URL:** $API_BASE_URL
**Test Timeout:** ${TEST_TIMEOUT}s
**Parallel Execution:** $PARALLEL_EXECUTION

## Summary

- **Total Scenarios:** $TOTAL_SCENARIOS
- **Passed:** $PASSED_SCENARIOS
- **Failed:** $FAILED_SCENARIOS
- **Skipped:** $SKIPPED_SCENARIOS
- **Success Rate:** $(( TOTAL_SCENARIOS > 0 ? PASSED_SCENARIOS * 100 / TOTAL_SCENARIOS : 0 ))%

## Scenario Results

EOF
    
    for result in "${SCENARIO_RESULTS[@]}"; do
        IFS='|' read -r status name duration message <<< "$result"
        local status_icon
        case "$status" in
            "PASS") status_icon="✅" ;;
            "FAIL") status_icon="❌" ;;
            "SKIP") status_icon="⏭️" ;;
        esac
        
        echo "### $status_icon $name" >> "$report_file"
        echo "- **Status:** $status" >> "$report_file"
        echo "- **Duration:** ${duration}s" >> "$report_file"
        echo "- **Message:** $message" >> "$report_file"
        echo "" >> "$report_file"
    done
    
    cat >> "$report_file" << EOF

## Configuration

- **API Base URL:** $API_BASE_URL
- **Test Timeout:** ${TEST_TIMEOUT}s
- **Verbose Output:** ${VERBOSE_OUTPUT:-false}
- **Skip Cleanup:** ${SKIP_CLEANUP:-false}
- **Parallel Execution:** $PARALLEL_EXECUTION

## Test Coverage

This test suite covers:

1. **Authentication and User Management**
   - User registration via webhooks
   - JWT token validation
   - Role-based access control
   - Profile management

2. **Advertiser Workflow**
   - Organization creation
   - Advertiser profile management
   - Campaign creation and management
   - Tracking link generation
   - Association invitation system
   - External provider integration

3. **Affiliate Workflow**
   - Affiliate organization setup
   - Affiliate profile management
   - Campaign discovery
   - Association management
   - Performance tracking

4. **Association Invitation System**
   - Invitation creation and management
   - Public invitation access
   - Usage tracking and analytics
   - Restriction enforcement
   - Expiration handling

5. **Integration and Error Handling**
   - External service integration
   - Webhook processing
   - Input validation
   - Error response consistency
   - Performance testing

## Next Steps

EOF
    
    if [[ $FAILED_SCENARIOS -gt 0 ]]; then
        cat >> "$report_file" << EOF
⚠️ **Action Required:** Some test scenarios failed. Please review the failed scenarios and address the issues before deploying to production.

EOF
    else
        cat >> "$report_file" << EOF
🎉 **All tests passed!** The system is ready for deployment.

EOF
    fi
    
    cat >> "$report_file" << EOF
For detailed logs of individual scenarios, check the individual report files in the reports directory.

---
*Report generated by E2E Test Suite Runner*
EOF
}

# Generate JSON report
generate_json_report() {
    local report_file="$1"
    
    cat > "$report_file" << EOF
{
  "generated": "$(date -Iseconds)",
  "api_base_url": "$API_BASE_URL",
  "test_timeout": $TEST_TIMEOUT,
  "parallel_execution": $PARALLEL_EXECUTION,
  "summary": {
    "total_scenarios": $TOTAL_SCENARIOS,
    "passed": $PASSED_SCENARIOS,
    "failed": $FAILED_SCENARIOS,
    "skipped": $SKIPPED_SCENARIOS,
    "success_rate": $(( TOTAL_SCENARIOS > 0 ? PASSED_SCENARIOS * 100 / TOTAL_SCENARIOS : 0 ))
  },
  "scenarios": [
EOF
    
    local first=true
    for result in "${SCENARIO_RESULTS[@]}"; do
        IFS='|' read -r status name duration message <<< "$result"
        
        if [[ "$first" == "true" ]]; then
            first=false
        else
            echo "," >> "$report_file"
        fi
        
        cat >> "$report_file" << EOF
    {
      "name": "$name",
      "status": "$status",
      "duration": $duration,
      "message": "$message"
    }
EOF
    done
    
    cat >> "$report_file" << EOF
  ],
  "configuration": {
    "api_base_url": "$API_BASE_URL",
    "test_timeout": $TEST_TIMEOUT,
    "verbose_output": ${VERBOSE_OUTPUT:-false},
    "skip_cleanup": ${SKIP_CLEANUP:-false},
    "parallel_execution": $PARALLEL_EXECUTION
  }
}
EOF
}

# Print final summary
print_final_summary() {
    echo -e "\n${CYAN}========================================${NC}"
    echo -e "${CYAN}E2E Test Suite Final Summary${NC}"
    echo -e "${CYAN}========================================${NC}"
    
    echo -e "Total Scenarios: $TOTAL_SCENARIOS"
    echo -e "${GREEN}Passed: $PASSED_SCENARIOS${NC}"
    echo -e "${RED}Failed: $FAILED_SCENARIOS${NC}"
    echo -e "${YELLOW}Skipped: $SKIPPED_SCENARIOS${NC}"
    
    if [[ $TOTAL_SCENARIOS -gt 0 ]]; then
        local success_rate=$(( PASSED_SCENARIOS * 100 / TOTAL_SCENARIOS ))
        echo -e "Success Rate: ${success_rate}%"
        
        if [[ $FAILED_SCENARIOS -eq 0 ]]; then
            echo -e "\n${GREEN}🎉 All test scenarios passed successfully!${NC}"
            echo -e "${GREEN}The system is ready for deployment.${NC}"
        else
            echo -e "\n${RED}❌ Some test scenarios failed.${NC}"
            echo -e "${RED}Please review the failed scenarios and address the issues.${NC}"
        fi
    fi
    
    echo -e "\n${BLUE}Test execution completed.${NC}"
}

# Main execution function
main() {
    # Parse command line arguments
    parse_arguments "$@"
    
    # Show configuration
    echo -e "${CYAN}E2E Test Suite Configuration:${NC}"
    echo -e "API Base URL: $API_BASE_URL"
    echo -e "Test Timeout: ${TEST_TIMEOUT}s"
    echo -e "Verbose Output: ${VERBOSE_OUTPUT:-false}"
    echo -e "Skip Cleanup: ${SKIP_CLEANUP:-false}"
    echo -e "Parallel Execution: $PARALLEL_EXECUTION"
    echo -e "Report Format: $REPORT_FORMAT"
    
    if [[ -n "$SCENARIO_TO_RUN" ]]; then
        echo -e "Running Specific Scenario: $SCENARIO_TO_RUN"
    fi
    
    echo ""
    
    # Run pre-flight checks
    run_preflight_checks
    
    # Record start time
    local suite_start_time=$(date +%s)
    
    # Run scenarios
    if [[ "$PARALLEL_EXECUTION" == "true" ]]; then
        run_scenarios_parallel
    else
        run_scenarios_sequential
    fi
    
    # Record end time
    local suite_end_time=$(date +%s)
    local suite_duration=$((suite_end_time - suite_start_time))
    
    echo -e "\n${BLUE}Total execution time: ${suite_duration}s${NC}"
    
    # Generate comprehensive report
    if [[ "$GENERATE_SUMMARY_REPORT" == "true" ]]; then
        generate_comprehensive_report
    fi
    
    # Print final summary
    print_final_summary
    
    # Exit with appropriate code
    if [[ $FAILED_SCENARIOS -eq 0 ]]; then
        exit 0
    else
        exit 1
    fi
}

# Run main function if script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi