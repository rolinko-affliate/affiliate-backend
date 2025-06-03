# UpdateAffiliateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the affiliate | 
**AccountStatus** | **string** | The affiliate&#39;s account status | 
**NetworkEmployeeId** | **int32** | The employee id of the affiliate&#39;s account manager | 
**InternalNotes** | Pointer to **string** | Internal notes for network usage | [optional] 
**DefaultCurrencyId** | Pointer to **string** | The affiliate&#39;s default currency (3-letter code) | [optional] 
**EnableMediaCostTrackingLinks** | Pointer to **bool** | Whether to allow affiliate to pass and override cost in their tracking links | [optional] 
**ReferrerId** | Pointer to **int32** | The id of the affiliate that referred the new affiliate | [optional] 
**IsContactAddressEnabled** | Pointer to **bool** | Whether to include a contact address for this affiliate | [optional] 
**NetworkAffiliateTierId** | Pointer to **int32** | The ID of the Affiliate Tier | [optional] 
**ContactAddress** | Pointer to [**ContactAddress**](ContactAddress.md) |  | [optional] 
**Labels** | Pointer to **[]string** | Labels to associate with the affiliate | [optional] 
**Billing** | Pointer to [**BillingInfo**](BillingInfo.md) |  | [optional] 

## Methods

### NewUpdateAffiliateRequest

`func NewUpdateAffiliateRequest(name string, accountStatus string, networkEmployeeId int32, ) *UpdateAffiliateRequest`

NewUpdateAffiliateRequest instantiates a new UpdateAffiliateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateAffiliateRequestWithDefaults

`func NewUpdateAffiliateRequestWithDefaults() *UpdateAffiliateRequest`

NewUpdateAffiliateRequestWithDefaults instantiates a new UpdateAffiliateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *UpdateAffiliateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *UpdateAffiliateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *UpdateAffiliateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetAccountStatus

`func (o *UpdateAffiliateRequest) GetAccountStatus() string`

GetAccountStatus returns the AccountStatus field if non-nil, zero value otherwise.

### GetAccountStatusOk

`func (o *UpdateAffiliateRequest) GetAccountStatusOk() (*string, bool)`

GetAccountStatusOk returns a tuple with the AccountStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountStatus

`func (o *UpdateAffiliateRequest) SetAccountStatus(v string)`

SetAccountStatus sets AccountStatus field to given value.


### GetNetworkEmployeeId

`func (o *UpdateAffiliateRequest) GetNetworkEmployeeId() int32`

GetNetworkEmployeeId returns the NetworkEmployeeId field if non-nil, zero value otherwise.

### GetNetworkEmployeeIdOk

`func (o *UpdateAffiliateRequest) GetNetworkEmployeeIdOk() (*int32, bool)`

GetNetworkEmployeeIdOk returns a tuple with the NetworkEmployeeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkEmployeeId

`func (o *UpdateAffiliateRequest) SetNetworkEmployeeId(v int32)`

SetNetworkEmployeeId sets NetworkEmployeeId field to given value.


### GetInternalNotes

`func (o *UpdateAffiliateRequest) GetInternalNotes() string`

GetInternalNotes returns the InternalNotes field if non-nil, zero value otherwise.

### GetInternalNotesOk

`func (o *UpdateAffiliateRequest) GetInternalNotesOk() (*string, bool)`

GetInternalNotesOk returns a tuple with the InternalNotes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternalNotes

`func (o *UpdateAffiliateRequest) SetInternalNotes(v string)`

SetInternalNotes sets InternalNotes field to given value.

### HasInternalNotes

`func (o *UpdateAffiliateRequest) HasInternalNotes() bool`

HasInternalNotes returns a boolean if a field has been set.

### GetDefaultCurrencyId

`func (o *UpdateAffiliateRequest) GetDefaultCurrencyId() string`

GetDefaultCurrencyId returns the DefaultCurrencyId field if non-nil, zero value otherwise.

### GetDefaultCurrencyIdOk

`func (o *UpdateAffiliateRequest) GetDefaultCurrencyIdOk() (*string, bool)`

GetDefaultCurrencyIdOk returns a tuple with the DefaultCurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultCurrencyId

`func (o *UpdateAffiliateRequest) SetDefaultCurrencyId(v string)`

SetDefaultCurrencyId sets DefaultCurrencyId field to given value.

### HasDefaultCurrencyId

`func (o *UpdateAffiliateRequest) HasDefaultCurrencyId() bool`

HasDefaultCurrencyId returns a boolean if a field has been set.

### GetEnableMediaCostTrackingLinks

`func (o *UpdateAffiliateRequest) GetEnableMediaCostTrackingLinks() bool`

GetEnableMediaCostTrackingLinks returns the EnableMediaCostTrackingLinks field if non-nil, zero value otherwise.

### GetEnableMediaCostTrackingLinksOk

`func (o *UpdateAffiliateRequest) GetEnableMediaCostTrackingLinksOk() (*bool, bool)`

GetEnableMediaCostTrackingLinksOk returns a tuple with the EnableMediaCostTrackingLinks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableMediaCostTrackingLinks

`func (o *UpdateAffiliateRequest) SetEnableMediaCostTrackingLinks(v bool)`

SetEnableMediaCostTrackingLinks sets EnableMediaCostTrackingLinks field to given value.

### HasEnableMediaCostTrackingLinks

`func (o *UpdateAffiliateRequest) HasEnableMediaCostTrackingLinks() bool`

HasEnableMediaCostTrackingLinks returns a boolean if a field has been set.

### GetReferrerId

`func (o *UpdateAffiliateRequest) GetReferrerId() int32`

GetReferrerId returns the ReferrerId field if non-nil, zero value otherwise.

### GetReferrerIdOk

`func (o *UpdateAffiliateRequest) GetReferrerIdOk() (*int32, bool)`

GetReferrerIdOk returns a tuple with the ReferrerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReferrerId

`func (o *UpdateAffiliateRequest) SetReferrerId(v int32)`

SetReferrerId sets ReferrerId field to given value.

### HasReferrerId

`func (o *UpdateAffiliateRequest) HasReferrerId() bool`

HasReferrerId returns a boolean if a field has been set.

### GetIsContactAddressEnabled

`func (o *UpdateAffiliateRequest) GetIsContactAddressEnabled() bool`

GetIsContactAddressEnabled returns the IsContactAddressEnabled field if non-nil, zero value otherwise.

### GetIsContactAddressEnabledOk

`func (o *UpdateAffiliateRequest) GetIsContactAddressEnabledOk() (*bool, bool)`

GetIsContactAddressEnabledOk returns a tuple with the IsContactAddressEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsContactAddressEnabled

`func (o *UpdateAffiliateRequest) SetIsContactAddressEnabled(v bool)`

SetIsContactAddressEnabled sets IsContactAddressEnabled field to given value.

### HasIsContactAddressEnabled

`func (o *UpdateAffiliateRequest) HasIsContactAddressEnabled() bool`

HasIsContactAddressEnabled returns a boolean if a field has been set.

### GetNetworkAffiliateTierId

`func (o *UpdateAffiliateRequest) GetNetworkAffiliateTierId() int32`

GetNetworkAffiliateTierId returns the NetworkAffiliateTierId field if non-nil, zero value otherwise.

### GetNetworkAffiliateTierIdOk

`func (o *UpdateAffiliateRequest) GetNetworkAffiliateTierIdOk() (*int32, bool)`

GetNetworkAffiliateTierIdOk returns a tuple with the NetworkAffiliateTierId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkAffiliateTierId

`func (o *UpdateAffiliateRequest) SetNetworkAffiliateTierId(v int32)`

SetNetworkAffiliateTierId sets NetworkAffiliateTierId field to given value.

### HasNetworkAffiliateTierId

`func (o *UpdateAffiliateRequest) HasNetworkAffiliateTierId() bool`

HasNetworkAffiliateTierId returns a boolean if a field has been set.

### GetContactAddress

`func (o *UpdateAffiliateRequest) GetContactAddress() ContactAddress`

GetContactAddress returns the ContactAddress field if non-nil, zero value otherwise.

### GetContactAddressOk

`func (o *UpdateAffiliateRequest) GetContactAddressOk() (*ContactAddress, bool)`

GetContactAddressOk returns a tuple with the ContactAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContactAddress

`func (o *UpdateAffiliateRequest) SetContactAddress(v ContactAddress)`

SetContactAddress sets ContactAddress field to given value.

### HasContactAddress

`func (o *UpdateAffiliateRequest) HasContactAddress() bool`

HasContactAddress returns a boolean if a field has been set.

### GetLabels

`func (o *UpdateAffiliateRequest) GetLabels() []string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *UpdateAffiliateRequest) GetLabelsOk() (*[]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *UpdateAffiliateRequest) SetLabels(v []string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *UpdateAffiliateRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetBilling

`func (o *UpdateAffiliateRequest) GetBilling() BillingInfo`

GetBilling returns the Billing field if non-nil, zero value otherwise.

### GetBillingOk

`func (o *UpdateAffiliateRequest) GetBillingOk() (*BillingInfo, bool)`

GetBillingOk returns a tuple with the Billing field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBilling

`func (o *UpdateAffiliateRequest) SetBilling(v BillingInfo)`

SetBilling sets Billing field to given value.

### HasBilling

`func (o *UpdateAffiliateRequest) HasBilling() bool`

HasBilling returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


