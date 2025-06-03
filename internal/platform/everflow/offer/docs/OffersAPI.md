# \OffersAPI

All URIs are relative to *https://api.eflow.team/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateOffer**](OffersAPI.md#CreateOffer) | **Post** /networks/offers | Create an offer
[**GetOfferById**](OffersAPI.md#GetOfferById) | **Get** /networks/offers/{offerId} | Find offer by ID
[**UpdateOffer**](OffersAPI.md#UpdateOffer) | **Put** /networks/offers/{offerId} | Update an offer



## CreateOffer

> OfferResponse CreateOffer(ctx).CreateOfferRequest(createOfferRequest).Execute()

Create an offer



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	createOfferRequest := *openapiclient.NewCreateOfferRequest(int32(123), "Name_example", "DestinationUrl_example", "OfferStatus_example", []openapiclient.PayoutRevenue{*openapiclient.NewPayoutRevenue("PayoutType_example", "RevenueType_example", false, false)}) // CreateOfferRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OffersAPI.CreateOffer(context.Background()).CreateOfferRequest(createOfferRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OffersAPI.CreateOffer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateOffer`: OfferResponse
	fmt.Fprintf(os.Stdout, "Response from `OffersAPI.CreateOffer`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateOfferRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createOfferRequest** | [**CreateOfferRequest**](CreateOfferRequest.md) |  | 

### Return type

[**OfferResponse**](OfferResponse.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOfferById

> OfferResponse GetOfferById(ctx, offerId).Relationship(relationship).Execute()

Find offer by ID



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	offerId := int32(56) // int32 | The ID of the offer you want to fetch
	relationship := "relationship_example" // string | Additional relationships to include (comma-separated) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OffersAPI.GetOfferById(context.Background(), offerId).Relationship(relationship).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OffersAPI.GetOfferById``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOfferById`: OfferResponse
	fmt.Fprintf(os.Stdout, "Response from `OffersAPI.GetOfferById`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**offerId** | **int32** | The ID of the offer you want to fetch | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetOfferByIdRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **relationship** | **string** | Additional relationships to include (comma-separated) | 

### Return type

[**OfferResponse**](OfferResponse.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateOffer

> OfferResponse UpdateOffer(ctx, offerId).UpdateOfferRequest(updateOfferRequest).Execute()

Update an offer



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	offerId := int32(56) // int32 | The ID of the offer you want to update
	updateOfferRequest := *openapiclient.NewUpdateOfferRequest(int32(123), "Name_example", "DestinationUrl_example", "OfferStatus_example", []openapiclient.PayoutRevenue{*openapiclient.NewPayoutRevenue("PayoutType_example", "RevenueType_example", false, false)}) // UpdateOfferRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OffersAPI.UpdateOffer(context.Background(), offerId).UpdateOfferRequest(updateOfferRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OffersAPI.UpdateOffer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateOffer`: OfferResponse
	fmt.Fprintf(os.Stdout, "Response from `OffersAPI.UpdateOffer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**offerId** | **int32** | The ID of the offer you want to update | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateOfferRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateOfferRequest** | [**UpdateOfferRequest**](UpdateOfferRequest.md) |  | 

### Return type

[**OfferResponse**](OfferResponse.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

