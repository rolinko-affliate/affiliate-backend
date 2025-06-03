# InternalRedirect

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RedirectNetworkOfferId** | Pointer to **int32** | Offer ID to redirect to | [optional] 
**RedirectNetworkOfferUrlId** | Pointer to **int32** | Offer URL ID (0 for default) | [optional] 
**RedirectNetworkOfferGroupId** | Pointer to **int32** | Offer group ID to redirect to | [optional] 
**RedirectNetworkCampaignId** | Pointer to **int32** | Campaign ID to redirect to | [optional] 
**RoutingValue** | Pointer to **int32** | Priority or weight value | [optional] 
**Ruleset** | Pointer to [**Ruleset**](Ruleset.md) |  | [optional] 
**Categories** | Pointer to **[]string** | Fail traffic categories | [optional] 
**IsPayAffiliate** | Pointer to **bool** | Whether to pay affiliate | [optional] [default to false]
**IsPassThrough** | Pointer to **bool** | Whether to pass through to destination | [optional] [default to false]
**IsApplySpecificAffiliates** | Pointer to **bool** | Whether to apply to specific affiliates | [optional] [default to false]
**NetworkAffiliateIds** | Pointer to **[]int32** | Specific affiliate IDs | [optional] 

## Methods

### NewInternalRedirect

`func NewInternalRedirect() *InternalRedirect`

NewInternalRedirect instantiates a new InternalRedirect object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInternalRedirectWithDefaults

`func NewInternalRedirectWithDefaults() *InternalRedirect`

NewInternalRedirectWithDefaults instantiates a new InternalRedirect object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRedirectNetworkOfferId

`func (o *InternalRedirect) GetRedirectNetworkOfferId() int32`

GetRedirectNetworkOfferId returns the RedirectNetworkOfferId field if non-nil, zero value otherwise.

### GetRedirectNetworkOfferIdOk

`func (o *InternalRedirect) GetRedirectNetworkOfferIdOk() (*int32, bool)`

GetRedirectNetworkOfferIdOk returns a tuple with the RedirectNetworkOfferId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectNetworkOfferId

`func (o *InternalRedirect) SetRedirectNetworkOfferId(v int32)`

SetRedirectNetworkOfferId sets RedirectNetworkOfferId field to given value.

### HasRedirectNetworkOfferId

`func (o *InternalRedirect) HasRedirectNetworkOfferId() bool`

HasRedirectNetworkOfferId returns a boolean if a field has been set.

### GetRedirectNetworkOfferUrlId

`func (o *InternalRedirect) GetRedirectNetworkOfferUrlId() int32`

GetRedirectNetworkOfferUrlId returns the RedirectNetworkOfferUrlId field if non-nil, zero value otherwise.

### GetRedirectNetworkOfferUrlIdOk

`func (o *InternalRedirect) GetRedirectNetworkOfferUrlIdOk() (*int32, bool)`

GetRedirectNetworkOfferUrlIdOk returns a tuple with the RedirectNetworkOfferUrlId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectNetworkOfferUrlId

`func (o *InternalRedirect) SetRedirectNetworkOfferUrlId(v int32)`

SetRedirectNetworkOfferUrlId sets RedirectNetworkOfferUrlId field to given value.

### HasRedirectNetworkOfferUrlId

`func (o *InternalRedirect) HasRedirectNetworkOfferUrlId() bool`

HasRedirectNetworkOfferUrlId returns a boolean if a field has been set.

### GetRedirectNetworkOfferGroupId

`func (o *InternalRedirect) GetRedirectNetworkOfferGroupId() int32`

GetRedirectNetworkOfferGroupId returns the RedirectNetworkOfferGroupId field if non-nil, zero value otherwise.

### GetRedirectNetworkOfferGroupIdOk

`func (o *InternalRedirect) GetRedirectNetworkOfferGroupIdOk() (*int32, bool)`

GetRedirectNetworkOfferGroupIdOk returns a tuple with the RedirectNetworkOfferGroupId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectNetworkOfferGroupId

`func (o *InternalRedirect) SetRedirectNetworkOfferGroupId(v int32)`

SetRedirectNetworkOfferGroupId sets RedirectNetworkOfferGroupId field to given value.

### HasRedirectNetworkOfferGroupId

`func (o *InternalRedirect) HasRedirectNetworkOfferGroupId() bool`

HasRedirectNetworkOfferGroupId returns a boolean if a field has been set.

### GetRedirectNetworkCampaignId

`func (o *InternalRedirect) GetRedirectNetworkCampaignId() int32`

GetRedirectNetworkCampaignId returns the RedirectNetworkCampaignId field if non-nil, zero value otherwise.

### GetRedirectNetworkCampaignIdOk

`func (o *InternalRedirect) GetRedirectNetworkCampaignIdOk() (*int32, bool)`

GetRedirectNetworkCampaignIdOk returns a tuple with the RedirectNetworkCampaignId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectNetworkCampaignId

`func (o *InternalRedirect) SetRedirectNetworkCampaignId(v int32)`

SetRedirectNetworkCampaignId sets RedirectNetworkCampaignId field to given value.

### HasRedirectNetworkCampaignId

`func (o *InternalRedirect) HasRedirectNetworkCampaignId() bool`

HasRedirectNetworkCampaignId returns a boolean if a field has been set.

### GetRoutingValue

`func (o *InternalRedirect) GetRoutingValue() int32`

GetRoutingValue returns the RoutingValue field if non-nil, zero value otherwise.

### GetRoutingValueOk

`func (o *InternalRedirect) GetRoutingValueOk() (*int32, bool)`

GetRoutingValueOk returns a tuple with the RoutingValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoutingValue

`func (o *InternalRedirect) SetRoutingValue(v int32)`

SetRoutingValue sets RoutingValue field to given value.

### HasRoutingValue

`func (o *InternalRedirect) HasRoutingValue() bool`

HasRoutingValue returns a boolean if a field has been set.

### GetRuleset

`func (o *InternalRedirect) GetRuleset() Ruleset`

GetRuleset returns the Ruleset field if non-nil, zero value otherwise.

### GetRulesetOk

`func (o *InternalRedirect) GetRulesetOk() (*Ruleset, bool)`

GetRulesetOk returns a tuple with the Ruleset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRuleset

`func (o *InternalRedirect) SetRuleset(v Ruleset)`

SetRuleset sets Ruleset field to given value.

### HasRuleset

`func (o *InternalRedirect) HasRuleset() bool`

HasRuleset returns a boolean if a field has been set.

### GetCategories

`func (o *InternalRedirect) GetCategories() []string`

GetCategories returns the Categories field if non-nil, zero value otherwise.

### GetCategoriesOk

`func (o *InternalRedirect) GetCategoriesOk() (*[]string, bool)`

GetCategoriesOk returns a tuple with the Categories field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCategories

`func (o *InternalRedirect) SetCategories(v []string)`

SetCategories sets Categories field to given value.

### HasCategories

`func (o *InternalRedirect) HasCategories() bool`

HasCategories returns a boolean if a field has been set.

### GetIsPayAffiliate

`func (o *InternalRedirect) GetIsPayAffiliate() bool`

GetIsPayAffiliate returns the IsPayAffiliate field if non-nil, zero value otherwise.

### GetIsPayAffiliateOk

`func (o *InternalRedirect) GetIsPayAffiliateOk() (*bool, bool)`

GetIsPayAffiliateOk returns a tuple with the IsPayAffiliate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPayAffiliate

`func (o *InternalRedirect) SetIsPayAffiliate(v bool)`

SetIsPayAffiliate sets IsPayAffiliate field to given value.

### HasIsPayAffiliate

`func (o *InternalRedirect) HasIsPayAffiliate() bool`

HasIsPayAffiliate returns a boolean if a field has been set.

### GetIsPassThrough

`func (o *InternalRedirect) GetIsPassThrough() bool`

GetIsPassThrough returns the IsPassThrough field if non-nil, zero value otherwise.

### GetIsPassThroughOk

`func (o *InternalRedirect) GetIsPassThroughOk() (*bool, bool)`

GetIsPassThroughOk returns a tuple with the IsPassThrough field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPassThrough

`func (o *InternalRedirect) SetIsPassThrough(v bool)`

SetIsPassThrough sets IsPassThrough field to given value.

### HasIsPassThrough

`func (o *InternalRedirect) HasIsPassThrough() bool`

HasIsPassThrough returns a boolean if a field has been set.

### GetIsApplySpecificAffiliates

`func (o *InternalRedirect) GetIsApplySpecificAffiliates() bool`

GetIsApplySpecificAffiliates returns the IsApplySpecificAffiliates field if non-nil, zero value otherwise.

### GetIsApplySpecificAffiliatesOk

`func (o *InternalRedirect) GetIsApplySpecificAffiliatesOk() (*bool, bool)`

GetIsApplySpecificAffiliatesOk returns a tuple with the IsApplySpecificAffiliates field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsApplySpecificAffiliates

`func (o *InternalRedirect) SetIsApplySpecificAffiliates(v bool)`

SetIsApplySpecificAffiliates sets IsApplySpecificAffiliates field to given value.

### HasIsApplySpecificAffiliates

`func (o *InternalRedirect) HasIsApplySpecificAffiliates() bool`

HasIsApplySpecificAffiliates returns a boolean if a field has been set.

### GetNetworkAffiliateIds

`func (o *InternalRedirect) GetNetworkAffiliateIds() []int32`

GetNetworkAffiliateIds returns the NetworkAffiliateIds field if non-nil, zero value otherwise.

### GetNetworkAffiliateIdsOk

`func (o *InternalRedirect) GetNetworkAffiliateIdsOk() (*[]int32, bool)`

GetNetworkAffiliateIdsOk returns a tuple with the NetworkAffiliateIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkAffiliateIds

`func (o *InternalRedirect) SetNetworkAffiliateIds(v []int32)`

SetNetworkAffiliateIds sets NetworkAffiliateIds field to given value.

### HasNetworkAffiliateIds

`func (o *InternalRedirect) HasNetworkAffiliateIds() bool`

HasNetworkAffiliateIds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


