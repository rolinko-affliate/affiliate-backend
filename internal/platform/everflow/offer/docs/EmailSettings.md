# EmailSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IsEnabled** | Pointer to **bool** | Email settings enabled | [optional] 
**SubjectLines** | Pointer to **string** | Approved subject lines | [optional] 
**FromLines** | Pointer to **string** | Approved from lines | [optional] 

## Methods

### NewEmailSettings

`func NewEmailSettings() *EmailSettings`

NewEmailSettings instantiates a new EmailSettings object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEmailSettingsWithDefaults

`func NewEmailSettingsWithDefaults() *EmailSettings`

NewEmailSettingsWithDefaults instantiates a new EmailSettings object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIsEnabled

`func (o *EmailSettings) GetIsEnabled() bool`

GetIsEnabled returns the IsEnabled field if non-nil, zero value otherwise.

### GetIsEnabledOk

`func (o *EmailSettings) GetIsEnabledOk() (*bool, bool)`

GetIsEnabledOk returns a tuple with the IsEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsEnabled

`func (o *EmailSettings) SetIsEnabled(v bool)`

SetIsEnabled sets IsEnabled field to given value.

### HasIsEnabled

`func (o *EmailSettings) HasIsEnabled() bool`

HasIsEnabled returns a boolean if a field has been set.

### GetSubjectLines

`func (o *EmailSettings) GetSubjectLines() string`

GetSubjectLines returns the SubjectLines field if non-nil, zero value otherwise.

### GetSubjectLinesOk

`func (o *EmailSettings) GetSubjectLinesOk() (*string, bool)`

GetSubjectLinesOk returns a tuple with the SubjectLines field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubjectLines

`func (o *EmailSettings) SetSubjectLines(v string)`

SetSubjectLines sets SubjectLines field to given value.

### HasSubjectLines

`func (o *EmailSettings) HasSubjectLines() bool`

HasSubjectLines returns a boolean if a field has been set.

### GetFromLines

`func (o *EmailSettings) GetFromLines() string`

GetFromLines returns the FromLines field if non-nil, zero value otherwise.

### GetFromLinesOk

`func (o *EmailSettings) GetFromLinesOk() (*string, bool)`

GetFromLinesOk returns a tuple with the FromLines field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFromLines

`func (o *EmailSettings) SetFromLines(v string)`

SetFromLines sets FromLines field to given value.

### HasFromLines

`func (o *EmailSettings) HasFromLines() bool`

HasFromLines returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


