# AffiliateWithRelationships

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NetworkAffiliateId** | Pointer to **int32** | The unique ID of the affiliate | [optional] 
**NetworkId** | Pointer to **int32** | The network ID | [optional] 
**Name** | Pointer to **string** | The name of the affiliate | [optional] 
**AccountStatus** | Pointer to **string** | The affiliate&#39;s account status | [optional] 
**NetworkEmployeeId** | Pointer to **int32** | The employee id of the account manager | [optional] 
**AccountManagerId** | Pointer to **int32** | The account manager ID | [optional] 
**AccountManagerName** | Pointer to **string** | The account manager&#39;s name | [optional] 
**AccountExecutiveId** | Pointer to **int32** | Account executive ID | [optional] 
**AccountExecutiveName** | Pointer to **string** | The account executive&#39;s name | [optional] 
**InternalNotes** | Pointer to **string** | Internal notes | [optional] 
**HasNotifications** | Pointer to **bool** | Whether the affiliate has notifications enabled | [optional] 
**NetworkTrafficSourceId** | Pointer to **int32** | Traffic source ID | [optional] 
**AdressId** | Pointer to **int32** | Address ID (note the typo in the field name) | [optional] 
**DefaultCurrencyId** | Pointer to **string** | Default currency code | [optional] 
**IsContactAddressEnabled** | Pointer to **bool** | Whether contact address is enabled | [optional] 
**EnableMediaCostTrackingLinks** | Pointer to **bool** | Whether media cost tracking links are enabled | [optional] 
**TodayRevenue** | Pointer to **string** | Today&#39;s revenue (formatted as currency string) | [optional] 
**TimeCreated** | Pointer to **int64** | Unix timestamp of creation | [optional] 
**TimeSaved** | Pointer to **int64** | Unix timestamp of last save | [optional] 
**Labels** | Pointer to **[]string** | Array of labels associated with the affiliate | [optional] 
**Balance** | Pointer to **float32** | The affiliate&#39;s balance | [optional] 
**LastLogin** | Pointer to **int64** | Unix timestamp of last login | [optional] 
**GlobalTrackingDomainUrl** | Pointer to **string** | Global tracking domain URL | [optional] 
**NetworkCountryCode** | Pointer to **string** | Network country code | [optional] 
**IsPayable** | Pointer to **bool** | Whether the affiliate is payable | [optional] 
**PaymentType** | Pointer to **string** | The payment type | [optional] 
**ReferrerId** | Pointer to **int32** | ID of referring affiliate | [optional] 
**Relationship** | Pointer to [**AffiliateWithRelationshipsAllOfRelationship**](AffiliateWithRelationshipsAllOfRelationship.md) |  | [optional] 

## Methods

### NewAffiliateWithRelationships

`func NewAffiliateWithRelationships() *AffiliateWithRelationships`

NewAffiliateWithRelationships instantiates a new AffiliateWithRelationships object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAffiliateWithRelationshipsWithDefaults

`func NewAffiliateWithRelationshipsWithDefaults() *AffiliateWithRelationships`

NewAffiliateWithRelationshipsWithDefaults instantiates a new AffiliateWithRelationships object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNetworkAffiliateId

`func (o *AffiliateWithRelationships) GetNetworkAffiliateId() int32`

GetNetworkAffiliateId returns the NetworkAffiliateId field if non-nil, zero value otherwise.

### GetNetworkAffiliateIdOk

`func (o *AffiliateWithRelationships) GetNetworkAffiliateIdOk() (*int32, bool)`

GetNetworkAffiliateIdOk returns a tuple with the NetworkAffiliateId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkAffiliateId

`func (o *AffiliateWithRelationships) SetNetworkAffiliateId(v int32)`

SetNetworkAffiliateId sets NetworkAffiliateId field to given value.

### HasNetworkAffiliateId

`func (o *AffiliateWithRelationships) HasNetworkAffiliateId() bool`

HasNetworkAffiliateId returns a boolean if a field has been set.

### GetNetworkId

`func (o *AffiliateWithRelationships) GetNetworkId() int32`

GetNetworkId returns the NetworkId field if non-nil, zero value otherwise.

### GetNetworkIdOk

`func (o *AffiliateWithRelationships) GetNetworkIdOk() (*int32, bool)`

GetNetworkIdOk returns a tuple with the NetworkId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkId

`func (o *AffiliateWithRelationships) SetNetworkId(v int32)`

SetNetworkId sets NetworkId field to given value.

### HasNetworkId

`func (o *AffiliateWithRelationships) HasNetworkId() bool`

HasNetworkId returns a boolean if a field has been set.

### GetName

`func (o *AffiliateWithRelationships) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AffiliateWithRelationships) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AffiliateWithRelationships) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *AffiliateWithRelationships) HasName() bool`

HasName returns a boolean if a field has been set.

### GetAccountStatus

`func (o *AffiliateWithRelationships) GetAccountStatus() string`

GetAccountStatus returns the AccountStatus field if non-nil, zero value otherwise.

### GetAccountStatusOk

`func (o *AffiliateWithRelationships) GetAccountStatusOk() (*string, bool)`

GetAccountStatusOk returns a tuple with the AccountStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountStatus

`func (o *AffiliateWithRelationships) SetAccountStatus(v string)`

SetAccountStatus sets AccountStatus field to given value.

### HasAccountStatus

`func (o *AffiliateWithRelationships) HasAccountStatus() bool`

HasAccountStatus returns a boolean if a field has been set.

### GetNetworkEmployeeId

`func (o *AffiliateWithRelationships) GetNetworkEmployeeId() int32`

GetNetworkEmployeeId returns the NetworkEmployeeId field if non-nil, zero value otherwise.

### GetNetworkEmployeeIdOk

`func (o *AffiliateWithRelationships) GetNetworkEmployeeIdOk() (*int32, bool)`

GetNetworkEmployeeIdOk returns a tuple with the NetworkEmployeeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkEmployeeId

`func (o *AffiliateWithRelationships) SetNetworkEmployeeId(v int32)`

SetNetworkEmployeeId sets NetworkEmployeeId field to given value.

### HasNetworkEmployeeId

`func (o *AffiliateWithRelationships) HasNetworkEmployeeId() bool`

HasNetworkEmployeeId returns a boolean if a field has been set.

### GetAccountManagerId

`func (o *AffiliateWithRelationships) GetAccountManagerId() int32`

GetAccountManagerId returns the AccountManagerId field if non-nil, zero value otherwise.

### GetAccountManagerIdOk

`func (o *AffiliateWithRelationships) GetAccountManagerIdOk() (*int32, bool)`

GetAccountManagerIdOk returns a tuple with the AccountManagerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountManagerId

`func (o *AffiliateWithRelationships) SetAccountManagerId(v int32)`

SetAccountManagerId sets AccountManagerId field to given value.

### HasAccountManagerId

`func (o *AffiliateWithRelationships) HasAccountManagerId() bool`

HasAccountManagerId returns a boolean if a field has been set.

### GetAccountManagerName

`func (o *AffiliateWithRelationships) GetAccountManagerName() string`

GetAccountManagerName returns the AccountManagerName field if non-nil, zero value otherwise.

### GetAccountManagerNameOk

`func (o *AffiliateWithRelationships) GetAccountManagerNameOk() (*string, bool)`

GetAccountManagerNameOk returns a tuple with the AccountManagerName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountManagerName

`func (o *AffiliateWithRelationships) SetAccountManagerName(v string)`

SetAccountManagerName sets AccountManagerName field to given value.

### HasAccountManagerName

`func (o *AffiliateWithRelationships) HasAccountManagerName() bool`

HasAccountManagerName returns a boolean if a field has been set.

### GetAccountExecutiveId

`func (o *AffiliateWithRelationships) GetAccountExecutiveId() int32`

GetAccountExecutiveId returns the AccountExecutiveId field if non-nil, zero value otherwise.

### GetAccountExecutiveIdOk

`func (o *AffiliateWithRelationships) GetAccountExecutiveIdOk() (*int32, bool)`

GetAccountExecutiveIdOk returns a tuple with the AccountExecutiveId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountExecutiveId

`func (o *AffiliateWithRelationships) SetAccountExecutiveId(v int32)`

SetAccountExecutiveId sets AccountExecutiveId field to given value.

### HasAccountExecutiveId

`func (o *AffiliateWithRelationships) HasAccountExecutiveId() bool`

HasAccountExecutiveId returns a boolean if a field has been set.

### GetAccountExecutiveName

`func (o *AffiliateWithRelationships) GetAccountExecutiveName() string`

GetAccountExecutiveName returns the AccountExecutiveName field if non-nil, zero value otherwise.

### GetAccountExecutiveNameOk

`func (o *AffiliateWithRelationships) GetAccountExecutiveNameOk() (*string, bool)`

GetAccountExecutiveNameOk returns a tuple with the AccountExecutiveName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountExecutiveName

`func (o *AffiliateWithRelationships) SetAccountExecutiveName(v string)`

SetAccountExecutiveName sets AccountExecutiveName field to given value.

### HasAccountExecutiveName

`func (o *AffiliateWithRelationships) HasAccountExecutiveName() bool`

HasAccountExecutiveName returns a boolean if a field has been set.

### GetInternalNotes

`func (o *AffiliateWithRelationships) GetInternalNotes() string`

GetInternalNotes returns the InternalNotes field if non-nil, zero value otherwise.

### GetInternalNotesOk

`func (o *AffiliateWithRelationships) GetInternalNotesOk() (*string, bool)`

GetInternalNotesOk returns a tuple with the InternalNotes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternalNotes

`func (o *AffiliateWithRelationships) SetInternalNotes(v string)`

SetInternalNotes sets InternalNotes field to given value.

### HasInternalNotes

`func (o *AffiliateWithRelationships) HasInternalNotes() bool`

HasInternalNotes returns a boolean if a field has been set.

### GetHasNotifications

`func (o *AffiliateWithRelationships) GetHasNotifications() bool`

GetHasNotifications returns the HasNotifications field if non-nil, zero value otherwise.

### GetHasNotificationsOk

`func (o *AffiliateWithRelationships) GetHasNotificationsOk() (*bool, bool)`

GetHasNotificationsOk returns a tuple with the HasNotifications field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasNotifications

`func (o *AffiliateWithRelationships) SetHasNotifications(v bool)`

SetHasNotifications sets HasNotifications field to given value.

### HasHasNotifications

`func (o *AffiliateWithRelationships) HasHasNotifications() bool`

HasHasNotifications returns a boolean if a field has been set.

### GetNetworkTrafficSourceId

`func (o *AffiliateWithRelationships) GetNetworkTrafficSourceId() int32`

GetNetworkTrafficSourceId returns the NetworkTrafficSourceId field if non-nil, zero value otherwise.

### GetNetworkTrafficSourceIdOk

`func (o *AffiliateWithRelationships) GetNetworkTrafficSourceIdOk() (*int32, bool)`

GetNetworkTrafficSourceIdOk returns a tuple with the NetworkTrafficSourceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkTrafficSourceId

`func (o *AffiliateWithRelationships) SetNetworkTrafficSourceId(v int32)`

SetNetworkTrafficSourceId sets NetworkTrafficSourceId field to given value.

### HasNetworkTrafficSourceId

`func (o *AffiliateWithRelationships) HasNetworkTrafficSourceId() bool`

HasNetworkTrafficSourceId returns a boolean if a field has been set.

### GetAdressId

`func (o *AffiliateWithRelationships) GetAdressId() int32`

GetAdressId returns the AdressId field if non-nil, zero value otherwise.

### GetAdressIdOk

`func (o *AffiliateWithRelationships) GetAdressIdOk() (*int32, bool)`

GetAdressIdOk returns a tuple with the AdressId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdressId

`func (o *AffiliateWithRelationships) SetAdressId(v int32)`

SetAdressId sets AdressId field to given value.

### HasAdressId

`func (o *AffiliateWithRelationships) HasAdressId() bool`

HasAdressId returns a boolean if a field has been set.

### GetDefaultCurrencyId

`func (o *AffiliateWithRelationships) GetDefaultCurrencyId() string`

GetDefaultCurrencyId returns the DefaultCurrencyId field if non-nil, zero value otherwise.

### GetDefaultCurrencyIdOk

`func (o *AffiliateWithRelationships) GetDefaultCurrencyIdOk() (*string, bool)`

GetDefaultCurrencyIdOk returns a tuple with the DefaultCurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultCurrencyId

`func (o *AffiliateWithRelationships) SetDefaultCurrencyId(v string)`

SetDefaultCurrencyId sets DefaultCurrencyId field to given value.

### HasDefaultCurrencyId

`func (o *AffiliateWithRelationships) HasDefaultCurrencyId() bool`

HasDefaultCurrencyId returns a boolean if a field has been set.

### GetIsContactAddressEnabled

`func (o *AffiliateWithRelationships) GetIsContactAddressEnabled() bool`

GetIsContactAddressEnabled returns the IsContactAddressEnabled field if non-nil, zero value otherwise.

### GetIsContactAddressEnabledOk

`func (o *AffiliateWithRelationships) GetIsContactAddressEnabledOk() (*bool, bool)`

GetIsContactAddressEnabledOk returns a tuple with the IsContactAddressEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsContactAddressEnabled

`func (o *AffiliateWithRelationships) SetIsContactAddressEnabled(v bool)`

SetIsContactAddressEnabled sets IsContactAddressEnabled field to given value.

### HasIsContactAddressEnabled

`func (o *AffiliateWithRelationships) HasIsContactAddressEnabled() bool`

HasIsContactAddressEnabled returns a boolean if a field has been set.

### GetEnableMediaCostTrackingLinks

`func (o *AffiliateWithRelationships) GetEnableMediaCostTrackingLinks() bool`

GetEnableMediaCostTrackingLinks returns the EnableMediaCostTrackingLinks field if non-nil, zero value otherwise.

### GetEnableMediaCostTrackingLinksOk

`func (o *AffiliateWithRelationships) GetEnableMediaCostTrackingLinksOk() (*bool, bool)`

GetEnableMediaCostTrackingLinksOk returns a tuple with the EnableMediaCostTrackingLinks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableMediaCostTrackingLinks

`func (o *AffiliateWithRelationships) SetEnableMediaCostTrackingLinks(v bool)`

SetEnableMediaCostTrackingLinks sets EnableMediaCostTrackingLinks field to given value.

### HasEnableMediaCostTrackingLinks

`func (o *AffiliateWithRelationships) HasEnableMediaCostTrackingLinks() bool`

HasEnableMediaCostTrackingLinks returns a boolean if a field has been set.

### GetTodayRevenue

`func (o *AffiliateWithRelationships) GetTodayRevenue() string`

GetTodayRevenue returns the TodayRevenue field if non-nil, zero value otherwise.

### GetTodayRevenueOk

`func (o *AffiliateWithRelationships) GetTodayRevenueOk() (*string, bool)`

GetTodayRevenueOk returns a tuple with the TodayRevenue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTodayRevenue

`func (o *AffiliateWithRelationships) SetTodayRevenue(v string)`

SetTodayRevenue sets TodayRevenue field to given value.

### HasTodayRevenue

`func (o *AffiliateWithRelationships) HasTodayRevenue() bool`

HasTodayRevenue returns a boolean if a field has been set.

### GetTimeCreated

`func (o *AffiliateWithRelationships) GetTimeCreated() int64`

GetTimeCreated returns the TimeCreated field if non-nil, zero value otherwise.

### GetTimeCreatedOk

`func (o *AffiliateWithRelationships) GetTimeCreatedOk() (*int64, bool)`

GetTimeCreatedOk returns a tuple with the TimeCreated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeCreated

`func (o *AffiliateWithRelationships) SetTimeCreated(v int64)`

SetTimeCreated sets TimeCreated field to given value.

### HasTimeCreated

`func (o *AffiliateWithRelationships) HasTimeCreated() bool`

HasTimeCreated returns a boolean if a field has been set.

### GetTimeSaved

`func (o *AffiliateWithRelationships) GetTimeSaved() int64`

GetTimeSaved returns the TimeSaved field if non-nil, zero value otherwise.

### GetTimeSavedOk

`func (o *AffiliateWithRelationships) GetTimeSavedOk() (*int64, bool)`

GetTimeSavedOk returns a tuple with the TimeSaved field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeSaved

`func (o *AffiliateWithRelationships) SetTimeSaved(v int64)`

SetTimeSaved sets TimeSaved field to given value.

### HasTimeSaved

`func (o *AffiliateWithRelationships) HasTimeSaved() bool`

HasTimeSaved returns a boolean if a field has been set.

### GetLabels

`func (o *AffiliateWithRelationships) GetLabels() []string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *AffiliateWithRelationships) GetLabelsOk() (*[]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *AffiliateWithRelationships) SetLabels(v []string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *AffiliateWithRelationships) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetBalance

`func (o *AffiliateWithRelationships) GetBalance() float32`

GetBalance returns the Balance field if non-nil, zero value otherwise.

### GetBalanceOk

`func (o *AffiliateWithRelationships) GetBalanceOk() (*float32, bool)`

GetBalanceOk returns a tuple with the Balance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBalance

`func (o *AffiliateWithRelationships) SetBalance(v float32)`

SetBalance sets Balance field to given value.

### HasBalance

`func (o *AffiliateWithRelationships) HasBalance() bool`

HasBalance returns a boolean if a field has been set.

### GetLastLogin

`func (o *AffiliateWithRelationships) GetLastLogin() int64`

GetLastLogin returns the LastLogin field if non-nil, zero value otherwise.

### GetLastLoginOk

`func (o *AffiliateWithRelationships) GetLastLoginOk() (*int64, bool)`

GetLastLoginOk returns a tuple with the LastLogin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastLogin

`func (o *AffiliateWithRelationships) SetLastLogin(v int64)`

SetLastLogin sets LastLogin field to given value.

### HasLastLogin

`func (o *AffiliateWithRelationships) HasLastLogin() bool`

HasLastLogin returns a boolean if a field has been set.

### GetGlobalTrackingDomainUrl

`func (o *AffiliateWithRelationships) GetGlobalTrackingDomainUrl() string`

GetGlobalTrackingDomainUrl returns the GlobalTrackingDomainUrl field if non-nil, zero value otherwise.

### GetGlobalTrackingDomainUrlOk

`func (o *AffiliateWithRelationships) GetGlobalTrackingDomainUrlOk() (*string, bool)`

GetGlobalTrackingDomainUrlOk returns a tuple with the GlobalTrackingDomainUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGlobalTrackingDomainUrl

`func (o *AffiliateWithRelationships) SetGlobalTrackingDomainUrl(v string)`

SetGlobalTrackingDomainUrl sets GlobalTrackingDomainUrl field to given value.

### HasGlobalTrackingDomainUrl

`func (o *AffiliateWithRelationships) HasGlobalTrackingDomainUrl() bool`

HasGlobalTrackingDomainUrl returns a boolean if a field has been set.

### GetNetworkCountryCode

`func (o *AffiliateWithRelationships) GetNetworkCountryCode() string`

GetNetworkCountryCode returns the NetworkCountryCode field if non-nil, zero value otherwise.

### GetNetworkCountryCodeOk

`func (o *AffiliateWithRelationships) GetNetworkCountryCodeOk() (*string, bool)`

GetNetworkCountryCodeOk returns a tuple with the NetworkCountryCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkCountryCode

`func (o *AffiliateWithRelationships) SetNetworkCountryCode(v string)`

SetNetworkCountryCode sets NetworkCountryCode field to given value.

### HasNetworkCountryCode

`func (o *AffiliateWithRelationships) HasNetworkCountryCode() bool`

HasNetworkCountryCode returns a boolean if a field has been set.

### GetIsPayable

`func (o *AffiliateWithRelationships) GetIsPayable() bool`

GetIsPayable returns the IsPayable field if non-nil, zero value otherwise.

### GetIsPayableOk

`func (o *AffiliateWithRelationships) GetIsPayableOk() (*bool, bool)`

GetIsPayableOk returns a tuple with the IsPayable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPayable

`func (o *AffiliateWithRelationships) SetIsPayable(v bool)`

SetIsPayable sets IsPayable field to given value.

### HasIsPayable

`func (o *AffiliateWithRelationships) HasIsPayable() bool`

HasIsPayable returns a boolean if a field has been set.

### GetPaymentType

`func (o *AffiliateWithRelationships) GetPaymentType() string`

GetPaymentType returns the PaymentType field if non-nil, zero value otherwise.

### GetPaymentTypeOk

`func (o *AffiliateWithRelationships) GetPaymentTypeOk() (*string, bool)`

GetPaymentTypeOk returns a tuple with the PaymentType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPaymentType

`func (o *AffiliateWithRelationships) SetPaymentType(v string)`

SetPaymentType sets PaymentType field to given value.

### HasPaymentType

`func (o *AffiliateWithRelationships) HasPaymentType() bool`

HasPaymentType returns a boolean if a field has been set.

### GetReferrerId

`func (o *AffiliateWithRelationships) GetReferrerId() int32`

GetReferrerId returns the ReferrerId field if non-nil, zero value otherwise.

### GetReferrerIdOk

`func (o *AffiliateWithRelationships) GetReferrerIdOk() (*int32, bool)`

GetReferrerIdOk returns a tuple with the ReferrerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReferrerId

`func (o *AffiliateWithRelationships) SetReferrerId(v int32)`

SetReferrerId sets ReferrerId field to given value.

### HasReferrerId

`func (o *AffiliateWithRelationships) HasReferrerId() bool`

HasReferrerId returns a boolean if a field has been set.

### GetRelationship

`func (o *AffiliateWithRelationships) GetRelationship() AffiliateWithRelationshipsAllOfRelationship`

GetRelationship returns the Relationship field if non-nil, zero value otherwise.

### GetRelationshipOk

`func (o *AffiliateWithRelationships) GetRelationshipOk() (*AffiliateWithRelationshipsAllOfRelationship, bool)`

GetRelationshipOk returns a tuple with the Relationship field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationship

`func (o *AffiliateWithRelationships) SetRelationship(v AffiliateWithRelationshipsAllOfRelationship)`

SetRelationship sets Relationship field to given value.

### HasRelationship

`func (o *AffiliateWithRelationships) HasRelationship() bool`

HasRelationship returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


