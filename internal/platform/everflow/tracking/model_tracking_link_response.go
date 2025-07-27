/*
Everflow Network API - Tracking

API for generating tracking links in the Everflow platform

API version: 1.0.0
*/

package tracking

import (
	"encoding/json"
)

// TrackingLinkResponse struct for TrackingLinkResponse
type TrackingLinkResponse struct {
	TrackingUrl            *string `json:"tracking_url,omitempty"`
	NetworkAffiliateId     *int32  `json:"network_affiliate_id,omitempty"`
	NetworkOfferId         *int32  `json:"network_offer_id,omitempty"`
	NetworkTrackingDomainId *int32 `json:"network_tracking_domain_id,omitempty"`
	NetworkOfferUrlId      *int32  `json:"network_offer_url_id,omitempty"`
	CreativeId             *int32  `json:"creative_id,omitempty"`
	NetworkTrafficSourceId *int32  `json:"network_traffic_source_id,omitempty"`
	SourceId               *string `json:"source_id,omitempty"`
	Sub1                   *string `json:"sub1,omitempty"`
	Sub2                   *string `json:"sub2,omitempty"`
	Sub3                   *string `json:"sub3,omitempty"`
	Sub4                   *string `json:"sub4,omitempty"`
	Sub5                   *string `json:"sub5,omitempty"`
	IsEncryptParameters    *bool   `json:"is_encrypt_parameters,omitempty"`
	IsRedirectLink         *bool   `json:"is_redirect_link,omitempty"`
}

// NewTrackingLinkResponse instantiates a new TrackingLinkResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTrackingLinkResponse() *TrackingLinkResponse {
	this := TrackingLinkResponse{}
	return &this
}

// NewTrackingLinkResponseWithDefaults instantiates a new TrackingLinkResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTrackingLinkResponseWithDefaults() *TrackingLinkResponse {
	this := TrackingLinkResponse{}
	return &this
}

// GetTrackingUrl returns the TrackingUrl field value if set, zero value otherwise.
func (o *TrackingLinkResponse) GetTrackingUrl() string {
	if o == nil || o.TrackingUrl == nil {
		var ret string
		return ret
	}
	return *o.TrackingUrl
}

// GetTrackingUrlOk returns a tuple with the TrackingUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TrackingLinkResponse) GetTrackingUrlOk() (*string, bool) {
	if o == nil || o.TrackingUrl == nil {
		return nil, false
	}
	return o.TrackingUrl, true
}

// HasTrackingUrl returns a boolean if a field has been set.
func (o *TrackingLinkResponse) HasTrackingUrl() bool {
	if o != nil && o.TrackingUrl != nil {
		return true
	}

	return false
}

// SetTrackingUrl gets a reference to the given string and assigns it to the TrackingUrl field.
func (o *TrackingLinkResponse) SetTrackingUrl(v string) {
	o.TrackingUrl = &v
}

// GetNetworkAffiliateId returns the NetworkAffiliateId field value if set, zero value otherwise.
func (o *TrackingLinkResponse) GetNetworkAffiliateId() int32 {
	if o == nil || o.NetworkAffiliateId == nil {
		var ret int32
		return ret
	}
	return *o.NetworkAffiliateId
}

// GetNetworkAffiliateIdOk returns a tuple with the NetworkAffiliateId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TrackingLinkResponse) GetNetworkAffiliateIdOk() (*int32, bool) {
	if o == nil || o.NetworkAffiliateId == nil {
		return nil, false
	}
	return o.NetworkAffiliateId, true
}

// HasNetworkAffiliateId returns a boolean if a field has been set.
func (o *TrackingLinkResponse) HasNetworkAffiliateId() bool {
	if o != nil && o.NetworkAffiliateId != nil {
		return true
	}

	return false
}

// SetNetworkAffiliateId gets a reference to the given int32 and assigns it to the NetworkAffiliateId field.
func (o *TrackingLinkResponse) SetNetworkAffiliateId(v int32) {
	o.NetworkAffiliateId = &v
}

// GetNetworkOfferId returns the NetworkOfferId field value if set, zero value otherwise.
func (o *TrackingLinkResponse) GetNetworkOfferId() int32 {
	if o == nil || o.NetworkOfferId == nil {
		var ret int32
		return ret
	}
	return *o.NetworkOfferId
}

// GetNetworkOfferIdOk returns a tuple with the NetworkOfferId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TrackingLinkResponse) GetNetworkOfferIdOk() (*int32, bool) {
	if o == nil || o.NetworkOfferId == nil {
		return nil, false
	}
	return o.NetworkOfferId, true
}

// HasNetworkOfferId returns a boolean if a field has been set.
func (o *TrackingLinkResponse) HasNetworkOfferId() bool {
	if o != nil && o.NetworkOfferId != nil {
		return true
	}

	return false
}

// SetNetworkOfferId gets a reference to the given int32 and assigns it to the NetworkOfferId field.
func (o *TrackingLinkResponse) SetNetworkOfferId(v int32) {
	o.NetworkOfferId = &v
}

// Additional getter/setter methods for other fields would follow the same pattern...
// For brevity, I'll include the MarshalJSON method

func (o TrackingLinkResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.TrackingUrl != nil {
		toSerialize["tracking_url"] = o.TrackingUrl
	}
	if o.NetworkAffiliateId != nil {
		toSerialize["network_affiliate_id"] = o.NetworkAffiliateId
	}
	if o.NetworkOfferId != nil {
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

type NullableTrackingLinkResponse struct {
	value *TrackingLinkResponse
	isSet bool
}

func (v NullableTrackingLinkResponse) Get() *TrackingLinkResponse {
	return v.value
}

func (v *NullableTrackingLinkResponse) Set(val *TrackingLinkResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableTrackingLinkResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableTrackingLinkResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTrackingLinkResponse(val *TrackingLinkResponse) *NullableTrackingLinkResponse {
	return &NullableTrackingLinkResponse{value: val, isSet: true}
}

func (v NullableTrackingLinkResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTrackingLinkResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}