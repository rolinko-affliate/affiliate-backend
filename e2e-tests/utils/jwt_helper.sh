#!/bin/bash

# JWT Helper functions for E2E tests
# This file contains functions for generating and managing JWT tokens for testing

# JWT Secret (should match the server configuration)
JWT_SECRET="${JWT_SECRET:-gDxsm/JerlPJiOObQLtfjViLBQF2ggmJpYCNW+9LPwL2QJksmiYlzRCJCKseCLxJtGysx+awZvoiS0MF0pLjnw==}"

# Test user configurations
declare -A TEST_USERS=(
    ["admin"]="43bad314-bdd3-49d3-9f85-5be4d019c2ae|admin@rolinko.com|Admin|User|1"
    ["advertiser_manager"]="22222222-2222-2222-2222-222222222222|advertiser@test.com|John|Advertiser|1000"
    ["affiliate_manager"]="33333333-3333-3333-3333-333333333333|affiliate@test.com|Jane|Affiliate|1001"
    ["regular_user"]="44444444-4444-4444-4444-444444444444|user@test.com|Regular|User|100000"
)

# Generate JWT token using Python (more reliable than bash)
generate_jwt_token() {
    local user_id="$1"
    local email="$2"
    local role="$3"
    local exp_hours="${4:-1}"  # Default 1 hour expiration
    
    python3 -c "
import jwt
import json
from datetime import datetime, timedelta

# Token payload
payload = {
    'sub': '$user_id',
    'email': '$email',
    'role': '$role',
    'exp': int((datetime.now() + timedelta(hours=$exp_hours)).timestamp()),
    'iat': int(datetime.now().timestamp())
}

# Generate token
secret = '$JWT_SECRET'
token = jwt.encode(payload, secret.encode(), algorithm='HS256')
print(token)
"
}

# Get predefined test user token
get_test_user_token() {
    local user_type="$1"
    
    if [[ -z "${TEST_USERS[$user_type]}" ]]; then
        log_error "Unknown test user type: $user_type"
        return 1
    fi
    
    IFS='|' read -r user_id email first_name last_name role_id <<< "${TEST_USERS[$user_type]}"
    
    # Map role_id to role name for JWT
    local role_name
    case "$role_id" in
        "1") role_name="Admin" ;;
        "1000") role_name="AdvertiserManager" ;;
        "1001") role_name="AffiliateManager" ;;
        "100000") role_name="User" ;;
        *) role_name="User" ;;
    esac
    
    generate_jwt_token "$user_id" "$email" "$role_name"
}

# Validate JWT token
validate_jwt_token() {
    local token="$1"
    
    python3 -c "
import jwt
import json

try:
    secret = '$JWT_SECRET'
    decoded = jwt.decode('$token', secret.encode(), algorithms=['HS256'])
    print(json.dumps(decoded, indent=2))
    exit(0)
except jwt.ExpiredSignatureError:
    print('Token has expired')
    exit(1)
except jwt.InvalidTokenError as e:
    print(f'Invalid token: {e}')
    exit(1)
"
}

# Get user info from test user type
get_test_user_info() {
    local user_type="$1"
    local field="$2"  # user_id, email, first_name, last_name, role_id
    
    if [[ -z "${TEST_USERS[$user_type]}" ]]; then
        log_error "Unknown test user type: $user_type"
        return 1
    fi
    
    IFS='|' read -r user_id email first_name last_name role_id <<< "${TEST_USERS[$user_type]}"
    
    case "$field" in
        "user_id") echo "$user_id" ;;
        "email") echo "$email" ;;
        "first_name") echo "$first_name" ;;
        "last_name") echo "$last_name" ;;
        "role_id") echo "$role_id" ;;
        "full_name") echo "$first_name $last_name" ;;
        *) 
            log_error "Unknown field: $field"
            return 1
            ;;
    esac
}

# Create authorization header
create_auth_header() {
    local user_type="$1"
    local token
    
    token=$(get_test_user_token "$user_type")
    if [[ $? -eq 0 ]]; then
        echo "Authorization: Bearer $token"
    else
        return 1
    fi
}

# Test JWT functionality
test_jwt_functions() {
    log_info "Testing JWT helper functions..."
    
    # Test token generation for each user type
    for user_type in "${!TEST_USERS[@]}"; do
        log_debug "Testing token generation for $user_type"
        
        local token
        token=$(get_test_user_token "$user_type")
        
        if [[ $? -eq 0 && -n "$token" ]]; then
            log_debug "✓ Token generated for $user_type"
            
            # Validate the token
            if validate_jwt_token "$token" >/dev/null 2>&1; then
                log_debug "✓ Token validation passed for $user_type"
            else
                log_error "✗ Token validation failed for $user_type"
                return 1
            fi
        else
            log_error "✗ Failed to generate token for $user_type"
            return 1
        fi
    done
    
    log_success "JWT helper functions test completed successfully"
    return 0
}

# Export functions for use in other scripts
export -f generate_jwt_token
export -f get_test_user_token
export -f validate_jwt_token
export -f get_test_user_info
export -f create_auth_header
export -f test_jwt_functions