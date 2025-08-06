#!/bin/bash

# Individual Scenario Test Runner
# 
# This script runs a single E2E test scenario with enhanced options
#
# Usage:
#   ./run_scenario.sh <scenario_name> [options]
#
# Examples:
#   ./run_scenario.sh authentication_flow
#   ./run_scenario.sh advertiser_workflow --verbose
#   ./run_scenario.sh association_invitation_system --skip-cleanup

set -e

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Load common utilities
source "$SCRIPT_DIR/utils/common.sh"

# Available scenarios
declare -A AVAILABLE_SCENARIOS=(
    ["authentication_flow"]="01_authentication_flow.sh|Authentication and User Management Flow"
    ["advertiser_workflow"]="02_advertiser_workflow.sh|Advertiser Complete Workflow"
    ["affiliate_workflow"]="03_affiliate_workflow.sh|Affiliate Complete Workflow"
    ["association_invitation_system"]="04_association_invitation_system.sh|Association Invitation System Complete Workflow"
    ["integration_and_error_handling"]="05_integration_and_error_handling.sh|Integration and Error Handling"
)

# Help function
show_help() {
    cat << EOF
Individual Scenario Test Runner

Usage: $0 <scenario_name> [options]

Available Scenarios:
EOF

    for scenario_name in "${!AVAILABLE_SCENARIOS[@]}"; do
        IFS='|' read -r script_name description <<< "${AVAILABLE_SCENARIOS[$scenario_name]}"
        echo "    $scenario_name  $description"
    done

    cat << EOF

Options:
    --verbose            Enable verbose output
    --skip-cleanup       Skip test data cleanup
    --timeout <seconds>  Set request timeout (default: 30)
    --api-url <url>      Set API base URL (default: http://localhost:8080)
    --help              Show this help message

Examples:
    $0 authentication_flow
    $0 advertiser_workflow --verbose
    $0 association_invitation_system --skip-cleanup
    $0 integration_and_error_handling --api-url http://staging.api.com

EOF
}

# Parse arguments
if [[ $# -eq 0 ]]; then
    echo "Error: No scenario specified"
    show_help
    exit 1
fi

SCENARIO_NAME="$1"
shift

# Check if scenario exists
if [[ -z "${AVAILABLE_SCENARIOS[$SCENARIO_NAME]}" ]]; then
    echo "Error: Unknown scenario '$SCENARIO_NAME'"
    echo ""
    echo "Available scenarios:"
    for name in "${!AVAILABLE_SCENARIOS[@]}"; do
        echo "  $name"
    done
    exit 1
fi

# Parse remaining arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --verbose)
            export VERBOSE_OUTPUT="true"
            shift
            ;;
        --skip-cleanup)
            export SKIP_CLEANUP="true"
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

# Get scenario details
IFS='|' read -r script_name description <<< "${AVAILABLE_SCENARIOS[$SCENARIO_NAME]}"
scenario_script="$SCRIPT_DIR/scenarios/$script_name"

# Check if script exists
if [[ ! -f "$scenario_script" ]]; then
    echo "Error: Scenario script not found: $scenario_script"
    exit 1
fi

# Show configuration
echo -e "${CYAN}Running Individual Scenario Test${NC}"
echo -e "${CYAN}================================${NC}"
echo -e "Scenario: $SCENARIO_NAME"
echo -e "Description: $description"
echo -e "Script: $script_name"
echo -e "API Base URL: ${API_BASE_URL:-http://localhost:8080}"
echo -e "Test Timeout: ${TEST_TIMEOUT:-30}s"
echo -e "Verbose Output: ${VERBOSE_OUTPUT:-false}"
echo -e "Skip Cleanup: ${SKIP_CLEANUP:-false}"
echo ""

# Make script executable
chmod +x "$scenario_script"

# Run the scenario
echo -e "${BLUE}Starting scenario execution...${NC}"
start_time=$(date +%s)

if bash "$scenario_script"; then
    end_time=$(date +%s)
    duration=$((end_time - start_time))
    echo -e "\n${GREEN}✓ Scenario '$SCENARIO_NAME' completed successfully in ${duration}s${NC}"
    exit 0
else
    end_time=$(date +%s)
    duration=$((end_time - start_time))
    echo -e "\n${RED}✗ Scenario '$SCENARIO_NAME' failed in ${duration}s${NC}"
    exit 1
fi