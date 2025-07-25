/*
Everflow Network API - Offers

API for managing offers in the Everflow platform

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package offer

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// checks if the Creative type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Creative{}

// Creative struct for Creative
type Creative struct {
	// Name of the creative
	Name string `json:"name"`
	// Type of creative
	CreativeType string `json:"creative_type"`
	// Whether creative is private
	IsPrivate *bool `json:"is_private,omitempty"`
	// Status of creative
	CreativeStatus string `json:"creative_status"`
	// HTML content (required for html/email types)
	HtmlCode *string `json:"html_code,omitempty"`
	// Width (required for html type)
	Width *int32 `json:"width,omitempty"`
	// Height (required for html type)
	Height *int32 `json:"height,omitempty"`
	// From field (required for email type)
	EmailFrom *string `json:"email_from,omitempty"`
	// Subject field (required for email type)
	EmailSubject *string       `json:"email_subject,omitempty"`
	ResourceFile *ResourceFile `json:"resource_file,omitempty"`
	HtmlFiles    []HtmlFile    `json:"html_files,omitempty"`
}

type _Creative Creative

// NewCreative instantiates a new Creative object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreative(name string, creativeType string, creativeStatus string) *Creative {
	this := Creative{}
	this.Name = name
	this.CreativeType = creativeType
	var isPrivate bool = false
	this.IsPrivate = &isPrivate
	this.CreativeStatus = creativeStatus
	return &this
}

// NewCreativeWithDefaults instantiates a new Creative object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreativeWithDefaults() *Creative {
	this := Creative{}
	var isPrivate bool = false
	this.IsPrivate = &isPrivate
	return &this
}

// GetName returns the Name field value
func (o *Creative) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *Creative) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *Creative) SetName(v string) {
	o.Name = v
}

// GetCreativeType returns the CreativeType field value
func (o *Creative) GetCreativeType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CreativeType
}

// GetCreativeTypeOk returns a tuple with the CreativeType field value
// and a boolean to check if the value has been set.
func (o *Creative) GetCreativeTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreativeType, true
}

// SetCreativeType sets field value
func (o *Creative) SetCreativeType(v string) {
	o.CreativeType = v
}

// GetIsPrivate returns the IsPrivate field value if set, zero value otherwise.
func (o *Creative) GetIsPrivate() bool {
	if o == nil || IsNil(o.IsPrivate) {
		var ret bool
		return ret
	}
	return *o.IsPrivate
}

// GetIsPrivateOk returns a tuple with the IsPrivate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Creative) GetIsPrivateOk() (*bool, bool) {
	if o == nil || IsNil(o.IsPrivate) {
		return nil, false
	}
	return o.IsPrivate, true
}

// HasIsPrivate returns a boolean if a field has been set.
func (o *Creative) HasIsPrivate() bool {
	if o != nil && !IsNil(o.IsPrivate) {
		return true
	}

	return false
}

// SetIsPrivate gets a reference to the given bool and assigns it to the IsPrivate field.
func (o *Creative) SetIsPrivate(v bool) {
	o.IsPrivate = &v
}

// GetCreativeStatus returns the CreativeStatus field value
func (o *Creative) GetCreativeStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CreativeStatus
}

// GetCreativeStatusOk returns a tuple with the CreativeStatus field value
// and a boolean to check if the value has been set.
func (o *Creative) GetCreativeStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreativeStatus, true
}

// SetCreativeStatus sets field value
func (o *Creative) SetCreativeStatus(v string) {
	o.CreativeStatus = v
}

// GetHtmlCode returns the HtmlCode field value if set, zero value otherwise.
func (o *Creative) GetHtmlCode() string {
	if o == nil || IsNil(o.HtmlCode) {
		var ret string
		return ret
	}
	return *o.HtmlCode
}

// GetHtmlCodeOk returns a tuple with the HtmlCode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Creative) GetHtmlCodeOk() (*string, bool) {
	if o == nil || IsNil(o.HtmlCode) {
		return nil, false
	}
	return o.HtmlCode, true
}

// HasHtmlCode returns a boolean if a field has been set.
func (o *Creative) HasHtmlCode() bool {
	if o != nil && !IsNil(o.HtmlCode) {
		return true
	}

	return false
}

// SetHtmlCode gets a reference to the given string and assigns it to the HtmlCode field.
func (o *Creative) SetHtmlCode(v string) {
	o.HtmlCode = &v
}

// GetWidth returns the Width field value if set, zero value otherwise.
func (o *Creative) GetWidth() int32 {
	if o == nil || IsNil(o.Width) {
		var ret int32
		return ret
	}
	return *o.Width
}

// GetWidthOk returns a tuple with the Width field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Creative) GetWidthOk() (*int32, bool) {
	if o == nil || IsNil(o.Width) {
		return nil, false
	}
	return o.Width, true
}

// HasWidth returns a boolean if a field has been set.
func (o *Creative) HasWidth() bool {
	if o != nil && !IsNil(o.Width) {
		return true
	}

	return false
}

// SetWidth gets a reference to the given int32 and assigns it to the Width field.
func (o *Creative) SetWidth(v int32) {
	o.Width = &v
}

// GetHeight returns the Height field value if set, zero value otherwise.
func (o *Creative) GetHeight() int32 {
	if o == nil || IsNil(o.Height) {
		var ret int32
		return ret
	}
	return *o.Height
}

// GetHeightOk returns a tuple with the Height field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Creative) GetHeightOk() (*int32, bool) {
	if o == nil || IsNil(o.Height) {
		return nil, false
	}
	return o.Height, true
}

// HasHeight returns a boolean if a field has been set.
func (o *Creative) HasHeight() bool {
	if o != nil && !IsNil(o.Height) {
		return true
	}

	return false
}

// SetHeight gets a reference to the given int32 and assigns it to the Height field.
func (o *Creative) SetHeight(v int32) {
	o.Height = &v
}

// GetEmailFrom returns the EmailFrom field value if set, zero value otherwise.
func (o *Creative) GetEmailFrom() string {
	if o == nil || IsNil(o.EmailFrom) {
		var ret string
		return ret
	}
	return *o.EmailFrom
}

// GetEmailFromOk returns a tuple with the EmailFrom field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Creative) GetEmailFromOk() (*string, bool) {
	if o == nil || IsNil(o.EmailFrom) {
		return nil, false
	}
	return o.EmailFrom, true
}

// HasEmailFrom returns a boolean if a field has been set.
func (o *Creative) HasEmailFrom() bool {
	if o != nil && !IsNil(o.EmailFrom) {
		return true
	}

	return false
}

// SetEmailFrom gets a reference to the given string and assigns it to the EmailFrom field.
func (o *Creative) SetEmailFrom(v string) {
	o.EmailFrom = &v
}

// GetEmailSubject returns the EmailSubject field value if set, zero value otherwise.
func (o *Creative) GetEmailSubject() string {
	if o == nil || IsNil(o.EmailSubject) {
		var ret string
		return ret
	}
	return *o.EmailSubject
}

// GetEmailSubjectOk returns a tuple with the EmailSubject field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Creative) GetEmailSubjectOk() (*string, bool) {
	if o == nil || IsNil(o.EmailSubject) {
		return nil, false
	}
	return o.EmailSubject, true
}

// HasEmailSubject returns a boolean if a field has been set.
func (o *Creative) HasEmailSubject() bool {
	if o != nil && !IsNil(o.EmailSubject) {
		return true
	}

	return false
}

// SetEmailSubject gets a reference to the given string and assigns it to the EmailSubject field.
func (o *Creative) SetEmailSubject(v string) {
	o.EmailSubject = &v
}

// GetResourceFile returns the ResourceFile field value if set, zero value otherwise.
func (o *Creative) GetResourceFile() ResourceFile {
	if o == nil || IsNil(o.ResourceFile) {
		var ret ResourceFile
		return ret
	}
	return *o.ResourceFile
}

// GetResourceFileOk returns a tuple with the ResourceFile field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Creative) GetResourceFileOk() (*ResourceFile, bool) {
	if o == nil || IsNil(o.ResourceFile) {
		return nil, false
	}
	return o.ResourceFile, true
}

// HasResourceFile returns a boolean if a field has been set.
func (o *Creative) HasResourceFile() bool {
	if o != nil && !IsNil(o.ResourceFile) {
		return true
	}

	return false
}

// SetResourceFile gets a reference to the given ResourceFile and assigns it to the ResourceFile field.
func (o *Creative) SetResourceFile(v ResourceFile) {
	o.ResourceFile = &v
}

// GetHtmlFiles returns the HtmlFiles field value if set, zero value otherwise.
func (o *Creative) GetHtmlFiles() []HtmlFile {
	if o == nil || IsNil(o.HtmlFiles) {
		var ret []HtmlFile
		return ret
	}
	return o.HtmlFiles
}

// GetHtmlFilesOk returns a tuple with the HtmlFiles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Creative) GetHtmlFilesOk() ([]HtmlFile, bool) {
	if o == nil || IsNil(o.HtmlFiles) {
		return nil, false
	}
	return o.HtmlFiles, true
}

// HasHtmlFiles returns a boolean if a field has been set.
func (o *Creative) HasHtmlFiles() bool {
	if o != nil && !IsNil(o.HtmlFiles) {
		return true
	}

	return false
}

// SetHtmlFiles gets a reference to the given []HtmlFile and assigns it to the HtmlFiles field.
func (o *Creative) SetHtmlFiles(v []HtmlFile) {
	o.HtmlFiles = v
}

func (o Creative) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Creative) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["creative_type"] = o.CreativeType
	if !IsNil(o.IsPrivate) {
		toSerialize["is_private"] = o.IsPrivate
	}
	toSerialize["creative_status"] = o.CreativeStatus
	if !IsNil(o.HtmlCode) {
		toSerialize["html_code"] = o.HtmlCode
	}
	if !IsNil(o.Width) {
		toSerialize["width"] = o.Width
	}
	if !IsNil(o.Height) {
		toSerialize["height"] = o.Height
	}
	if !IsNil(o.EmailFrom) {
		toSerialize["email_from"] = o.EmailFrom
	}
	if !IsNil(o.EmailSubject) {
		toSerialize["email_subject"] = o.EmailSubject
	}
	if !IsNil(o.ResourceFile) {
		toSerialize["resource_file"] = o.ResourceFile
	}
	if !IsNil(o.HtmlFiles) {
		toSerialize["html_files"] = o.HtmlFiles
	}
	return toSerialize, nil
}

func (o *Creative) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
		"creative_type",
		"creative_status",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varCreative := _Creative{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varCreative)

	if err != nil {
		return err
	}

	*o = Creative(varCreative)

	return err
}

type NullableCreative struct {
	value *Creative
	isSet bool
}

func (v NullableCreative) Get() *Creative {
	return v.value
}

func (v *NullableCreative) Set(val *Creative) {
	v.value = val
	v.isSet = true
}

func (v NullableCreative) IsSet() bool {
	return v.isSet
}

func (v *NullableCreative) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreative(val *Creative) *NullableCreative {
	return &NullableCreative{value: val, isSet: true}
}

func (v NullableCreative) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreative) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
