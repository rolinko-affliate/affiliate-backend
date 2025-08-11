/*
Everflow Network API - Tracking

API for generating tracking links in the Everflow platform

API version: 1.0.0
*/

package tracking

import (
	"encoding/json"
)

// CreateTrackingLinkRequest struct for CreateTrackingLinkRequest
type CreateTrackingLinkRequest struct {
	NetworkAffiliateId     int32  `json:"network_affiliate_id"`
	NetworkOfferId         int32  `json:"network_offer_id"`
	NetworkTrackingDomainId *int32 `json:"network_tracking_domain_id,omitempty"`
	NetworkOfferUrlId      *int32 `json:"network_offer_url_id,omitempty"`
	CreativeId             *int32 `json:"creative_id,omitempty"`
	NetworkTrafficSourceId *int32 `json:"network_traffic_source_id,omitempty"`
	SourceId               *string `json:"source_id,omitempty"`
	Sub1                   *string `json:"sub1,omitempty"`
	Sub2                   *string `json:"sub2,omitempty"`
	Sub3                   *string `json:"sub3,omitempty"`
	Sub4                   *string `json:"sub4,omitempty"`
	Sub5                   *string `json:"sub5,omitempty"`
	IsEncryptParameters    *bool   `json:"is_encrypt_parameters,omitempty"`
	IsRedirectLink         *bool   `json:"is_redirect_link,omitempty"`
}

// NewCreateTrackingLinkRequest instantiates a new CreateTrackingLinkRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateTrackingLinkRequest(networkAffiliateId int32, networkOfferId int32) *CreateTrackingLinkRequest {
	this := CreateTrackingLinkRequest{}
	this.NetworkAffiliateId = networkAffiliateId
	this.NetworkOfferId = networkOfferId
	return &this
}

// NewCreateTrackingLinkRequestWithDefaults instantiates a new CreateTrackingLinkRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateTrackingLinkRequestWithDefaults() *CreateTrackingLinkRequest {
	this := CreateTrackingLinkRequest{}
	return &this
}

// GetNetworkAffiliateId returns the NetworkAffiliateId field value
func (o *CreateTrackingLinkRequest) GetNetworkAffiliateId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.NetworkAffiliateId
}

// GetNetworkAffiliateIdOk returns a tuple with the NetworkAffiliateId field value
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetNetworkAffiliateIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NetworkAffiliateId, true
}

// SetNetworkAffiliateId sets field value
func (o *CreateTrackingLinkRequest) SetNetworkAffiliateId(v int32) {
	o.NetworkAffiliateId = v
}

// GetNetworkOfferId returns the NetworkOfferId field value
func (o *CreateTrackingLinkRequest) GetNetworkOfferId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.NetworkOfferId
}

// GetNetworkOfferIdOk returns a tuple with the NetworkOfferId field value
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetNetworkOfferIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NetworkOfferId, true
}

// SetNetworkOfferId sets field value
func (o *CreateTrackingLinkRequest) SetNetworkOfferId(v int32) {
	o.NetworkOfferId = v
}

// GetNetworkTrackingDomainId returns the NetworkTrackingDomainId field value if set, zero value otherwise.
func (o *CreateTrackingLinkRequest) GetNetworkTrackingDomainId() int32 {
	if o == nil || o.NetworkTrackingDomainId == nil {
		var ret int32
		return ret
	}
	return *o.NetworkTrackingDomainId
}

// GetNetworkTrackingDomainIdOk returns a tuple with the NetworkTrackingDomainId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetNetworkTrackingDomainIdOk() (*int32, bool) {
	if o == nil || o.NetworkTrackingDomainId == nil {
		return nil, false
	}
	return o.NetworkTrackingDomainId, true
}

// HasNetworkTrackingDomainId returns a boolean if a field has been set.
func (o *CreateTrackingLinkRequest) HasNetworkTrackingDomainId() bool {
	if o != nil && o.NetworkTrackingDomainId != nil {
		return true
	}

	return false
}

// SetNetworkTrackingDomainId gets a reference to the given int32 and assigns it to the NetworkTrackingDomainId field.
func (o *CreateTrackingLinkRequest) SetNetworkTrackingDomainId(v int32) {
	o.NetworkTrackingDomainId = &v
}

// GetNetworkOfferUrlId returns the NetworkOfferUrlId field value if set, zero value otherwise.
func (o *CreateTrackingLinkRequest) GetNetworkOfferUrlId() int32 {
	if o == nil || o.NetworkOfferUrlId == nil {
		var ret int32
		return ret
	}
	return *o.NetworkOfferUrlId
}

// GetNetworkOfferUrlIdOk returns a tuple with the NetworkOfferUrlId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetNetworkOfferUrlIdOk() (*int32, bool) {
	if o == nil || o.NetworkOfferUrlId == nil {
		return nil, false
	}
	return o.NetworkOfferUrlId, true
}

// HasNetworkOfferUrlId returns a boolean if a field has been set.
func (o *CreateTrackingLinkRequest) HasNetworkOfferUrlId() bool {
	if o != nil && o.NetworkOfferUrlId != nil {
		return true
	}

	return false
}

// SetNetworkOfferUrlId gets a reference to the given int32 and assigns it to the NetworkOfferUrlId field.
func (o *CreateTrackingLinkRequest) SetNetworkOfferUrlId(v int32) {
	o.NetworkOfferUrlId = &v
}

// GetCreativeId returns the CreativeId field value if set, zero value otherwise.
func (o *CreateTrackingLinkRequest) GetCreativeId() int32 {
	if o == nil || o.CreativeId == nil {
		var ret int32
		return ret
	}
	return *o.CreativeId
}

// GetCreativeIdOk returns a tuple with the CreativeId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetCreativeIdOk() (*int32, bool) {
	if o == nil || o.CreativeId == nil {
		return nil, false
	}
	return o.CreativeId, true
}

// HasCreativeId returns a boolean if a field has been set.
func (o *CreateTrackingLinkRequest) HasCreativeId() bool {
	if o != nil && o.CreativeId != nil {
		return true
	}

	return false
}

// SetCreativeId gets a reference to the given int32 and assigns it to the CreativeId field.
func (o *CreateTrackingLinkRequest) SetCreativeId(v int32) {
	o.CreativeId = &v
}

// GetNetworkTrafficSourceId returns the NetworkTrafficSourceId field value if set, zero value otherwise.
func (o *CreateTrackingLinkRequest) GetNetworkTrafficSourceId() int32 {
	if o == nil || o.NetworkTrafficSourceId == nil {
		var ret int32
		return ret
	}
	return *o.NetworkTrafficSourceId
}

// GetNetworkTrafficSourceIdOk returns a tuple with the NetworkTrafficSourceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetNetworkTrafficSourceIdOk() (*int32, bool) {
	if o == nil || o.NetworkTrafficSourceId == nil {
		return nil, false
	}
	return o.NetworkTrafficSourceId, true
}

// HasNetworkTrafficSourceId returns a boolean if a field has been set.
func (o *CreateTrackingLinkRequest) HasNetworkTrafficSourceId() bool {
	if o != nil && o.NetworkTrafficSourceId != nil {
		return true
	}

	return false
}

// SetNetworkTrafficSourceId gets a reference to the given int32 and assigns it to the NetworkTrafficSourceId field.
func (o *CreateTrackingLinkRequest) SetNetworkTrafficSourceId(v int32) {
	o.NetworkTrafficSourceId = &v
}

// GetSourceId returns the SourceId field value if set, zero value otherwise.
func (o *CreateTrackingLinkRequest) GetSourceId() string {
	if o == nil || o.SourceId == nil {
		var ret string
		return ret
	}
	return *o.SourceId
}

// GetSourceIdOk returns a tuple with the SourceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetSourceIdOk() (*string, bool) {
	if o == nil || o.SourceId == nil {
		return nil, false
	}
	return o.SourceId, true
}

// HasSourceId returns a boolean if a field has been set.
func (o *CreateTrackingLinkRequest) HasSourceId() bool {
	if o != nil && o.SourceId != nil {
		return true
	}

	return false
}

// SetSourceId gets a reference to the given string and assigns it to the SourceId field.
func (o *CreateTrackingLinkRequest) SetSourceId(v string) {
	o.SourceId = &v
}

// GetSub1 returns the Sub1 field value if set, zero value otherwise.
func (o *CreateTrackingLinkRequest) GetSub1() string {
	if o == nil || o.Sub1 == nil {
		var ret string
		return ret
	}
	return *o.Sub1
}

// GetSub1Ok returns a tuple with the Sub1 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetSub1Ok() (*string, bool) {
	if o == nil || o.Sub1 == nil {
		return nil, false
	}
	return o.Sub1, true
}

// HasSub1 returns a boolean if a field has been set.
func (o *CreateTrackingLinkRequest) HasSub1() bool {
	if o != nil && o.Sub1 != nil {
		return true
	}

	return false
}

// SetSub1 gets a reference to the given string and assigns it to the Sub1 field.
func (o *CreateTrackingLinkRequest) SetSub1(v string) {
	o.Sub1 = &v
}

// GetSub2 returns the Sub2 field value if set, zero value otherwise.
func (o *CreateTrackingLinkRequest) GetSub2() string {
	if o == nil || o.Sub2 == nil {
		var ret string
		return ret
	}
	return *o.Sub2
}

// GetSub2Ok returns a tuple with the Sub2 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetSub2Ok() (*string, bool) {
	if o == nil || o.Sub2 == nil {
		return nil, false
	}
	return o.Sub2, true
}

// HasSub2 returns a boolean if a field has been set.
func (o *CreateTrackingLinkRequest) HasSub2() bool {
	if o != nil && o.Sub2 != nil {
		return true
	}

	return false
}

// SetSub2 gets a reference to the given string and assigns it to the Sub2 field.
func (o *CreateTrackingLinkRequest) SetSub2(v string) {
	o.Sub2 = &v
}

// GetSub3 returns the Sub3 field value if set, zero value otherwise.
func (o *CreateTrackingLinkRequest) GetSub3() string {
	if o == nil || o.Sub3 == nil {
		var ret string
		return ret
	}
	return *o.Sub3
}

// GetSub3Ok returns a tuple with the Sub3 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetSub3Ok() (*string, bool) {
	if o == nil || o.Sub3 == nil {
		return nil, false
	}
	return o.Sub3, true
}

// HasSub3 returns a boolean if a field has been set.
func (o *CreateTrackingLinkRequest) HasSub3() bool {
	if o != nil && o.Sub3 != nil {
		return true
	}

	return false
}

// SetSub3 gets a reference to the given string and assigns it to the Sub3 field.
func (o *CreateTrackingLinkRequest) SetSub3(v string) {
	o.Sub3 = &v
}

// GetSub4 returns the Sub4 field value if set, zero value otherwise.
func (o *CreateTrackingLinkRequest) GetSub4() string {
	if o == nil || o.Sub4 == nil {
		var ret string
		return ret
	}
	return *o.Sub4
}

// GetSub4Ok returns a tuple with the Sub4 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetSub4Ok() (*string, bool) {
	if o == nil || o.Sub4 == nil {
		return nil, false
	}
	return o.Sub4, true
}

// HasSub4 returns a boolean if a field has been set.
func (o *CreateTrackingLinkRequest) HasSub4() bool {
	if o != nil && o.Sub4 != nil {
		return true
	}

	return false
}

// SetSub4 gets a reference to the given string and assigns it to the Sub4 field.
func (o *CreateTrackingLinkRequest) SetSub4(v string) {
	o.Sub4 = &v
}

// GetSub5 returns the Sub5 field value if set, zero value otherwise.
func (o *CreateTrackingLinkRequest) GetSub5() string {
	if o == nil || o.Sub5 == nil {
		var ret string
		return ret
	}
	return *o.Sub5
}

// GetSub5Ok returns a tuple with the Sub5 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetSub5Ok() (*string, bool) {
	if o == nil || o.Sub5 == nil {
		return nil, false
	}
	return o.Sub5, true
}

// HasSub5 returns a boolean if a field has been set.
func (o *CreateTrackingLinkRequest) HasSub5() bool {
	if o != nil && o.Sub5 != nil {
		return true
	}

	return false
}

// SetSub5 gets a reference to the given string and assigns it to the Sub5 field.
func (o *CreateTrackingLinkRequest) SetSub5(v string) {
	o.Sub5 = &v
}

// GetIsEncryptParameters returns the IsEncryptParameters field value if set, zero value otherwise.
func (o *CreateTrackingLinkRequest) GetIsEncryptParameters() bool {
	if o == nil || o.IsEncryptParameters == nil {
		var ret bool
		return ret
	}
	return *o.IsEncryptParameters
}

// GetIsEncryptParametersOk returns a tuple with the IsEncryptParameters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetIsEncryptParametersOk() (*bool, bool) {
	if o == nil || o.IsEncryptParameters == nil {
		return nil, false
	}
	return o.IsEncryptParameters, true
}

// HasIsEncryptParameters returns a boolean if a field has been set.
func (o *CreateTrackingLinkRequest) HasIsEncryptParameters() bool {
	if o != nil && o.IsEncryptParameters != nil {
		return true
	}

	return false
}

// SetIsEncryptParameters gets a reference to the given bool and assigns it to the IsEncryptParameters field.
func (o *CreateTrackingLinkRequest) SetIsEncryptParameters(v bool) {
	o.IsEncryptParameters = &v
}

// GetIsRedirectLink returns the IsRedirectLink field value if set, zero value otherwise.
func (o *CreateTrackingLinkRequest) GetIsRedirectLink() bool {
	if o == nil || o.IsRedirectLink == nil {
		var ret bool
		return ret
	}
	return *o.IsRedirectLink
}

// GetIsRedirectLinkOk returns a tuple with the IsRedirectLink field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateTrackingLinkRequest) GetIsRedirectLinkOk() (*bool, bool) {
	if o == nil || o.IsRedirectLink == nil {
		return nil, false
	}
	return o.IsRedirectLink, true
}

// HasIsRedirectLink returns a boolean if a field has been set.
func (o *CreateTrackingLinkRequest) HasIsRedirectLink() bool {
	if o != nil && o.IsRedirectLink != nil {
		return true
	}

	return false
}

// SetIsRedirectLink gets a reference to the given bool and assigns it to the IsRedirectLink field.
func (o *CreateTrackingLinkRequest) SetIsRedirectLink(v bool) {
	o.IsRedirectLink = &v
}

func (o CreateTrackingLinkRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["network_affiliate_id"] = o.NetworkAffiliateId
	}
	if true {
		toSerialize["network_offer_id"] = o.NetworkOfferId
	}
	if o.NetworkTrackingDomainId != nil {
		toSerialize["network_tracking_domain_id"] = o.NetworkTrackingDomainId
	}
	if o.NetworkOfferUrlId != nil {
		toSerialize["network_offer_url_id"] = o.NetworkOfferUrlId
	}
	if o.CreativeId != nil {
		toSerialize["creative_id"] = o.CreativeId
	}
	if o.NetworkTrafficSourceId != nil {
		toSerialize["network_traffic_source_id"] = o.NetworkTrafficSourceId
	}
	if o.SourceId != nil {
		toSerialize["source_id"] = o.SourceId
	}
	if o.Sub1 != nil {
		toSerialize["sub1"] = o.Sub1
	}
	if o.Sub2 != nil {
		toSerialize["sub2"] = o.Sub2
	}
	if o.Sub3 != nil {
		toSerialize["sub3"] = o.Sub3
	}
	if o.Sub4 != nil {
		toSerialize["sub4"] = o.Sub4
	}
	if o.Sub5 != nil {
		toSerialize["sub5"] = o.Sub5
	}
	if o.IsEncryptParameters != nil {
		toSerialize["is_encrypt_parameters"] = o.IsEncryptParameters
	}
	if o.IsRedirectLink != nil {
		toSerialize["is_redirect_link"] = o.IsRedirectLink
	}
	return json.Marshal(toSerialize)
}

type NullableCreateTrackingLinkRequest struct {
	value *CreateTrackingLinkRequest
	isSet bool
}

func (v NullableCreateTrackingLinkRequest) Get() *CreateTrackingLinkRequest {
	return v.value
}

func (v *NullableCreateTrackingLinkRequest) Set(val *CreateTrackingLinkRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateTrackingLinkRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateTrackingLinkRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateTrackingLinkRequest(val *CreateTrackingLinkRequest) *NullableCreateTrackingLinkRequest {
	return &NullableCreateTrackingLinkRequest{value: val, isSet: true}
}

func (v NullableCreateTrackingLinkRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateTrackingLinkRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}