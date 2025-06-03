# PaymentDetails

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PaxumId** | Pointer to **string** | The affiliate&#39;s paxum id (when payment_type is paxum) | [optional] 
**PaypalReceptionMethod** | Pointer to **string** | Reception method for PayPal | [optional] 
**ReceptionIdentifier** | Pointer to **string** | PayPal reception identifier | [optional] 
**Email** | Pointer to **string** | Email for payoneer or veem | [optional] 
**IsExistingPayee** | Pointer to **bool** | Whether to assign existing payee id (tipalti) | [optional] 
**Idap** | Pointer to **string** | Payee&#39;s IDAP (Payee ID) for tipalti | [optional] 
**FirstName** | Pointer to **string** | First name for veem | [optional] 
**LastName** | Pointer to **string** | Last name for veem | [optional] 
**Phone** | Pointer to **string** | Phone number in international format (veem) | [optional] 
**CountryIso** | Pointer to **string** | Country ISO code (veem) | [optional] 
**BankName** | Pointer to **string** | Bank name (wire/direct_deposit) | [optional] 
**BankAddress** | Pointer to **string** | Bank address (wire/direct_deposit) | [optional] 
**AccountName** | Pointer to **string** | Bank account name (wire/direct_deposit) | [optional] 
**AccountNumber** | Pointer to **string** | Bank account number (wire/direct_deposit) | [optional] 
**RoutingNumber** | Pointer to **string** | Bank routing number (wire/direct_deposit) | [optional] 
**SwiftCode** | Pointer to **string** | SWIFT code (wire/direct_deposit) | [optional] 

## Methods

### NewPaymentDetails

`func NewPaymentDetails() *PaymentDetails`

NewPaymentDetails instantiates a new PaymentDetails object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPaymentDetailsWithDefaults

`func NewPaymentDetailsWithDefaults() *PaymentDetails`

NewPaymentDetailsWithDefaults instantiates a new PaymentDetails object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPaxumId

`func (o *PaymentDetails) GetPaxumId() string`

GetPaxumId returns the PaxumId field if non-nil, zero value otherwise.

### GetPaxumIdOk

`func (o *PaymentDetails) GetPaxumIdOk() (*string, bool)`

GetPaxumIdOk returns a tuple with the PaxumId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPaxumId

`func (o *PaymentDetails) SetPaxumId(v string)`

SetPaxumId sets PaxumId field to given value.

### HasPaxumId

`func (o *PaymentDetails) HasPaxumId() bool`

HasPaxumId returns a boolean if a field has been set.

### GetPaypalReceptionMethod

`func (o *PaymentDetails) GetPaypalReceptionMethod() string`

GetPaypalReceptionMethod returns the PaypalReceptionMethod field if non-nil, zero value otherwise.

### GetPaypalReceptionMethodOk

`func (o *PaymentDetails) GetPaypalReceptionMethodOk() (*string, bool)`

GetPaypalReceptionMethodOk returns a tuple with the PaypalReceptionMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPaypalReceptionMethod

`func (o *PaymentDetails) SetPaypalReceptionMethod(v string)`

SetPaypalReceptionMethod sets PaypalReceptionMethod field to given value.

### HasPaypalReceptionMethod

`func (o *PaymentDetails) HasPaypalReceptionMethod() bool`

HasPaypalReceptionMethod returns a boolean if a field has been set.

### GetReceptionIdentifier

`func (o *PaymentDetails) GetReceptionIdentifier() string`

GetReceptionIdentifier returns the ReceptionIdentifier field if non-nil, zero value otherwise.

### GetReceptionIdentifierOk

`func (o *PaymentDetails) GetReceptionIdentifierOk() (*string, bool)`

GetReceptionIdentifierOk returns a tuple with the ReceptionIdentifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReceptionIdentifier

`func (o *PaymentDetails) SetReceptionIdentifier(v string)`

SetReceptionIdentifier sets ReceptionIdentifier field to given value.

### HasReceptionIdentifier

`func (o *PaymentDetails) HasReceptionIdentifier() bool`

HasReceptionIdentifier returns a boolean if a field has been set.

### GetEmail

`func (o *PaymentDetails) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *PaymentDetails) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *PaymentDetails) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *PaymentDetails) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetIsExistingPayee

`func (o *PaymentDetails) GetIsExistingPayee() bool`

GetIsExistingPayee returns the IsExistingPayee field if non-nil, zero value otherwise.

### GetIsExistingPayeeOk

`func (o *PaymentDetails) GetIsExistingPayeeOk() (*bool, bool)`

GetIsExistingPayeeOk returns a tuple with the IsExistingPayee field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsExistingPayee

`func (o *PaymentDetails) SetIsExistingPayee(v bool)`

SetIsExistingPayee sets IsExistingPayee field to given value.

### HasIsExistingPayee

`func (o *PaymentDetails) HasIsExistingPayee() bool`

HasIsExistingPayee returns a boolean if a field has been set.

### GetIdap

`func (o *PaymentDetails) GetIdap() string`

GetIdap returns the Idap field if non-nil, zero value otherwise.

### GetIdapOk

`func (o *PaymentDetails) GetIdapOk() (*string, bool)`

GetIdapOk returns a tuple with the Idap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdap

`func (o *PaymentDetails) SetIdap(v string)`

SetIdap sets Idap field to given value.

### HasIdap

`func (o *PaymentDetails) HasIdap() bool`

HasIdap returns a boolean if a field has been set.

### GetFirstName

`func (o *PaymentDetails) GetFirstName() string`

GetFirstName returns the FirstName field if non-nil, zero value otherwise.

### GetFirstNameOk

`func (o *PaymentDetails) GetFirstNameOk() (*string, bool)`

GetFirstNameOk returns a tuple with the FirstName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirstName

`func (o *PaymentDetails) SetFirstName(v string)`

SetFirstName sets FirstName field to given value.

### HasFirstName

`func (o *PaymentDetails) HasFirstName() bool`

HasFirstName returns a boolean if a field has been set.

### GetLastName

`func (o *PaymentDetails) GetLastName() string`

GetLastName returns the LastName field if non-nil, zero value otherwise.

### GetLastNameOk

`func (o *PaymentDetails) GetLastNameOk() (*string, bool)`

GetLastNameOk returns a tuple with the LastName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastName

`func (o *PaymentDetails) SetLastName(v string)`

SetLastName sets LastName field to given value.

### HasLastName

`func (o *PaymentDetails) HasLastName() bool`

HasLastName returns a boolean if a field has been set.

### GetPhone

`func (o *PaymentDetails) GetPhone() string`

GetPhone returns the Phone field if non-nil, zero value otherwise.

### GetPhoneOk

`func (o *PaymentDetails) GetPhoneOk() (*string, bool)`

GetPhoneOk returns a tuple with the Phone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhone

`func (o *PaymentDetails) SetPhone(v string)`

SetPhone sets Phone field to given value.

### HasPhone

`func (o *PaymentDetails) HasPhone() bool`

HasPhone returns a boolean if a field has been set.

### GetCountryIso

`func (o *PaymentDetails) GetCountryIso() string`

GetCountryIso returns the CountryIso field if non-nil, zero value otherwise.

### GetCountryIsoOk

`func (o *PaymentDetails) GetCountryIsoOk() (*string, bool)`

GetCountryIsoOk returns a tuple with the CountryIso field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryIso

`func (o *PaymentDetails) SetCountryIso(v string)`

SetCountryIso sets CountryIso field to given value.

### HasCountryIso

`func (o *PaymentDetails) HasCountryIso() bool`

HasCountryIso returns a boolean if a field has been set.

### GetBankName

`func (o *PaymentDetails) GetBankName() string`

GetBankName returns the BankName field if non-nil, zero value otherwise.

### GetBankNameOk

`func (o *PaymentDetails) GetBankNameOk() (*string, bool)`

GetBankNameOk returns a tuple with the BankName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBankName

`func (o *PaymentDetails) SetBankName(v string)`

SetBankName sets BankName field to given value.

### HasBankName

`func (o *PaymentDetails) HasBankName() bool`

HasBankName returns a boolean if a field has been set.

### GetBankAddress

`func (o *PaymentDetails) GetBankAddress() string`

GetBankAddress returns the BankAddress field if non-nil, zero value otherwise.

### GetBankAddressOk

`func (o *PaymentDetails) GetBankAddressOk() (*string, bool)`

GetBankAddressOk returns a tuple with the BankAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBankAddress

`func (o *PaymentDetails) SetBankAddress(v string)`

SetBankAddress sets BankAddress field to given value.

### HasBankAddress

`func (o *PaymentDetails) HasBankAddress() bool`

HasBankAddress returns a boolean if a field has been set.

### GetAccountName

`func (o *PaymentDetails) GetAccountName() string`

GetAccountName returns the AccountName field if non-nil, zero value otherwise.

### GetAccountNameOk

`func (o *PaymentDetails) GetAccountNameOk() (*string, bool)`

GetAccountNameOk returns a tuple with the AccountName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountName

`func (o *PaymentDetails) SetAccountName(v string)`

SetAccountName sets AccountName field to given value.

### HasAccountName

`func (o *PaymentDetails) HasAccountName() bool`

HasAccountName returns a boolean if a field has been set.

### GetAccountNumber

`func (o *PaymentDetails) GetAccountNumber() string`

GetAccountNumber returns the AccountNumber field if non-nil, zero value otherwise.

### GetAccountNumberOk

`func (o *PaymentDetails) GetAccountNumberOk() (*string, bool)`

GetAccountNumberOk returns a tuple with the AccountNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountNumber

`func (o *PaymentDetails) SetAccountNumber(v string)`

SetAccountNumber sets AccountNumber field to given value.

### HasAccountNumber

`func (o *PaymentDetails) HasAccountNumber() bool`

HasAccountNumber returns a boolean if a field has been set.

### GetRoutingNumber

`func (o *PaymentDetails) GetRoutingNumber() string`

GetRoutingNumber returns the RoutingNumber field if non-nil, zero value otherwise.

### GetRoutingNumberOk

`func (o *PaymentDetails) GetRoutingNumberOk() (*string, bool)`

GetRoutingNumberOk returns a tuple with the RoutingNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoutingNumber

`func (o *PaymentDetails) SetRoutingNumber(v string)`

SetRoutingNumber sets RoutingNumber field to given value.

### HasRoutingNumber

`func (o *PaymentDetails) HasRoutingNumber() bool`

HasRoutingNumber returns a boolean if a field has been set.

### GetSwiftCode

`func (o *PaymentDetails) GetSwiftCode() string`

GetSwiftCode returns the SwiftCode field if non-nil, zero value otherwise.

### GetSwiftCodeOk

`func (o *PaymentDetails) GetSwiftCodeOk() (*string, bool)`

GetSwiftCodeOk returns a tuple with the SwiftCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwiftCode

`func (o *PaymentDetails) SetSwiftCode(v string)`

SetSwiftCode sets SwiftCode field to given value.

### HasSwiftCode

`func (o *PaymentDetails) HasSwiftCode() bool`

HasSwiftCode returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


