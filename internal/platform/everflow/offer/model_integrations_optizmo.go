/*
Everflow Network API - Offers

API for managing offers in the Everflow platform

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package offer

import (
	"encoding/json"
)

// checks if the IntegrationsOptizmo type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &IntegrationsOptizmo{}

// IntegrationsOptizmo struct for IntegrationsOptizmo
type IntegrationsOptizmo struct {
	// Optizmo optout list ID
	OptoutlistId *string `json:"optoutlist_id,omitempty"`
}

// NewIntegrationsOptizmo instantiates a new IntegrationsOptizmo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewIntegrationsOptizmo() *IntegrationsOptizmo {
	this := IntegrationsOptizmo{}
	return &this
}

// NewIntegrationsOptizmoWithDefaults instantiates a new IntegrationsOptizmo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewIntegrationsOptizmoWithDefaults() *IntegrationsOptizmo {
	this := IntegrationsOptizmo{}
	return &this
}

// GetOptoutlistId returns the OptoutlistId field value if set, zero value otherwise.
func (o *IntegrationsOptizmo) GetOptoutlistId() string {
	if o == nil || IsNil(o.OptoutlistId) {
		var ret string
		return ret
	}
	return *o.OptoutlistId
}

// GetOptoutlistIdOk returns a tuple with the OptoutlistId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IntegrationsOptizmo) GetOptoutlistIdOk() (*string, bool) {
	if o == nil || IsNil(o.OptoutlistId) {
		return nil, false
	}
	return o.OptoutlistId, true
}

// HasOptoutlistId returns a boolean if a field has been set.
func (o *IntegrationsOptizmo) HasOptoutlistId() bool {
	if o != nil && !IsNil(o.OptoutlistId) {
		return true
	}

	return false
}

// SetOptoutlistId gets a reference to the given string and assigns it to the OptoutlistId field.
func (o *IntegrationsOptizmo) SetOptoutlistId(v string) {
	o.OptoutlistId = &v
}

func (o IntegrationsOptizmo) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o IntegrationsOptizmo) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.OptoutlistId) {
		toSerialize["optoutlist_id"] = o.OptoutlistId
	}
	return toSerialize, nil
}

type NullableIntegrationsOptizmo struct {
	value *IntegrationsOptizmo
	isSet bool
}

func (v NullableIntegrationsOptizmo) Get() *IntegrationsOptizmo {
	return v.value
}

func (v *NullableIntegrationsOptizmo) Set(val *IntegrationsOptizmo) {
	v.value = val
	v.isSet = true
}

func (v NullableIntegrationsOptizmo) IsSet() bool {
	return v.isSet
}

func (v *NullableIntegrationsOptizmo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableIntegrationsOptizmo(val *IntegrationsOptizmo) *NullableIntegrationsOptizmo {
	return &NullableIntegrationsOptizmo{value: val, isSet: true}
}

func (v NullableIntegrationsOptizmo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableIntegrationsOptizmo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
