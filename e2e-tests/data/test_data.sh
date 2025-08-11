#!/bin/bash

# Test data definitions for E2E tests
# This file contains predefined test data used across different test scenarios

# Test Organizations
declare -A TEST_ORGANIZATIONS=(
    ["advertiser_org"]="Test Advertiser Organization|advertiser"
    ["affiliate_org"]="Test Affiliate Organization|affiliate"
    ["agency_org"]="Test Agency Organization|agency"
    ["advertiser_org_2"]="Second Advertiser Org|advertiser"
    ["affiliate_org_2"]="Second Affiliate Org|affiliate"
)

# Test Advertisers
declare -A TEST_ADVERTISERS=(
    ["advertiser_1"]="Test Advertiser 1|test-advertiser-1|A comprehensive test advertiser for E2E testing|Technology|https://test-advertiser-1.com|advertiser1@test.com"
    ["advertiser_2"]="Test Advertiser 2|test-advertiser-2|Another test advertiser for multi-tenant testing|E-commerce|https://test-advertiser-2.com|advertiser2@test.com"
)

# Test Affiliates
declare -A TEST_AFFILIATES=(
    ["affiliate_1"]="Test Affiliate 1|test-affiliate-1|A test affiliate partner|Blog|https://test-affiliate-1.com"
    ["affiliate_2"]="Test Affiliate 2|test-affiliate-2|Another test affiliate partner|Social Media|https://test-affiliate-2.com"
)

# Test Campaigns
declare -A TEST_CAMPAIGNS=(
    ["campaign_1"]="Summer Sale Campaign|summer-sale-2025|Promote our summer sale with great discounts|active|cpa|50.00"
    ["campaign_2"]="Product Launch Campaign|product-launch-2025|Launch campaign for new product line|active|cpc|2.50"
    ["campaign_3"]="Holiday Special|holiday-special-2025|Special holiday promotion campaign|paused|cpa|75.00"
)

# Test Tracking Links
declare -A TEST_TRACKING_LINKS=(
    ["link_1"]="Summer Sale Link|https://example.com/summer-sale|Test tracking link for summer sale"
    ["link_2"]="Product Launch Link|https://example.com/product-launch|Test tracking link for product launch"
)

# Test Invitations
declare -A TEST_INVITATIONS=(
    ["invitation_1"]="Partner Program Invitation|Join our exclusive partner program|10|30"
    ["invitation_2"]="Limited Time Offer|Special invitation for limited partners|5|7"
    ["invitation_3"]="Open Invitation|Open invitation for all affiliates|100|90"
)

# Test Messages
declare -A TEST_MESSAGES=(
    ["message_1"]="Welcome Message|Welcome to our affiliate program! We're excited to work with you."
    ["message_2"]="Campaign Update|New campaign available with higher payouts. Check it out!"
    ["message_3"]="Performance Review|Great job this month! Your performance has been outstanding."
)

# Helper functions to get test data
get_test_org_data() {
    local org_key="$1"
    local field="$2"  # name, type
    
    if [[ -z "${TEST_ORGANIZATIONS[$org_key]}" ]]; then
        log_error "Unknown test organization: $org_key"
        return 1
    fi
    
    IFS='|' read -r name type <<< "${TEST_ORGANIZATIONS[$org_key]}"
    
    case "$field" in
        "name") echo "$name" ;;
        "type") echo "$type" ;;
        "json") echo "{\"name\":\"$name\",\"type\":\"$type\"}" ;;
        *) 
            log_error "Unknown field: $field"
            return 1
            ;;
    esac
}

get_test_advertiser_data() {
    local adv_key="$1"
    local field="$2"  # name, slug, description, industry, website, contact_email
    
    if [[ -z "${TEST_ADVERTISERS[$adv_key]}" ]]; then
        log_error "Unknown test advertiser: $adv_key"
        return 1
    fi
    
    IFS='|' read -r name slug description industry website contact_email <<< "${TEST_ADVERTISERS[$adv_key]}"
    
    case "$field" in
        "name") echo "$name" ;;
        "slug") echo "$slug" ;;
        "description") echo "$description" ;;
        "industry") echo "$industry" ;;
        "website") echo "$website" ;;
        "contact_email") echo "$contact_email" ;;
        "json") echo "{\"name\":\"$name\",\"slug\":\"$slug\",\"description\":\"$description\",\"industry\":\"$industry\",\"website\":\"$website\",\"contact_email\":\"$contact_email\"}" ;;
        *) 
            log_error "Unknown field: $field"
            return 1
            ;;
    esac
}

get_test_affiliate_data() {
    local aff_key="$1"
    local field="$2"  # name, slug, description, traffic_source, website
    
    if [[ -z "${TEST_AFFILIATES[$aff_key]}" ]]; then
        log_error "Unknown test affiliate: $aff_key"
        return 1
    fi
    
    IFS='|' read -r name slug description traffic_source website <<< "${TEST_AFFILIATES[$aff_key]}"
    
    case "$field" in
        "name") echo "$name" ;;
        "slug") echo "$slug" ;;
        "description") echo "$description" ;;
        "traffic_source") echo "$traffic_source" ;;
        "website") echo "$website" ;;
        "json") echo "{\"name\":\"$name\",\"slug\":\"$slug\",\"description\":\"$description\",\"traffic_source\":\"$traffic_source\",\"website\":\"$website\",\"status\":\"active\"}" ;;
        *) 
            log_error "Unknown field: $field"
            return 1
            ;;
    esac
}

get_test_campaign_data() {
    local camp_key="$1"
    local field="$2"  # name, slug, description, status, payout_type, payout_amount
    
    if [[ -z "${TEST_CAMPAIGNS[$camp_key]}" ]]; then
        log_error "Unknown test campaign: $camp_key"
        return 1
    fi
    
    IFS='|' read -r name slug description status payout_type payout_amount <<< "${TEST_CAMPAIGNS[$camp_key]}"
    
    case "$field" in
        "name") echo "$name" ;;
        "slug") echo "$slug" ;;
        "description") echo "$description" ;;
        "status") echo "$status" ;;
        "payout_type") echo "$payout_type" ;;
        "payout_amount") echo "$payout_amount" ;;
        "json") echo "{\"name\":\"$name\",\"description\":\"$description\",\"status\":\"$status\"}" ;;
        *) 
            log_error "Unknown field: $field"
            return 1
            ;;
    esac
}

get_test_invitation_data() {
    local inv_key="$1"
    local field="$2"  # name, description, max_uses, expires_days
    
    if [[ -z "${TEST_INVITATIONS[$inv_key]}" ]]; then
        log_error "Unknown test invitation: $inv_key"
        return 1
    fi
    
    IFS='|' read -r name description max_uses expires_days <<< "${TEST_INVITATIONS[$inv_key]}"
    
    case "$field" in
        "name") echo "$name" ;;
        "description") echo "$description" ;;
        "max_uses") echo "$max_uses" ;;
        "expires_days") echo "$expires_days" ;;
        "expires_at") 
            # Calculate expiration date
            local expires_at
            expires_at=$(date -d "+${expires_days} days" -Iseconds)
            echo "$expires_at"
            ;;
        "json") 
            local expires_at
            expires_at=$(date -d "+${expires_days} days" -Iseconds)
            echo "{\"name\":\"$name\",\"description\":\"$description\",\"max_uses\":$max_uses,\"expires_at\":\"$expires_at\"}"
            ;;
        *) 
            log_error "Unknown field: $field"
            return 1
            ;;
    esac
}

# Generate random test data
generate_random_org_name() {
    local org_type="$1"
    local timestamp=$(date +%s)
    echo "Test ${org_type^} Org $timestamp"
}

generate_random_email() {
    local prefix="$1"
    local timestamp=$(date +%s)
    echo "${prefix}-${timestamp}@test.example.com"
}

generate_random_slug() {
    local prefix="$1"
    local timestamp=$(date +%s)
    echo "${prefix}-${timestamp}"
}

# Validation data
VALID_ORG_TYPES=("advertiser" "affiliate" "agency")
VALID_CAMPAIGN_STATUSES=("active" "paused" "completed" "draft")
VALID_PAYOUT_TYPES=("cpa" "cpc" "cpm" "revenue_share")
VALID_TRAFFIC_SOURCES=("Blog" "Social Media" "Email" "PPC" "SEO" "Display")
VALID_INDUSTRIES=("Technology" "E-commerce" "Finance" "Healthcare" "Education" "Entertainment")

# Validation functions
is_valid_org_type() {
    local org_type="$1"
    for valid_type in "${VALID_ORG_TYPES[@]}"; do
        if [[ "$org_type" == "$valid_type" ]]; then
            return 0
        fi
    done
    return 1
}

is_valid_campaign_status() {
    local status="$1"
    for valid_status in "${VALID_CAMPAIGN_STATUSES[@]}"; do
        if [[ "$status" == "$valid_status" ]]; then
            return 0
        fi
    done
    return 1
}

# Export functions for use in other scripts
export -f get_test_org_data
export -f get_test_advertiser_data
export -f get_test_affiliate_data
export -f get_test_campaign_data
export -f get_test_invitation_data
export -f generate_random_org_name
export -f generate_random_email
export -f generate_random_slug
export -f is_valid_org_type
export -f is_valid_campaign_status