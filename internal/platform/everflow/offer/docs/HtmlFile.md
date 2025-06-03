# HtmlFile

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TempUrl** | **string** | Temporary URL from file upload | 
**OriginalFileName** | **string** | Filename used for macro generation | 

## Methods

### NewHtmlFile

`func NewHtmlFile(tempUrl string, originalFileName string, ) *HtmlFile`

NewHtmlFile instantiates a new HtmlFile object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHtmlFileWithDefaults

`func NewHtmlFileWithDefaults() *HtmlFile`

NewHtmlFileWithDefaults instantiates a new HtmlFile object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTempUrl

`func (o *HtmlFile) GetTempUrl() string`

GetTempUrl returns the TempUrl field if non-nil, zero value otherwise.

### GetTempUrlOk

`func (o *HtmlFile) GetTempUrlOk() (*string, bool)`

GetTempUrlOk returns a tuple with the TempUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTempUrl

`func (o *HtmlFile) SetTempUrl(v string)`

SetTempUrl sets TempUrl field to given value.


### GetOriginalFileName

`func (o *HtmlFile) GetOriginalFileName() string`

GetOriginalFileName returns the OriginalFileName field if non-nil, zero value otherwise.

### GetOriginalFileNameOk

`func (o *HtmlFile) GetOriginalFileNameOk() (*string, bool)`

GetOriginalFileNameOk returns a tuple with the OriginalFileName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOriginalFileName

`func (o *HtmlFile) SetOriginalFileName(v string)`

SetOriginalFileName sets OriginalFileName field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


