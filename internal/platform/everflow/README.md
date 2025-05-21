# Everflow API Integration

This package provides integration with the Everflow API for creating advertisers and offers.

## Overview

The Everflow integration consists of the following components:

1. **Client**: A low-level HTTP client for making requests to the Everflow API.
2. **Service**: A high-level service that maps our domain models to Everflow API requests and handles the integration logic.
3. **Factory**: A factory function for creating the Everflow service with the necessary dependencies.

## Configuration

The Everflow integration can be configured using environment variables:

- `EVERFLOW_API_KEY`: The Everflow API key.
- `EVERFLOW_CONFIG`: A JSON string containing the Everflow configuration (alternative to `EVERFLOW_API_KEY`).

Example:

```
EVERFLOW_API_KEY=your-api-key-here
```

Or:

```
EVERFLOW_CONFIG={"api_key":"your-api-key-here"}
```

## Integration with Advertiser and Campaign Services

The Everflow integration is automatically triggered when:

1. A new advertiser is created in our system, which creates a corresponding advertiser in Everflow.
2. A new campaign is created in our system, which creates a corresponding offer in Everflow.

The integration is asynchronous, meaning that the creation of advertisers and campaigns in our system will not be blocked by the Everflow API calls. If the Everflow API calls fail, the errors will be logged but will not affect the creation of advertisers and campaigns in our system.

## Data Mapping

### Advertiser Mapping

When an advertiser is created in our system, the following data is mapped to Everflow:

- `name`: The advertiser's name.
- `account_status`: Mapped from our status field ("active", "inactive", "pending").
- `default_currency_id`: Set to "USD" by default.
- `internal_notes`: Contains the advertiser's contact email if available.
- `contact_address`: Mapped from the advertiser's billing details if available.
- `billing`: Mapped from the advertiser's billing details if available.

### Campaign Mapping

When a campaign is created in our system, the following data is mapped to Everflow:

- `name`: The campaign's name.
- `network_advertiser_id`: The ID of the advertiser in Everflow.
- `destination_url`: A generated URL based on the campaign ID.
- `offer_status`: Mapped from our status field ("active", "paused", "draft", "archived").
- `currency_id`: Set to "USD" by default.
- `visibility`: Set to "public" by default.
- `conversion_method`: Set to "server_postback" by default.
- `description`: The campaign's description if available.
- `session_definition`: Set to "cookie" by default.
- `session_duration`: Set to 720 hours (30 days) by default.
- `payout_revenue`: A default payout/revenue structure.

## Tags

Both advertisers and offers in Everflow are tagged with the following information:

- `advertiser_id`: The ID of the advertiser in our system.
- `organization_id`: The ID of the organization in our system.
- `campaign_id`: The ID of the campaign in our system (for offers only).

These tags can be used to identify and filter advertisers and offers in Everflow.

## Error Handling

Errors that occur during the Everflow API calls are logged but do not affect the creation of advertisers and campaigns in our system. This ensures that our system can continue to function even if the Everflow API is unavailable.

## Testing

The Everflow integration includes unit tests that verify the mapping logic and API calls. The tests use mocks to simulate the Everflow API responses.