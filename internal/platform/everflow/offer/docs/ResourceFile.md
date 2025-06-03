# ResourceFile

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TempUrl** | **string** | Temporary URL from file upload | 
**OriginalFileName** | **string** | Original filename | 

## Methods

### NewResourceFile

`func NewResourceFile(tempUrl string, originalFileName string, ) *ResourceFile`

NewResourceFile instantiates a new ResourceFile object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewResourceFileWithDefaults

`func NewResourceFileWithDefaults() *ResourceFile`

NewResourceFileWithDefaults instantiates a new ResourceFile object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTempUrl

`func (o *ResourceFile) GetTempUrl() string`

GetTempUrl returns the TempUrl field if non-nil, zero value otherwise.

### GetTempUrlOk

`func (o *ResourceFile) GetTempUrlOk() (*string, bool)`

GetTempUrlOk returns a tuple with the TempUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTempUrl

`func (o *ResourceFile) SetTempUrl(v string)`

SetTempUrl sets TempUrl field to given value.


### GetOriginalFileName

`func (o *ResourceFile) GetOriginalFileName() string`

GetOriginalFileName returns the OriginalFileName field if non-nil, zero value otherwise.

### GetOriginalFileNameOk

`func (o *ResourceFile) GetOriginalFileNameOk() (*string, bool)`

GetOriginalFileNameOk returns a tuple with the OriginalFileName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOriginalFileName

`func (o *ResourceFile) SetOriginalFileName(v string)`

SetOriginalFileName sets OriginalFileName field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


