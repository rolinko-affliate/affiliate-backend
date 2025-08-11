# Everflow API Key Permissions Issue

## Current Status

The current Everflow API key (`3ytTLxgNTZW3rpRwz7ormw`) has **READ-ONLY** permissions.

### What Works ✅
- GET requests to retrieve existing entities:
  - `GET /v1/networks/advertisers` - Returns existing advertisers
  - `GET /v1/networks/affiliates` - Returns existing affiliates  
  - `GET /v1/networks/offers` - Returns existing offers
  - `GET /v1/networks/tracking_links` - Returns existing tracking links

### What Doesn't Work ❌
- POST requests to create new entities:
  - `POST /v1/networks/advertisers` - Returns 400 "Invalid parameters"
  - `POST /v1/networks/affiliates` - Returns 400 "Invalid parameters"
  - All creation attempts fail regardless of payload content

## Root Cause

The API key lacks write permissions. Even minimal payloads that match existing entity structures fail with "Invalid parameters", indicating a permissions issue rather than a data validation issue.

## Solutions

### Option 1: Get Read-Write API Key (Recommended)
Contact the Everflow API provider to:
1. Request a read-write API key for testing
2. Verify the key has permissions for:
   - Creating advertisers
   - Creating affiliates  
   - Creating offers
   - Creating tracking links

### Option 2: Use Mock Mode for Testing
The integration tests can be run in mock mode to simulate successful creation:
```bash
export EVERFLOW_MOCK_MODE=true
go test ./integration_tests/ -v
```

### Option 3: Test with Existing Entities
Focus on synchronization testing with existing entities:
- Verify existing advertisers are properly mapped
- Test data consistency between our platform and Everflow
- Validate read operations and data transformation

## Next Steps

1. **Immediate**: Document the limitation and create mock tests
2. **Short-term**: Request proper API key from Everflow team
3. **Long-term**: Implement full integration tests with read-write permissions

## Test Implementation Status

- ✅ Read operations tested and working
- ✅ API connectivity verified
- ✅ Data mapping logic implemented
- ❌ Write operations blocked by API key permissions
- ✅ Mock test framework ready for implementation