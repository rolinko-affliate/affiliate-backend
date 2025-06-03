# \DefaultAPI

All URIs are relative to *https://api.eflow.team*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateAdvertiser**](DefaultAPI.md#CreateAdvertiser) | **Post** /v1/networks/advertisers | Create Advertiser
[**GetAdvertiserById**](DefaultAPI.md#GetAdvertiserById) | **Get** /v1/networks/advertisers/{advertiserId} | Get Advertiser by ID
[**UpdateAdvertiser**](DefaultAPI.md#UpdateAdvertiser) | **Put** /v1/networks/advertisers/{advertiserId} | Update Advertiser



## CreateAdvertiser

> Advertiser CreateAdvertiser(ctx).CreateAdvertiserRequest(createAdvertiserRequest).Execute()

Create Advertiser



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
	createAdvertiserRequest := *openapiclient.NewCreateAdvertiserRequest("Name_example", "AccountStatus_example", int32(123), "DefaultCurrencyId_example", int32(123), "AttributionMethod_example", "EmailAttributionMethod_example", "AttributionPriority_example") // CreateAdvertiserRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.CreateAdvertiser(context.Background()).CreateAdvertiserRequest(createAdvertiserRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.CreateAdvertiser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateAdvertiser`: Advertiser
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.CreateAdvertiser`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateAdvertiserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createAdvertiserRequest** | [**CreateAdvertiserRequest**](CreateAdvertiserRequest.md) |  | 

### Return type

[**Advertiser**](Advertiser.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAdvertiserById

> Advertiser GetAdvertiserById(ctx, advertiserId).Relationship(relationship).Execute()

Get Advertiser by ID



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
	advertiserId := int32(56) // int32 | The ID of the advertiser to retrieve
	relationship := []string{"Relationship_example"} // []string | Additional relationships to include in the response (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetAdvertiserById(context.Background(), advertiserId).Relationship(relationship).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetAdvertiserById``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAdvertiserById`: Advertiser
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetAdvertiserById`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**advertiserId** | **int32** | The ID of the advertiser to retrieve | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAdvertiserByIdRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **relationship** | **[]string** | Additional relationships to include in the response | 

### Return type

[**Advertiser**](Advertiser.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateAdvertiser

> Advertiser UpdateAdvertiser(ctx, advertiserId).UpdateAdvertiserRequest(updateAdvertiserRequest).Execute()

Update Advertiser



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
	advertiserId := int32(56) // int32 | The ID of the advertiser to update
	updateAdvertiserRequest := *openapiclient.NewUpdateAdvertiserRequest("Name_example", "AccountStatus_example", int32(123), "DefaultCurrencyId_example", int32(123), "AttributionMethod_example", "EmailAttributionMethod_example", "AttributionPriority_example") // UpdateAdvertiserRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.UpdateAdvertiser(context.Background(), advertiserId).UpdateAdvertiserRequest(updateAdvertiserRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.UpdateAdvertiser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateAdvertiser`: Advertiser
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.UpdateAdvertiser`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**advertiserId** | **int32** | The ID of the advertiser to update | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateAdvertiserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateAdvertiserRequest** | [**UpdateAdvertiserRequest**](UpdateAdvertiserRequest.md) |  | 

### Return type

[**Advertiser**](Advertiser.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

