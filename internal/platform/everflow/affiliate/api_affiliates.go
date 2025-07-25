/*
Everflow Affiliates API

API for managing affiliates in the Everflow platform

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package affiliate

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// AffiliatesAPIService AffiliatesAPI service
type AffiliatesAPIService service

type ApiCreateAffiliateRequest struct {
	ctx                    context.Context
	ApiService             *AffiliatesAPIService
	createAffiliateRequest *CreateAffiliateRequest
}

func (r ApiCreateAffiliateRequest) CreateAffiliateRequest(createAffiliateRequest CreateAffiliateRequest) ApiCreateAffiliateRequest {
	r.createAffiliateRequest = &createAffiliateRequest
	return r
}

func (r ApiCreateAffiliateRequest) Execute() (*Affiliate, *http.Response, error) {
	return r.ApiService.CreateAffiliateExecute(r)
}

/*
CreateAffiliate Create Affiliate

Creates a new affiliate in the network

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiCreateAffiliateRequest
*/
func (a *AffiliatesAPIService) CreateAffiliate(ctx context.Context) ApiCreateAffiliateRequest {
	return ApiCreateAffiliateRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//
//	@return Affiliate
func (a *AffiliatesAPIService) CreateAffiliateExecute(r ApiCreateAffiliateRequest) (*Affiliate, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Affiliate
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AffiliatesAPIService.CreateAffiliate")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/networks/affiliates"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.createAffiliateRequest == nil {
		return localVarReturnValue, nil, reportError("createAffiliateRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.createAffiliateRequest
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["ApiKeyAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["X-Eflow-API-Key"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiGetAffiliateByIdRequest struct {
	ctx          context.Context
	ApiService   *AffiliatesAPIService
	affiliateId  int32
	relationship *[]string
}

// Additional relationships to include in the response
func (r ApiGetAffiliateByIdRequest) Relationship(relationship []string) ApiGetAffiliateByIdRequest {
	r.relationship = &relationship
	return r
}

func (r ApiGetAffiliateByIdRequest) Execute() (*AffiliateWithRelationships, *http.Response, error) {
	return r.ApiService.GetAffiliateByIdExecute(r)
}

/*
GetAffiliateById Find Affiliate By ID

Retrieves a single affiliate by its ID

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param affiliateId The ID of the affiliate you want to fetch
	@return ApiGetAffiliateByIdRequest
*/
func (a *AffiliatesAPIService) GetAffiliateById(ctx context.Context, affiliateId int32) ApiGetAffiliateByIdRequest {
	return ApiGetAffiliateByIdRequest{
		ApiService:  a,
		ctx:         ctx,
		affiliateId: affiliateId,
	}
}

// Execute executes the request
//
//	@return AffiliateWithRelationships
func (a *AffiliatesAPIService) GetAffiliateByIdExecute(r ApiGetAffiliateByIdRequest) (*AffiliateWithRelationships, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *AffiliateWithRelationships
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AffiliatesAPIService.GetAffiliateById")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/networks/affiliates/{affiliateId}"
	localVarPath = strings.Replace(localVarPath, "{"+"affiliateId"+"}", url.PathEscape(parameterValueToString(r.affiliateId, "affiliateId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.relationship != nil {
		t := *r.relationship
		if reflect.TypeOf(t).Kind() == reflect.Slice {
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				parameterAddToHeaderOrQuery(localVarQueryParams, "relationship", s.Index(i).Interface(), "form", "multi")
			}
		} else {
			parameterAddToHeaderOrQuery(localVarQueryParams, "relationship", t, "form", "multi")
		}
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["ApiKeyAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["X-Eflow-API-Key"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiUpdateAffiliateRequest struct {
	ctx                    context.Context
	ApiService             *AffiliatesAPIService
	affiliateId            int32
	updateAffiliateRequest *UpdateAffiliateRequest
}

func (r ApiUpdateAffiliateRequest) UpdateAffiliateRequest(updateAffiliateRequest UpdateAffiliateRequest) ApiUpdateAffiliateRequest {
	r.updateAffiliateRequest = &updateAffiliateRequest
	return r
}

func (r ApiUpdateAffiliateRequest) Execute() (*Affiliate, *http.Response, error) {
	return r.ApiService.UpdateAffiliateExecute(r)
}

/*
UpdateAffiliate Update Affiliate

Updates an existing affiliate. All fields must be specified, not only the ones you wish to update.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param affiliateId The ID of the affiliate you want to update
	@return ApiUpdateAffiliateRequest
*/
func (a *AffiliatesAPIService) UpdateAffiliate(ctx context.Context, affiliateId int32) ApiUpdateAffiliateRequest {
	return ApiUpdateAffiliateRequest{
		ApiService:  a,
		ctx:         ctx,
		affiliateId: affiliateId,
	}
}

// Execute executes the request
//
//	@return Affiliate
func (a *AffiliatesAPIService) UpdateAffiliateExecute(r ApiUpdateAffiliateRequest) (*Affiliate, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Affiliate
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AffiliatesAPIService.UpdateAffiliate")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/networks/affiliates/{affiliateId}"
	localVarPath = strings.Replace(localVarPath, "{"+"affiliateId"+"}", url.PathEscape(parameterValueToString(r.affiliateId, "affiliateId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.updateAffiliateRequest == nil {
		return localVarReturnValue, nil, reportError("updateAffiliateRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.updateAffiliateRequest
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["ApiKeyAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["X-Eflow-API-Key"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
