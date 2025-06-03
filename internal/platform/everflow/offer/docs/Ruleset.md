# Ruleset

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Platforms** | Pointer to **[]map[string]interface{}** |  | [optional] 
**DeviceTypes** | Pointer to **[]map[string]interface{}** |  | [optional] 
**OsVersions** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Browsers** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Languages** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Countries** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Regions** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Cities** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Dmas** | Pointer to **[]map[string]interface{}** |  | [optional] 
**MobileCarriers** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ConnectionTypes** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Ips** | Pointer to **[]map[string]interface{}** |  | [optional] 
**IsBlockProxy** | Pointer to **bool** | Block proxy traffic | [optional] [default to false]
**IsUseDayParting** | Pointer to **bool** | Enable day parting | [optional] [default to false]
**DayPartingApplyTo** | Pointer to **string** | Day parting timezone setting | [optional] 
**DayPartingTimezoneId** | Pointer to **int32** | Timezone ID for day parting | [optional] 
**DaysParting** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Isps** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Brands** | Pointer to **[]map[string]interface{}** |  | [optional] 
**PostalCodes** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Methods

### NewRuleset

`func NewRuleset() *Ruleset`

NewRuleset instantiates a new Ruleset object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRulesetWithDefaults

`func NewRulesetWithDefaults() *Ruleset`

NewRulesetWithDefaults instantiates a new Ruleset object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPlatforms

`func (o *Ruleset) GetPlatforms() []map[string]interface{}`

GetPlatforms returns the Platforms field if non-nil, zero value otherwise.

### GetPlatformsOk

`func (o *Ruleset) GetPlatformsOk() (*[]map[string]interface{}, bool)`

GetPlatformsOk returns a tuple with the Platforms field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlatforms

`func (o *Ruleset) SetPlatforms(v []map[string]interface{})`

SetPlatforms sets Platforms field to given value.

### HasPlatforms

`func (o *Ruleset) HasPlatforms() bool`

HasPlatforms returns a boolean if a field has been set.

### GetDeviceTypes

`func (o *Ruleset) GetDeviceTypes() []map[string]interface{}`

GetDeviceTypes returns the DeviceTypes field if non-nil, zero value otherwise.

### GetDeviceTypesOk

`func (o *Ruleset) GetDeviceTypesOk() (*[]map[string]interface{}, bool)`

GetDeviceTypesOk returns a tuple with the DeviceTypes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceTypes

`func (o *Ruleset) SetDeviceTypes(v []map[string]interface{})`

SetDeviceTypes sets DeviceTypes field to given value.

### HasDeviceTypes

`func (o *Ruleset) HasDeviceTypes() bool`

HasDeviceTypes returns a boolean if a field has been set.

### GetOsVersions

`func (o *Ruleset) GetOsVersions() []map[string]interface{}`

GetOsVersions returns the OsVersions field if non-nil, zero value otherwise.

### GetOsVersionsOk

`func (o *Ruleset) GetOsVersionsOk() (*[]map[string]interface{}, bool)`

GetOsVersionsOk returns a tuple with the OsVersions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOsVersions

`func (o *Ruleset) SetOsVersions(v []map[string]interface{})`

SetOsVersions sets OsVersions field to given value.

### HasOsVersions

`func (o *Ruleset) HasOsVersions() bool`

HasOsVersions returns a boolean if a field has been set.

### GetBrowsers

`func (o *Ruleset) GetBrowsers() []map[string]interface{}`

GetBrowsers returns the Browsers field if non-nil, zero value otherwise.

### GetBrowsersOk

`func (o *Ruleset) GetBrowsersOk() (*[]map[string]interface{}, bool)`

GetBrowsersOk returns a tuple with the Browsers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBrowsers

`func (o *Ruleset) SetBrowsers(v []map[string]interface{})`

SetBrowsers sets Browsers field to given value.

### HasBrowsers

`func (o *Ruleset) HasBrowsers() bool`

HasBrowsers returns a boolean if a field has been set.

### GetLanguages

`func (o *Ruleset) GetLanguages() []map[string]interface{}`

GetLanguages returns the Languages field if non-nil, zero value otherwise.

### GetLanguagesOk

`func (o *Ruleset) GetLanguagesOk() (*[]map[string]interface{}, bool)`

GetLanguagesOk returns a tuple with the Languages field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLanguages

`func (o *Ruleset) SetLanguages(v []map[string]interface{})`

SetLanguages sets Languages field to given value.

### HasLanguages

`func (o *Ruleset) HasLanguages() bool`

HasLanguages returns a boolean if a field has been set.

### GetCountries

`func (o *Ruleset) GetCountries() []map[string]interface{}`

GetCountries returns the Countries field if non-nil, zero value otherwise.

### GetCountriesOk

`func (o *Ruleset) GetCountriesOk() (*[]map[string]interface{}, bool)`

GetCountriesOk returns a tuple with the Countries field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountries

`func (o *Ruleset) SetCountries(v []map[string]interface{})`

SetCountries sets Countries field to given value.

### HasCountries

`func (o *Ruleset) HasCountries() bool`

HasCountries returns a boolean if a field has been set.

### GetRegions

`func (o *Ruleset) GetRegions() []map[string]interface{}`

GetRegions returns the Regions field if non-nil, zero value otherwise.

### GetRegionsOk

`func (o *Ruleset) GetRegionsOk() (*[]map[string]interface{}, bool)`

GetRegionsOk returns a tuple with the Regions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegions

`func (o *Ruleset) SetRegions(v []map[string]interface{})`

SetRegions sets Regions field to given value.

### HasRegions

`func (o *Ruleset) HasRegions() bool`

HasRegions returns a boolean if a field has been set.

### GetCities

`func (o *Ruleset) GetCities() []map[string]interface{}`

GetCities returns the Cities field if non-nil, zero value otherwise.

### GetCitiesOk

`func (o *Ruleset) GetCitiesOk() (*[]map[string]interface{}, bool)`

GetCitiesOk returns a tuple with the Cities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCities

`func (o *Ruleset) SetCities(v []map[string]interface{})`

SetCities sets Cities field to given value.

### HasCities

`func (o *Ruleset) HasCities() bool`

HasCities returns a boolean if a field has been set.

### GetDmas

`func (o *Ruleset) GetDmas() []map[string]interface{}`

GetDmas returns the Dmas field if non-nil, zero value otherwise.

### GetDmasOk

`func (o *Ruleset) GetDmasOk() (*[]map[string]interface{}, bool)`

GetDmasOk returns a tuple with the Dmas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDmas

`func (o *Ruleset) SetDmas(v []map[string]interface{})`

SetDmas sets Dmas field to given value.

### HasDmas

`func (o *Ruleset) HasDmas() bool`

HasDmas returns a boolean if a field has been set.

### GetMobileCarriers

`func (o *Ruleset) GetMobileCarriers() []map[string]interface{}`

GetMobileCarriers returns the MobileCarriers field if non-nil, zero value otherwise.

### GetMobileCarriersOk

`func (o *Ruleset) GetMobileCarriersOk() (*[]map[string]interface{}, bool)`

GetMobileCarriersOk returns a tuple with the MobileCarriers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMobileCarriers

`func (o *Ruleset) SetMobileCarriers(v []map[string]interface{})`

SetMobileCarriers sets MobileCarriers field to given value.

### HasMobileCarriers

`func (o *Ruleset) HasMobileCarriers() bool`

HasMobileCarriers returns a boolean if a field has been set.

### GetConnectionTypes

`func (o *Ruleset) GetConnectionTypes() []map[string]interface{}`

GetConnectionTypes returns the ConnectionTypes field if non-nil, zero value otherwise.

### GetConnectionTypesOk

`func (o *Ruleset) GetConnectionTypesOk() (*[]map[string]interface{}, bool)`

GetConnectionTypesOk returns a tuple with the ConnectionTypes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectionTypes

`func (o *Ruleset) SetConnectionTypes(v []map[string]interface{})`

SetConnectionTypes sets ConnectionTypes field to given value.

### HasConnectionTypes

`func (o *Ruleset) HasConnectionTypes() bool`

HasConnectionTypes returns a boolean if a field has been set.

### GetIps

`func (o *Ruleset) GetIps() []map[string]interface{}`

GetIps returns the Ips field if non-nil, zero value otherwise.

### GetIpsOk

`func (o *Ruleset) GetIpsOk() (*[]map[string]interface{}, bool)`

GetIpsOk returns a tuple with the Ips field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIps

`func (o *Ruleset) SetIps(v []map[string]interface{})`

SetIps sets Ips field to given value.

### HasIps

`func (o *Ruleset) HasIps() bool`

HasIps returns a boolean if a field has been set.

### GetIsBlockProxy

`func (o *Ruleset) GetIsBlockProxy() bool`

GetIsBlockProxy returns the IsBlockProxy field if non-nil, zero value otherwise.

### GetIsBlockProxyOk

`func (o *Ruleset) GetIsBlockProxyOk() (*bool, bool)`

GetIsBlockProxyOk returns a tuple with the IsBlockProxy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsBlockProxy

`func (o *Ruleset) SetIsBlockProxy(v bool)`

SetIsBlockProxy sets IsBlockProxy field to given value.

### HasIsBlockProxy

`func (o *Ruleset) HasIsBlockProxy() bool`

HasIsBlockProxy returns a boolean if a field has been set.

### GetIsUseDayParting

`func (o *Ruleset) GetIsUseDayParting() bool`

GetIsUseDayParting returns the IsUseDayParting field if non-nil, zero value otherwise.

### GetIsUseDayPartingOk

`func (o *Ruleset) GetIsUseDayPartingOk() (*bool, bool)`

GetIsUseDayPartingOk returns a tuple with the IsUseDayParting field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsUseDayParting

`func (o *Ruleset) SetIsUseDayParting(v bool)`

SetIsUseDayParting sets IsUseDayParting field to given value.

### HasIsUseDayParting

`func (o *Ruleset) HasIsUseDayParting() bool`

HasIsUseDayParting returns a boolean if a field has been set.

### GetDayPartingApplyTo

`func (o *Ruleset) GetDayPartingApplyTo() string`

GetDayPartingApplyTo returns the DayPartingApplyTo field if non-nil, zero value otherwise.

### GetDayPartingApplyToOk

`func (o *Ruleset) GetDayPartingApplyToOk() (*string, bool)`

GetDayPartingApplyToOk returns a tuple with the DayPartingApplyTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDayPartingApplyTo

`func (o *Ruleset) SetDayPartingApplyTo(v string)`

SetDayPartingApplyTo sets DayPartingApplyTo field to given value.

### HasDayPartingApplyTo

`func (o *Ruleset) HasDayPartingApplyTo() bool`

HasDayPartingApplyTo returns a boolean if a field has been set.

### GetDayPartingTimezoneId

`func (o *Ruleset) GetDayPartingTimezoneId() int32`

GetDayPartingTimezoneId returns the DayPartingTimezoneId field if non-nil, zero value otherwise.

### GetDayPartingTimezoneIdOk

`func (o *Ruleset) GetDayPartingTimezoneIdOk() (*int32, bool)`

GetDayPartingTimezoneIdOk returns a tuple with the DayPartingTimezoneId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDayPartingTimezoneId

`func (o *Ruleset) SetDayPartingTimezoneId(v int32)`

SetDayPartingTimezoneId sets DayPartingTimezoneId field to given value.

### HasDayPartingTimezoneId

`func (o *Ruleset) HasDayPartingTimezoneId() bool`

HasDayPartingTimezoneId returns a boolean if a field has been set.

### GetDaysParting

`func (o *Ruleset) GetDaysParting() []map[string]interface{}`

GetDaysParting returns the DaysParting field if non-nil, zero value otherwise.

### GetDaysPartingOk

`func (o *Ruleset) GetDaysPartingOk() (*[]map[string]interface{}, bool)`

GetDaysPartingOk returns a tuple with the DaysParting field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDaysParting

`func (o *Ruleset) SetDaysParting(v []map[string]interface{})`

SetDaysParting sets DaysParting field to given value.

### HasDaysParting

`func (o *Ruleset) HasDaysParting() bool`

HasDaysParting returns a boolean if a field has been set.

### GetIsps

`func (o *Ruleset) GetIsps() []map[string]interface{}`

GetIsps returns the Isps field if non-nil, zero value otherwise.

### GetIspsOk

`func (o *Ruleset) GetIspsOk() (*[]map[string]interface{}, bool)`

GetIspsOk returns a tuple with the Isps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsps

`func (o *Ruleset) SetIsps(v []map[string]interface{})`

SetIsps sets Isps field to given value.

### HasIsps

`func (o *Ruleset) HasIsps() bool`

HasIsps returns a boolean if a field has been set.

### GetBrands

`func (o *Ruleset) GetBrands() []map[string]interface{}`

GetBrands returns the Brands field if non-nil, zero value otherwise.

### GetBrandsOk

`func (o *Ruleset) GetBrandsOk() (*[]map[string]interface{}, bool)`

GetBrandsOk returns a tuple with the Brands field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBrands

`func (o *Ruleset) SetBrands(v []map[string]interface{})`

SetBrands sets Brands field to given value.

### HasBrands

`func (o *Ruleset) HasBrands() bool`

HasBrands returns a boolean if a field has been set.

### GetPostalCodes

`func (o *Ruleset) GetPostalCodes() []map[string]interface{}`

GetPostalCodes returns the PostalCodes field if non-nil, zero value otherwise.

### GetPostalCodesOk

`func (o *Ruleset) GetPostalCodesOk() (*[]map[string]interface{}, bool)`

GetPostalCodesOk returns a tuple with the PostalCodes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostalCodes

`func (o *Ruleset) SetPostalCodes(v []map[string]interface{})`

SetPostalCodes sets PostalCodes field to given value.

### HasPostalCodes

`func (o *Ruleset) HasPostalCodes() bool`

HasPostalCodes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


