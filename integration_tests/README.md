# Everflow Integration Tests

This directory contains comprehensive integration tests for synchronizing core entities between our affiliate platform and the Everflow network.

## 🎯 Test Coverage

### Core Entity Synchronization Tests

| Test | Entity | Everflow Counterpart | Status |
|------|--------|---------------------|---------|
| `TestAdvertiserSynchronization` | Advertiser | Advertiser | ⚠️ Limited by API key |
| `TestAffiliateSynchronization` | Affiliate | Partner | 🚧 Ready for implementation |
| `TestCampaignSynchronization` | Campaign | Offer | 🚧 Ready for implementation |
| `TestTrackingLinkSynchronization` | Tracking Link | Tracking Link | 🚧 Ready for implementation |
| `TestFullSynchronizationWorkflow` | All entities | Complete workflow | 🚧 Pending individual tests |

### Direct API Tests

| Test | Purpose | Status |
|------|---------|---------|
| `TestEverflowDirectAPI` | Direct Everflow API testing | ✅ Working |

## 🔑 API Key Limitations

**Current Issue**: The Everflow API key has **read-only permissions**.

- ✅ **GET requests work**: Can retrieve existing advertisers, affiliates, offers, tracking links
- ❌ **POST requests fail**: Cannot create new entities (returns 400 "Invalid parameters")

See [API_KEY_PERMISSIONS.md](./API_KEY_PERMISSIONS.md) for detailed analysis and solutions.

## 🚀 Running Tests

### Prerequisites

1. **Server running**: Start the affiliate backend server without `--mock` flag
   ```bash
   cd /workspace && go run cmd/api/main.go --auto-migrate
   ```

2. **Environment variables**: Set the Everflow API key
   ```bash
   export EVERFLOW_API_KEY="3ytTLxgNTZW3rpRwz7ormw"
   ```

### Test Execution

```bash
# Run all integration tests
go test ./integration_tests/ -v

# Run specific test
go test ./integration_tests/ -run TestAdvertiserSynchronization -v

# Run in mock mode (when proper API key is available)
export EVERFLOW_MOCK_MODE=true
go test ./integration_tests/ -v
```

## 📁 File Structure

```
integration_tests/
├── README.md                           # This file
├── API_KEY_PERMISSIONS.md              # API key limitation analysis
├── test_helpers.go                     # Core test utilities and helpers
├── cleanup_helpers.go                  # Entity cleanup and tracking
├── everflow_sync_integration_test.go   # Main synchronization tests
└── everflow_direct_test.go            # Direct API testing
```

## 🛠️ Test Framework Features

### Test Helpers (`test_helpers.go`)
- JWT token generation for authentication
- HTTP request utilities for platform and Everflow APIs
- Response parsing and validation
- Test data generation (names, emails, URLs)

### Cleanup System (`cleanup_helpers.go`)
- Automatic tracking of created entities
- Cleanup on test completion or failure
- Support for platform and Everflow entity cleanup

### Integration Tests (`everflow_sync_integration_test.go`)
- End-to-end synchronization testing
- API key permission handling
- Mock mode support for development
- Comprehensive entity lifecycle testing

## 🔄 Test Workflow

Each synchronization test follows this pattern:

1. **Setup**: Create test user profile and organization
2. **Create**: Create entity via our platform API
3. **Verify**: Check entity was created successfully
4. **Sync**: Wait for or trigger synchronization to Everflow
5. **Validate**: Verify entity exists in Everflow with correct attributes
6. **Cleanup**: Remove all created entities

## ⚠️ Current Limitations

1. **API Key Permissions**: Read-only access limits full testing
2. **Sync Implementation**: Some entity types may not have sync implemented yet
3. **Mock Mode**: Not yet fully implemented (ready for proper API key)

## 🎯 Next Steps

1. **Get Read-Write API Key**: Contact Everflow team for proper permissions
2. **Complete Sync Implementation**: Implement missing entity synchronization
3. **Add Mock Mode**: Full mock implementation for development testing
4. **Error Handling**: Add comprehensive error scenario testing
5. **Performance Testing**: Add load and performance tests

## 📊 Test Results

With current read-only API key:
- ✅ Platform entity creation works
- ✅ Everflow API connectivity confirmed
- ✅ Read operations successful
- ❌ Write operations blocked by permissions
- ✅ Test framework fully functional

## 🤝 Contributing

When adding new tests:

1. Follow the established test pattern
2. Use the cleanup system to track entities
3. Add appropriate error handling
4. Update this README with new test coverage
5. Ensure tests work in both normal and mock modes