# Billing

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BillingFrequency** | Pointer to **string** | The advertiser&#39;s invoicing frequency | [optional] 
**InvoiceAmountThreshold** | Pointer to **float64** | Minimal amount required for invoice generation | [optional] [default to 0]
**TaxId** | Pointer to **string** | The advertiser&#39;s tax id | [optional] 
**IsInvoiceCreationAuto** | Pointer to **bool** | Configures automatic invoice creations | [optional] [default to false]
**AutoInvoiceStartDate** | Pointer to **string** | Automatic invoice creation start date (YYYY-mm-dd) | [optional] 
**DefaultInvoiceIsHidden** | Pointer to **bool** | Whether invoices are hidden from advertiser by default | [optional] [default to false]
**InvoiceGenerationDaysDelay** | Pointer to **int32** | Days to wait for invoice generation after billing period | [optional] [default to 0]
**DefaultPaymentTerms** | Pointer to **int32** | Number of days for payment terms on invoices | [optional] [default to 0]
**Details** | Pointer to [**BillingDetails**](BillingDetails.md) |  | [optional] 

## Methods

### NewBilling

`func NewBilling() *Billing`

NewBilling instantiates a new Billing object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBillingWithDefaults

`func NewBillingWithDefaults() *Billing`

NewBillingWithDefaults instantiates a new Billing object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBillingFrequency

`func (o *Billing) GetBillingFrequency() string`

GetBillingFrequency returns the BillingFrequency field if non-nil, zero value otherwise.

### GetBillingFrequencyOk

`func (o *Billing) GetBillingFrequencyOk() (*string, bool)`

GetBillingFrequencyOk returns a tuple with the BillingFrequency field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBillingFrequency

`func (o *Billing) SetBillingFrequency(v string)`

SetBillingFrequency sets BillingFrequency field to given value.

### HasBillingFrequency

`func (o *Billing) HasBillingFrequency() bool`

HasBillingFrequency returns a boolean if a field has been set.

### GetInvoiceAmountThreshold

`func (o *Billing) GetInvoiceAmountThreshold() float64`

GetInvoiceAmountThreshold returns the InvoiceAmountThreshold field if non-nil, zero value otherwise.

### GetInvoiceAmountThresholdOk

`func (o *Billing) GetInvoiceAmountThresholdOk() (*float64, bool)`

GetInvoiceAmountThresholdOk returns a tuple with the InvoiceAmountThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInvoiceAmountThreshold

`func (o *Billing) SetInvoiceAmountThreshold(v float64)`

SetInvoiceAmountThreshold sets InvoiceAmountThreshold field to given value.

### HasInvoiceAmountThreshold

`func (o *Billing) HasInvoiceAmountThreshold() bool`

HasInvoiceAmountThreshold returns a boolean if a field has been set.

### GetTaxId

`func (o *Billing) GetTaxId() string`

GetTaxId returns the TaxId field if non-nil, zero value otherwise.

### GetTaxIdOk

`func (o *Billing) GetTaxIdOk() (*string, bool)`

GetTaxIdOk returns a tuple with the TaxId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTaxId

`func (o *Billing) SetTaxId(v string)`

SetTaxId sets TaxId field to given value.

### HasTaxId

`func (o *Billing) HasTaxId() bool`

HasTaxId returns a boolean if a field has been set.

### GetIsInvoiceCreationAuto

`func (o *Billing) GetIsInvoiceCreationAuto() bool`

GetIsInvoiceCreationAuto returns the IsInvoiceCreationAuto field if non-nil, zero value otherwise.

### GetIsInvoiceCreationAutoOk

`func (o *Billing) GetIsInvoiceCreationAutoOk() (*bool, bool)`

GetIsInvoiceCreationAutoOk returns a tuple with the IsInvoiceCreationAuto field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsInvoiceCreationAuto

`func (o *Billing) SetIsInvoiceCreationAuto(v bool)`

SetIsInvoiceCreationAuto sets IsInvoiceCreationAuto field to given value.

### HasIsInvoiceCreationAuto

`func (o *Billing) HasIsInvoiceCreationAuto() bool`

HasIsInvoiceCreationAuto returns a boolean if a field has been set.

### GetAutoInvoiceStartDate

`func (o *Billing) GetAutoInvoiceStartDate() string`

GetAutoInvoiceStartDate returns the AutoInvoiceStartDate field if non-nil, zero value otherwise.

### GetAutoInvoiceStartDateOk

`func (o *Billing) GetAutoInvoiceStartDateOk() (*string, bool)`

GetAutoInvoiceStartDateOk returns a tuple with the AutoInvoiceStartDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAutoInvoiceStartDate

`func (o *Billing) SetAutoInvoiceStartDate(v string)`

SetAutoInvoiceStartDate sets AutoInvoiceStartDate field to given value.

### HasAutoInvoiceStartDate

`func (o *Billing) HasAutoInvoiceStartDate() bool`

HasAutoInvoiceStartDate returns a boolean if a field has been set.

### GetDefaultInvoiceIsHidden

`func (o *Billing) GetDefaultInvoiceIsHidden() bool`

GetDefaultInvoiceIsHidden returns the DefaultInvoiceIsHidden field if non-nil, zero value otherwise.

### GetDefaultInvoiceIsHiddenOk

`func (o *Billing) GetDefaultInvoiceIsHiddenOk() (*bool, bool)`

GetDefaultInvoiceIsHiddenOk returns a tuple with the DefaultInvoiceIsHidden field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultInvoiceIsHidden

`func (o *Billing) SetDefaultInvoiceIsHidden(v bool)`

SetDefaultInvoiceIsHidden sets DefaultInvoiceIsHidden field to given value.

### HasDefaultInvoiceIsHidden

`func (o *Billing) HasDefaultInvoiceIsHidden() bool`

HasDefaultInvoiceIsHidden returns a boolean if a field has been set.

### GetInvoiceGenerationDaysDelay

`func (o *Billing) GetInvoiceGenerationDaysDelay() int32`

GetInvoiceGenerationDaysDelay returns the InvoiceGenerationDaysDelay field if non-nil, zero value otherwise.

### GetInvoiceGenerationDaysDelayOk

`func (o *Billing) GetInvoiceGenerationDaysDelayOk() (*int32, bool)`

GetInvoiceGenerationDaysDelayOk returns a tuple with the InvoiceGenerationDaysDelay field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInvoiceGenerationDaysDelay

`func (o *Billing) SetInvoiceGenerationDaysDelay(v int32)`

SetInvoiceGenerationDaysDelay sets InvoiceGenerationDaysDelay field to given value.

### HasInvoiceGenerationDaysDelay

`func (o *Billing) HasInvoiceGenerationDaysDelay() bool`

HasInvoiceGenerationDaysDelay returns a boolean if a field has been set.

### GetDefaultPaymentTerms

`func (o *Billing) GetDefaultPaymentTerms() int32`

GetDefaultPaymentTerms returns the DefaultPaymentTerms field if non-nil, zero value otherwise.

### GetDefaultPaymentTermsOk

`func (o *Billing) GetDefaultPaymentTermsOk() (*int32, bool)`

GetDefaultPaymentTermsOk returns a tuple with the DefaultPaymentTerms field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultPaymentTerms

`func (o *Billing) SetDefaultPaymentTerms(v int32)`

SetDefaultPaymentTerms sets DefaultPaymentTerms field to given value.

### HasDefaultPaymentTerms

`func (o *Billing) HasDefaultPaymentTerms() bool`

HasDefaultPaymentTerms returns a boolean if a field has been set.

### GetDetails

`func (o *Billing) GetDetails() BillingDetails`

GetDetails returns the Details field if non-nil, zero value otherwise.

### GetDetailsOk

`func (o *Billing) GetDetailsOk() (*BillingDetails, bool)`

GetDetailsOk returns a tuple with the Details field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetails

`func (o *Billing) SetDetails(v BillingDetails)`

SetDetails sets Details field to given value.

### HasDetails

`func (o *Billing) HasDetails() bool`

HasDetails returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


