# AffiliateUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FirstName** | **string** | The affiliate user&#39;s first name | 
**LastName** | **string** | The affiliate user&#39;s last name | 
**Email** | **string** | The affiliate user&#39;s email (must be unique) | 
**AccountStatus** | **string** | The affiliate user&#39;s account status | 
**Title** | Pointer to **string** | The affiliate user&#39;s title | [optional] 
**WorkPhone** | Pointer to **string** | The affiliate user&#39;s work phone number | [optional] 
**CellPhone** | Pointer to **string** | The affiliate user&#39;s cell phone number | [optional] 
**InstantMessagingId** | Pointer to **int32** | The id of an instant messaging platform | [optional] 
**InstantMessagingIdentifier** | Pointer to **string** | The affiliate user&#39;s instant messaging identifier | [optional] 
**LanguageId** | Pointer to **int32** | The affiliate user&#39;s language id (1 &#x3D; English) | [optional] 
**TimezoneId** | Pointer to **int32** | The affiliate user&#39;s timezone id | [optional] 
**CurrencyId** | Pointer to **string** | The affiliate user&#39;s currency id | [optional] 
**InitialPassword** | Pointer to **string** | The affiliate user&#39;s login password (min 8 chars, 1 non-alphanumeric, 1 uppercase, 1 lowercase) | [optional] 

## Methods

### NewAffiliateUser

`func NewAffiliateUser(firstName string, lastName string, email string, accountStatus string, ) *AffiliateUser`

NewAffiliateUser instantiates a new AffiliateUser object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAffiliateUserWithDefaults

`func NewAffiliateUserWithDefaults() *AffiliateUser`

NewAffiliateUserWithDefaults instantiates a new AffiliateUser object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFirstName

`func (o *AffiliateUser) GetFirstName() string`

GetFirstName returns the FirstName field if non-nil, zero value otherwise.

### GetFirstNameOk

`func (o *AffiliateUser) GetFirstNameOk() (*string, bool)`

GetFirstNameOk returns a tuple with the FirstName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirstName

`func (o *AffiliateUser) SetFirstName(v string)`

SetFirstName sets FirstName field to given value.


### GetLastName

`func (o *AffiliateUser) GetLastName() string`

GetLastName returns the LastName field if non-nil, zero value otherwise.

### GetLastNameOk

`func (o *AffiliateUser) GetLastNameOk() (*string, bool)`

GetLastNameOk returns a tuple with the LastName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastName

`func (o *AffiliateUser) SetLastName(v string)`

SetLastName sets LastName field to given value.


### GetEmail

`func (o *AffiliateUser) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *AffiliateUser) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *AffiliateUser) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetAccountStatus

`func (o *AffiliateUser) GetAccountStatus() string`

GetAccountStatus returns the AccountStatus field if non-nil, zero value otherwise.

### GetAccountStatusOk

`func (o *AffiliateUser) GetAccountStatusOk() (*string, bool)`

GetAccountStatusOk returns a tuple with the AccountStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountStatus

`func (o *AffiliateUser) SetAccountStatus(v string)`

SetAccountStatus sets AccountStatus field to given value.


### GetTitle

`func (o *AffiliateUser) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *AffiliateUser) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *AffiliateUser) SetTitle(v string)`

SetTitle sets Title field to given value.

### HasTitle

`func (o *AffiliateUser) HasTitle() bool`

HasTitle returns a boolean if a field has been set.

### GetWorkPhone

`func (o *AffiliateUser) GetWorkPhone() string`

GetWorkPhone returns the WorkPhone field if non-nil, zero value otherwise.

### GetWorkPhoneOk

`func (o *AffiliateUser) GetWorkPhoneOk() (*string, bool)`

GetWorkPhoneOk returns a tuple with the WorkPhone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkPhone

`func (o *AffiliateUser) SetWorkPhone(v string)`

SetWorkPhone sets WorkPhone field to given value.

### HasWorkPhone

`func (o *AffiliateUser) HasWorkPhone() bool`

HasWorkPhone returns a boolean if a field has been set.

### GetCellPhone

`func (o *AffiliateUser) GetCellPhone() string`

GetCellPhone returns the CellPhone field if non-nil, zero value otherwise.

### GetCellPhoneOk

`func (o *AffiliateUser) GetCellPhoneOk() (*string, bool)`

GetCellPhoneOk returns a tuple with the CellPhone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCellPhone

`func (o *AffiliateUser) SetCellPhone(v string)`

SetCellPhone sets CellPhone field to given value.

### HasCellPhone

`func (o *AffiliateUser) HasCellPhone() bool`

HasCellPhone returns a boolean if a field has been set.

### GetInstantMessagingId

`func (o *AffiliateUser) GetInstantMessagingId() int32`

GetInstantMessagingId returns the InstantMessagingId field if non-nil, zero value otherwise.

### GetInstantMessagingIdOk

`func (o *AffiliateUser) GetInstantMessagingIdOk() (*int32, bool)`

GetInstantMessagingIdOk returns a tuple with the InstantMessagingId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstantMessagingId

`func (o *AffiliateUser) SetInstantMessagingId(v int32)`

SetInstantMessagingId sets InstantMessagingId field to given value.

### HasInstantMessagingId

`func (o *AffiliateUser) HasInstantMessagingId() bool`

HasInstantMessagingId returns a boolean if a field has been set.

### GetInstantMessagingIdentifier

`func (o *AffiliateUser) GetInstantMessagingIdentifier() string`

GetInstantMessagingIdentifier returns the InstantMessagingIdentifier field if non-nil, zero value otherwise.

### GetInstantMessagingIdentifierOk

`func (o *AffiliateUser) GetInstantMessagingIdentifierOk() (*string, bool)`

GetInstantMessagingIdentifierOk returns a tuple with the InstantMessagingIdentifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstantMessagingIdentifier

`func (o *AffiliateUser) SetInstantMessagingIdentifier(v string)`

SetInstantMessagingIdentifier sets InstantMessagingIdentifier field to given value.

### HasInstantMessagingIdentifier

`func (o *AffiliateUser) HasInstantMessagingIdentifier() bool`

HasInstantMessagingIdentifier returns a boolean if a field has been set.

### GetLanguageId

`func (o *AffiliateUser) GetLanguageId() int32`

GetLanguageId returns the LanguageId field if non-nil, zero value otherwise.

### GetLanguageIdOk

`func (o *AffiliateUser) GetLanguageIdOk() (*int32, bool)`

GetLanguageIdOk returns a tuple with the LanguageId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLanguageId

`func (o *AffiliateUser) SetLanguageId(v int32)`

SetLanguageId sets LanguageId field to given value.

### HasLanguageId

`func (o *AffiliateUser) HasLanguageId() bool`

HasLanguageId returns a boolean if a field has been set.

### GetTimezoneId

`func (o *AffiliateUser) GetTimezoneId() int32`

GetTimezoneId returns the TimezoneId field if non-nil, zero value otherwise.

### GetTimezoneIdOk

`func (o *AffiliateUser) GetTimezoneIdOk() (*int32, bool)`

GetTimezoneIdOk returns a tuple with the TimezoneId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimezoneId

`func (o *AffiliateUser) SetTimezoneId(v int32)`

SetTimezoneId sets TimezoneId field to given value.

### HasTimezoneId

`func (o *AffiliateUser) HasTimezoneId() bool`

HasTimezoneId returns a boolean if a field has been set.

### GetCurrencyId

`func (o *AffiliateUser) GetCurrencyId() string`

GetCurrencyId returns the CurrencyId field if non-nil, zero value otherwise.

### GetCurrencyIdOk

`func (o *AffiliateUser) GetCurrencyIdOk() (*string, bool)`

GetCurrencyIdOk returns a tuple with the CurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyId

`func (o *AffiliateUser) SetCurrencyId(v string)`

SetCurrencyId sets CurrencyId field to given value.

### HasCurrencyId

`func (o *AffiliateUser) HasCurrencyId() bool`

HasCurrencyId returns a boolean if a field has been set.

### GetInitialPassword

`func (o *AffiliateUser) GetInitialPassword() string`

GetInitialPassword returns the InitialPassword field if non-nil, zero value otherwise.

### GetInitialPasswordOk

`func (o *AffiliateUser) GetInitialPasswordOk() (*string, bool)`

GetInitialPasswordOk returns a tuple with the InitialPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitialPassword

`func (o *AffiliateUser) SetInitialPassword(v string)`

SetInitialPassword sets InitialPassword field to given value.

### HasInitialPassword

`func (o *AffiliateUser) HasInitialPassword() bool`

HasInitialPassword returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


