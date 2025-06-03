# Creative

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name of the creative | 
**CreativeType** | **string** | Type of creative | 
**IsPrivate** | Pointer to **bool** | Whether creative is private | [optional] [default to false]
**CreativeStatus** | **string** | Status of creative | 
**HtmlCode** | Pointer to **string** | HTML content (required for html/email types) | [optional] 
**Width** | Pointer to **int32** | Width (required for html type) | [optional] 
**Height** | Pointer to **int32** | Height (required for html type) | [optional] 
**EmailFrom** | Pointer to **string** | From field (required for email type) | [optional] 
**EmailSubject** | Pointer to **string** | Subject field (required for email type) | [optional] 
**ResourceFile** | Pointer to [**ResourceFile**](ResourceFile.md) |  | [optional] 
**HtmlFiles** | Pointer to [**[]HtmlFile**](HtmlFile.md) |  | [optional] 

## Methods

### NewCreative

`func NewCreative(name string, creativeType string, creativeStatus string, ) *Creative`

NewCreative instantiates a new Creative object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreativeWithDefaults

`func NewCreativeWithDefaults() *Creative`

NewCreativeWithDefaults instantiates a new Creative object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *Creative) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Creative) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Creative) SetName(v string)`

SetName sets Name field to given value.


### GetCreativeType

`func (o *Creative) GetCreativeType() string`

GetCreativeType returns the CreativeType field if non-nil, zero value otherwise.

### GetCreativeTypeOk

`func (o *Creative) GetCreativeTypeOk() (*string, bool)`

GetCreativeTypeOk returns a tuple with the CreativeType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreativeType

`func (o *Creative) SetCreativeType(v string)`

SetCreativeType sets CreativeType field to given value.


### GetIsPrivate

`func (o *Creative) GetIsPrivate() bool`

GetIsPrivate returns the IsPrivate field if non-nil, zero value otherwise.

### GetIsPrivateOk

`func (o *Creative) GetIsPrivateOk() (*bool, bool)`

GetIsPrivateOk returns a tuple with the IsPrivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPrivate

`func (o *Creative) SetIsPrivate(v bool)`

SetIsPrivate sets IsPrivate field to given value.

### HasIsPrivate

`func (o *Creative) HasIsPrivate() bool`

HasIsPrivate returns a boolean if a field has been set.

### GetCreativeStatus

`func (o *Creative) GetCreativeStatus() string`

GetCreativeStatus returns the CreativeStatus field if non-nil, zero value otherwise.

### GetCreativeStatusOk

`func (o *Creative) GetCreativeStatusOk() (*string, bool)`

GetCreativeStatusOk returns a tuple with the CreativeStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreativeStatus

`func (o *Creative) SetCreativeStatus(v string)`

SetCreativeStatus sets CreativeStatus field to given value.


### GetHtmlCode

`func (o *Creative) GetHtmlCode() string`

GetHtmlCode returns the HtmlCode field if non-nil, zero value otherwise.

### GetHtmlCodeOk

`func (o *Creative) GetHtmlCodeOk() (*string, bool)`

GetHtmlCodeOk returns a tuple with the HtmlCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHtmlCode

`func (o *Creative) SetHtmlCode(v string)`

SetHtmlCode sets HtmlCode field to given value.

### HasHtmlCode

`func (o *Creative) HasHtmlCode() bool`

HasHtmlCode returns a boolean if a field has been set.

### GetWidth

`func (o *Creative) GetWidth() int32`

GetWidth returns the Width field if non-nil, zero value otherwise.

### GetWidthOk

`func (o *Creative) GetWidthOk() (*int32, bool)`

GetWidthOk returns a tuple with the Width field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWidth

`func (o *Creative) SetWidth(v int32)`

SetWidth sets Width field to given value.

### HasWidth

`func (o *Creative) HasWidth() bool`

HasWidth returns a boolean if a field has been set.

### GetHeight

`func (o *Creative) GetHeight() int32`

GetHeight returns the Height field if non-nil, zero value otherwise.

### GetHeightOk

`func (o *Creative) GetHeightOk() (*int32, bool)`

GetHeightOk returns a tuple with the Height field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHeight

`func (o *Creative) SetHeight(v int32)`

SetHeight sets Height field to given value.

### HasHeight

`func (o *Creative) HasHeight() bool`

HasHeight returns a boolean if a field has been set.

### GetEmailFrom

`func (o *Creative) GetEmailFrom() string`

GetEmailFrom returns the EmailFrom field if non-nil, zero value otherwise.

### GetEmailFromOk

`func (o *Creative) GetEmailFromOk() (*string, bool)`

GetEmailFromOk returns a tuple with the EmailFrom field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmailFrom

`func (o *Creative) SetEmailFrom(v string)`

SetEmailFrom sets EmailFrom field to given value.

### HasEmailFrom

`func (o *Creative) HasEmailFrom() bool`

HasEmailFrom returns a boolean if a field has been set.

### GetEmailSubject

`func (o *Creative) GetEmailSubject() string`

GetEmailSubject returns the EmailSubject field if non-nil, zero value otherwise.

### GetEmailSubjectOk

`func (o *Creative) GetEmailSubjectOk() (*string, bool)`

GetEmailSubjectOk returns a tuple with the EmailSubject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmailSubject

`func (o *Creative) SetEmailSubject(v string)`

SetEmailSubject sets EmailSubject field to given value.

### HasEmailSubject

`func (o *Creative) HasEmailSubject() bool`

HasEmailSubject returns a boolean if a field has been set.

### GetResourceFile

`func (o *Creative) GetResourceFile() ResourceFile`

GetResourceFile returns the ResourceFile field if non-nil, zero value otherwise.

### GetResourceFileOk

`func (o *Creative) GetResourceFileOk() (*ResourceFile, bool)`

GetResourceFileOk returns a tuple with the ResourceFile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceFile

`func (o *Creative) SetResourceFile(v ResourceFile)`

SetResourceFile sets ResourceFile field to given value.

### HasResourceFile

`func (o *Creative) HasResourceFile() bool`

HasResourceFile returns a boolean if a field has been set.

### GetHtmlFiles

`func (o *Creative) GetHtmlFiles() []HtmlFile`

GetHtmlFiles returns the HtmlFiles field if non-nil, zero value otherwise.

### GetHtmlFilesOk

`func (o *Creative) GetHtmlFilesOk() (*[]HtmlFile, bool)`

GetHtmlFilesOk returns a tuple with the HtmlFiles field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHtmlFiles

`func (o *Creative) SetHtmlFiles(v []HtmlFile)`

SetHtmlFiles sets HtmlFiles field to given value.

### HasHtmlFiles

`func (o *Creative) HasHtmlFiles() bool`

HasHtmlFiles returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


