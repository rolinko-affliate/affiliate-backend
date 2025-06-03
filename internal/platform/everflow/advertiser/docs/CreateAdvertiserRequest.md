# CreateAdvertiserRequest

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

## Methods

### NewCreateAdvertiserRequest

`func NewCreateAdvertiserRequest(name string, accountStatus string, networkEmployeeId int32, defaultCurrencyId string, reportingTimezoneId int32, attributionMethod string, emailAttributionMethod string, attributionPriority string, ) *CreateAdvertiserRequest`

NewCreateAdvertiserRequest instantiates a new CreateAdvertiserRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateAdvertiserRequestWithDefaults

`func NewCreateAdvertiserRequestWithDefaults() *CreateAdvertiserRequest`

NewCreateAdvertiserRequestWithDefaults instantiates a new CreateAdvertiserRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *CreateAdvertiserRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateAdvertiserRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateAdvertiserRequest) SetName(v string)`

SetName sets Name field to given value.


### GetAccountStatus

`func (o *CreateAdvertiserRequest) GetAccountStatus() string`

GetAccountStatus returns the AccountStatus field if non-nil, zero value otherwise.

### GetAccountStatusOk

`func (o *CreateAdvertiserRequest) GetAccountStatusOk() (*string, bool)`

GetAccountStatusOk returns a tuple with the AccountStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountStatus

`func (o *CreateAdvertiserRequest) SetAccountStatus(v string)`

SetAccountStatus sets AccountStatus field to given value.


### GetNetworkEmployeeId

`func (o *CreateAdvertiserRequest) GetNetworkEmployeeId() int32`

GetNetworkEmployeeId returns the NetworkEmployeeId field if non-nil, zero value otherwise.

### GetNetworkEmployeeIdOk

`func (o *CreateAdvertiserRequest) GetNetworkEmployeeIdOk() (*int32, bool)`

GetNetworkEmployeeIdOk returns a tuple with the NetworkEmployeeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkEmployeeId

`func (o *CreateAdvertiserRequest) SetNetworkEmployeeId(v int32)`

SetNetworkEmployeeId sets NetworkEmployeeId field to given value.


### GetInternalNotes

`func (o *CreateAdvertiserRequest) GetInternalNotes() string`

GetInternalNotes returns the InternalNotes field if non-nil, zero value otherwise.

### GetInternalNotesOk

`func (o *CreateAdvertiserRequest) GetInternalNotesOk() (*string, bool)`

GetInternalNotesOk returns a tuple with the InternalNotes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternalNotes

`func (o *CreateAdvertiserRequest) SetInternalNotes(v string)`

SetInternalNotes sets InternalNotes field to given value.

### HasInternalNotes

`func (o *CreateAdvertiserRequest) HasInternalNotes() bool`

HasInternalNotes returns a boolean if a field has been set.

### GetAddressId

`func (o *CreateAdvertiserRequest) GetAddressId() int32`

GetAddressId returns the AddressId field if non-nil, zero value otherwise.

### GetAddressIdOk

`func (o *CreateAdvertiserRequest) GetAddressIdOk() (*int32, bool)`

GetAddressIdOk returns a tuple with the AddressId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddressId

`func (o *CreateAdvertiserRequest) SetAddressId(v int32)`

SetAddressId sets AddressId field to given value.

### HasAddressId

`func (o *CreateAdvertiserRequest) HasAddressId() bool`

HasAddressId returns a boolean if a field has been set.

### GetIsContactAddressEnabled

`func (o *CreateAdvertiserRequest) GetIsContactAddressEnabled() bool`

GetIsContactAddressEnabled returns the IsContactAddressEnabled field if non-nil, zero value otherwise.

### GetIsContactAddressEnabledOk

`func (o *CreateAdvertiserRequest) GetIsContactAddressEnabledOk() (*bool, bool)`

GetIsContactAddressEnabledOk returns a tuple with the IsContactAddressEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsContactAddressEnabled

`func (o *CreateAdvertiserRequest) SetIsContactAddressEnabled(v bool)`

SetIsContactAddressEnabled sets IsContactAddressEnabled field to given value.

### HasIsContactAddressEnabled

`func (o *CreateAdvertiserRequest) HasIsContactAddressEnabled() bool`

HasIsContactAddressEnabled returns a boolean if a field has been set.

### GetSalesManagerId

`func (o *CreateAdvertiserRequest) GetSalesManagerId() int32`

GetSalesManagerId returns the SalesManagerId field if non-nil, zero value otherwise.

### GetSalesManagerIdOk

`func (o *CreateAdvertiserRequest) GetSalesManagerIdOk() (*int32, bool)`

GetSalesManagerIdOk returns a tuple with the SalesManagerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSalesManagerId

`func (o *CreateAdvertiserRequest) SetSalesManagerId(v int32)`

SetSalesManagerId sets SalesManagerId field to given value.

### HasSalesManagerId

`func (o *CreateAdvertiserRequest) HasSalesManagerId() bool`

HasSalesManagerId returns a boolean if a field has been set.

### GetDefaultCurrencyId

`func (o *CreateAdvertiserRequest) GetDefaultCurrencyId() string`

GetDefaultCurrencyId returns the DefaultCurrencyId field if non-nil, zero value otherwise.

### GetDefaultCurrencyIdOk

`func (o *CreateAdvertiserRequest) GetDefaultCurrencyIdOk() (*string, bool)`

GetDefaultCurrencyIdOk returns a tuple with the DefaultCurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultCurrencyId

`func (o *CreateAdvertiserRequest) SetDefaultCurrencyId(v string)`

SetDefaultCurrencyId sets DefaultCurrencyId field to given value.


### GetPlatformName

`func (o *CreateAdvertiserRequest) GetPlatformName() string`

GetPlatformName returns the PlatformName field if non-nil, zero value otherwise.

### GetPlatformNameOk

`func (o *CreateAdvertiserRequest) GetPlatformNameOk() (*string, bool)`

GetPlatformNameOk returns a tuple with the PlatformName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlatformName

`func (o *CreateAdvertiserRequest) SetPlatformName(v string)`

SetPlatformName sets PlatformName field to given value.

### HasPlatformName

`func (o *CreateAdvertiserRequest) HasPlatformName() bool`

HasPlatformName returns a boolean if a field has been set.

### GetPlatformUrl

`func (o *CreateAdvertiserRequest) GetPlatformUrl() string`

GetPlatformUrl returns the PlatformUrl field if non-nil, zero value otherwise.

### GetPlatformUrlOk

`func (o *CreateAdvertiserRequest) GetPlatformUrlOk() (*string, bool)`

GetPlatformUrlOk returns a tuple with the PlatformUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlatformUrl

`func (o *CreateAdvertiserRequest) SetPlatformUrl(v string)`

SetPlatformUrl sets PlatformUrl field to given value.

### HasPlatformUrl

`func (o *CreateAdvertiserRequest) HasPlatformUrl() bool`

HasPlatformUrl returns a boolean if a field has been set.

### GetPlatformUsername

`func (o *CreateAdvertiserRequest) GetPlatformUsername() string`

GetPlatformUsername returns the PlatformUsername field if non-nil, zero value otherwise.

### GetPlatformUsernameOk

`func (o *CreateAdvertiserRequest) GetPlatformUsernameOk() (*string, bool)`

GetPlatformUsernameOk returns a tuple with the PlatformUsername field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlatformUsername

`func (o *CreateAdvertiserRequest) SetPlatformUsername(v string)`

SetPlatformUsername sets PlatformUsername field to given value.

### HasPlatformUsername

`func (o *CreateAdvertiserRequest) HasPlatformUsername() bool`

HasPlatformUsername returns a boolean if a field has been set.

### GetReportingTimezoneId

`func (o *CreateAdvertiserRequest) GetReportingTimezoneId() int32`

GetReportingTimezoneId returns the ReportingTimezoneId field if non-nil, zero value otherwise.

### GetReportingTimezoneIdOk

`func (o *CreateAdvertiserRequest) GetReportingTimezoneIdOk() (*int32, bool)`

GetReportingTimezoneIdOk returns a tuple with the ReportingTimezoneId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReportingTimezoneId

`func (o *CreateAdvertiserRequest) SetReportingTimezoneId(v int32)`

SetReportingTimezoneId sets ReportingTimezoneId field to given value.


### GetAttributionMethod

`func (o *CreateAdvertiserRequest) GetAttributionMethod() string`

GetAttributionMethod returns the AttributionMethod field if non-nil, zero value otherwise.

### GetAttributionMethodOk

`func (o *CreateAdvertiserRequest) GetAttributionMethodOk() (*string, bool)`

GetAttributionMethodOk returns a tuple with the AttributionMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributionMethod

`func (o *CreateAdvertiserRequest) SetAttributionMethod(v string)`

SetAttributionMethod sets AttributionMethod field to given value.


### GetEmailAttributionMethod

`func (o *CreateAdvertiserRequest) GetEmailAttributionMethod() string`

GetEmailAttributionMethod returns the EmailAttributionMethod field if non-nil, zero value otherwise.

### GetEmailAttributionMethodOk

`func (o *CreateAdvertiserRequest) GetEmailAttributionMethodOk() (*string, bool)`

GetEmailAttributionMethodOk returns a tuple with the EmailAttributionMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmailAttributionMethod

`func (o *CreateAdvertiserRequest) SetEmailAttributionMethod(v string)`

SetEmailAttributionMethod sets EmailAttributionMethod field to given value.


### GetAttributionPriority

`func (o *CreateAdvertiserRequest) GetAttributionPriority() string`

GetAttributionPriority returns the AttributionPriority field if non-nil, zero value otherwise.

### GetAttributionPriorityOk

`func (o *CreateAdvertiserRequest) GetAttributionPriorityOk() (*string, bool)`

GetAttributionPriorityOk returns a tuple with the AttributionPriority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributionPriority

`func (o *CreateAdvertiserRequest) SetAttributionPriority(v string)`

SetAttributionPriority sets AttributionPriority field to given value.


### GetAccountingContactEmail

`func (o *CreateAdvertiserRequest) GetAccountingContactEmail() string`

GetAccountingContactEmail returns the AccountingContactEmail field if non-nil, zero value otherwise.

### GetAccountingContactEmailOk

`func (o *CreateAdvertiserRequest) GetAccountingContactEmailOk() (*string, bool)`

GetAccountingContactEmailOk returns a tuple with the AccountingContactEmail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountingContactEmail

`func (o *CreateAdvertiserRequest) SetAccountingContactEmail(v string)`

SetAccountingContactEmail sets AccountingContactEmail field to given value.

### HasAccountingContactEmail

`func (o *CreateAdvertiserRequest) HasAccountingContactEmail() bool`

HasAccountingContactEmail returns a boolean if a field has been set.

### GetVerificationToken

`func (o *CreateAdvertiserRequest) GetVerificationToken() string`

GetVerificationToken returns the VerificationToken field if non-nil, zero value otherwise.

### GetVerificationTokenOk

`func (o *CreateAdvertiserRequest) GetVerificationTokenOk() (*string, bool)`

GetVerificationTokenOk returns a tuple with the VerificationToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVerificationToken

`func (o *CreateAdvertiserRequest) SetVerificationToken(v string)`

SetVerificationToken sets VerificationToken field to given value.

### HasVerificationToken

`func (o *CreateAdvertiserRequest) HasVerificationToken() bool`

HasVerificationToken returns a boolean if a field has been set.

### GetOfferIdMacro

`func (o *CreateAdvertiserRequest) GetOfferIdMacro() string`

GetOfferIdMacro returns the OfferIdMacro field if non-nil, zero value otherwise.

### GetOfferIdMacroOk

`func (o *CreateAdvertiserRequest) GetOfferIdMacroOk() (*string, bool)`

GetOfferIdMacroOk returns a tuple with the OfferIdMacro field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOfferIdMacro

`func (o *CreateAdvertiserRequest) SetOfferIdMacro(v string)`

SetOfferIdMacro sets OfferIdMacro field to given value.

### HasOfferIdMacro

`func (o *CreateAdvertiserRequest) HasOfferIdMacro() bool`

HasOfferIdMacro returns a boolean if a field has been set.

### GetAffiliateIdMacro

`func (o *CreateAdvertiserRequest) GetAffiliateIdMacro() string`

GetAffiliateIdMacro returns the AffiliateIdMacro field if non-nil, zero value otherwise.

### GetAffiliateIdMacroOk

`func (o *CreateAdvertiserRequest) GetAffiliateIdMacroOk() (*string, bool)`

GetAffiliateIdMacroOk returns a tuple with the AffiliateIdMacro field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAffiliateIdMacro

`func (o *CreateAdvertiserRequest) SetAffiliateIdMacro(v string)`

SetAffiliateIdMacro sets AffiliateIdMacro field to given value.

### HasAffiliateIdMacro

`func (o *CreateAdvertiserRequest) HasAffiliateIdMacro() bool`

HasAffiliateIdMacro returns a boolean if a field has been set.

### GetLabels

`func (o *CreateAdvertiserRequest) GetLabels() []string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *CreateAdvertiserRequest) GetLabelsOk() (*[]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *CreateAdvertiserRequest) SetLabels(v []string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *CreateAdvertiserRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetUsers

`func (o *CreateAdvertiserRequest) GetUsers() []AdvertiserUser`

GetUsers returns the Users field if non-nil, zero value otherwise.

### GetUsersOk

`func (o *CreateAdvertiserRequest) GetUsersOk() (*[]AdvertiserUser, bool)`

GetUsersOk returns a tuple with the Users field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsers

`func (o *CreateAdvertiserRequest) SetUsers(v []AdvertiserUser)`

SetUsers sets Users field to given value.

### HasUsers

`func (o *CreateAdvertiserRequest) HasUsers() bool`

HasUsers returns a boolean if a field has been set.

### GetContactAddress

`func (o *CreateAdvertiserRequest) GetContactAddress() ContactAddress`

GetContactAddress returns the ContactAddress field if non-nil, zero value otherwise.

### GetContactAddressOk

`func (o *CreateAdvertiserRequest) GetContactAddressOk() (*ContactAddress, bool)`

GetContactAddressOk returns a tuple with the ContactAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContactAddress

`func (o *CreateAdvertiserRequest) SetContactAddress(v ContactAddress)`

SetContactAddress sets ContactAddress field to given value.

### HasContactAddress

`func (o *CreateAdvertiserRequest) HasContactAddress() bool`

HasContactAddress returns a boolean if a field has been set.

### GetBilling

`func (o *CreateAdvertiserRequest) GetBilling() Billing`

GetBilling returns the Billing field if non-nil, zero value otherwise.

### GetBillingOk

`func (o *CreateAdvertiserRequest) GetBillingOk() (*Billing, bool)`

GetBillingOk returns a tuple with the Billing field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBilling

`func (o *CreateAdvertiserRequest) SetBilling(v Billing)`

SetBilling sets Billing field to given value.

### HasBilling

`func (o *CreateAdvertiserRequest) HasBilling() bool`

HasBilling returns a boolean if a field has been set.

### GetSettings

`func (o *CreateAdvertiserRequest) GetSettings() Settings`

GetSettings returns the Settings field if non-nil, zero value otherwise.

### GetSettingsOk

`func (o *CreateAdvertiserRequest) GetSettingsOk() (*Settings, bool)`

GetSettingsOk returns a tuple with the Settings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSettings

`func (o *CreateAdvertiserRequest) SetSettings(v Settings)`

SetSettings sets Settings field to given value.

### HasSettings

`func (o *CreateAdvertiserRequest) HasSettings() bool`

HasSettings returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


