# PayoutRevenue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**EntryName** | Pointer to **string** | Event name | [optional] [default to "Base"]
**PayoutType** | **string** | Payout type | 
**PayoutAmount** | Pointer to **float64** | Payout amount | [optional] 
**PayoutPercentage** | Pointer to **int32** | Payout percentage | [optional] 
**RevenueType** | **string** | Revenue type | 
**RevenueAmount** | Pointer to **float64** | Revenue amount | [optional] 
**RevenuePercentage** | Pointer to **int32** | Revenue percentage | [optional] 
**IsDefault** | **bool** | Is base conversion | 
**IsPrivate** | **bool** | Is private event | 
**IsPostbackDisabled** | Pointer to **bool** | Disable partner postback | [optional] [default to false]
**GlobalAdvertiserEventId** | Pointer to **int32** | Global advertiser event ID | [optional] [default to 0]
**IsMustApproveConversion** | Pointer to **bool** | Require conversion approval | [optional] [default to false]
**IsAllowDuplicateConversion** | Pointer to **bool** | Allow duplicate conversions | [optional] [default to true]

## Methods

### NewPayoutRevenue

`func NewPayoutRevenue(payoutType string, revenueType string, isDefault bool, isPrivate bool, ) *PayoutRevenue`

NewPayoutRevenue instantiates a new PayoutRevenue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPayoutRevenueWithDefaults

`func NewPayoutRevenueWithDefaults() *PayoutRevenue`

NewPayoutRevenueWithDefaults instantiates a new PayoutRevenue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEntryName

`func (o *PayoutRevenue) GetEntryName() string`

GetEntryName returns the EntryName field if non-nil, zero value otherwise.

### GetEntryNameOk

`func (o *PayoutRevenue) GetEntryNameOk() (*string, bool)`

GetEntryNameOk returns a tuple with the EntryName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEntryName

`func (o *PayoutRevenue) SetEntryName(v string)`

SetEntryName sets EntryName field to given value.

### HasEntryName

`func (o *PayoutRevenue) HasEntryName() bool`

HasEntryName returns a boolean if a field has been set.

### GetPayoutType

`func (o *PayoutRevenue) GetPayoutType() string`

GetPayoutType returns the PayoutType field if non-nil, zero value otherwise.

### GetPayoutTypeOk

`func (o *PayoutRevenue) GetPayoutTypeOk() (*string, bool)`

GetPayoutTypeOk returns a tuple with the PayoutType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPayoutType

`func (o *PayoutRevenue) SetPayoutType(v string)`

SetPayoutType sets PayoutType field to given value.


### GetPayoutAmount

`func (o *PayoutRevenue) GetPayoutAmount() float64`

GetPayoutAmount returns the PayoutAmount field if non-nil, zero value otherwise.

### GetPayoutAmountOk

`func (o *PayoutRevenue) GetPayoutAmountOk() (*float64, bool)`

GetPayoutAmountOk returns a tuple with the PayoutAmount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPayoutAmount

`func (o *PayoutRevenue) SetPayoutAmount(v float64)`

SetPayoutAmount sets PayoutAmount field to given value.

### HasPayoutAmount

`func (o *PayoutRevenue) HasPayoutAmount() bool`

HasPayoutAmount returns a boolean if a field has been set.

### GetPayoutPercentage

`func (o *PayoutRevenue) GetPayoutPercentage() int32`

GetPayoutPercentage returns the PayoutPercentage field if non-nil, zero value otherwise.

### GetPayoutPercentageOk

`func (o *PayoutRevenue) GetPayoutPercentageOk() (*int32, bool)`

GetPayoutPercentageOk returns a tuple with the PayoutPercentage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPayoutPercentage

`func (o *PayoutRevenue) SetPayoutPercentage(v int32)`

SetPayoutPercentage sets PayoutPercentage field to given value.

### HasPayoutPercentage

`func (o *PayoutRevenue) HasPayoutPercentage() bool`

HasPayoutPercentage returns a boolean if a field has been set.

### GetRevenueType

`func (o *PayoutRevenue) GetRevenueType() string`

GetRevenueType returns the RevenueType field if non-nil, zero value otherwise.

### GetRevenueTypeOk

`func (o *PayoutRevenue) GetRevenueTypeOk() (*string, bool)`

GetRevenueTypeOk returns a tuple with the RevenueType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevenueType

`func (o *PayoutRevenue) SetRevenueType(v string)`

SetRevenueType sets RevenueType field to given value.


### GetRevenueAmount

`func (o *PayoutRevenue) GetRevenueAmount() float64`

GetRevenueAmount returns the RevenueAmount field if non-nil, zero value otherwise.

### GetRevenueAmountOk

`func (o *PayoutRevenue) GetRevenueAmountOk() (*float64, bool)`

GetRevenueAmountOk returns a tuple with the RevenueAmount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevenueAmount

`func (o *PayoutRevenue) SetRevenueAmount(v float64)`

SetRevenueAmount sets RevenueAmount field to given value.

### HasRevenueAmount

`func (o *PayoutRevenue) HasRevenueAmount() bool`

HasRevenueAmount returns a boolean if a field has been set.

### GetRevenuePercentage

`func (o *PayoutRevenue) GetRevenuePercentage() int32`

GetRevenuePercentage returns the RevenuePercentage field if non-nil, zero value otherwise.

### GetRevenuePercentageOk

`func (o *PayoutRevenue) GetRevenuePercentageOk() (*int32, bool)`

GetRevenuePercentageOk returns a tuple with the RevenuePercentage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevenuePercentage

`func (o *PayoutRevenue) SetRevenuePercentage(v int32)`

SetRevenuePercentage sets RevenuePercentage field to given value.

### HasRevenuePercentage

`func (o *PayoutRevenue) HasRevenuePercentage() bool`

HasRevenuePercentage returns a boolean if a field has been set.

### GetIsDefault

`func (o *PayoutRevenue) GetIsDefault() bool`

GetIsDefault returns the IsDefault field if non-nil, zero value otherwise.

### GetIsDefaultOk

`func (o *PayoutRevenue) GetIsDefaultOk() (*bool, bool)`

GetIsDefaultOk returns a tuple with the IsDefault field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsDefault

`func (o *PayoutRevenue) SetIsDefault(v bool)`

SetIsDefault sets IsDefault field to given value.


### GetIsPrivate

`func (o *PayoutRevenue) GetIsPrivate() bool`

GetIsPrivate returns the IsPrivate field if non-nil, zero value otherwise.

### GetIsPrivateOk

`func (o *PayoutRevenue) GetIsPrivateOk() (*bool, bool)`

GetIsPrivateOk returns a tuple with the IsPrivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPrivate

`func (o *PayoutRevenue) SetIsPrivate(v bool)`

SetIsPrivate sets IsPrivate field to given value.


### GetIsPostbackDisabled

`func (o *PayoutRevenue) GetIsPostbackDisabled() bool`

GetIsPostbackDisabled returns the IsPostbackDisabled field if non-nil, zero value otherwise.

### GetIsPostbackDisabledOk

`func (o *PayoutRevenue) GetIsPostbackDisabledOk() (*bool, bool)`

GetIsPostbackDisabledOk returns a tuple with the IsPostbackDisabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPostbackDisabled

`func (o *PayoutRevenue) SetIsPostbackDisabled(v bool)`

SetIsPostbackDisabled sets IsPostbackDisabled field to given value.

### HasIsPostbackDisabled

`func (o *PayoutRevenue) HasIsPostbackDisabled() bool`

HasIsPostbackDisabled returns a boolean if a field has been set.

### GetGlobalAdvertiserEventId

`func (o *PayoutRevenue) GetGlobalAdvertiserEventId() int32`

GetGlobalAdvertiserEventId returns the GlobalAdvertiserEventId field if non-nil, zero value otherwise.

### GetGlobalAdvertiserEventIdOk

`func (o *PayoutRevenue) GetGlobalAdvertiserEventIdOk() (*int32, bool)`

GetGlobalAdvertiserEventIdOk returns a tuple with the GlobalAdvertiserEventId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGlobalAdvertiserEventId

`func (o *PayoutRevenue) SetGlobalAdvertiserEventId(v int32)`

SetGlobalAdvertiserEventId sets GlobalAdvertiserEventId field to given value.

### HasGlobalAdvertiserEventId

`func (o *PayoutRevenue) HasGlobalAdvertiserEventId() bool`

HasGlobalAdvertiserEventId returns a boolean if a field has been set.

### GetIsMustApproveConversion

`func (o *PayoutRevenue) GetIsMustApproveConversion() bool`

GetIsMustApproveConversion returns the IsMustApproveConversion field if non-nil, zero value otherwise.

### GetIsMustApproveConversionOk

`func (o *PayoutRevenue) GetIsMustApproveConversionOk() (*bool, bool)`

GetIsMustApproveConversionOk returns a tuple with the IsMustApproveConversion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsMustApproveConversion

`func (o *PayoutRevenue) SetIsMustApproveConversion(v bool)`

SetIsMustApproveConversion sets IsMustApproveConversion field to given value.

### HasIsMustApproveConversion

`func (o *PayoutRevenue) HasIsMustApproveConversion() bool`

HasIsMustApproveConversion returns a boolean if a field has been set.

### GetIsAllowDuplicateConversion

`func (o *PayoutRevenue) GetIsAllowDuplicateConversion() bool`

GetIsAllowDuplicateConversion returns the IsAllowDuplicateConversion field if non-nil, zero value otherwise.

### GetIsAllowDuplicateConversionOk

`func (o *PayoutRevenue) GetIsAllowDuplicateConversionOk() (*bool, bool)`

GetIsAllowDuplicateConversionOk returns a tuple with the IsAllowDuplicateConversion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsAllowDuplicateConversion

`func (o *PayoutRevenue) SetIsAllowDuplicateConversion(v bool)`

SetIsAllowDuplicateConversion sets IsAllowDuplicateConversion field to given value.

### HasIsAllowDuplicateConversion

`func (o *PayoutRevenue) HasIsAllowDuplicateConversion() bool`

HasIsAllowDuplicateConversion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


