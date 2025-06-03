# ContactAddress

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Address1** | **string** | The address first line | 
**Address2** | Pointer to **string** | The address second line | [optional] 
**City** | **string** | The city name | 
**RegionCode** | **string** | The region code | 
**CountryCode** | **string** | The country code | 
**CountryId** | Pointer to **int32** | The country ID (numeric identifier) | [optional] 
**ZipPostalCode** | **string** | The ZIP or Postal code | 

## Methods

### NewContactAddress

`func NewContactAddress(address1 string, city string, regionCode string, countryCode string, zipPostalCode string, ) *ContactAddress`

NewContactAddress instantiates a new ContactAddress object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewContactAddressWithDefaults

`func NewContactAddressWithDefaults() *ContactAddress`

NewContactAddressWithDefaults instantiates a new ContactAddress object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAddress1

`func (o *ContactAddress) GetAddress1() string`

GetAddress1 returns the Address1 field if non-nil, zero value otherwise.

### GetAddress1Ok

`func (o *ContactAddress) GetAddress1Ok() (*string, bool)`

GetAddress1Ok returns a tuple with the Address1 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress1

`func (o *ContactAddress) SetAddress1(v string)`

SetAddress1 sets Address1 field to given value.


### GetAddress2

`func (o *ContactAddress) GetAddress2() string`

GetAddress2 returns the Address2 field if non-nil, zero value otherwise.

### GetAddress2Ok

`func (o *ContactAddress) GetAddress2Ok() (*string, bool)`

GetAddress2Ok returns a tuple with the Address2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress2

`func (o *ContactAddress) SetAddress2(v string)`

SetAddress2 sets Address2 field to given value.

### HasAddress2

`func (o *ContactAddress) HasAddress2() bool`

HasAddress2 returns a boolean if a field has been set.

### GetCity

`func (o *ContactAddress) GetCity() string`

GetCity returns the City field if non-nil, zero value otherwise.

### GetCityOk

`func (o *ContactAddress) GetCityOk() (*string, bool)`

GetCityOk returns a tuple with the City field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCity

`func (o *ContactAddress) SetCity(v string)`

SetCity sets City field to given value.


### GetRegionCode

`func (o *ContactAddress) GetRegionCode() string`

GetRegionCode returns the RegionCode field if non-nil, zero value otherwise.

### GetRegionCodeOk

`func (o *ContactAddress) GetRegionCodeOk() (*string, bool)`

GetRegionCodeOk returns a tuple with the RegionCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegionCode

`func (o *ContactAddress) SetRegionCode(v string)`

SetRegionCode sets RegionCode field to given value.


### GetCountryCode

`func (o *ContactAddress) GetCountryCode() string`

GetCountryCode returns the CountryCode field if non-nil, zero value otherwise.

### GetCountryCodeOk

`func (o *ContactAddress) GetCountryCodeOk() (*string, bool)`

GetCountryCodeOk returns a tuple with the CountryCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryCode

`func (o *ContactAddress) SetCountryCode(v string)`

SetCountryCode sets CountryCode field to given value.


### GetCountryId

`func (o *ContactAddress) GetCountryId() int32`

GetCountryId returns the CountryId field if non-nil, zero value otherwise.

### GetCountryIdOk

`func (o *ContactAddress) GetCountryIdOk() (*int32, bool)`

GetCountryIdOk returns a tuple with the CountryId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryId

`func (o *ContactAddress) SetCountryId(v int32)`

SetCountryId sets CountryId field to given value.

### HasCountryId

`func (o *ContactAddress) HasCountryId() bool`

HasCountryId returns a boolean if a field has been set.

### GetZipPostalCode

`func (o *ContactAddress) GetZipPostalCode() string`

GetZipPostalCode returns the ZipPostalCode field if non-nil, zero value otherwise.

### GetZipPostalCodeOk

`func (o *ContactAddress) GetZipPostalCodeOk() (*string, bool)`

GetZipPostalCodeOk returns a tuple with the ZipPostalCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZipPostalCode

`func (o *ContactAddress) SetZipPostalCode(v string)`

SetZipPostalCode sets ZipPostalCode field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


