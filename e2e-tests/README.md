# E2E Test Suite for Affiliate Platform

## Overview

This directory contains comprehensive End-to-End (E2E) tests for the Affiliate Platform API. The tests are written as bash scripts that simulate real-world user workflows and validate the complete system functionality.

## Test Structure

```
e2e-tests/
├── scenarios/          # Individual test scenario scripts
├── utils/             # Shared utilities and helper functions
├── data/              # Test data and fixtures
├── reports/           # Test execution reports
└── README.md          # This file
```

## Test Categories

### 1. Authentication & User Management
- User registration and profile creation
- JWT token validation
- Role-based access control

### 2. Organization Management
- Organization creation and management
- Multi-tenant data isolation
- Organization associations

### 3. Advertiser Workflows
- Advertiser creation and management
- Campaign management
- Tracking link generation
- Association invitation system

### 4. Affiliate Workflows
- Affiliate registration and management
- Campaign discovery
- Invitation acceptance
- Performance tracking

### 5. Integration Testing
- External provider synchronization (Everflow)
- Webhook processing
- Data consistency validation

### 6. Error Handling
- Invalid input validation
- Permission denied scenarios
- Resource not found handling
- External service failure simulation

## Running Tests

### Prerequisites
- Server running on localhost:8080
- PostgreSQL database with test data
- Valid JWT tokens for different user roles

### Run All Tests
```bash
./run_all_tests.sh
```

### Run Specific Test Category
```bash
./scenarios/01_authentication_flow.sh
./scenarios/02_advertiser_workflow.sh
./scenarios/03_affiliate_workflow.sh
```

### Run Individual Test
```bash
./scenarios/01_authentication_flow.sh --test user_registration
```

## Test Data

Test scenarios use predefined test data including:
- Test organizations (advertiser, affiliate, agency)
- Test user profiles with different roles
- Sample campaigns and tracking links
- Mock external provider data

## Reporting

Test results are automatically generated in the `reports/` directory:
- `test_summary.html` - Overall test execution summary
- `detailed_results.json` - Detailed test results in JSON format
- `error_logs.txt` - Error logs and debugging information

## Configuration

Tests can be configured via environment variables:
- `API_BASE_URL` - Base URL for API (default: http://localhost:8080)
- `TEST_TIMEOUT` - Request timeout in seconds (default: 30)
- `VERBOSE_OUTPUT` - Enable verbose logging (default: false)
- `SKIP_CLEANUP` - Skip test data cleanup (default: false)

## Best Practices

1. **Isolation**: Each test scenario is independent and can run in isolation
2. **Cleanup**: Tests clean up their data after execution
3. **Idempotency**: Tests can be run multiple times safely
4. **Validation**: Comprehensive response validation and error checking
5. **Documentation**: Each test includes detailed scenario descriptions

## Contributing

When adding new tests:
1. Follow the existing naming convention
2. Include comprehensive scenario documentation
3. Add proper error handling and cleanup
4. Update this README with new test categories
5. Ensure tests are idempotent and isolated