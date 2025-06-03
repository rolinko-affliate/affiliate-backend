# EmailOptoutSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IsEnabled** | Pointer to **bool** | Email optout enabled | [optional] 
**SuppressionFileLink** | Pointer to **string** | Suppression file URL | [optional] 
**UnsubLink** | Pointer to **string** | Unsubscribe URL | [optional] 

## Methods

### NewEmailOptoutSettings

`func NewEmailOptoutSettings() *EmailOptoutSettings`

NewEmailOptoutSettings instantiates a new EmailOptoutSettings object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEmailOptoutSettingsWithDefaults

`func NewEmailOptoutSettingsWithDefaults() *EmailOptoutSettings`

NewEmailOptoutSettingsWithDefaults instantiates a new EmailOptoutSettings object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIsEnabled

`func (o *EmailOptoutSettings) GetIsEnabled() bool`

GetIsEnabled returns the IsEnabled field if non-nil, zero value otherwise.

### GetIsEnabledOk

`func (o *EmailOptoutSettings) GetIsEnabledOk() (*bool, bool)`

GetIsEnabledOk returns a tuple with the IsEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsEnabled

`func (o *EmailOptoutSettings) SetIsEnabled(v bool)`

SetIsEnabled sets IsEnabled field to given value.

### HasIsEnabled

`func (o *EmailOptoutSettings) HasIsEnabled() bool`

HasIsEnabled returns a boolean if a field has been set.

### GetSuppressionFileLink

`func (o *EmailOptoutSettings) GetSuppressionFileLink() string`

GetSuppressionFileLink returns the SuppressionFileLink field if non-nil, zero value otherwise.

### GetSuppressionFileLinkOk

`func (o *EmailOptoutSettings) GetSuppressionFileLinkOk() (*string, bool)`

GetSuppressionFileLinkOk returns a tuple with the SuppressionFileLink field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuppressionFileLink

`func (o *EmailOptoutSettings) SetSuppressionFileLink(v string)`

SetSuppressionFileLink sets SuppressionFileLink field to given value.

### HasSuppressionFileLink

`func (o *EmailOptoutSettings) HasSuppressionFileLink() bool`

HasSuppressionFileLink returns a boolean if a field has been set.

### GetUnsubLink

`func (o *EmailOptoutSettings) GetUnsubLink() string`

GetUnsubLink returns the UnsubLink field if non-nil, zero value otherwise.

### GetUnsubLinkOk

`func (o *EmailOptoutSettings) GetUnsubLinkOk() (*string, bool)`

GetUnsubLinkOk returns a tuple with the UnsubLink field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnsubLink

`func (o *EmailOptoutSettings) SetUnsubLink(v string)`

SetUnsubLink sets UnsubLink field to given value.

### HasUnsubLink

`func (o *EmailOptoutSettings) HasUnsubLink() bool`

HasUnsubLink returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


