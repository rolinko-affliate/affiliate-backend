/*
Everflow Network API - Tracking

API for generating tracking links in the Everflow platform

API version: 1.0.0
*/

package tracking

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)

// TrackingAPIService TrackingAPI service
type TrackingAPIService service

type ApiCreateTrackingLinkRequest struct {
	ctx                        context.Context
	apiService                 *TrackingAPIService
	createTrackingLinkRequest  *CreateTrackingLinkRequest
}

func (r ApiCreateTrackingLinkRequest) CreateTrackingLinkRequest(createTrackingLinkRequest CreateTrackingLinkRequest) ApiCreateTrackingLinkRequest {
	r.createTrackingLinkRequest = &createTrackingLinkRequest
	return r
}

func (r ApiCreateTrackingLinkRequest) Execute() (*TrackingLinkResponse, *http.Response, error) {
	return r.apiService.CreateTrackingLinkExecute(r)
}

/*
CreateTrackingLink Generate a tracking link

Generate a tracking link for an affiliate and offer combination

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiCreateTrackingLinkRequest
*/
func (a *TrackingAPIService) CreateTrackingLink(ctx context.Context) ApiCreateTrackingLinkRequest {
	return ApiCreateTrackingLinkRequest{
		apiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//  @return TrackingLinkResponse
func (a *TrackingAPIService) CreateTrackingLinkExecute(r ApiCreateTrackingLinkRequest) (*TrackingLinkResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *TrackingLinkResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TrackingAPIService.CreateTrackingLink")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/networks/tracking/offers/clicks"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.createTrackingLinkRequest == nil {
		return localVarReturnValue, nil, reportError("createTrackingLinkRequest is required and must be specified")
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
	localVarPostBody = r.createTrackingLinkRequest
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

