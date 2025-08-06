# Everflow Integration

This package provides integration with the Everflow affiliate marketing platform.

## Overview

The Everflow integration allows you to:
- Create affiliates (partners) in Everflow
- Sync affiliate data between your system and Everflow
- Maintain provider mappings to track synchronization status

## Configuration

### API Credentials

To use the Everflow integration, you need to configure the API clients with your Everflow credentials:

```go
import (
    "github.com/affiliate-backend/internal/platform/everflow"
    "github.com/affiliate-backend/internal/platform/everflow/affiliate"
)

// Configure the affiliate client
config := affiliate.NewConfiguration()
config.Host = "api.eflow.team" // Your Everflow API host
config.DefaultHeader["X-Eflow-API-Key"] = "your-api-key-here"

// Create the client
client := affiliate.NewAPIClient(config)
```

### Integration Service

Create an integration service instance:

```go
service := everflow.NewIntegrationService(
    advertiserClient,
    affiliateClient,
    offerClient,
    advertiserRepo,
    affiliateRepo,
    campaignRepo,
    advertiserProviderMappingRepo,
    affiliateProviderMappingRepo,
    campaignProviderMappingRepo,
)
```

## Usage

### Creating an Affiliate

```go
ctx := context.Background()

// Create your affiliate domain object
affiliate := domain.Affiliate{
    AffiliateID:      12345,
    OrganizationID:   1,
    Name:             "Example Affiliate",
    ContactEmail:     "affiliate@example.com",
    Status:           "active",
    InternalNotes:    stringPtr("High-performing affiliate"),
    DefaultCurrency:  stringPtr("USD"),
    InvoiceAmountThreshold: float64Ptr(1000.00),
    DefaultPaymentTerms:    int32Ptr(30),
    ContactAddress: &domain.Address{
        Address1:    "123 Main St",
        City:        "New York",
        RegionCode:  "NY",
        CountryCode: "US",
        ZipCode:     "10001",
    },
    BillingInfo: &domain.BillingInfo{
        CompanyName: "Example Affiliate LLC",
        TaxID:       "12-3456789",
    },
    PaymentDetails: &domain.PaymentDetails{
        BankAccount:   "123456789",
        RoutingNumber: "987654321",
        PaymentMethod: "bank_transfer",
    },
    Labels: []string{"premium", "high-volume", "trusted"},
}

// Create the affiliate in Everflow
result, err := service.CreateAffiliate(ctx, affiliate)
if err != nil {
    log.Printf("Failed to create affiliate: %v", err)
    return
}

fmt.Printf("Affiliate created successfully: %+v\n", result)
```

## Data Mapping

The integration handles mapping between your internal domain models and Everflow's API format:

### Internal → Everflow Mapping

- `domain.Affiliate` → `affiliate.CreateAffiliateRequest`
- Status mapping: `active` → `active`, `pending` → `pending`, etc.
- Address and billing information mapping
- Payment details mapping
- Labels converted to JSON array format

### Everflow → Internal Mapping

- `affiliate.Affiliate` response → `domain.Affiliate` updates
- Provider-specific data stored in `domain.EverflowProviderData`
- Provider mappings track synchronization status

## Provider Mappings

The integration automatically creates and manages provider mappings to track:
- Synchronization status (`synced`, `failed`, `pending`)
- Provider-specific IDs and data
- Request/response payloads for debugging
- Last sync timestamps

## Error Handling

The integration provides comprehensive error handling:
- API communication errors
- Data validation errors
- Mapping errors
- Provider mapping creation/update errors

All errors include detailed context to help with debugging.

## Testing

Run the test suite:

```bash
go test ./internal/platform/everflow/... -v
```

The test suite includes:
- Unit tests for all mapper functions
- Integration tests with mock HTTP servers
- Round-trip data mapping tests
- Error scenario testing

## Debugging

The integration includes debug logging that shows:
- Outbound API requests
- API responses
- Provider mapping operations

Enable debug output by checking the console for `DEBUG:` prefixed messages.

## Example Test Program

See `cmd/test_affiliate_creation/main.go` for a complete example of how to use the integration.

## API Reference

### Key Functions

- `CreateAffiliate(ctx, affiliate)` - Creates an affiliate in Everflow
- `UpdateAffiliate(ctx, affiliate)` - Updates an existing affiliate
- `GetAffiliate(ctx, affiliateID)` - Retrieves affiliate data

### Key Types

- `IntegrationService` - Main service interface
- `AffiliateProviderMapper` - Handles data mapping
- `domain.Affiliate` - Internal affiliate representation
- `domain.AffiliateProviderMapping` - Provider sync tracking

## Troubleshooting

### Common Issues

1. **API Authentication Errors**
   - Verify your API key is correct
   - Check that the API host is properly configured
   - Ensure your API key has the necessary permissions

2. **Mapping Errors**
   - Check that required fields are populated
   - Verify data types match expected formats
   - Review debug logs for detailed error information

3. **Network Issues**
   - Verify connectivity to Everflow API endpoints
   - Check firewall and proxy settings
   - Review HTTP response codes and bodies

### Debug Steps

1. Enable debug logging to see request/response data
2. Check provider mapping status in your database
3. Review error messages for specific failure points
4. Test with the example program to isolate issues

## Implementation Status

✅ **COMPLETED**: The Everflow affiliate integration is fully implemented and tested.

### Key Features:
- ✅ Complete affiliate creation in Everflow
- ✅ Proper API key configuration via factory pattern
- ✅ Comprehensive data mapping between domain and Everflow models
- ✅ Provider mapping tracking for sync status
- ✅ Enhanced error handling with detailed API responses
- ✅ Full test coverage with unit and integration tests

### Resolved Issues:
- **API Endpoint URL**: Fixed affiliate API to use correct `/v1/networks/affiliates` endpoint
- **Server Configuration**: Properly configured server URL as `https://api.eflow.team/v1` for affiliate client
- **API Authentication**: Verified API key configuration through `X-Eflow-API-Key` header

### Verification:
- ✅ All unit tests pass
- ✅ Integration tests demonstrate successful API calls
- ✅ Manual testing confirms affiliate creation in Everflow (returns network_affiliate_id)
- ✅ API key from .env file is properly loaded and used
- ✅ Production-ready implementation