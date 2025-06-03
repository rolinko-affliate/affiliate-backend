# AdvertiserRelationship

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Labels** | Pointer to [**AdvertiserRelationshipLabels**](AdvertiserRelationshipLabels.md) |  | [optional] 
**AccountManager** | Pointer to [**Employee**](Employee.md) |  | [optional] 
**SalesManager** | Pointer to [**Employee**](Employee.md) |  | [optional] 
**Reporting** | Pointer to [**ReportingData**](ReportingData.md) |  | [optional] 
**ApiKeys** | Pointer to [**AdvertiserRelationshipApiKeys**](AdvertiserRelationshipApiKeys.md) |  | [optional] 
**ApiWhitelistIps** | Pointer to [**AdvertiserRelationshipApiKeys**](AdvertiserRelationshipApiKeys.md) |  | [optional] 
**Billing** | Pointer to [**Billing**](Billing.md) |  | [optional] 
**Settings** | Pointer to [**Settings**](Settings.md) |  | [optional] 

## Methods

### NewAdvertiserRelationship

`func NewAdvertiserRelationship() *AdvertiserRelationship`

NewAdvertiserRelationship instantiates a new AdvertiserRelationship object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdvertiserRelationshipWithDefaults

`func NewAdvertiserRelationshipWithDefaults() *AdvertiserRelationship`

NewAdvertiserRelationshipWithDefaults instantiates a new AdvertiserRelationship object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLabels

`func (o *AdvertiserRelationship) GetLabels() AdvertiserRelationshipLabels`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *AdvertiserRelationship) GetLabelsOk() (*AdvertiserRelationshipLabels, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *AdvertiserRelationship) SetLabels(v AdvertiserRelationshipLabels)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *AdvertiserRelationship) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetAccountManager

`func (o *AdvertiserRelationship) GetAccountManager() Employee`

GetAccountManager returns the AccountManager field if non-nil, zero value otherwise.

### GetAccountManagerOk

`func (o *AdvertiserRelationship) GetAccountManagerOk() (*Employee, bool)`

GetAccountManagerOk returns a tuple with the AccountManager field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountManager

`func (o *AdvertiserRelationship) SetAccountManager(v Employee)`

SetAccountManager sets AccountManager field to given value.

### HasAccountManager

`func (o *AdvertiserRelationship) HasAccountManager() bool`

HasAccountManager returns a boolean if a field has been set.

### GetSalesManager

`func (o *AdvertiserRelationship) GetSalesManager() Employee`

GetSalesManager returns the SalesManager field if non-nil, zero value otherwise.

### GetSalesManagerOk

`func (o *AdvertiserRelationship) GetSalesManagerOk() (*Employee, bool)`

GetSalesManagerOk returns a tuple with the SalesManager field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSalesManager

`func (o *AdvertiserRelationship) SetSalesManager(v Employee)`

SetSalesManager sets SalesManager field to given value.

### HasSalesManager

`func (o *AdvertiserRelationship) HasSalesManager() bool`

HasSalesManager returns a boolean if a field has been set.

### GetReporting

`func (o *AdvertiserRelationship) GetReporting() ReportingData`

GetReporting returns the Reporting field if non-nil, zero value otherwise.

### GetReportingOk

`func (o *AdvertiserRelationship) GetReportingOk() (*ReportingData, bool)`

GetReportingOk returns a tuple with the Reporting field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReporting

`func (o *AdvertiserRelationship) SetReporting(v ReportingData)`

SetReporting sets Reporting field to given value.

### HasReporting

`func (o *AdvertiserRelationship) HasReporting() bool`

HasReporting returns a boolean if a field has been set.

### GetApiKeys

`func (o *AdvertiserRelationship) GetApiKeys() AdvertiserRelationshipApiKeys`

GetApiKeys returns the ApiKeys field if non-nil, zero value otherwise.

### GetApiKeysOk

`func (o *AdvertiserRelationship) GetApiKeysOk() (*AdvertiserRelationshipApiKeys, bool)`

GetApiKeysOk returns a tuple with the ApiKeys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiKeys

`func (o *AdvertiserRelationship) SetApiKeys(v AdvertiserRelationshipApiKeys)`

SetApiKeys sets ApiKeys field to given value.

### HasApiKeys

`func (o *AdvertiserRelationship) HasApiKeys() bool`

HasApiKeys returns a boolean if a field has been set.

### GetApiWhitelistIps

`func (o *AdvertiserRelationship) GetApiWhitelistIps() AdvertiserRelationshipApiKeys`

GetApiWhitelistIps returns the ApiWhitelistIps field if non-nil, zero value otherwise.

### GetApiWhitelistIpsOk

`func (o *AdvertiserRelationship) GetApiWhitelistIpsOk() (*AdvertiserRelationshipApiKeys, bool)`

GetApiWhitelistIpsOk returns a tuple with the ApiWhitelistIps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiWhitelistIps

`func (o *AdvertiserRelationship) SetApiWhitelistIps(v AdvertiserRelationshipApiKeys)`

SetApiWhitelistIps sets ApiWhitelistIps field to given value.

### HasApiWhitelistIps

`func (o *AdvertiserRelationship) HasApiWhitelistIps() bool`

HasApiWhitelistIps returns a boolean if a field has been set.

### GetBilling

`func (o *AdvertiserRelationship) GetBilling() Billing`

GetBilling returns the Billing field if non-nil, zero value otherwise.

### GetBillingOk

`func (o *AdvertiserRelationship) GetBillingOk() (*Billing, bool)`

GetBillingOk returns a tuple with the Billing field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBilling

`func (o *AdvertiserRelationship) SetBilling(v Billing)`

SetBilling sets Billing field to given value.

### HasBilling

`func (o *AdvertiserRelationship) HasBilling() bool`

HasBilling returns a boolean if a field has been set.

### GetSettings

`func (o *AdvertiserRelationship) GetSettings() Settings`

GetSettings returns the Settings field if non-nil, zero value otherwise.

### GetSettingsOk

`func (o *AdvertiserRelationship) GetSettingsOk() (*Settings, bool)`

GetSettingsOk returns a tuple with the Settings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSettings

`func (o *AdvertiserRelationship) SetSettings(v Settings)`

SetSettings sets Settings field to given value.

### HasSettings

`func (o *AdvertiserRelationship) HasSettings() bool`

HasSettings returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


