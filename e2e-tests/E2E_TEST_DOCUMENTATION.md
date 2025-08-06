# E2E Test Suite Documentation

## Overview

This comprehensive End-to-End (E2E) test suite validates the complete functionality of the Affiliate Platform API. The tests are designed to simulate real-world user workflows and ensure system reliability, security, and performance.

## Architecture

### Test Structure
```
e2e-tests/
├── scenarios/              # Individual test scenario scripts
│   ├── 01_authentication_flow.sh
│   ├── 02_advertiser_workflow.sh
│   ├── 03_affiliate_workflow.sh
│   ├── 04_association_invitation_system.sh
│   └── 05_integration_and_error_handling.sh
├── utils/                  # Shared utilities and helper functions
│   ├── common.sh          # Common test utilities
│   ├── jwt_helper.sh      # JWT token management
│   └── test_data.sh       # Test data definitions
├── data/                   # Test data and fixtures
├── reports/               # Generated test reports
├── run_all_tests.sh       # Main test suite runner
├── run_scenario.sh        # Individual scenario runner
└── README.md              # Basic usage instructions
```

### Design Principles

1. **Isolation**: Each test scenario is independent and can run in isolation
2. **Idempotency**: Tests can be run multiple times safely without side effects
3. **Cleanup**: Automatic cleanup of test data after execution
4. **Comprehensive Coverage**: Tests cover all major API endpoints and workflows
5. **Error Handling**: Robust error handling and validation
6. **Documentation**: Detailed scenario descriptions and explanations

## Test Scenarios

### 1. Authentication and User Management Flow (`01_authentication_flow.sh`)

**Purpose**: Validates the complete authentication and user management system.

**Test Coverage**:
- Health check endpoint validation
- Supabase webhook processing for user creation
- JWT token generation and validation for different user roles
- User profile retrieval with authentication
- Profile creation, update, and management
- Unauthorized access rejection
- Invalid token handling
- Role-based access control (RBAC) enforcement
- Profile upsert functionality

**Key Tests**:
- `test_health_check`: Verifies API server health
- `test_supabase_webhook`: Tests user creation via Supabase webhook
- `test_jwt_token_generation`: Validates JWT token generation for all user types
- `test_get_user_profile`: Tests authenticated profile retrieval
- `test_profile_creation`: Validates profile creation workflow
- `test_profile_update`: Tests profile modification
- `test_unauthorized_access`: Ensures proper rejection of unauthenticated requests
- `test_invalid_token_access`: Validates invalid token rejection
- `test_rbac_enforcement`: Tests role-based access control
- `test_profile_upsert`: Validates upsert (create or update) functionality

**Expected Results**: All authentication mechanisms work correctly, RBAC is properly enforced, and user management operations function as expected.

### 2. Advertiser Complete Workflow (`02_advertiser_workflow.sh`)

**Purpose**: Validates the complete advertiser journey from organization creation to campaign management and affiliate relationship building.

**Test Coverage**:
- Advertiser organization creation and management
- Advertiser profile creation, retrieval, and updates
- Campaign creation and management
- Tracking link generation and QR code creation
- Association invitation system
- External provider synchronization (Everflow mock)
- Analytics data retrieval
- Error handling and validation
- Permission enforcement

**Key Tests**:
- `test_create_advertiser_organization`: Creates advertiser organization
- `test_list_organizations`: Tests organization listing
- `test_get_organization_details`: Validates organization retrieval
- `test_create_advertiser`: Creates advertiser profile
- `test_get_advertiser`: Tests advertiser retrieval
- `test_update_advertiser`: Validates advertiser updates
- `test_create_campaign`: Creates marketing campaign
- `test_list_campaigns_by_advertiser`: Lists advertiser campaigns
- `test_create_tracking_link`: Creates tracking links
- `test_generate_tracking_link`: Tests dynamic link generation
- `test_get_tracking_link_qr`: Generates QR codes for links
- `test_create_association_invitation`: Creates affiliate invitations
- `test_generate_invitation_link`: Generates shareable invitation links
- `test_list_association_invitations`: Lists all invitations
- `test_everflow_sync`: Tests external provider synchronization
- `test_advertiser_analytics`: Retrieves analytics data
- `test_error_handling_invalid_data`: Tests error handling
- `test_permission_denied`: Validates permission enforcement

**Expected Results**: Complete advertiser workflow functions correctly, all CRUD operations work, external integrations are functional, and proper error handling is in place.

### 3. Affiliate Complete Workflow (`03_affiliate_workflow.sh`)

**Purpose**: Validates the complete affiliate journey from organization creation to campaign discovery and performance tracking.

**Test Coverage**:
- Affiliate organization creation and management
- Affiliate profile creation, retrieval, and updates
- Campaign search and discovery
- Association invitation acceptance
- Tracking link generation for affiliates
- Performance analytics and reporting
- Publisher messaging system
- Error handling and validation
- Permission enforcement

**Key Tests**:
- `test_create_affiliate_organization`: Creates affiliate organization
- `test_create_affiliate`: Creates affiliate profile
- `test_get_affiliate`: Tests affiliate retrieval
- `test_update_affiliate`: Validates affiliate updates
- `test_list_affiliates_by_organization`: Lists organization affiliates
- `test_search_affiliates`: Tests affiliate search functionality
- `test_get_visible_campaigns`: Retrieves visible campaigns
- `test_create_association_request`: Creates association requests
- `test_list_associations`: Lists organization associations
- `test_use_association_invitation`: Uses advertiser invitations
- `test_generate_affiliate_tracking_links`: Generates affiliate tracking links
- `test_affiliate_analytics`: Retrieves affiliate analytics
- `test_create_messaging_conversation`: Creates messaging conversations
- `test_add_message_to_conversation`: Adds messages to conversations
- `test_get_messaging_conversations`: Retrieves conversations
- `test_error_handling_invalid_affiliate_data`: Tests error handling
- `test_advertiser_manager_permission`: Validates permission enforcement

**Expected Results**: Complete affiliate workflow functions correctly, campaign discovery works, association management is functional, and messaging system operates properly.

### 4. Association Invitation System Complete Workflow (`04_association_invitation_system.sh`)

**Purpose**: Validates the complete association invitation system that enables advertisers to invite affiliates through shareable links.

**Test Coverage**:
- Invitation creation with various configurations
- Invitation management (list, update, delete)
- Public invitation access without authentication
- Invitation usage by affiliates
- Usage tracking and analytics
- Link generation and sharing
- Expiration and usage limit handling
- Restriction enforcement (allowed affiliate organizations)
- Error scenarios and edge cases
- Security and access control

**Key Tests**:
- `test_setup_test_organizations`: Creates test organizations
- `test_create_basic_invitation`: Creates basic invitation
- `test_create_restricted_invitation`: Creates restricted invitation
- `test_list_invitations`: Lists all invitations
- `test_get_invitation_details`: Retrieves invitation details
- `test_generate_invitation_link`: Generates shareable links
- `test_public_invitation_access`: Tests public access without auth
- `test_use_invitation_successfully`: Uses invitation to create association
- `test_use_invitation_duplicate`: Tests duplicate usage handling
- `test_use_restricted_invitation_allowed`: Uses restricted invitation with allowed affiliate
- `test_use_restricted_invitation_disallowed`: Tests restriction enforcement
- `test_get_invitation_usage_history`: Retrieves usage history
- `test_update_invitation`: Updates invitation details
- `test_create_expired_invitation`: Creates expired invitation
- `test_use_expired_invitation`: Tests expired invitation rejection
- `test_invalid_invitation_token`: Tests invalid token handling
- `test_missing_required_fields`: Tests validation
- `test_affiliate_manager_cannot_create_invitations`: Tests permission enforcement

**Expected Results**: Complete invitation system works correctly, restrictions are enforced, expiration is handled properly, and security measures are in place.

### 5. Integration and Error Handling (`05_integration_and_error_handling.sh`)

**Purpose**: Validates system integration points and comprehensive error handling across the platform.

**Test Coverage**:
- External provider integration (Everflow mock)
- Webhook processing and validation
- Billing system integration
- Analytics and reporting integration
- Comprehensive error scenarios
- Edge cases and boundary conditions
- Performance and timeout handling
- Data consistency validation
- Input validation edge cases
- Authentication edge cases

**Key Tests**:
- `test_supabase_webhook_integration`: Tests Supabase webhook processing
- `test_stripe_webhook_integration`: Tests Stripe webhook processing
- `test_everflow_integration`: Tests Everflow provider integration
- `test_analytics_integration`: Tests analytics endpoints
- `test_billing_integration`: Tests billing system
- `test_input_validation_edge_cases`: Tests input validation with edge cases
- `test_authentication_edge_cases`: Tests authentication edge cases
- `test_rate_limiting_and_performance`: Tests system performance
- `test_data_consistency`: Validates data consistency
- `test_error_response_consistency`: Tests error response consistency

**Expected Results**: All integrations work correctly, error handling is comprehensive, system performance is acceptable, and data consistency is maintained.

## Utilities and Helpers

### Common Utilities (`utils/common.sh`)

**Functions**:
- `log_info`, `log_success`, `log_warning`, `log_error`, `log_debug`: Logging functions
- `start_test`, `end_test`: Test execution tracking
- `make_request`: HTTP request wrapper with validation
- `parse_json`: JSON parsing utility
- `validate_response_field`, `validate_response_contains`, `validate_response_not_empty`: Response validation
- `create_test_organization`, `create_test_profile`: Test data creation
- `cleanup_test_data`: Test data cleanup
- `generate_test_report`: Report generation
- `print_test_summary`: Summary printing

### JWT Helper (`utils/jwt_helper.sh`)

**Functions**:
- `generate_jwt_token`: Generates JWT tokens for testing
- `get_test_user_token`: Gets predefined test user tokens
- `validate_jwt_token`: Validates JWT token structure
- `get_test_user_info`: Retrieves test user information
- `create_auth_header`: Creates authorization headers
- `test_jwt_functions`: Tests JWT functionality

**Test Users**:
- `admin`: Full system administrator
- `advertiser_manager`: Advertiser management role
- `affiliate_manager`: Affiliate management role
- `regular_user`: Basic user role

### Test Data (`data/test_data.sh`)

**Data Sets**:
- Test organizations (advertiser, affiliate, agency)
- Test advertisers with various industries
- Test affiliates with different traffic sources
- Test campaigns with different payout models
- Test invitations with various configurations
- Test messages for communication testing

**Functions**:
- `get_test_org_data`, `get_test_advertiser_data`, `get_test_affiliate_data`: Data retrieval
- `get_test_campaign_data`, `get_test_invitation_data`: Campaign and invitation data
- `generate_random_org_name`, `generate_random_email`, `generate_random_slug`: Random data generation
- `is_valid_org_type`, `is_valid_campaign_status`: Validation functions

## Running Tests

### Prerequisites

1. **Server Running**: API server must be running on the configured URL (default: http://localhost:8080)
2. **Database**: PostgreSQL database with proper schema and test data
3. **Dependencies**: Required tools (curl, python3, jq) and Python packages (pyjwt, requests)
4. **Configuration**: Proper environment variables set

### Environment Variables

```bash
# API Configuration
export API_BASE_URL="http://localhost:8080"
export TEST_TIMEOUT="30"

# Test Configuration
export VERBOSE_OUTPUT="false"
export SKIP_CLEANUP="false"

# JWT Configuration
export JWT_SECRET="your-jwt-secret-here"
```

### Running All Tests

```bash
# Run all test scenarios
./run_all_tests.sh

# Run with verbose output
./run_all_tests.sh --verbose

# Run without cleanup
./run_all_tests.sh --skip-cleanup

# Run against different environment
./run_all_tests.sh --api-url http://staging.api.com

# Run in parallel (experimental)
./run_all_tests.sh --parallel

# Generate HTML report
./run_all_tests.sh --report-format html
```

### Running Individual Scenarios

```bash
# Run specific scenario
./run_scenario.sh authentication_flow

# Run with options
./run_scenario.sh advertiser_workflow --verbose --skip-cleanup

# Run against different environment
./run_scenario.sh association_invitation_system --api-url http://staging.api.com
```

### Running Individual Tests

```bash
# Run specific test within a scenario
cd scenarios
./01_authentication_flow.sh --test user_registration
```

## Test Reports

### Report Types

1. **Markdown Reports**: Human-readable reports with detailed results
2. **JSON Reports**: Machine-readable reports for CI/CD integration
3. **HTML Reports**: Web-friendly reports with styling
4. **Individual Scenario Reports**: Detailed reports for each scenario

### Report Contents

- Test execution summary
- Individual test results with timing
- Error logs and debugging information
- Configuration details
- Coverage analysis
- Recommendations and next steps

### Report Locations

```
e2e-tests/reports/
├── comprehensive_test_report_YYYYMMDD_HHMMSS.md
├── authentication_flow_report.md
├── advertiser_workflow_report.md
├── affiliate_workflow_report.md
├── association_invitation_system_report.md
└── integration_and_error_handling_report.md
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: E2E Tests
on: [push, pull_request]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v2
      
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.23
          
      - name: Start API Server
        run: |
          make build
          ./bin/affiliate-backend &
          sleep 10
          
      - name: Install Test Dependencies
        run: |
          pip3 install pyjwt requests
          
      - name: Run E2E Tests
        run: |
          cd e2e-tests
          ./run_all_tests.sh --report-format json
          
      - name: Upload Test Results
        uses: actions/upload-artifact@v2
        if: always()
        with:
          name: e2e-test-results
          path: e2e-tests/reports/
```

### Jenkins Pipeline Example

```groovy
pipeline {
    agent any
    
    stages {
        stage('Setup') {
            steps {
                sh 'make build'
                sh './bin/affiliate-backend &'
                sh 'sleep 10'
            }
        }
        
        stage('E2E Tests') {
            steps {
                dir('e2e-tests') {
                    sh './run_all_tests.sh --report-format json'
                }
            }
            post {
                always {
                    archiveArtifacts artifacts: 'e2e-tests/reports/**/*', fingerprint: true
                    publishHTML([
                        allowMissing: false,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
                        reportDir: 'e2e-tests/reports',
                        reportFiles: '*.html',
                        reportName: 'E2E Test Report'
                    ])
                }
            }
        }
    }
}
```

## Troubleshooting

### Common Issues

1. **Server Not Running**
   - Error: "API server is not accessible"
   - Solution: Ensure the API server is running on the configured URL

2. **JWT Token Issues**
   - Error: "JWT token generation failed"
   - Solution: Check JWT_SECRET configuration and Python dependencies

3. **Database Connection Issues**
   - Error: Database connection failures in tests
   - Solution: Ensure PostgreSQL is running and properly configured

4. **Permission Denied**
   - Error: Tests fail with 403 errors
   - Solution: Check user roles and RBAC configuration

5. **Test Data Conflicts**
   - Error: Tests fail due to existing data
   - Solution: Run with --skip-cleanup=false or clean database

### Debugging

1. **Enable Verbose Output**:
   ```bash
   ./run_all_tests.sh --verbose
   ```

2. **Skip Cleanup for Investigation**:
   ```bash
   ./run_all_tests.sh --skip-cleanup
   ```

3. **Run Individual Tests**:
   ```bash
   ./run_scenario.sh authentication_flow --verbose
   ```

4. **Check Logs**:
   - Review individual scenario reports in `reports/` directory
   - Check API server logs for backend issues
   - Review database logs for data issues

### Performance Tuning

1. **Increase Timeout**:
   ```bash
   ./run_all_tests.sh --timeout 60
   ```

2. **Parallel Execution**:
   ```bash
   ./run_all_tests.sh --parallel
   ```

3. **Optimize Database**:
   - Ensure proper indexing
   - Use connection pooling
   - Optimize test data size

## Best Practices

### Test Development

1. **Isolation**: Each test should be independent
2. **Cleanup**: Always clean up test data
3. **Validation**: Comprehensive response validation
4. **Error Handling**: Robust error handling
5. **Documentation**: Clear test descriptions

### Maintenance

1. **Regular Updates**: Keep tests updated with API changes
2. **Data Management**: Maintain test data consistency
3. **Performance Monitoring**: Monitor test execution times
4. **Coverage Analysis**: Ensure comprehensive coverage

### Security

1. **Test Data**: Use non-sensitive test data
2. **Credentials**: Secure test credentials
3. **Environment Isolation**: Use separate test environments
4. **Access Control**: Proper test user permissions

## Contributing

### Adding New Tests

1. **Create Test Function**: Follow naming convention `test_<functionality>`
2. **Add Documentation**: Include detailed test description
3. **Add Cleanup**: Ensure proper cleanup of test data
4. **Update Documentation**: Update this documentation
5. **Test Thoroughly**: Test in isolation and with full suite

### Modifying Existing Tests

1. **Understand Impact**: Consider impact on other tests
2. **Maintain Compatibility**: Ensure backward compatibility
3. **Update Documentation**: Update relevant documentation
4. **Test Changes**: Run full test suite to verify changes

### Code Review Checklist

- [ ] Test is properly documented
- [ ] Test includes proper cleanup
- [ ] Test is idempotent
- [ ] Test handles errors appropriately
- [ ] Test follows naming conventions
- [ ] Test is independent of other tests
- [ ] Documentation is updated

## Conclusion

This E2E test suite provides comprehensive coverage of the Affiliate Platform API, ensuring system reliability, security, and performance. The modular design allows for easy maintenance and extension, while the detailed reporting provides valuable insights into system health and test results.

For questions or issues, please refer to the troubleshooting section or contact the development team.