#!/bin/bash

# Test Account Creation Runner Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Default values
API_URL="http://localhost:8080"
ADMIN_TOKEN=""
VERBOSE=false
DRY_RUN=false

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_highlight() {
    echo -e "${CYAN}[HIGHLIGHT]${NC} $1"
}

# Function to show usage
show_usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Create test accounts for the affiliate platform.

OPTIONS:
    -u, --api-url URL       API base URL (default: http://localhost:8080)
    -t, --token TOKEN       Admin JWT token for authentication
    -v, --verbose           Enable verbose output
    -d, --dry-run           Only check API health, don't create accounts
    -h, --help              Show this help message

EXAMPLES:
    $0 -t "your-token"
    $0 -u "http://localhost:3000" -t "your-token"
    $0 -t "your-token" -v

EOF
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check if Python 3 is available
    if ! command -v python3 >/dev/null 2>&1; then
        print_error "Python 3 is required but not installed"
        return 1
    fi
    
    # Check if requests library is available
    if ! python3 -c "import requests" >/dev/null 2>&1; then
        print_warning "Python requests library not found, attempting to install..."
        if command -v pip3 >/dev/null 2>&1; then
            pip3 install requests
        elif command -v pip >/dev/null 2>&1; then
            pip install requests
        else
            print_error "Cannot install requests library. Please install it manually:"
            print_error "  pip3 install requests"
            return 1
        fi
    fi
    
    # Check if the Python script exists
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PYTHON_SCRIPT="$SCRIPT_DIR/create_test_accounts.py"
    
    if [[ ! -f "$PYTHON_SCRIPT" ]]; then
        print_error "Python script not found: $PYTHON_SCRIPT"
        return 1
    fi
    
    print_success "All prerequisites satisfied"
    return 0
}

# Function to check API connectivity
check_api_connectivity() {
    print_status "Checking API connectivity..."
    
    if command -v curl >/dev/null 2>&1; then
        if curl -s --connect-timeout 5 "$API_URL/health" >/dev/null 2>&1; then
            print_success "API is reachable at $API_URL"
            return 0
        else
            print_error "Cannot reach API at $API_URL"
            print_error "Make sure the API server is running: make run"
            return 1
        fi
    else
        print_warning "curl not available, skipping connectivity check"
        return 0
    fi
}

# Function to run the Python script
run_python_script() {
    print_status "Running idempotent test account creation script..."
    
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PYTHON_SCRIPT="$SCRIPT_DIR/create_test_accounts.py"
    
    # Build Python command arguments
    PYTHON_ARGS=()
    PYTHON_ARGS+=("-u" "$API_URL")
    
    if [[ -n "$ADMIN_TOKEN" ]]; then
        PYTHON_ARGS+=("-t" "$ADMIN_TOKEN")
    fi
    
    if [[ "$VERBOSE" == "true" ]]; then
        PYTHON_ARGS+=("-v")
    fi
    
    if [[ "$DRY_RUN" == "true" ]]; then
        PYTHON_ARGS+=("--dry-run")
    fi
    
    # Run the Python script
    if [[ "$VERBOSE" == "true" ]]; then
        print_status "Executing: python3 $PYTHON_SCRIPT ${PYTHON_ARGS[*]}"
    fi
    
    python3 "$PYTHON_SCRIPT" "${PYTHON_ARGS[@]}"
    return $?
}

# Function to run verification
run_verification() {
    print_status "Running verification..."
    
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    VERIFY_SCRIPT="$SCRIPT_DIR/verify_test_accounts.py"
    
    if [[ ! -f "$VERIFY_SCRIPT" ]]; then
        print_warning "Verification script not found, skipping verification"
        return 0
    fi
    
    # Build verification command arguments
    VERIFY_ARGS=()
    VERIFY_ARGS+=("-u" "$API_URL")
    
    if [[ -n "$ADMIN_TOKEN" ]]; then
        VERIFY_ARGS+=("-t" "$ADMIN_TOKEN")
    fi
    
    if [[ "$VERBOSE" == "true" ]]; then
        VERIFY_ARGS+=("-v")
    fi
    
    # Run the verification script
    python3 "$VERIFY_SCRIPT" "${VERIFY_ARGS[@]}"
    return $?
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -u|--api-url)
            API_URL="$2"
            shift 2
            ;;
        -t|--token)
            ADMIN_TOKEN="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Get values from environment if not provided
if [[ -z "$API_URL" ]]; then
    API_URL="${API_BASE_URL:-http://localhost:8080}"
fi

if [[ -z "$ADMIN_TOKEN" ]]; then
    ADMIN_TOKEN="${ADMIN_JWT_TOKEN:-}"
fi

# Validate required parameters
if [[ -z "$ADMIN_TOKEN" && "$DRY_RUN" != "true" ]]; then
    print_error "Admin JWT token is required for creating test accounts"
    print_error "Provide it via -t option or ADMIN_JWT_TOKEN environment variable"
    print_error "Use -d for dry run (health check only)"
    show_usage
    exit 1
fi

# Main execution
echo ""
print_highlight "🚀 TEST ACCOUNT CREATION"
print_status "API URL: $API_URL"

if [[ -n "$ADMIN_TOKEN" ]]; then
    print_status "Authentication: JWT token provided"
else
    print_warning "Authentication: No token provided (dry run mode)"
fi

if [[ "$DRY_RUN" == "true" ]]; then
    print_highlight "Mode: DRY RUN (health check only)"
else
    print_highlight "Mode: FULL CREATION"
fi

echo ""

# Check prerequisites
if ! check_prerequisites; then
    print_error "Prerequisites check failed"
    exit 1
fi

echo ""

# Check API connectivity
if ! check_api_connectivity; then
    print_error "API connectivity check failed"
    exit 1
fi

echo ""

# Run the Python script
if ! run_python_script; then
    print_error "Test account creation script failed"
    exit 1
fi

# Run verification if not in dry run mode
if [[ "$DRY_RUN" != "true" ]]; then
    echo ""
    print_status "Running post-creation verification..."
    
    if run_verification; then
        print_success "Verification completed successfully"
    else
        print_warning "Verification completed with some issues"
        print_warning "Check the verification output above for details"
    fi
fi

echo ""
print_success "🎉 Test account creation process completed!"