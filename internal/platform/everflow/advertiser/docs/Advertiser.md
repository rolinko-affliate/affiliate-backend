# Advertiser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NetworkAdvertiserId** | Pointer to **int32** | The unique ID of the advertiser | [optional] 
**NetworkId** | Pointer to **int32** | The network ID | [optional] 
**Name** | Pointer to **string** | The name of the advertiser | [optional] 
**AccountStatus** | Pointer to **string** | Status of the advertiser | [optional] 
**NetworkEmployeeId** | Pointer to **int32** | The employee id of the advertiser&#39;s account manager | [optional] 
**InternalNotes** | Pointer to **string** | Internal notes for the advertiser | [optional] 
**AddressId** | Pointer to **int32** | The address id of the advertiser | [optional] 
**IsContactAddressEnabled** | Pointer to **bool** | Whether contact address is enabled | [optional] 
**SalesManagerId** | Pointer to **int32** | The employee id of the advertiser&#39;s sales manager | [optional] 
**IsExposePublisherReportingData** | Pointer to **NullableBool** | Whether to expose publisher reporting data | [optional] 
**DefaultCurrencyId** | Pointer to **string** | The advertiser&#39;s default currency | [optional] 
**PlatformName** | Pointer to **string** | The name of the shopping cart or attribution platform | [optional] 
**PlatformUrl** | Pointer to **string** | The URL for logging into the advertiser&#39;s platform | [optional] 
**PlatformUsername** | Pointer to **string** | The username for logging into the advertiser&#39;s platform | [optional] 
**ReportingTimezoneId** | Pointer to **int32** | The timezone used in the advertiser&#39;s platform reporting | [optional] 
**AccountingContactEmail** | Pointer to **string** | The email address of the accounting contact | [optional] 
**VerificationToken** | Pointer to **string** | Verification token for incoming postbacks | [optional] 
**OfferIdMacro** | Pointer to **string** | The string used for the offer id macro | [optional] 
**AffiliateIdMacro** | Pointer to **string** | The string used for the affiliate id macro | [optional] 
**AttributionMethod** | Pointer to **string** | How attribution works for this advertiser | [optional] 
**EmailAttributionMethod** | Pointer to **string** | How email attribution works for this advertiser | [optional] 
**AttributionPriority** | Pointer to **string** | Attribution priority between click and coupon code | [optional] 
**TimeCreated** | Pointer to **int32** | Creation timestamp | [optional] 
**TimeSaved** | Pointer to **int32** | Last save timestamp | [optional] 
**Relationship** | Pointer to [**AdvertiserRelationship**](AdvertiserRelationship.md) |  | [optional] 

## Methods

### NewAdvertiser

`func NewAdvertiser() *Advertiser`

NewAdvertiser instantiates a new Advertiser object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdvertiserWithDefaults

`func NewAdvertiserWithDefaults() *Advertiser`

NewAdvertiserWithDefaults instantiates a new Advertiser object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNetworkAdvertiserId

`func (o *Advertiser) GetNetworkAdvertiserId() int32`

GetNetworkAdvertiserId returns the NetworkAdvertiserId field if non-nil, zero value otherwise.

### GetNetworkAdvertiserIdOk

`func (o *Advertiser) GetNetworkAdvertiserIdOk() (*int32, bool)`

GetNetworkAdvertiserIdOk returns a tuple with the NetworkAdvertiserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkAdvertiserId

`func (o *Advertiser) SetNetworkAdvertiserId(v int32)`

SetNetworkAdvertiserId sets NetworkAdvertiserId field to given value.

### HasNetworkAdvertiserId

`func (o *Advertiser) HasNetworkAdvertiserId() bool`

HasNetworkAdvertiserId returns a boolean if a field has been set.

### GetNetworkId

`func (o *Advertiser) GetNetworkId() int32`

GetNetworkId returns the NetworkId field if non-nil, zero value otherwise.

### GetNetworkIdOk

`func (o *Advertiser) GetNetworkIdOk() (*int32, bool)`

GetNetworkIdOk returns a tuple with the NetworkId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkId

`func (o *Advertiser) SetNetworkId(v int32)`

SetNetworkId sets NetworkId field to given value.

### HasNetworkId

`func (o *Advertiser) HasNetworkId() bool`

HasNetworkId returns a boolean if a field has been set.

### GetName

`func (o *Advertiser) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Advertiser) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Advertiser) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *Advertiser) HasName() bool`

HasName returns a boolean if a field has been set.

### GetAccountStatus

`func (o *Advertiser) GetAccountStatus() string`

GetAccountStatus returns the AccountStatus field if non-nil, zero value otherwise.

### GetAccountStatusOk

`func (o *Advertiser) GetAccountStatusOk() (*string, bool)`

GetAccountStatusOk returns a tuple with the AccountStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountStatus

`func (o *Advertiser) SetAccountStatus(v string)`

SetAccountStatus sets AccountStatus field to given value.

### HasAccountStatus

`func (o *Advertiser) HasAccountStatus() bool`

HasAccountStatus returns a boolean if a field has been set.

### GetNetworkEmployeeId

`func (o *Advertiser) GetNetworkEmployeeId() int32`

GetNetworkEmployeeId returns the NetworkEmployeeId field if non-nil, zero value otherwise.

### GetNetworkEmployeeIdOk

`func (o *Advertiser) GetNetworkEmployeeIdOk() (*int32, bool)`

GetNetworkEmployeeIdOk returns a tuple with the NetworkEmployeeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkEmployeeId

`func (o *Advertiser) SetNetworkEmployeeId(v int32)`

SetNetworkEmployeeId sets NetworkEmployeeId field to given value.

### HasNetworkEmployeeId

`func (o *Advertiser) HasNetworkEmployeeId() bool`

HasNetworkEmployeeId returns a boolean if a field has been set.

### GetInternalNotes

`func (o *Advertiser) GetInternalNotes() string`

GetInternalNotes returns the InternalNotes field if non-nil, zero value otherwise.

### GetInternalNotesOk

`func (o *Advertiser) GetInternalNotesOk() (*string, bool)`

GetInternalNotesOk returns a tuple with the InternalNotes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternalNotes

`func (o *Advertiser) SetInternalNotes(v string)`

SetInternalNotes sets InternalNotes field to given value.

### HasInternalNotes

`func (o *Advertiser) HasInternalNotes() bool`

HasInternalNotes returns a boolean if a field has been set.

### GetAddressId

`func (o *Advertiser) GetAddressId() int32`

GetAddressId returns the AddressId field if non-nil, zero value otherwise.

### GetAddressIdOk

`func (o *Advertiser) GetAddressIdOk() (*int32, bool)`

GetAddressIdOk returns a tuple with the AddressId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddressId

`func (o *Advertiser) SetAddressId(v int32)`

SetAddressId sets AddressId field to given value.

### HasAddressId

`func (o *Advertiser) HasAddressId() bool`

HasAddressId returns a boolean if a field has been set.

### GetIsContactAddressEnabled

`func (o *Advertiser) GetIsContactAddressEnabled() bool`

GetIsContactAddressEnabled returns the IsContactAddressEnabled field if non-nil, zero value otherwise.

### GetIsContactAddressEnabledOk

`func (o *Advertiser) GetIsContactAddressEnabledOk() (*bool, bool)`

GetIsContactAddressEnabledOk returns a tuple with the IsContactAddressEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsContactAddressEnabled

`func (o *Advertiser) SetIsContactAddressEnabled(v bool)`

SetIsContactAddressEnabled sets IsContactAddressEnabled field to given value.

### HasIsContactAddressEnabled

`func (o *Advertiser) HasIsContactAddressEnabled() bool`

HasIsContactAddressEnabled returns a boolean if a field has been set.

### GetSalesManagerId

`func (o *Advertiser) GetSalesManagerId() int32`

GetSalesManagerId returns the SalesManagerId field if non-nil, zero value otherwise.

### GetSalesManagerIdOk

`func (o *Advertiser) GetSalesManagerIdOk() (*int32, bool)`

GetSalesManagerIdOk returns a tuple with the SalesManagerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSalesManagerId

`func (o *Advertiser) SetSalesManagerId(v int32)`

SetSalesManagerId sets SalesManagerId field to given value.

### HasSalesManagerId

`func (o *Advertiser) HasSalesManagerId() bool`

HasSalesManagerId returns a boolean if a field has been set.

### GetIsExposePublisherReportingData

`func (o *Advertiser) GetIsExposePublisherReportingData() bool`

GetIsExposePublisherReportingData returns the IsExposePublisherReportingData field if non-nil, zero value otherwise.

### GetIsExposePublisherReportingDataOk

`func (o *Advertiser) GetIsExposePublisherReportingDataOk() (*bool, bool)`

GetIsExposePublisherReportingDataOk returns a tuple with the IsExposePublisherReportingData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsExposePublisherReportingData

`func (o *Advertiser) SetIsExposePublisherReportingData(v bool)`

SetIsExposePublisherReportingData sets IsExposePublisherReportingData field to given value.

### HasIsExposePublisherReportingData

`func (o *Advertiser) HasIsExposePublisherReportingData() bool`

HasIsExposePublisherReportingData returns a boolean if a field has been set.

### SetIsExposePublisherReportingDataNil

`func (o *Advertiser) SetIsExposePublisherReportingDataNil(b bool)`

 SetIsExposePublisherReportingDataNil sets the value for IsExposePublisherReportingData to be an explicit nil

### UnsetIsExposePublisherReportingData
`func (o *Advertiser) UnsetIsExposePublisherReportingData()`

UnsetIsExposePublisherReportingData ensures that no value is present for IsExposePublisherReportingData, not even an explicit nil
### GetDefaultCurrencyId

`func (o *Advertiser) GetDefaultCurrencyId() string`

GetDefaultCurrencyId returns the DefaultCurrencyId field if non-nil, zero value otherwise.

### GetDefaultCurrencyIdOk

`func (o *Advertiser) GetDefaultCurrencyIdOk() (*string, bool)`

GetDefaultCurrencyIdOk returns a tuple with the DefaultCurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultCurrencyId

`func (o *Advertiser) SetDefaultCurrencyId(v string)`

SetDefaultCurrencyId sets DefaultCurrencyId field to given value.

### HasDefaultCurrencyId

`func (o *Advertiser) HasDefaultCurrencyId() bool`

HasDefaultCurrencyId returns a boolean if a field has been set.

### GetPlatformName

`func (o *Advertiser) GetPlatformName() string`

GetPlatformName returns the PlatformName field if non-nil, zero value otherwise.

### GetPlatformNameOk

`func (o *Advertiser) GetPlatformNameOk() (*string, bool)`

GetPlatformNameOk returns a tuple with the PlatformName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlatformName

`func (o *Advertiser) SetPlatformName(v string)`

SetPlatformName sets PlatformName field to given value.

### HasPlatformName

`func (o *Advertiser) HasPlatformName() bool`

HasPlatformName returns a boolean if a field has been set.

### GetPlatformUrl

`func (o *Advertiser) GetPlatformUrl() string`

GetPlatformUrl returns the PlatformUrl field if non-nil, zero value otherwise.

### GetPlatformUrlOk

`func (o *Advertiser) GetPlatformUrlOk() (*string, bool)`

GetPlatformUrlOk returns a tuple with the PlatformUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlatformUrl

`func (o *Advertiser) SetPlatformUrl(v string)`

SetPlatformUrl sets PlatformUrl field to given value.

### HasPlatformUrl

`func (o *Advertiser) HasPlatformUrl() bool`

HasPlatformUrl returns a boolean if a field has been set.

### GetPlatformUsername

`func (o *Advertiser) GetPlatformUsername() string`

GetPlatformUsername returns the PlatformUsername field if non-nil, zero value otherwise.

### GetPlatformUsernameOk

`func (o *Advertiser) GetPlatformUsernameOk() (*string, bool)`

GetPlatformUsernameOk returns a tuple with the PlatformUsername field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlatformUsername

`func (o *Advertiser) SetPlatformUsername(v string)`

SetPlatformUsername sets PlatformUsername field to given value.

### HasPlatformUsername

`func (o *Advertiser) HasPlatformUsername() bool`

HasPlatformUsername returns a boolean if a field has been set.

### GetReportingTimezoneId

`func (o *Advertiser) GetReportingTimezoneId() int32`

GetReportingTimezoneId returns the ReportingTimezoneId field if non-nil, zero value otherwise.

### GetReportingTimezoneIdOk

`func (o *Advertiser) GetReportingTimezoneIdOk() (*int32, bool)`

GetReportingTimezoneIdOk returns a tuple with the ReportingTimezoneId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReportingTimezoneId

`func (o *Advertiser) SetReportingTimezoneId(v int32)`

SetReportingTimezoneId sets ReportingTimezoneId field to given value.

### HasReportingTimezoneId

`func (o *Advertiser) HasReportingTimezoneId() bool`

HasReportingTimezoneId returns a boolean if a field has been set.

### GetAccountingContactEmail

`func (o *Advertiser) GetAccountingContactEmail() string`

GetAccountingContactEmail returns the AccountingContactEmail field if non-nil, zero value otherwise.

### GetAccountingContactEmailOk

`func (o *Advertiser) GetAccountingContactEmailOk() (*string, bool)`

GetAccountingContactEmailOk returns a tuple with the AccountingContactEmail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountingContactEmail

`func (o *Advertiser) SetAccountingContactEmail(v string)`

SetAccountingContactEmail sets AccountingContactEmail field to given value.

### HasAccountingContactEmail

`func (o *Advertiser) HasAccountingContactEmail() bool`

HasAccountingContactEmail returns a boolean if a field has been set.

### GetVerificationToken

`func (o *Advertiser) GetVerificationToken() string`

GetVerificationToken returns the VerificationToken field if non-nil, zero value otherwise.

### GetVerificationTokenOk

`func (o *Advertiser) GetVerificationTokenOk() (*string, bool)`

GetVerificationTokenOk returns a tuple with the VerificationToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVerificationToken

`func (o *Advertiser) SetVerificationToken(v string)`

SetVerificationToken sets VerificationToken field to given value.

### HasVerificationToken

`func (o *Advertiser) HasVerificationToken() bool`

HasVerificationToken returns a boolean if a field has been set.

### GetOfferIdMacro

`func (o *Advertiser) GetOfferIdMacro() string`

GetOfferIdMacro returns the OfferIdMacro field if non-nil, zero value otherwise.

### GetOfferIdMacroOk

`func (o *Advertiser) GetOfferIdMacroOk() (*string, bool)`

GetOfferIdMacroOk returns a tuple with the OfferIdMacro field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOfferIdMacro

`func (o *Advertiser) SetOfferIdMacro(v string)`

SetOfferIdMacro sets OfferIdMacro field to given value.

### HasOfferIdMacro

`func (o *Advertiser) HasOfferIdMacro() bool`

HasOfferIdMacro returns a boolean if a field has been set.

### GetAffiliateIdMacro

`func (o *Advertiser) GetAffiliateIdMacro() string`

GetAffiliateIdMacro returns the AffiliateIdMacro field if non-nil, zero value otherwise.

### GetAffiliateIdMacroOk

`func (o *Advertiser) GetAffiliateIdMacroOk() (*string, bool)`

GetAffiliateIdMacroOk returns a tuple with the AffiliateIdMacro field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAffiliateIdMacro

`func (o *Advertiser) SetAffiliateIdMacro(v string)`

SetAffiliateIdMacro sets AffiliateIdMacro field to given value.

### HasAffiliateIdMacro

`func (o *Advertiser) HasAffiliateIdMacro() bool`

HasAffiliateIdMacro returns a boolean if a field has been set.

### GetAttributionMethod

`func (o *Advertiser) GetAttributionMethod() string`

GetAttributionMethod returns the AttributionMethod field if non-nil, zero value otherwise.

### GetAttributionMethodOk

`func (o *Advertiser) GetAttributionMethodOk() (*string, bool)`

GetAttributionMethodOk returns a tuple with the AttributionMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributionMethod

`func (o *Advertiser) SetAttributionMethod(v string)`

SetAttributionMethod sets AttributionMethod field to given value.

### HasAttributionMethod

`func (o *Advertiser) HasAttributionMethod() bool`

HasAttributionMethod returns a boolean if a field has been set.

### GetEmailAttributionMethod

`func (o *Advertiser) GetEmailAttributionMethod() string`

GetEmailAttributionMethod returns the EmailAttributionMethod field if non-nil, zero value otherwise.

### GetEmailAttributionMethodOk

`func (o *Advertiser) GetEmailAttributionMethodOk() (*string, bool)`

GetEmailAttributionMethodOk returns a tuple with the EmailAttributionMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmailAttributionMethod

`func (o *Advertiser) SetEmailAttributionMethod(v string)`

SetEmailAttributionMethod sets EmailAttributionMethod field to given value.

### HasEmailAttributionMethod

`func (o *Advertiser) HasEmailAttributionMethod() bool`

HasEmailAttributionMethod returns a boolean if a field has been set.

### GetAttributionPriority

`func (o *Advertiser) GetAttributionPriority() string`

GetAttributionPriority returns the AttributionPriority field if non-nil, zero value otherwise.

### GetAttributionPriorityOk

`func (o *Advertiser) GetAttributionPriorityOk() (*string, bool)`

GetAttributionPriorityOk returns a tuple with the AttributionPriority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributionPriority

`func (o *Advertiser) SetAttributionPriority(v string)`

SetAttributionPriority sets AttributionPriority field to given value.

### HasAttributionPriority

`func (o *Advertiser) HasAttributionPriority() bool`

HasAttributionPriority returns a boolean if a field has been set.

### GetTimeCreated

`func (o *Advertiser) GetTimeCreated() int32`

GetTimeCreated returns the TimeCreated field if non-nil, zero value otherwise.

### GetTimeCreatedOk

`func (o *Advertiser) GetTimeCreatedOk() (*int32, bool)`

GetTimeCreatedOk returns a tuple with the TimeCreated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeCreated

`func (o *Advertiser) SetTimeCreated(v int32)`

SetTimeCreated sets TimeCreated field to given value.

### HasTimeCreated

`func (o *Advertiser) HasTimeCreated() bool`

HasTimeCreated returns a boolean if a field has been set.

### GetTimeSaved

`func (o *Advertiser) GetTimeSaved() int32`

GetTimeSaved returns the TimeSaved field if non-nil, zero value otherwise.

### GetTimeSavedOk

`func (o *Advertiser) GetTimeSavedOk() (*int32, bool)`

GetTimeSavedOk returns a tuple with the TimeSaved field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeSaved

`func (o *Advertiser) SetTimeSaved(v int32)`

SetTimeSaved sets TimeSaved field to given value.

### HasTimeSaved

`func (o *Advertiser) HasTimeSaved() bool`

HasTimeSaved returns a boolean if a field has been set.

### GetRelationship

`func (o *Advertiser) GetRelationship() AdvertiserRelationship`

GetRelationship returns the Relationship field if non-nil, zero value otherwise.

### GetRelationshipOk

`func (o *Advertiser) GetRelationshipOk() (*AdvertiserRelationship, bool)`

GetRelationshipOk returns a tuple with the Relationship field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationship

`func (o *Advertiser) SetRelationship(v AdvertiserRelationship)`

SetRelationship sets Relationship field to given value.

### HasRelationship

`func (o *Advertiser) HasRelationship() bool`

HasRelationship returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


