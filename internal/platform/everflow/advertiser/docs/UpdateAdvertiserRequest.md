# UpdateAdvertiserRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the advertiser | 
**AccountStatus** | **string** | Status of the advertiser | 
**NetworkEmployeeId** | **int32** | The employee id of the advertiser&#39;s account manager | 
**InternalNotes** | Pointer to **string** | Internal notes for the advertiser | [optional] 
**AddressId** | Pointer to **int32** | The address id of the advertiser | [optional] 
**IsContactAddressEnabled** | Pointer to **bool** | Whether or not to include a contact address for this advertiser | [optional] [default to false]
**SalesManagerId** | Pointer to **int32** | The employee id of the advertiser&#39;s sales manager | [optional] 
**DefaultCurrencyId** | **string** | The advertiser&#39;s default currency | 
**PlatformName** | Pointer to **string** | The name of the shopping cart or attribution platform | [optional] 
**PlatformUrl** | Pointer to **string** | The URL for logging into the advertiser&#39;s platform | [optional] 
**PlatformUsername** | Pointer to **string** | The username for logging into the advertiser&#39;s platform | [optional] 
**ReportingTimezoneId** | **int32** | The timezone used in the advertiser&#39;s platform reporting | 
**AttributionMethod** | **string** | Determines how attribution works for this advertiser | 
**EmailAttributionMethod** | **string** | Determines how email attribution works for this advertiser | 
**AttributionPriority** | **string** | Determines attribution priority between click and coupon code | 
**AccountingContactEmail** | Pointer to **string** | The email address of the accounting contact | [optional] 
**VerificationToken** | Pointer to **string** | Verification token for incoming postbacks | [optional] 
**OfferIdMacro** | Pointer to **string** | The string used for the offer id macro | [optional] 
**AffiliateIdMacro** | Pointer to **string** | The string used for the affiliate id macro | [optional] 
**Labels** | Pointer to **[]string** | The list of labels associated with the advertiser | [optional] 
**Users** | Pointer to [**[]AdvertiserUser**](AdvertiserUser.md) | List of advertiser users (maximum one) | [optional] 
**ContactAddress** | Pointer to [**ContactAddress**](ContactAddress.md) |  | [optional] 
**Billing** | Pointer to [**Billing**](Billing.md) |  | [optional] 
**Settings** | Pointer to [**Settings**](Settings.md) |  | [optional] 
**NetworkAdvertiserId** | Pointer to **int32** | The ID of the advertiser (read-only in update) | [optional] 
**NetworkId** | Pointer to **int32** | The network ID (read-only in update) | [optional] 
**TimeCreated** | Pointer to **int32** | Creation timestamp (read-only in update) | [optional] 
**TimeSaved** | Pointer to **int32** | Last save timestamp (read-only in update) | [optional] 

## Methods

### NewUpdateAdvertiserRequest

`func NewUpdateAdvertiserRequest(name string, accountStatus string, networkEmployeeId int32, defaultCurrencyId string, reportingTimezoneId int32, attributionMethod string, emailAttributionMethod string, attributionPriority string, ) *UpdateAdvertiserRequest`

NewUpdateAdvertiserRequest instantiates a new UpdateAdvertiserRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateAdvertiserRequestWithDefaults

`func NewUpdateAdvertiserRequestWithDefaults() *UpdateAdvertiserRequest`

NewUpdateAdvertiserRequestWithDefaults instantiates a new UpdateAdvertiserRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *UpdateAdvertiserRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *UpdateAdvertiserRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *UpdateAdvertiserRequest) SetName(v string)`

SetName sets Name field to given value.


### GetAccountStatus

`func (o *UpdateAdvertiserRequest) GetAccountStatus() string`

GetAccountStatus returns the AccountStatus field if non-nil, zero value otherwise.

### GetAccountStatusOk

`func (o *UpdateAdvertiserRequest) GetAccountStatusOk() (*string, bool)`

GetAccountStatusOk returns a tuple with the AccountStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountStatus

`func (o *UpdateAdvertiserRequest) SetAccountStatus(v string)`

SetAccountStatus sets AccountStatus field to given value.


### GetNetworkEmployeeId

`func (o *UpdateAdvertiserRequest) GetNetworkEmployeeId() int32`

GetNetworkEmployeeId returns the NetworkEmployeeId field if non-nil, zero value otherwise.

### GetNetworkEmployeeIdOk

`func (o *UpdateAdvertiserRequest) GetNetworkEmployeeIdOk() (*int32, bool)`

GetNetworkEmployeeIdOk returns a tuple with the NetworkEmployeeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkEmployeeId

`func (o *UpdateAdvertiserRequest) SetNetworkEmployeeId(v int32)`

SetNetworkEmployeeId sets NetworkEmployeeId field to given value.


### GetInternalNotes

`func (o *UpdateAdvertiserRequest) GetInternalNotes() string`

GetInternalNotes returns the InternalNotes field if non-nil, zero value otherwise.

### GetInternalNotesOk

`func (o *UpdateAdvertiserRequest) GetInternalNotesOk() (*string, bool)`

GetInternalNotesOk returns a tuple with the InternalNotes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternalNotes

`func (o *UpdateAdvertiserRequest) SetInternalNotes(v string)`

SetInternalNotes sets InternalNotes field to given value.

### HasInternalNotes

`func (o *UpdateAdvertiserRequest) HasInternalNotes() bool`

HasInternalNotes returns a boolean if a field has been set.

### GetAddressId

`func (o *UpdateAdvertiserRequest) GetAddressId() int32`

GetAddressId returns the AddressId field if non-nil, zero value otherwise.

### GetAddressIdOk

`func (o *UpdateAdvertiserRequest) GetAddressIdOk() (*int32, bool)`

GetAddressIdOk returns a tuple with the AddressId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddressId

`func (o *UpdateAdvertiserRequest) SetAddressId(v int32)`

SetAddressId sets AddressId field to given value.

### HasAddressId

`func (o *UpdateAdvertiserRequest) HasAddressId() bool`

HasAddressId returns a boolean if a field has been set.

### GetIsContactAddressEnabled

`func (o *UpdateAdvertiserRequest) GetIsContactAddressEnabled() bool`

GetIsContactAddressEnabled returns the IsContactAddressEnabled field if non-nil, zero value otherwise.

### GetIsContactAddressEnabledOk

`func (o *UpdateAdvertiserRequest) GetIsContactAddressEnabledOk() (*bool, bool)`

GetIsContactAddressEnabledOk returns a tuple with the IsContactAddressEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsContactAddressEnabled

`func (o *UpdateAdvertiserRequest) SetIsContactAddressEnabled(v bool)`

SetIsContactAddressEnabled sets IsContactAddressEnabled field to given value.

### HasIsContactAddressEnabled

`func (o *UpdateAdvertiserRequest) HasIsContactAddressEnabled() bool`

HasIsContactAddressEnabled returns a boolean if a field has been set.

### GetSalesManagerId

`func (o *UpdateAdvertiserRequest) GetSalesManagerId() int32`

GetSalesManagerId returns the SalesManagerId field if non-nil, zero value otherwise.

### GetSalesManagerIdOk

`func (o *UpdateAdvertiserRequest) GetSalesManagerIdOk() (*int32, bool)`

GetSalesManagerIdOk returns a tuple with the SalesManagerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSalesManagerId

`func (o *UpdateAdvertiserRequest) SetSalesManagerId(v int32)`

SetSalesManagerId sets SalesManagerId field to given value.

### HasSalesManagerId

`func (o *UpdateAdvertiserRequest) HasSalesManagerId() bool`

HasSalesManagerId returns a boolean if a field has been set.

### GetDefaultCurrencyId

`func (o *UpdateAdvertiserRequest) GetDefaultCurrencyId() string`

GetDefaultCurrencyId returns the DefaultCurrencyId field if non-nil, zero value otherwise.

### GetDefaultCurrencyIdOk

`func (o *UpdateAdvertiserRequest) GetDefaultCurrencyIdOk() (*string, bool)`

GetDefaultCurrencyIdOk returns a tuple with the DefaultCurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultCurrencyId

`func (o *UpdateAdvertiserRequest) SetDefaultCurrencyId(v string)`

SetDefaultCurrencyId sets DefaultCurrencyId field to given value.


### GetPlatformName

`func (o *UpdateAdvertiserRequest) GetPlatformName() string`

GetPlatformName returns the PlatformName field if non-nil, zero value otherwise.

### GetPlatformNameOk

`func (o *UpdateAdvertiserRequest) GetPlatformNameOk() (*string, bool)`

GetPlatformNameOk returns a tuple with the PlatformName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlatformName

`func (o *UpdateAdvertiserRequest) SetPlatformName(v string)`

SetPlatformName sets PlatformName field to given value.

### HasPlatformName

`func (o *UpdateAdvertiserRequest) HasPlatformName() bool`

HasPlatformName returns a boolean if a field has been set.

### GetPlatformUrl

`func (o *UpdateAdvertiserRequest) GetPlatformUrl() string`

GetPlatformUrl returns the PlatformUrl field if non-nil, zero value otherwise.

### GetPlatformUrlOk

`func (o *UpdateAdvertiserRequest) GetPlatformUrlOk() (*string, bool)`

GetPlatformUrlOk returns a tuple with the PlatformUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlatformUrl

`func (o *UpdateAdvertiserRequest) SetPlatformUrl(v string)`

SetPlatformUrl sets PlatformUrl field to given value.

### HasPlatformUrl

`func (o *UpdateAdvertiserRequest) HasPlatformUrl() bool`

HasPlatformUrl returns a boolean if a field has been set.

### GetPlatformUsername

`func (o *UpdateAdvertiserRequest) GetPlatformUsername() string`

GetPlatformUsername returns the PlatformUsername field if non-nil, zero value otherwise.

### GetPlatformUsernameOk

`func (o *UpdateAdvertiserRequest) GetPlatformUsernameOk() (*string, bool)`

GetPlatformUsernameOk returns a tuple with the PlatformUsername field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlatformUsername

`func (o *UpdateAdvertiserRequest) SetPlatformUsername(v string)`

SetPlatformUsername sets PlatformUsername field to given value.

### HasPlatformUsername

`func (o *UpdateAdvertiserRequest) HasPlatformUsername() bool`

HasPlatformUsername returns a boolean if a field has been set.

### GetReportingTimezoneId

`func (o *UpdateAdvertiserRequest) GetReportingTimezoneId() int32`

GetReportingTimezoneId returns the ReportingTimezoneId field if non-nil, zero value otherwise.

### GetReportingTimezoneIdOk

`func (o *UpdateAdvertiserRequest) GetReportingTimezoneIdOk() (*int32, bool)`

GetReportingTimezoneIdOk returns a tuple with the ReportingTimezoneId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReportingTimezoneId

`func (o *UpdateAdvertiserRequest) SetReportingTimezoneId(v int32)`

SetReportingTimezoneId sets ReportingTimezoneId field to given value.


### GetAttributionMethod

`func (o *UpdateAdvertiserRequest) GetAttributionMethod() string`

GetAttributionMethod returns the AttributionMethod field if non-nil, zero value otherwise.

### GetAttributionMethodOk

`func (o *UpdateAdvertiserRequest) GetAttributionMethodOk() (*string, bool)`

GetAttributionMethodOk returns a tuple with the AttributionMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributionMethod

`func (o *UpdateAdvertiserRequest) SetAttributionMethod(v string)`

SetAttributionMethod sets AttributionMethod field to given value.


### GetEmailAttributionMethod

`func (o *UpdateAdvertiserRequest) GetEmailAttributionMethod() string`

GetEmailAttributionMethod returns the EmailAttributionMethod field if non-nil, zero value otherwise.

### GetEmailAttributionMethodOk

`func (o *UpdateAdvertiserRequest) GetEmailAttributionMethodOk() (*string, bool)`

GetEmailAttributionMethodOk returns a tuple with the EmailAttributionMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmailAttributionMethod

`func (o *UpdateAdvertiserRequest) SetEmailAttributionMethod(v string)`

SetEmailAttributionMethod sets EmailAttributionMethod field to given value.


### GetAttributionPriority

`func (o *UpdateAdvertiserRequest) GetAttributionPriority() string`

GetAttributionPriority returns the AttributionPriority field if non-nil, zero value otherwise.

### GetAttributionPriorityOk

`func (o *UpdateAdvertiserRequest) GetAttributionPriorityOk() (*string, bool)`

GetAttributionPriorityOk returns a tuple with the AttributionPriority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributionPriority

`func (o *UpdateAdvertiserRequest) SetAttributionPriority(v string)`

SetAttributionPriority sets AttributionPriority field to given value.


### GetAccountingContactEmail

`func (o *UpdateAdvertiserRequest) GetAccountingContactEmail() string`

GetAccountingContactEmail returns the AccountingContactEmail field if non-nil, zero value otherwise.

### GetAccountingContactEmailOk

`func (o *UpdateAdvertiserRequest) GetAccountingContactEmailOk() (*string, bool)`

GetAccountingContactEmailOk returns a tuple with the AccountingContactEmail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountingContactEmail

`func (o *UpdateAdvertiserRequest) SetAccountingContactEmail(v string)`

SetAccountingContactEmail sets AccountingContactEmail field to given value.

### HasAccountingContactEmail

`func (o *UpdateAdvertiserRequest) HasAccountingContactEmail() bool`

HasAccountingContactEmail returns a boolean if a field has been set.

### GetVerificationToken

`func (o *UpdateAdvertiserRequest) GetVerificationToken() string`

GetVerificationToken returns the VerificationToken field if non-nil, zero value otherwise.

### GetVerificationTokenOk

`func (o *UpdateAdvertiserRequest) GetVerificationTokenOk() (*string, bool)`

GetVerificationTokenOk returns a tuple with the VerificationToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVerificationToken

`func (o *UpdateAdvertiserRequest) SetVerificationToken(v string)`

SetVerificationToken sets VerificationToken field to given value.

### HasVerificationToken

`func (o *UpdateAdvertiserRequest) HasVerificationToken() bool`

HasVerificationToken returns a boolean if a field has been set.

### GetOfferIdMacro

`func (o *UpdateAdvertiserRequest) GetOfferIdMacro() string`

GetOfferIdMacro returns the OfferIdMacro field if non-nil, zero value otherwise.

### GetOfferIdMacroOk

`func (o *UpdateAdvertiserRequest) GetOfferIdMacroOk() (*string, bool)`

GetOfferIdMacroOk returns a tuple with the OfferIdMacro field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOfferIdMacro

`func (o *UpdateAdvertiserRequest) SetOfferIdMacro(v string)`

SetOfferIdMacro sets OfferIdMacro field to given value.

### HasOfferIdMacro

`func (o *UpdateAdvertiserRequest) HasOfferIdMacro() bool`

HasOfferIdMacro returns a boolean if a field has been set.

### GetAffiliateIdMacro

`func (o *UpdateAdvertiserRequest) GetAffiliateIdMacro() string`

GetAffiliateIdMacro returns the AffiliateIdMacro field if non-nil, zero value otherwise.

### GetAffiliateIdMacroOk

`func (o *UpdateAdvertiserRequest) GetAffiliateIdMacroOk() (*string, bool)`

GetAffiliateIdMacroOk returns a tuple with the AffiliateIdMacro field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAffiliateIdMacro

`func (o *UpdateAdvertiserRequest) SetAffiliateIdMacro(v string)`

SetAffiliateIdMacro sets AffiliateIdMacro field to given value.

### HasAffiliateIdMacro

`func (o *UpdateAdvertiserRequest) HasAffiliateIdMacro() bool`

HasAffiliateIdMacro returns a boolean if a field has been set.

### GetLabels

`func (o *UpdateAdvertiserRequest) GetLabels() []string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *UpdateAdvertiserRequest) GetLabelsOk() (*[]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *UpdateAdvertiserRequest) SetLabels(v []string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *UpdateAdvertiserRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetUsers

`func (o *UpdateAdvertiserRequest) GetUsers() []AdvertiserUser`

GetUsers returns the Users field if non-nil, zero value otherwise.

### GetUsersOk

`func (o *UpdateAdvertiserRequest) GetUsersOk() (*[]AdvertiserUser, bool)`

GetUsersOk returns a tuple with the Users field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsers

`func (o *UpdateAdvertiserRequest) SetUsers(v []AdvertiserUser)`

SetUsers sets Users field to given value.

### HasUsers

`func (o *UpdateAdvertiserRequest) HasUsers() bool`

HasUsers returns a boolean if a field has been set.

### GetContactAddress

`func (o *UpdateAdvertiserRequest) GetContactAddress() ContactAddress`

GetContactAddress returns the ContactAddress field if non-nil, zero value otherwise.

### GetContactAddressOk

`func (o *UpdateAdvertiserRequest) GetContactAddressOk() (*ContactAddress, bool)`

GetContactAddressOk returns a tuple with the ContactAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContactAddress

`func (o *UpdateAdvertiserRequest) SetContactAddress(v ContactAddress)`

SetContactAddress sets ContactAddress field to given value.

### HasContactAddress

`func (o *UpdateAdvertiserRequest) HasContactAddress() bool`

HasContactAddress returns a boolean if a field has been set.

### GetBilling

`func (o *UpdateAdvertiserRequest) GetBilling() Billing`

GetBilling returns the Billing field if non-nil, zero value otherwise.

### GetBillingOk

`func (o *UpdateAdvertiserRequest) GetBillingOk() (*Billing, bool)`

GetBillingOk returns a tuple with the Billing field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBilling

`func (o *UpdateAdvertiserRequest) SetBilling(v Billing)`

SetBilling sets Billing field to given value.

### HasBilling

`func (o *UpdateAdvertiserRequest) HasBilling() bool`

HasBilling returns a boolean if a field has been set.

### GetSettings

`func (o *UpdateAdvertiserRequest) GetSettings() Settings`

GetSettings returns the Settings field if non-nil, zero value otherwise.

### GetSettingsOk

`func (o *UpdateAdvertiserRequest) GetSettingsOk() (*Settings, bool)`

GetSettingsOk returns a tuple with the Settings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSettings

`func (o *UpdateAdvertiserRequest) SetSettings(v Settings)`

SetSettings sets Settings field to given value.

### HasSettings

`func (o *UpdateAdvertiserRequest) HasSettings() bool`

HasSettings returns a boolean if a field has been set.

### GetNetworkAdvertiserId

`func (o *UpdateAdvertiserRequest) GetNetworkAdvertiserId() int32`

GetNetworkAdvertiserId returns the NetworkAdvertiserId field if non-nil, zero value otherwise.

### GetNetworkAdvertiserIdOk

`func (o *UpdateAdvertiserRequest) GetNetworkAdvertiserIdOk() (*int32, bool)`

GetNetworkAdvertiserIdOk returns a tuple with the NetworkAdvertiserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkAdvertiserId

`func (o *UpdateAdvertiserRequest) SetNetworkAdvertiserId(v int32)`

SetNetworkAdvertiserId sets NetworkAdvertiserId field to given value.

### HasNetworkAdvertiserId

`func (o *UpdateAdvertiserRequest) HasNetworkAdvertiserId() bool`

HasNetworkAdvertiserId returns a boolean if a field has been set.

### GetNetworkId

`func (o *UpdateAdvertiserRequest) GetNetworkId() int32`

GetNetworkId returns the NetworkId field if non-nil, zero value otherwise.

### GetNetworkIdOk

`func (o *UpdateAdvertiserRequest) GetNetworkIdOk() (*int32, bool)`

GetNetworkIdOk returns a tuple with the NetworkId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkId

`func (o *UpdateAdvertiserRequest) SetNetworkId(v int32)`

SetNetworkId sets NetworkId field to given value.

### HasNetworkId

`func (o *UpdateAdvertiserRequest) HasNetworkId() bool`

HasNetworkId returns a boolean if a field has been set.

### GetTimeCreated

`func (o *UpdateAdvertiserRequest) GetTimeCreated() int32`

GetTimeCreated returns the TimeCreated field if non-nil, zero value otherwise.

### GetTimeCreatedOk

`func (o *UpdateAdvertiserRequest) GetTimeCreatedOk() (*int32, bool)`

GetTimeCreatedOk returns a tuple with the TimeCreated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeCreated

`func (o *UpdateAdvertiserRequest) SetTimeCreated(v int32)`

SetTimeCreated sets TimeCreated field to given value.

### HasTimeCreated

`func (o *UpdateAdvertiserRequest) HasTimeCreated() bool`

HasTimeCreated returns a boolean if a field has been set.

### GetTimeSaved

`func (o *UpdateAdvertiserRequest) GetTimeSaved() int32`

GetTimeSaved returns the TimeSaved field if non-nil, zero value otherwise.

### GetTimeSavedOk

`func (o *UpdateAdvertiserRequest) GetTimeSavedOk() (*int32, bool)`

GetTimeSavedOk returns a tuple with the TimeSaved field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeSaved

`func (o *UpdateAdvertiserRequest) SetTimeSaved(v int32)`

SetTimeSaved sets TimeSaved field to given value.

### HasTimeSaved

`func (o *UpdateAdvertiserRequest) HasTimeSaved() bool`

HasTimeSaved returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


