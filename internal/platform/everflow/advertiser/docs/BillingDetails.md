# BillingDetails

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DayOfWeek** | Pointer to **int32** | Day of the week (for weekly frequency) | [optional] 
**DayOfMonthOne** | Pointer to **int32** | First day of the month (for bimonthly frequency) | [optional] 
**DayOfMonthTwo** | Pointer to **int32** | Second day of the month (for bimonthly frequency) | [optional] 
**DayOfMonth** | Pointer to **int32** | Day of the month (for monthly, two_months, quarterly) | [optional] 
**StartingMonth** | Pointer to **int32** | Starting month for cycle (for two_months, quarterly) | [optional] 

## Methods

### NewBillingDetails

`func NewBillingDetails() *BillingDetails`

NewBillingDetails instantiates a new BillingDetails object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBillingDetailsWithDefaults

`func NewBillingDetailsWithDefaults() *BillingDetails`

NewBillingDetailsWithDefaults instantiates a new BillingDetails object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDayOfWeek

`func (o *BillingDetails) GetDayOfWeek() int32`

GetDayOfWeek returns the DayOfWeek field if non-nil, zero value otherwise.

### GetDayOfWeekOk

`func (o *BillingDetails) GetDayOfWeekOk() (*int32, bool)`

GetDayOfWeekOk returns a tuple with the DayOfWeek field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDayOfWeek

`func (o *BillingDetails) SetDayOfWeek(v int32)`

SetDayOfWeek sets DayOfWeek field to given value.

### HasDayOfWeek

`func (o *BillingDetails) HasDayOfWeek() bool`

HasDayOfWeek returns a boolean if a field has been set.

### GetDayOfMonthOne

`func (o *BillingDetails) GetDayOfMonthOne() int32`

GetDayOfMonthOne returns the DayOfMonthOne field if non-nil, zero value otherwise.

### GetDayOfMonthOneOk

`func (o *BillingDetails) GetDayOfMonthOneOk() (*int32, bool)`

GetDayOfMonthOneOk returns a tuple with the DayOfMonthOne field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDayOfMonthOne

`func (o *BillingDetails) SetDayOfMonthOne(v int32)`

SetDayOfMonthOne sets DayOfMonthOne field to given value.

### HasDayOfMonthOne

`func (o *BillingDetails) HasDayOfMonthOne() bool`

HasDayOfMonthOne returns a boolean if a field has been set.

### GetDayOfMonthTwo

`func (o *BillingDetails) GetDayOfMonthTwo() int32`

GetDayOfMonthTwo returns the DayOfMonthTwo field if non-nil, zero value otherwise.

### GetDayOfMonthTwoOk

`func (o *BillingDetails) GetDayOfMonthTwoOk() (*int32, bool)`

GetDayOfMonthTwoOk returns a tuple with the DayOfMonthTwo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDayOfMonthTwo

`func (o *BillingDetails) SetDayOfMonthTwo(v int32)`

SetDayOfMonthTwo sets DayOfMonthTwo field to given value.

### HasDayOfMonthTwo

`func (o *BillingDetails) HasDayOfMonthTwo() bool`

HasDayOfMonthTwo returns a boolean if a field has been set.

### GetDayOfMonth

`func (o *BillingDetails) GetDayOfMonth() int32`

GetDayOfMonth returns the DayOfMonth field if non-nil, zero value otherwise.

### GetDayOfMonthOk

`func (o *BillingDetails) GetDayOfMonthOk() (*int32, bool)`

GetDayOfMonthOk returns a tuple with the DayOfMonth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDayOfMonth

`func (o *BillingDetails) SetDayOfMonth(v int32)`

SetDayOfMonth sets DayOfMonth field to given value.

### HasDayOfMonth

`func (o *BillingDetails) HasDayOfMonth() bool`

HasDayOfMonth returns a boolean if a field has been set.

### GetStartingMonth

`func (o *BillingDetails) GetStartingMonth() int32`

GetStartingMonth returns the StartingMonth field if non-nil, zero value otherwise.

### GetStartingMonthOk

`func (o *BillingDetails) GetStartingMonthOk() (*int32, bool)`

GetStartingMonthOk returns a tuple with the StartingMonth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartingMonth

`func (o *BillingDetails) SetStartingMonth(v int32)`

SetStartingMonth sets StartingMonth field to given value.

### HasStartingMonth

`func (o *BillingDetails) HasStartingMonth() bool`

HasStartingMonth returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


