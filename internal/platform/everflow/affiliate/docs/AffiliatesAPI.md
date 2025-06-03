# \AffiliatesAPI

All URIs are relative to *https://api.eflow.team/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateAffiliate**](AffiliatesAPI.md#CreateAffiliate) | **Post** /networks/affiliates | Create Affiliate
[**GetAffiliateById**](AffiliatesAPI.md#GetAffiliateById) | **Get** /networks/affiliates/{affiliateId} | Find Affiliate By ID
[**UpdateAffiliate**](AffiliatesAPI.md#UpdateAffiliate) | **Put** /networks/affiliates/{affiliateId} | Update Affiliate



## CreateAffiliate

> Affiliate CreateAffiliate(ctx).CreateAffiliateRequest(createAffiliateRequest).Execute()

Create Affiliate



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
	createAffiliateRequest := *openapiclient.NewCreateAffiliateRequest("Name_example", "AccountStatus_example", int32(123)) // CreateAffiliateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AffiliatesAPI.CreateAffiliate(context.Background()).CreateAffiliateRequest(createAffiliateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AffiliatesAPI.CreateAffiliate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateAffiliate`: Affiliate
	fmt.Fprintf(os.Stdout, "Response from `AffiliatesAPI.CreateAffiliate`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateAffiliateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createAffiliateRequest** | [**CreateAffiliateRequest**](CreateAffiliateRequest.md) |  | 

### Return type

[**Affiliate**](Affiliate.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAffiliateById

> AffiliateWithRelationships GetAffiliateById(ctx, affiliateId).Relationship(relationship).Execute()

Find Affiliate By ID



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
	affiliateId := int32(56) // int32 | The ID of the affiliate you want to fetch
	relationship := []string{"Relationship_example"} // []string | Additional relationships to include in the response (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AffiliatesAPI.GetAffiliateById(context.Background(), affiliateId).Relationship(relationship).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AffiliatesAPI.GetAffiliateById``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAffiliateById`: AffiliateWithRelationships
	fmt.Fprintf(os.Stdout, "Response from `AffiliatesAPI.GetAffiliateById`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**affiliateId** | **int32** | The ID of the affiliate you want to fetch | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAffiliateByIdRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **relationship** | **[]string** | Additional relationships to include in the response | 

### Return type

[**AffiliateWithRelationships**](AffiliateWithRelationships.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateAffiliate

> Affiliate UpdateAffiliate(ctx, affiliateId).UpdateAffiliateRequest(updateAffiliateRequest).Execute()

Update Affiliate



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
	affiliateId := int32(56) // int32 | The ID of the affiliate you want to update
	updateAffiliateRequest := *openapiclient.NewUpdateAffiliateRequest("Name_example", "AccountStatus_example", int32(123)) // UpdateAffiliateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AffiliatesAPI.UpdateAffiliate(context.Background(), affiliateId).UpdateAffiliateRequest(updateAffiliateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AffiliatesAPI.UpdateAffiliate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateAffiliate`: Affiliate
	fmt.Fprintf(os.Stdout, "Response from `AffiliatesAPI.UpdateAffiliate`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**affiliateId** | **int32** | The ID of the affiliate you want to update | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateAffiliateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateAffiliateRequest** | [**UpdateAffiliateRequest**](UpdateAffiliateRequest.md) |  | 

### Return type

[**Affiliate**](Affiliate.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

