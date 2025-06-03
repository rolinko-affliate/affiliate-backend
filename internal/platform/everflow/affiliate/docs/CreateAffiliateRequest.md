# CreateAffiliateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the affiliate | 
**AccountStatus** | **string** | The affiliate&#39;s account status | 
**NetworkEmployeeId** | **int32** | The employee id of the affiliate&#39;s account manager | 
**InternalNotes** | Pointer to **string** | Internal notes for network usage | [optional] 
**DefaultCurrencyId** | Pointer to **string** | The affiliate&#39;s default currency (3-letter code) | [optional] 
**EnableMediaCostTrackingLinks** | Pointer to **bool** | Whether to allow affiliate to pass and override cost in their tracking links | [optional] [default to false]
**ReferrerId** | Pointer to **int32** | The id of the affiliate that referred the new affiliate | [optional] [default to 0]
**IsContactAddressEnabled** | Pointer to **bool** | Whether to include a contact address for this affiliate | [optional] [default to false]
**NetworkAffiliateTierId** | Pointer to **int32** | The ID of the Affiliate Tier | [optional] 
**ContactAddress** | Pointer to [**ContactAddress**](ContactAddress.md) |  | [optional] 
**Labels** | Pointer to **[]string** | Labels to associate with the affiliate | [optional] 
**Users** | Pointer to [**[]AffiliateUser**](AffiliateUser.md) | List of affiliate users to be created | [optional] 
**Billing** | Pointer to [**BillingInfo**](BillingInfo.md) |  | [optional] 

## Methods

### NewCreateAffiliateRequest

`func NewCreateAffiliateRequest(name string, accountStatus string, networkEmployeeId int32, ) *CreateAffiliateRequest`

NewCreateAffiliateRequest instantiates a new CreateAffiliateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateAffiliateRequestWithDefaults

`func NewCreateAffiliateRequestWithDefaults() *CreateAffiliateRequest`

NewCreateAffiliateRequestWithDefaults instantiates a new CreateAffiliateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *CreateAffiliateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateAffiliateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateAffiliateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetAccountStatus

`func (o *CreateAffiliateRequest) GetAccountStatus() string`

GetAccountStatus returns the AccountStatus field if non-nil, zero value otherwise.

### GetAccountStatusOk

`func (o *CreateAffiliateRequest) GetAccountStatusOk() (*string, bool)`

GetAccountStatusOk returns a tuple with the AccountStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountStatus

`func (o *CreateAffiliateRequest) SetAccountStatus(v string)`

SetAccountStatus sets AccountStatus field to given value.


### GetNetworkEmployeeId

`func (o *CreateAffiliateRequest) GetNetworkEmployeeId() int32`

GetNetworkEmployeeId returns the NetworkEmployeeId field if non-nil, zero value otherwise.

### GetNetworkEmployeeIdOk

`func (o *CreateAffiliateRequest) GetNetworkEmployeeIdOk() (*int32, bool)`

GetNetworkEmployeeIdOk returns a tuple with the NetworkEmployeeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkEmployeeId

`func (o *CreateAffiliateRequest) SetNetworkEmployeeId(v int32)`

SetNetworkEmployeeId sets NetworkEmployeeId field to given value.


### GetInternalNotes

`func (o *CreateAffiliateRequest) GetInternalNotes() string`

GetInternalNotes returns the InternalNotes field if non-nil, zero value otherwise.

### GetInternalNotesOk

`func (o *CreateAffiliateRequest) GetInternalNotesOk() (*string, bool)`

GetInternalNotesOk returns a tuple with the InternalNotes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternalNotes

`func (o *CreateAffiliateRequest) SetInternalNotes(v string)`

SetInternalNotes sets InternalNotes field to given value.

### HasInternalNotes

`func (o *CreateAffiliateRequest) HasInternalNotes() bool`

HasInternalNotes returns a boolean if a field has been set.

### GetDefaultCurrencyId

`func (o *CreateAffiliateRequest) GetDefaultCurrencyId() string`

GetDefaultCurrencyId returns the DefaultCurrencyId field if non-nil, zero value otherwise.

### GetDefaultCurrencyIdOk

`func (o *CreateAffiliateRequest) GetDefaultCurrencyIdOk() (*string, bool)`

GetDefaultCurrencyIdOk returns a tuple with the DefaultCurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultCurrencyId

`func (o *CreateAffiliateRequest) SetDefaultCurrencyId(v string)`

SetDefaultCurrencyId sets DefaultCurrencyId field to given value.

### HasDefaultCurrencyId

`func (o *CreateAffiliateRequest) HasDefaultCurrencyId() bool`

HasDefaultCurrencyId returns a boolean if a field has been set.

### GetEnableMediaCostTrackingLinks

`func (o *CreateAffiliateRequest) GetEnableMediaCostTrackingLinks() bool`

GetEnableMediaCostTrackingLinks returns the EnableMediaCostTrackingLinks field if non-nil, zero value otherwise.

### GetEnableMediaCostTrackingLinksOk

`func (o *CreateAffiliateRequest) GetEnableMediaCostTrackingLinksOk() (*bool, bool)`

GetEnableMediaCostTrackingLinksOk returns a tuple with the EnableMediaCostTrackingLinks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableMediaCostTrackingLinks

`func (o *CreateAffiliateRequest) SetEnableMediaCostTrackingLinks(v bool)`

SetEnableMediaCostTrackingLinks sets EnableMediaCostTrackingLinks field to given value.

### HasEnableMediaCostTrackingLinks

`func (o *CreateAffiliateRequest) HasEnableMediaCostTrackingLinks() bool`

HasEnableMediaCostTrackingLinks returns a boolean if a field has been set.

### GetReferrerId

`func (o *CreateAffiliateRequest) GetReferrerId() int32`

GetReferrerId returns the ReferrerId field if non-nil, zero value otherwise.

### GetReferrerIdOk

`func (o *CreateAffiliateRequest) GetReferrerIdOk() (*int32, bool)`

GetReferrerIdOk returns a tuple with the ReferrerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReferrerId

`func (o *CreateAffiliateRequest) SetReferrerId(v int32)`

SetReferrerId sets ReferrerId field to given value.

### HasReferrerId

`func (o *CreateAffiliateRequest) HasReferrerId() bool`

HasReferrerId returns a boolean if a field has been set.

### GetIsContactAddressEnabled

`func (o *CreateAffiliateRequest) GetIsContactAddressEnabled() bool`

GetIsContactAddressEnabled returns the IsContactAddressEnabled field if non-nil, zero value otherwise.

### GetIsContactAddressEnabledOk

`func (o *CreateAffiliateRequest) GetIsContactAddressEnabledOk() (*bool, bool)`

GetIsContactAddressEnabledOk returns a tuple with the IsContactAddressEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsContactAddressEnabled

`func (o *CreateAffiliateRequest) SetIsContactAddressEnabled(v bool)`

SetIsContactAddressEnabled sets IsContactAddressEnabled field to given value.

### HasIsContactAddressEnabled

`func (o *CreateAffiliateRequest) HasIsContactAddressEnabled() bool`

HasIsContactAddressEnabled returns a boolean if a field has been set.

### GetNetworkAffiliateTierId

`func (o *CreateAffiliateRequest) GetNetworkAffiliateTierId() int32`

GetNetworkAffiliateTierId returns the NetworkAffiliateTierId field if non-nil, zero value otherwise.

### GetNetworkAffiliateTierIdOk

`func (o *CreateAffiliateRequest) GetNetworkAffiliateTierIdOk() (*int32, bool)`

GetNetworkAffiliateTierIdOk returns a tuple with the NetworkAffiliateTierId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkAffiliateTierId

`func (o *CreateAffiliateRequest) SetNetworkAffiliateTierId(v int32)`

SetNetworkAffiliateTierId sets NetworkAffiliateTierId field to given value.

### HasNetworkAffiliateTierId

`func (o *CreateAffiliateRequest) HasNetworkAffiliateTierId() bool`

HasNetworkAffiliateTierId returns a boolean if a field has been set.

### GetContactAddress

`func (o *CreateAffiliateRequest) GetContactAddress() ContactAddress`

GetContactAddress returns the ContactAddress field if non-nil, zero value otherwise.

### GetContactAddressOk

`func (o *CreateAffiliateRequest) GetContactAddressOk() (*ContactAddress, bool)`

GetContactAddressOk returns a tuple with the ContactAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContactAddress

`func (o *CreateAffiliateRequest) SetContactAddress(v ContactAddress)`

SetContactAddress sets ContactAddress field to given value.

### HasContactAddress

`func (o *CreateAffiliateRequest) HasContactAddress() bool`

HasContactAddress returns a boolean if a field has been set.

### GetLabels

`func (o *CreateAffiliateRequest) GetLabels() []string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *CreateAffiliateRequest) GetLabelsOk() (*[]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *CreateAffiliateRequest) SetLabels(v []string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *CreateAffiliateRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetUsers

`func (o *CreateAffiliateRequest) GetUsers() []AffiliateUser`

GetUsers returns the Users field if non-nil, zero value otherwise.

### GetUsersOk

`func (o *CreateAffiliateRequest) GetUsersOk() (*[]AffiliateUser, bool)`

GetUsersOk returns a tuple with the Users field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsers

`func (o *CreateAffiliateRequest) SetUsers(v []AffiliateUser)`

SetUsers sets Users field to given value.

### HasUsers

`func (o *CreateAffiliateRequest) HasUsers() bool`

HasUsers returns a boolean if a field has been set.

### GetBilling

`func (o *CreateAffiliateRequest) GetBilling() BillingInfo`

GetBilling returns the Billing field if non-nil, zero value otherwise.

### GetBillingOk

`func (o *CreateAffiliateRequest) GetBillingOk() (*BillingInfo, bool)`

GetBillingOk returns a tuple with the Billing field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBilling

`func (o *CreateAffiliateRequest) SetBilling(v BillingInfo)`

SetBilling sets Billing field to given value.

### HasBilling

`func (o *CreateAffiliateRequest) HasBilling() bool`

HasBilling returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


