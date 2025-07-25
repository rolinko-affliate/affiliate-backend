/*
Everflow Network API - Advertisers

API for managing advertisers in the Everflow platform

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package advertiser

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// checks if the AdvertiserUser type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AdvertiserUser{}

// AdvertiserUser struct for AdvertiserUser
type AdvertiserUser struct {
	// The advertiser user's first name
	FirstName string `json:"first_name"`
	// The advertiser user's last name
	LastName string `json:"last_name"`
	// The advertiser user's email (must be unique)
	Email string `json:"email"`
	// The advertiser user's account status
	AccountStatus string `json:"account_status"`
	// The advertiser user's title
	Title *string `json:"title,omitempty"`
	// The advertiser user's work phone number
	WorkPhone *string `json:"work_phone,omitempty"`
	// The advertiser user's cell phone number
	CellPhone *string `json:"cell_phone,omitempty"`
	// The id of an instant messaging platform
	InstantMessagingId *int32 `json:"instant_messaging_id,omitempty"`
	// The advertiser user's instant messaging identifier
	InstantMessagingIdentifier *string `json:"instant_messaging_identifier,omitempty"`
	// The advertiser user's language id (limited to 1 for English)
	LanguageId int32 `json:"language_id"`
	// The advertiser user's timezone id
	TimezoneId int32 `json:"timezone_id"`
	// The advertiser user's currency id
	CurrencyId string `json:"currency_id"`
	// The advertiser user's login password (optional)
	InitialPassword *string `json:"initial_password,omitempty"`
}

type _AdvertiserUser AdvertiserUser

// NewAdvertiserUser instantiates a new AdvertiserUser object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAdvertiserUser(firstName string, lastName string, email string, accountStatus string, languageId int32, timezoneId int32, currencyId string) *AdvertiserUser {
	this := AdvertiserUser{}
	this.FirstName = firstName
	this.LastName = lastName
	this.Email = email
	this.AccountStatus = accountStatus
	this.LanguageId = languageId
	this.TimezoneId = timezoneId
	this.CurrencyId = currencyId
	return &this
}

// NewAdvertiserUserWithDefaults instantiates a new AdvertiserUser object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAdvertiserUserWithDefaults() *AdvertiserUser {
	this := AdvertiserUser{}
	return &this
}

// GetFirstName returns the FirstName field value
func (o *AdvertiserUser) GetFirstName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FirstName
}

// GetFirstNameOk returns a tuple with the FirstName field value
// and a boolean to check if the value has been set.
func (o *AdvertiserUser) GetFirstNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FirstName, true
}

// SetFirstName sets field value
func (o *AdvertiserUser) SetFirstName(v string) {
	o.FirstName = v
}

// GetLastName returns the LastName field value
func (o *AdvertiserUser) GetLastName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.LastName
}

// GetLastNameOk returns a tuple with the LastName field value
// and a boolean to check if the value has been set.
func (o *AdvertiserUser) GetLastNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LastName, true
}

// SetLastName sets field value
func (o *AdvertiserUser) SetLastName(v string) {
	o.LastName = v
}

// GetEmail returns the Email field value
func (o *AdvertiserUser) GetEmail() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Email
}

// GetEmailOk returns a tuple with the Email field value
// and a boolean to check if the value has been set.
func (o *AdvertiserUser) GetEmailOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Email, true
}

// SetEmail sets field value
func (o *AdvertiserUser) SetEmail(v string) {
	o.Email = v
}

// GetAccountStatus returns the AccountStatus field value
func (o *AdvertiserUser) GetAccountStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AccountStatus
}

// GetAccountStatusOk returns a tuple with the AccountStatus field value
// and a boolean to check if the value has been set.
func (o *AdvertiserUser) GetAccountStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AccountStatus, true
}

// SetAccountStatus sets field value
func (o *AdvertiserUser) SetAccountStatus(v string) {
	o.AccountStatus = v
}

// GetTitle returns the Title field value if set, zero value otherwise.
func (o *AdvertiserUser) GetTitle() string {
	if o == nil || IsNil(o.Title) {
		var ret string
		return ret
	}
	return *o.Title
}

// GetTitleOk returns a tuple with the Title field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdvertiserUser) GetTitleOk() (*string, bool) {
	if o == nil || IsNil(o.Title) {
		return nil, false
	}
	return o.Title, true
}

// HasTitle returns a boolean if a field has been set.
func (o *AdvertiserUser) HasTitle() bool {
	if o != nil && !IsNil(o.Title) {
		return true
	}

	return false
}

// SetTitle gets a reference to the given string and assigns it to the Title field.
func (o *AdvertiserUser) SetTitle(v string) {
	o.Title = &v
}

// GetWorkPhone returns the WorkPhone field value if set, zero value otherwise.
func (o *AdvertiserUser) GetWorkPhone() string {
	if o == nil || IsNil(o.WorkPhone) {
		var ret string
		return ret
	}
	return *o.WorkPhone
}

// GetWorkPhoneOk returns a tuple with the WorkPhone field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdvertiserUser) GetWorkPhoneOk() (*string, bool) {
	if o == nil || IsNil(o.WorkPhone) {
		return nil, false
	}
	return o.WorkPhone, true
}

// HasWorkPhone returns a boolean if a field has been set.
func (o *AdvertiserUser) HasWorkPhone() bool {
	if o != nil && !IsNil(o.WorkPhone) {
		return true
	}

	return false
}

// SetWorkPhone gets a reference to the given string and assigns it to the WorkPhone field.
func (o *AdvertiserUser) SetWorkPhone(v string) {
	o.WorkPhone = &v
}

// GetCellPhone returns the CellPhone field value if set, zero value otherwise.
func (o *AdvertiserUser) GetCellPhone() string {
	if o == nil || IsNil(o.CellPhone) {
		var ret string
		return ret
	}
	return *o.CellPhone
}

// GetCellPhoneOk returns a tuple with the CellPhone field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdvertiserUser) GetCellPhoneOk() (*string, bool) {
	if o == nil || IsNil(o.CellPhone) {
		return nil, false
	}
	return o.CellPhone, true
}

// HasCellPhone returns a boolean if a field has been set.
func (o *AdvertiserUser) HasCellPhone() bool {
	if o != nil && !IsNil(o.CellPhone) {
		return true
	}

	return false
}

// SetCellPhone gets a reference to the given string and assigns it to the CellPhone field.
func (o *AdvertiserUser) SetCellPhone(v string) {
	o.CellPhone = &v
}

// GetInstantMessagingId returns the InstantMessagingId field value if set, zero value otherwise.
func (o *AdvertiserUser) GetInstantMessagingId() int32 {
	if o == nil || IsNil(o.InstantMessagingId) {
		var ret int32
		return ret
	}
	return *o.InstantMessagingId
}

// GetInstantMessagingIdOk returns a tuple with the InstantMessagingId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdvertiserUser) GetInstantMessagingIdOk() (*int32, bool) {
	if o == nil || IsNil(o.InstantMessagingId) {
		return nil, false
	}
	return o.InstantMessagingId, true
}

// HasInstantMessagingId returns a boolean if a field has been set.
func (o *AdvertiserUser) HasInstantMessagingId() bool {
	if o != nil && !IsNil(o.InstantMessagingId) {
		return true
	}

	return false
}

// SetInstantMessagingId gets a reference to the given int32 and assigns it to the InstantMessagingId field.
func (o *AdvertiserUser) SetInstantMessagingId(v int32) {
	o.InstantMessagingId = &v
}

// GetInstantMessagingIdentifier returns the InstantMessagingIdentifier field value if set, zero value otherwise.
func (o *AdvertiserUser) GetInstantMessagingIdentifier() string {
	if o == nil || IsNil(o.InstantMessagingIdentifier) {
		var ret string
		return ret
	}
	return *o.InstantMessagingIdentifier
}

// GetInstantMessagingIdentifierOk returns a tuple with the InstantMessagingIdentifier field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdvertiserUser) GetInstantMessagingIdentifierOk() (*string, bool) {
	if o == nil || IsNil(o.InstantMessagingIdentifier) {
		return nil, false
	}
	return o.InstantMessagingIdentifier, true
}

// HasInstantMessagingIdentifier returns a boolean if a field has been set.
func (o *AdvertiserUser) HasInstantMessagingIdentifier() bool {
	if o != nil && !IsNil(o.InstantMessagingIdentifier) {
		return true
	}

	return false
}

// SetInstantMessagingIdentifier gets a reference to the given string and assigns it to the InstantMessagingIdentifier field.
func (o *AdvertiserUser) SetInstantMessagingIdentifier(v string) {
	o.InstantMessagingIdentifier = &v
}

// GetLanguageId returns the LanguageId field value
func (o *AdvertiserUser) GetLanguageId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.LanguageId
}

// GetLanguageIdOk returns a tuple with the LanguageId field value
// and a boolean to check if the value has been set.
func (o *AdvertiserUser) GetLanguageIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LanguageId, true
}

// SetLanguageId sets field value
func (o *AdvertiserUser) SetLanguageId(v int32) {
	o.LanguageId = v
}

// GetTimezoneId returns the TimezoneId field value
func (o *AdvertiserUser) GetTimezoneId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TimezoneId
}

// GetTimezoneIdOk returns a tuple with the TimezoneId field value
// and a boolean to check if the value has been set.
func (o *AdvertiserUser) GetTimezoneIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TimezoneId, true
}

// SetTimezoneId sets field value
func (o *AdvertiserUser) SetTimezoneId(v int32) {
	o.TimezoneId = v
}

// GetCurrencyId returns the CurrencyId field value
func (o *AdvertiserUser) GetCurrencyId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CurrencyId
}

// GetCurrencyIdOk returns a tuple with the CurrencyId field value
// and a boolean to check if the value has been set.
func (o *AdvertiserUser) GetCurrencyIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CurrencyId, true
}

// SetCurrencyId sets field value
func (o *AdvertiserUser) SetCurrencyId(v string) {
	o.CurrencyId = v
}

// GetInitialPassword returns the InitialPassword field value if set, zero value otherwise.
func (o *AdvertiserUser) GetInitialPassword() string {
	if o == nil || IsNil(o.InitialPassword) {
		var ret string
		return ret
	}
	return *o.InitialPassword
}

// GetInitialPasswordOk returns a tuple with the InitialPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdvertiserUser) GetInitialPasswordOk() (*string, bool) {
	if o == nil || IsNil(o.InitialPassword) {
		return nil, false
	}
	return o.InitialPassword, true
}

// HasInitialPassword returns a boolean if a field has been set.
func (o *AdvertiserUser) HasInitialPassword() bool {
	if o != nil && !IsNil(o.InitialPassword) {
		return true
	}

	return false
}

// SetInitialPassword gets a reference to the given string and assigns it to the InitialPassword field.
func (o *AdvertiserUser) SetInitialPassword(v string) {
	o.InitialPassword = &v
}

func (o AdvertiserUser) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AdvertiserUser) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["first_name"] = o.FirstName
	toSerialize["last_name"] = o.LastName
	toSerialize["email"] = o.Email
	toSerialize["account_status"] = o.AccountStatus
	if !IsNil(o.Title) {
		toSerialize["title"] = o.Title
	}
	if !IsNil(o.WorkPhone) {
		toSerialize["work_phone"] = o.WorkPhone
	}
	if !IsNil(o.CellPhone) {
		toSerialize["cell_phone"] = o.CellPhone
	}
	if !IsNil(o.InstantMessagingId) {
		toSerialize["instant_messaging_id"] = o.InstantMessagingId
	}
	if !IsNil(o.InstantMessagingIdentifier) {
		toSerialize["instant_messaging_identifier"] = o.InstantMessagingIdentifier
	}
	toSerialize["language_id"] = o.LanguageId
	toSerialize["timezone_id"] = o.TimezoneId
	toSerialize["currency_id"] = o.CurrencyId
	if !IsNil(o.InitialPassword) {
		toSerialize["initial_password"] = o.InitialPassword
	}
	return toSerialize, nil
}

func (o *AdvertiserUser) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"first_name",
		"last_name",
		"email",
		"account_status",
		"language_id",
		"timezone_id",
		"currency_id",
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

	varAdvertiserUser := _AdvertiserUser{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varAdvertiserUser)

	if err != nil {
		return err
	}

	*o = AdvertiserUser(varAdvertiserUser)

	return err
}

type NullableAdvertiserUser struct {
	value *AdvertiserUser
	isSet bool
}

func (v NullableAdvertiserUser) Get() *AdvertiserUser {
	return v.value
}

func (v *NullableAdvertiserUser) Set(val *AdvertiserUser) {
	v.value = val
	v.isSet = true
}

func (v NullableAdvertiserUser) IsSet() bool {
	return v.isSet
}

func (v *NullableAdvertiserUser) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAdvertiserUser(val *AdvertiserUser) *NullableAdvertiserUser {
	return &NullableAdvertiserUser{value: val, isSet: true}
}

func (v NullableAdvertiserUser) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAdvertiserUser) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
