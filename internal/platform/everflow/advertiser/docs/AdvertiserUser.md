# AdvertiserUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FirstName** | **string** | The advertiser user&#39;s first name | 
**LastName** | **string** | The advertiser user&#39;s last name | 
**Email** | **string** | The advertiser user&#39;s email (must be unique) | 
**AccountStatus** | **string** | The advertiser user&#39;s account status | 
**Title** | Pointer to **string** | The advertiser user&#39;s title | [optional] 
**WorkPhone** | Pointer to **string** | The advertiser user&#39;s work phone number | [optional] 
**CellPhone** | Pointer to **string** | The advertiser user&#39;s cell phone number | [optional] 
**InstantMessagingId** | Pointer to **int32** | The id of an instant messaging platform | [optional] 
**InstantMessagingIdentifier** | Pointer to **string** | The advertiser user&#39;s instant messaging identifier | [optional] 
**LanguageId** | **int32** | The advertiser user&#39;s language id (limited to 1 for English) | 
**TimezoneId** | **int32** | The advertiser user&#39;s timezone id | 
**CurrencyId** | **string** | The advertiser user&#39;s currency id | 
**InitialPassword** | Pointer to **string** | The advertiser user&#39;s login password (optional) | [optional] 

## Methods

### NewAdvertiserUser

`func NewAdvertiserUser(firstName string, lastName string, email string, accountStatus string, languageId int32, timezoneId int32, currencyId string, ) *AdvertiserUser`

NewAdvertiserUser instantiates a new AdvertiserUser object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAdvertiserUserWithDefaults

`func NewAdvertiserUserWithDefaults() *AdvertiserUser`

NewAdvertiserUserWithDefaults instantiates a new AdvertiserUser object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFirstName

`func (o *AdvertiserUser) GetFirstName() string`

GetFirstName returns the FirstName field if non-nil, zero value otherwise.

### GetFirstNameOk

`func (o *AdvertiserUser) GetFirstNameOk() (*string, bool)`

GetFirstNameOk returns a tuple with the FirstName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirstName

`func (o *AdvertiserUser) SetFirstName(v string)`

SetFirstName sets FirstName field to given value.


### GetLastName

`func (o *AdvertiserUser) GetLastName() string`

GetLastName returns the LastName field if non-nil, zero value otherwise.

### GetLastNameOk

`func (o *AdvertiserUser) GetLastNameOk() (*string, bool)`

GetLastNameOk returns a tuple with the LastName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastName

`func (o *AdvertiserUser) SetLastName(v string)`

SetLastName sets LastName field to given value.


### GetEmail

`func (o *AdvertiserUser) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *AdvertiserUser) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *AdvertiserUser) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetAccountStatus

`func (o *AdvertiserUser) GetAccountStatus() string`

GetAccountStatus returns the AccountStatus field if non-nil, zero value otherwise.

### GetAccountStatusOk

`func (o *AdvertiserUser) GetAccountStatusOk() (*string, bool)`

GetAccountStatusOk returns a tuple with the AccountStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountStatus

`func (o *AdvertiserUser) SetAccountStatus(v string)`

SetAccountStatus sets AccountStatus field to given value.


### GetTitle

`func (o *AdvertiserUser) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *AdvertiserUser) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *AdvertiserUser) SetTitle(v string)`

SetTitle sets Title field to given value.

### HasTitle

`func (o *AdvertiserUser) HasTitle() bool`

HasTitle returns a boolean if a field has been set.

### GetWorkPhone

`func (o *AdvertiserUser) GetWorkPhone() string`

GetWorkPhone returns the WorkPhone field if non-nil, zero value otherwise.

### GetWorkPhoneOk

`func (o *AdvertiserUser) GetWorkPhoneOk() (*string, bool)`

GetWorkPhoneOk returns a tuple with the WorkPhone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkPhone

`func (o *AdvertiserUser) SetWorkPhone(v string)`

SetWorkPhone sets WorkPhone field to given value.

### HasWorkPhone

`func (o *AdvertiserUser) HasWorkPhone() bool`

HasWorkPhone returns a boolean if a field has been set.

### GetCellPhone

`func (o *AdvertiserUser) GetCellPhone() string`

GetCellPhone returns the CellPhone field if non-nil, zero value otherwise.

### GetCellPhoneOk

`func (o *AdvertiserUser) GetCellPhoneOk() (*string, bool)`

GetCellPhoneOk returns a tuple with the CellPhone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCellPhone

`func (o *AdvertiserUser) SetCellPhone(v string)`

SetCellPhone sets CellPhone field to given value.

### HasCellPhone

`func (o *AdvertiserUser) HasCellPhone() bool`

HasCellPhone returns a boolean if a field has been set.

### GetInstantMessagingId

`func (o *AdvertiserUser) GetInstantMessagingId() int32`

GetInstantMessagingId returns the InstantMessagingId field if non-nil, zero value otherwise.

### GetInstantMessagingIdOk

`func (o *AdvertiserUser) GetInstantMessagingIdOk() (*int32, bool)`

GetInstantMessagingIdOk returns a tuple with the InstantMessagingId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstantMessagingId

`func (o *AdvertiserUser) SetInstantMessagingId(v int32)`

SetInstantMessagingId sets InstantMessagingId field to given value.

### HasInstantMessagingId

`func (o *AdvertiserUser) HasInstantMessagingId() bool`

HasInstantMessagingId returns a boolean if a field has been set.

### GetInstantMessagingIdentifier

`func (o *AdvertiserUser) GetInstantMessagingIdentifier() string`

GetInstantMessagingIdentifier returns the InstantMessagingIdentifier field if non-nil, zero value otherwise.

### GetInstantMessagingIdentifierOk

`func (o *AdvertiserUser) GetInstantMessagingIdentifierOk() (*string, bool)`

GetInstantMessagingIdentifierOk returns a tuple with the InstantMessagingIdentifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstantMessagingIdentifier

`func (o *AdvertiserUser) SetInstantMessagingIdentifier(v string)`

SetInstantMessagingIdentifier sets InstantMessagingIdentifier field to given value.

### HasInstantMessagingIdentifier

`func (o *AdvertiserUser) HasInstantMessagingIdentifier() bool`

HasInstantMessagingIdentifier returns a boolean if a field has been set.

### GetLanguageId

`func (o *AdvertiserUser) GetLanguageId() int32`

GetLanguageId returns the LanguageId field if non-nil, zero value otherwise.

### GetLanguageIdOk

`func (o *AdvertiserUser) GetLanguageIdOk() (*int32, bool)`

GetLanguageIdOk returns a tuple with the LanguageId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLanguageId

`func (o *AdvertiserUser) SetLanguageId(v int32)`

SetLanguageId sets LanguageId field to given value.


### GetTimezoneId

`func (o *AdvertiserUser) GetTimezoneId() int32`

GetTimezoneId returns the TimezoneId field if non-nil, zero value otherwise.

### GetTimezoneIdOk

`func (o *AdvertiserUser) GetTimezoneIdOk() (*int32, bool)`

GetTimezoneIdOk returns a tuple with the TimezoneId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimezoneId

`func (o *AdvertiserUser) SetTimezoneId(v int32)`

SetTimezoneId sets TimezoneId field to given value.


### GetCurrencyId

`func (o *AdvertiserUser) GetCurrencyId() string`

GetCurrencyId returns the CurrencyId field if non-nil, zero value otherwise.

### GetCurrencyIdOk

`func (o *AdvertiserUser) GetCurrencyIdOk() (*string, bool)`

GetCurrencyIdOk returns a tuple with the CurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyId

`func (o *AdvertiserUser) SetCurrencyId(v string)`

SetCurrencyId sets CurrencyId field to given value.


### GetInitialPassword

`func (o *AdvertiserUser) GetInitialPassword() string`

GetInitialPassword returns the InitialPassword field if non-nil, zero value otherwise.

### GetInitialPasswordOk

`func (o *AdvertiserUser) GetInitialPasswordOk() (*string, bool)`

GetInitialPasswordOk returns a tuple with the InitialPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitialPassword

`func (o *AdvertiserUser) SetInitialPassword(v string)`

SetInitialPassword sets InitialPassword field to given value.

### HasInitialPassword

`func (o *AdvertiserUser) HasInitialPassword() bool`

HasInitialPassword returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


