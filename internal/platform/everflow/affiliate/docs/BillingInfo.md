# BillingInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BillingFrequency** | Pointer to **string** | The affiliate&#39;s invoicing frequency | [optional] 
**PaymentType** | Pointer to **string** | The affiliate&#39;s payment type | [optional] 
**TaxId** | Pointer to **string** | The affiliate&#39;s tax id | [optional] 
**IsInvoiceCreationAuto** | Pointer to **bool** | Configures automatic invoice creations | [optional] [default to false]
**AutoInvoiceStartDate** | Pointer to **string** | Automatic invoice creation start date (YYYY-mm-dd) | [optional] 
**DefaultInvoiceIsHidden** | Pointer to **bool** | Whether invoices are hidden from the affiliate by default | [optional] [default to false]
**InvoiceGenerationDaysDelay** | Pointer to **int32** | Days to wait for invoice generation after billing period | [optional] [default to 0]
**InvoiceAmountThreshold** | Pointer to **float64** | Minimal amount required for invoice generation | [optional] [default to 0]
**DefaultPaymentTerms** | Pointer to **int32** | Number of days for payment terms on invoices | [optional] [default to 0]
**Details** | Pointer to [**BillingDetails**](BillingDetails.md) |  | [optional] 
**Payment** | Pointer to [**PaymentDetails**](PaymentDetails.md) |  | [optional] 

## Methods

### NewBillingInfo

`func NewBillingInfo() *BillingInfo`

NewBillingInfo instantiates a new BillingInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBillingInfoWithDefaults

`func NewBillingInfoWithDefaults() *BillingInfo`

NewBillingInfoWithDefaults instantiates a new BillingInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBillingFrequency

`func (o *BillingInfo) GetBillingFrequency() string`

GetBillingFrequency returns the BillingFrequency field if non-nil, zero value otherwise.

### GetBillingFrequencyOk

`func (o *BillingInfo) GetBillingFrequencyOk() (*string, bool)`

GetBillingFrequencyOk returns a tuple with the BillingFrequency field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBillingFrequency

`func (o *BillingInfo) SetBillingFrequency(v string)`

SetBillingFrequency sets BillingFrequency field to given value.

### HasBillingFrequency

`func (o *BillingInfo) HasBillingFrequency() bool`

HasBillingFrequency returns a boolean if a field has been set.

### GetPaymentType

`func (o *BillingInfo) GetPaymentType() string`

GetPaymentType returns the PaymentType field if non-nil, zero value otherwise.

### GetPaymentTypeOk

`func (o *BillingInfo) GetPaymentTypeOk() (*string, bool)`

GetPaymentTypeOk returns a tuple with the PaymentType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPaymentType

`func (o *BillingInfo) SetPaymentType(v string)`

SetPaymentType sets PaymentType field to given value.

### HasPaymentType

`func (o *BillingInfo) HasPaymentType() bool`

HasPaymentType returns a boolean if a field has been set.

### GetTaxId

`func (o *BillingInfo) GetTaxId() string`

GetTaxId returns the TaxId field if non-nil, zero value otherwise.

### GetTaxIdOk

`func (o *BillingInfo) GetTaxIdOk() (*string, bool)`

GetTaxIdOk returns a tuple with the TaxId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTaxId

`func (o *BillingInfo) SetTaxId(v string)`

SetTaxId sets TaxId field to given value.

### HasTaxId

`func (o *BillingInfo) HasTaxId() bool`

HasTaxId returns a boolean if a field has been set.

### GetIsInvoiceCreationAuto

`func (o *BillingInfo) GetIsInvoiceCreationAuto() bool`

GetIsInvoiceCreationAuto returns the IsInvoiceCreationAuto field if non-nil, zero value otherwise.

### GetIsInvoiceCreationAutoOk

`func (o *BillingInfo) GetIsInvoiceCreationAutoOk() (*bool, bool)`

GetIsInvoiceCreationAutoOk returns a tuple with the IsInvoiceCreationAuto field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsInvoiceCreationAuto

`func (o *BillingInfo) SetIsInvoiceCreationAuto(v bool)`

SetIsInvoiceCreationAuto sets IsInvoiceCreationAuto field to given value.

### HasIsInvoiceCreationAuto

`func (o *BillingInfo) HasIsInvoiceCreationAuto() bool`

HasIsInvoiceCreationAuto returns a boolean if a field has been set.

### GetAutoInvoiceStartDate

`func (o *BillingInfo) GetAutoInvoiceStartDate() string`

GetAutoInvoiceStartDate returns the AutoInvoiceStartDate field if non-nil, zero value otherwise.

### GetAutoInvoiceStartDateOk

`func (o *BillingInfo) GetAutoInvoiceStartDateOk() (*string, bool)`

GetAutoInvoiceStartDateOk returns a tuple with the AutoInvoiceStartDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAutoInvoiceStartDate

`func (o *BillingInfo) SetAutoInvoiceStartDate(v string)`

SetAutoInvoiceStartDate sets AutoInvoiceStartDate field to given value.

### HasAutoInvoiceStartDate

`func (o *BillingInfo) HasAutoInvoiceStartDate() bool`

HasAutoInvoiceStartDate returns a boolean if a field has been set.

### GetDefaultInvoiceIsHidden

`func (o *BillingInfo) GetDefaultInvoiceIsHidden() bool`

GetDefaultInvoiceIsHidden returns the DefaultInvoiceIsHidden field if non-nil, zero value otherwise.

### GetDefaultInvoiceIsHiddenOk

`func (o *BillingInfo) GetDefaultInvoiceIsHiddenOk() (*bool, bool)`

GetDefaultInvoiceIsHiddenOk returns a tuple with the DefaultInvoiceIsHidden field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultInvoiceIsHidden

`func (o *BillingInfo) SetDefaultInvoiceIsHidden(v bool)`

SetDefaultInvoiceIsHidden sets DefaultInvoiceIsHidden field to given value.

### HasDefaultInvoiceIsHidden

`func (o *BillingInfo) HasDefaultInvoiceIsHidden() bool`

HasDefaultInvoiceIsHidden returns a boolean if a field has been set.

### GetInvoiceGenerationDaysDelay

`func (o *BillingInfo) GetInvoiceGenerationDaysDelay() int32`

GetInvoiceGenerationDaysDelay returns the InvoiceGenerationDaysDelay field if non-nil, zero value otherwise.

### GetInvoiceGenerationDaysDelayOk

`func (o *BillingInfo) GetInvoiceGenerationDaysDelayOk() (*int32, bool)`

GetInvoiceGenerationDaysDelayOk returns a tuple with the InvoiceGenerationDaysDelay field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInvoiceGenerationDaysDelay

`func (o *BillingInfo) SetInvoiceGenerationDaysDelay(v int32)`

SetInvoiceGenerationDaysDelay sets InvoiceGenerationDaysDelay field to given value.

### HasInvoiceGenerationDaysDelay

`func (o *BillingInfo) HasInvoiceGenerationDaysDelay() bool`

HasInvoiceGenerationDaysDelay returns a boolean if a field has been set.

### GetInvoiceAmountThreshold

`func (o *BillingInfo) GetInvoiceAmountThreshold() float64`

GetInvoiceAmountThreshold returns the InvoiceAmountThreshold field if non-nil, zero value otherwise.

### GetInvoiceAmountThresholdOk

`func (o *BillingInfo) GetInvoiceAmountThresholdOk() (*float64, bool)`

GetInvoiceAmountThresholdOk returns a tuple with the InvoiceAmountThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInvoiceAmountThreshold

`func (o *BillingInfo) SetInvoiceAmountThreshold(v float64)`

SetInvoiceAmountThreshold sets InvoiceAmountThreshold field to given value.

### HasInvoiceAmountThreshold

`func (o *BillingInfo) HasInvoiceAmountThreshold() bool`

HasInvoiceAmountThreshold returns a boolean if a field has been set.

### GetDefaultPaymentTerms

`func (o *BillingInfo) GetDefaultPaymentTerms() int32`

GetDefaultPaymentTerms returns the DefaultPaymentTerms field if non-nil, zero value otherwise.

### GetDefaultPaymentTermsOk

`func (o *BillingInfo) GetDefaultPaymentTermsOk() (*int32, bool)`

GetDefaultPaymentTermsOk returns a tuple with the DefaultPaymentTerms field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultPaymentTerms

`func (o *BillingInfo) SetDefaultPaymentTerms(v int32)`

SetDefaultPaymentTerms sets DefaultPaymentTerms field to given value.

### HasDefaultPaymentTerms

`func (o *BillingInfo) HasDefaultPaymentTerms() bool`

HasDefaultPaymentTerms returns a boolean if a field has been set.

### GetDetails

`func (o *BillingInfo) GetDetails() BillingDetails`

GetDetails returns the Details field if non-nil, zero value otherwise.

### GetDetailsOk

`func (o *BillingInfo) GetDetailsOk() (*BillingDetails, bool)`

GetDetailsOk returns a tuple with the Details field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetails

`func (o *BillingInfo) SetDetails(v BillingDetails)`

SetDetails sets Details field to given value.

### HasDetails

`func (o *BillingInfo) HasDetails() bool`

HasDetails returns a boolean if a field has been set.

### GetPayment

`func (o *BillingInfo) GetPayment() PaymentDetails`

GetPayment returns the Payment field if non-nil, zero value otherwise.

### GetPaymentOk

`func (o *BillingInfo) GetPaymentOk() (*PaymentDetails, bool)`

GetPaymentOk returns a tuple with the Payment field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPayment

`func (o *BillingInfo) SetPayment(v PaymentDetails)`

SetPayment sets Payment field to given value.

### HasPayment

`func (o *BillingInfo) HasPayment() bool`

HasPayment returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


