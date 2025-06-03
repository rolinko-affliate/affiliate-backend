# CreateOfferRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NetworkAdvertiserId** | **int32** | ID of the advertiser submitting the offer | 
**NetworkOfferGroupId** | Pointer to **int32** | ID of the offer group associated with the offer | [optional] 
**Name** | **string** | Displayed name of the offer | 
**ThumbnailUrl** | Pointer to **string** | URL of the image thumbnail associated with the offer | [optional] 
**NetworkCategoryId** | Pointer to **int32** | ID of the category type associated with the offer | [optional] 
**InternalNotes** | Pointer to **string** | Notes on the offer for network employees | [optional] 
**DestinationUrl** | **string** | URL of the final landing page associated with the offer | 
**ServerSideUrl** | Pointer to **string** | Server-side URL that will be asynchronously fired by Everflow | [optional] 
**IsViewThroughEnabled** | Pointer to **bool** | Whether conversions can be generated from impressions | [optional] [default to false]
**ViewThroughDestinationUrl** | Pointer to **string** | URL of the final landing page when redirected from an impression | [optional] 
**PreviewUrl** | Pointer to **string** | URL of a preview of the offer landing page | [optional] 
**OfferStatus** | **string** | Status of the offer | 
**CurrencyId** | Pointer to **string** | Currency used to compute payouts, costs and revenues | [optional] [default to "USD"]
**CapsTimezoneId** | Pointer to **int32** | ID of the timezone used for caps | [optional] [default to 0]
**ProjectId** | Pointer to **string** | ID for the advertiser campaign or an Insertion Order | [optional] 
**DateLiveUntil** | Pointer to **string** | Date until when the offer can be run (yyyy-MM-dd) | [optional] 
**HtmlDescription** | Pointer to **string** | Description of the offer for affiliates (HTML accepted) | [optional] 
**IsUsingExplicitTermsAndConditions** | Pointer to **bool** | Whether the offer is using specific Terms and Conditions | [optional] [default to false]
**TermsAndConditions** | Pointer to **string** | Text listing the specific Terms and Conditions | [optional] 
**IsForceTermsAndConditions** | Pointer to **bool** | Whether affiliates are required to accept the offer&#39;s Terms and Conditions | [optional] [default to false]
**IsCapsEnabled** | Pointer to **bool** | Whether caps are enabled | [optional] [default to false]
**DailyConversionCap** | Pointer to **int32** | Limit to the number of unique conversions in one day | [optional] [default to 0]
**WeeklyConversionCap** | Pointer to **int32** | Limit to the number of unique conversions in one week | [optional] [default to 0]
**MonthlyConversionCap** | Pointer to **int32** | Limit to the number of unique conversions in one month | [optional] [default to 0]
**GlobalConversionCap** | Pointer to **int32** | Limit to the total number of unique conversions | [optional] [default to 0]
**DailyPayoutCap** | Pointer to **int32** | Limit to the affiliate&#39;s payout for one day | [optional] [default to 0]
**WeeklyPayoutCap** | Pointer to **int32** | Limit to the affiliate&#39;s payout for one week | [optional] [default to 0]
**MonthlyPayoutCap** | Pointer to **int32** | Limit to the affiliate&#39;s payout for one month | [optional] [default to 0]
**GlobalPayoutCap** | Pointer to **int32** | Limit to the affiliate&#39;s total payout | [optional] [default to 0]
**DailyRevenueCap** | Pointer to **int32** | Limit to the network&#39;s revenue for one day | [optional] [default to 0]
**WeeklyRevenueCap** | Pointer to **int32** | Limit to the network&#39;s revenue for one week | [optional] [default to 0]
**MonthlyRevenueCap** | Pointer to **int32** | Limit to the network&#39;s revenue for one month | [optional] [default to 0]
**GlobalRevenueCap** | Pointer to **int32** | Limit to the network&#39;s total revenue | [optional] [default to 0]
**DailyClickCap** | Pointer to **int32** | Limit to the number of unique clicks in one day | [optional] [default to 0]
**WeeklyClickCap** | Pointer to **int32** | Limit to the number of unique clicks in one week | [optional] [default to 0]
**MonthlyClickCap** | Pointer to **int32** | Limit to the number of unique clicks in one month | [optional] [default to 0]
**GlobalClickCap** | Pointer to **int32** | Limit to the total number of unique clicks | [optional] [default to 0]
**RedirectMode** | Pointer to **string** | Setting used to obscure referrer URLs from advertisers | [optional] [default to "standard"]
**IsUsingSuppressionList** | Pointer to **bool** | Whether an email suppression list is used | [optional] [default to false]
**SuppressionListId** | Pointer to **int32** | ID of the suppression list | [optional] [default to 0]
**IsDuplicateFilterEnabled** | Pointer to **bool** | Whether duplicate clicks are filtered | [optional] [default to false]
**DuplicateFilterTargetingAction** | Pointer to **string** | Action for duplicate clicks | [optional] 
**NetworkTrackingDomainId** | Pointer to **int32** | ID of the tracking domain | [optional] 
**IsUseSecureLink** | Pointer to **bool** | Whether tracking links use HTTPS | [optional] [default to false]
**IsAllowDeepLink** | Pointer to **bool** | Whether affiliates can send traffic to target URLs | [optional] [default to false]
**IsSessionTrackingEnabled** | Pointer to **bool** | Whether conversions are blocked based on time intervals | [optional] [default to false]
**SessionTrackingLifespanHour** | Pointer to **int32** | Maximum interval between click and conversion | [optional] [default to 0]
**SessionTrackingMinimumLifespanSecond** | Pointer to **int32** | Minimum interval between click and conversion | [optional] [default to 0]
**IsViewThroughSessionTrackingEnabled** | Pointer to **bool** | Whether conversions from impressions are time-tracked | [optional] [default to false]
**ViewThroughSessionTrackingLifespanMinute** | Pointer to **int32** | Maximum interval between impression and conversion | [optional] [default to 0]
**ViewThroughSessionTrackingMinimalLifespanSecond** | Pointer to **int32** | Minimum interval between impression and conversion | [optional] [default to 0]
**IsBlockAlreadyConverted** | Pointer to **bool** | Whether to block clicks from already-converted users | [optional] [default to false]
**AlreadyConvertedAction** | Pointer to **string** | Action for already-converted users | [optional] 
**Visibility** | Pointer to **string** | Offer visibility for affiliates | [optional] [default to "public"]
**ConversionMethod** | Pointer to **string** | Method used by advertiser to fire tracking data | [optional] [default to "server_postback"]
**IsWhitelistCheckEnabled** | Pointer to **bool** | Whether to check conversion postback origin | [optional] [default to false]
**IsUseScrubRate** | Pointer to **bool** | Whether to throttle conversions | [optional] [default to false]
**ScrubRateStatus** | Pointer to **string** | Status for throttled conversions | [optional] 
**ScrubRatePercentage** | Pointer to **int32** | Percentage of conversions to throttle | [optional] [default to 0]
**SessionDefinition** | Pointer to **string** | Method for determining unique clicks | [optional] [default to "cookie"]
**SessionDuration** | Pointer to **int32** | Duration of active session in hours | [optional] [default to 24]
**AppIdentifier** | Pointer to **string** | Bundle ID for iOS/Android Apps | [optional] 
**IsDescriptionPlainText** | Pointer to **bool** | Whether description is plain text | [optional] [default to false]
**IsUseDirectLinking** | Pointer to **bool** | Whether offer uses Direct Linking | [optional] [default to false]
**IsFailTrafficEnabled** | Pointer to **bool** | Whether invalid clicks are redirected | [optional] [default to false]
**RedirectRoutingMethod** | Pointer to **string** | How fail traffic is handled | [optional] [default to "internal"]
**RedirectInternalRoutingType** | Pointer to **string** | Redirect distribution mechanism | [optional] [default to "priority"]
**Meta** | Pointer to [**Meta**](Meta.md) |  | [optional] 
**Creatives** | Pointer to [**[]Creative**](Creative.md) |  | [optional] 
**InternalRedirects** | Pointer to [**[]InternalRedirect**](InternalRedirect.md) |  | [optional] 
**Ruleset** | Pointer to [**Ruleset**](Ruleset.md) |  | [optional] 
**TrafficFilters** | Pointer to [**[]TrafficFilter**](TrafficFilter.md) |  | [optional] 
**Email** | Pointer to [**EmailSettings**](EmailSettings.md) |  | [optional] 
**EmailOptout** | Pointer to [**EmailOptoutSettings**](EmailOptoutSettings.md) |  | [optional] 
**Labels** | Pointer to **[]string** | Labels for organizing offers | [optional] 
**SourceNames** | Pointer to **[]string** | Names of the source | [optional] 
**PayoutRevenue** | [**[]PayoutRevenue**](PayoutRevenue.md) |  | 
**ThumbnailFile** | Pointer to [**ThumbnailFile**](ThumbnailFile.md) |  | [optional] 
**Integrations** | Pointer to [**Integrations**](Integrations.md) |  | [optional] 
**Channels** | Pointer to [**[]Channel**](Channel.md) |  | [optional] 
**RequirementKpis** | Pointer to [**[]RequirementKPI**](RequirementKPI.md) |  | [optional] 
**RequirementTrackingParameters** | Pointer to [**[]RequirementTrackingParameter**](RequirementTrackingParameter.md) |  | [optional] 
**EmailAttributionMethod** | Pointer to **string** | Email attribution method | [optional] 
**AttributionMethod** | Pointer to **string** | Attribution method | [optional] 

## Methods

### NewCreateOfferRequest

`func NewCreateOfferRequest(networkAdvertiserId int32, name string, destinationUrl string, offerStatus string, payoutRevenue []PayoutRevenue, ) *CreateOfferRequest`

NewCreateOfferRequest instantiates a new CreateOfferRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateOfferRequestWithDefaults

`func NewCreateOfferRequestWithDefaults() *CreateOfferRequest`

NewCreateOfferRequestWithDefaults instantiates a new CreateOfferRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNetworkAdvertiserId

`func (o *CreateOfferRequest) GetNetworkAdvertiserId() int32`

GetNetworkAdvertiserId returns the NetworkAdvertiserId field if non-nil, zero value otherwise.

### GetNetworkAdvertiserIdOk

`func (o *CreateOfferRequest) GetNetworkAdvertiserIdOk() (*int32, bool)`

GetNetworkAdvertiserIdOk returns a tuple with the NetworkAdvertiserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkAdvertiserId

`func (o *CreateOfferRequest) SetNetworkAdvertiserId(v int32)`

SetNetworkAdvertiserId sets NetworkAdvertiserId field to given value.


### GetNetworkOfferGroupId

`func (o *CreateOfferRequest) GetNetworkOfferGroupId() int32`

GetNetworkOfferGroupId returns the NetworkOfferGroupId field if non-nil, zero value otherwise.

### GetNetworkOfferGroupIdOk

`func (o *CreateOfferRequest) GetNetworkOfferGroupIdOk() (*int32, bool)`

GetNetworkOfferGroupIdOk returns a tuple with the NetworkOfferGroupId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkOfferGroupId

`func (o *CreateOfferRequest) SetNetworkOfferGroupId(v int32)`

SetNetworkOfferGroupId sets NetworkOfferGroupId field to given value.

### HasNetworkOfferGroupId

`func (o *CreateOfferRequest) HasNetworkOfferGroupId() bool`

HasNetworkOfferGroupId returns a boolean if a field has been set.

### GetName

`func (o *CreateOfferRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateOfferRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateOfferRequest) SetName(v string)`

SetName sets Name field to given value.


### GetThumbnailUrl

`func (o *CreateOfferRequest) GetThumbnailUrl() string`

GetThumbnailUrl returns the ThumbnailUrl field if non-nil, zero value otherwise.

### GetThumbnailUrlOk

`func (o *CreateOfferRequest) GetThumbnailUrlOk() (*string, bool)`

GetThumbnailUrlOk returns a tuple with the ThumbnailUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThumbnailUrl

`func (o *CreateOfferRequest) SetThumbnailUrl(v string)`

SetThumbnailUrl sets ThumbnailUrl field to given value.

### HasThumbnailUrl

`func (o *CreateOfferRequest) HasThumbnailUrl() bool`

HasThumbnailUrl returns a boolean if a field has been set.

### GetNetworkCategoryId

`func (o *CreateOfferRequest) GetNetworkCategoryId() int32`

GetNetworkCategoryId returns the NetworkCategoryId field if non-nil, zero value otherwise.

### GetNetworkCategoryIdOk

`func (o *CreateOfferRequest) GetNetworkCategoryIdOk() (*int32, bool)`

GetNetworkCategoryIdOk returns a tuple with the NetworkCategoryId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkCategoryId

`func (o *CreateOfferRequest) SetNetworkCategoryId(v int32)`

SetNetworkCategoryId sets NetworkCategoryId field to given value.

### HasNetworkCategoryId

`func (o *CreateOfferRequest) HasNetworkCategoryId() bool`

HasNetworkCategoryId returns a boolean if a field has been set.

### GetInternalNotes

`func (o *CreateOfferRequest) GetInternalNotes() string`

GetInternalNotes returns the InternalNotes field if non-nil, zero value otherwise.

### GetInternalNotesOk

`func (o *CreateOfferRequest) GetInternalNotesOk() (*string, bool)`

GetInternalNotesOk returns a tuple with the InternalNotes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternalNotes

`func (o *CreateOfferRequest) SetInternalNotes(v string)`

SetInternalNotes sets InternalNotes field to given value.

### HasInternalNotes

`func (o *CreateOfferRequest) HasInternalNotes() bool`

HasInternalNotes returns a boolean if a field has been set.

### GetDestinationUrl

`func (o *CreateOfferRequest) GetDestinationUrl() string`

GetDestinationUrl returns the DestinationUrl field if non-nil, zero value otherwise.

### GetDestinationUrlOk

`func (o *CreateOfferRequest) GetDestinationUrlOk() (*string, bool)`

GetDestinationUrlOk returns a tuple with the DestinationUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinationUrl

`func (o *CreateOfferRequest) SetDestinationUrl(v string)`

SetDestinationUrl sets DestinationUrl field to given value.


### GetServerSideUrl

`func (o *CreateOfferRequest) GetServerSideUrl() string`

GetServerSideUrl returns the ServerSideUrl field if non-nil, zero value otherwise.

### GetServerSideUrlOk

`func (o *CreateOfferRequest) GetServerSideUrlOk() (*string, bool)`

GetServerSideUrlOk returns a tuple with the ServerSideUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServerSideUrl

`func (o *CreateOfferRequest) SetServerSideUrl(v string)`

SetServerSideUrl sets ServerSideUrl field to given value.

### HasServerSideUrl

`func (o *CreateOfferRequest) HasServerSideUrl() bool`

HasServerSideUrl returns a boolean if a field has been set.

### GetIsViewThroughEnabled

`func (o *CreateOfferRequest) GetIsViewThroughEnabled() bool`

GetIsViewThroughEnabled returns the IsViewThroughEnabled field if non-nil, zero value otherwise.

### GetIsViewThroughEnabledOk

`func (o *CreateOfferRequest) GetIsViewThroughEnabledOk() (*bool, bool)`

GetIsViewThroughEnabledOk returns a tuple with the IsViewThroughEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsViewThroughEnabled

`func (o *CreateOfferRequest) SetIsViewThroughEnabled(v bool)`

SetIsViewThroughEnabled sets IsViewThroughEnabled field to given value.

### HasIsViewThroughEnabled

`func (o *CreateOfferRequest) HasIsViewThroughEnabled() bool`

HasIsViewThroughEnabled returns a boolean if a field has been set.

### GetViewThroughDestinationUrl

`func (o *CreateOfferRequest) GetViewThroughDestinationUrl() string`

GetViewThroughDestinationUrl returns the ViewThroughDestinationUrl field if non-nil, zero value otherwise.

### GetViewThroughDestinationUrlOk

`func (o *CreateOfferRequest) GetViewThroughDestinationUrlOk() (*string, bool)`

GetViewThroughDestinationUrlOk returns a tuple with the ViewThroughDestinationUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetViewThroughDestinationUrl

`func (o *CreateOfferRequest) SetViewThroughDestinationUrl(v string)`

SetViewThroughDestinationUrl sets ViewThroughDestinationUrl field to given value.

### HasViewThroughDestinationUrl

`func (o *CreateOfferRequest) HasViewThroughDestinationUrl() bool`

HasViewThroughDestinationUrl returns a boolean if a field has been set.

### GetPreviewUrl

`func (o *CreateOfferRequest) GetPreviewUrl() string`

GetPreviewUrl returns the PreviewUrl field if non-nil, zero value otherwise.

### GetPreviewUrlOk

`func (o *CreateOfferRequest) GetPreviewUrlOk() (*string, bool)`

GetPreviewUrlOk returns a tuple with the PreviewUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreviewUrl

`func (o *CreateOfferRequest) SetPreviewUrl(v string)`

SetPreviewUrl sets PreviewUrl field to given value.

### HasPreviewUrl

`func (o *CreateOfferRequest) HasPreviewUrl() bool`

HasPreviewUrl returns a boolean if a field has been set.

### GetOfferStatus

`func (o *CreateOfferRequest) GetOfferStatus() string`

GetOfferStatus returns the OfferStatus field if non-nil, zero value otherwise.

### GetOfferStatusOk

`func (o *CreateOfferRequest) GetOfferStatusOk() (*string, bool)`

GetOfferStatusOk returns a tuple with the OfferStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOfferStatus

`func (o *CreateOfferRequest) SetOfferStatus(v string)`

SetOfferStatus sets OfferStatus field to given value.


### GetCurrencyId

`func (o *CreateOfferRequest) GetCurrencyId() string`

GetCurrencyId returns the CurrencyId field if non-nil, zero value otherwise.

### GetCurrencyIdOk

`func (o *CreateOfferRequest) GetCurrencyIdOk() (*string, bool)`

GetCurrencyIdOk returns a tuple with the CurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyId

`func (o *CreateOfferRequest) SetCurrencyId(v string)`

SetCurrencyId sets CurrencyId field to given value.

### HasCurrencyId

`func (o *CreateOfferRequest) HasCurrencyId() bool`

HasCurrencyId returns a boolean if a field has been set.

### GetCapsTimezoneId

`func (o *CreateOfferRequest) GetCapsTimezoneId() int32`

GetCapsTimezoneId returns the CapsTimezoneId field if non-nil, zero value otherwise.

### GetCapsTimezoneIdOk

`func (o *CreateOfferRequest) GetCapsTimezoneIdOk() (*int32, bool)`

GetCapsTimezoneIdOk returns a tuple with the CapsTimezoneId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapsTimezoneId

`func (o *CreateOfferRequest) SetCapsTimezoneId(v int32)`

SetCapsTimezoneId sets CapsTimezoneId field to given value.

### HasCapsTimezoneId

`func (o *CreateOfferRequest) HasCapsTimezoneId() bool`

HasCapsTimezoneId returns a boolean if a field has been set.

### GetProjectId

`func (o *CreateOfferRequest) GetProjectId() string`

GetProjectId returns the ProjectId field if non-nil, zero value otherwise.

### GetProjectIdOk

`func (o *CreateOfferRequest) GetProjectIdOk() (*string, bool)`

GetProjectIdOk returns a tuple with the ProjectId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProjectId

`func (o *CreateOfferRequest) SetProjectId(v string)`

SetProjectId sets ProjectId field to given value.

### HasProjectId

`func (o *CreateOfferRequest) HasProjectId() bool`

HasProjectId returns a boolean if a field has been set.

### GetDateLiveUntil

`func (o *CreateOfferRequest) GetDateLiveUntil() string`

GetDateLiveUntil returns the DateLiveUntil field if non-nil, zero value otherwise.

### GetDateLiveUntilOk

`func (o *CreateOfferRequest) GetDateLiveUntilOk() (*string, bool)`

GetDateLiveUntilOk returns a tuple with the DateLiveUntil field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDateLiveUntil

`func (o *CreateOfferRequest) SetDateLiveUntil(v string)`

SetDateLiveUntil sets DateLiveUntil field to given value.

### HasDateLiveUntil

`func (o *CreateOfferRequest) HasDateLiveUntil() bool`

HasDateLiveUntil returns a boolean if a field has been set.

### GetHtmlDescription

`func (o *CreateOfferRequest) GetHtmlDescription() string`

GetHtmlDescription returns the HtmlDescription field if non-nil, zero value otherwise.

### GetHtmlDescriptionOk

`func (o *CreateOfferRequest) GetHtmlDescriptionOk() (*string, bool)`

GetHtmlDescriptionOk returns a tuple with the HtmlDescription field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHtmlDescription

`func (o *CreateOfferRequest) SetHtmlDescription(v string)`

SetHtmlDescription sets HtmlDescription field to given value.

### HasHtmlDescription

`func (o *CreateOfferRequest) HasHtmlDescription() bool`

HasHtmlDescription returns a boolean if a field has been set.

### GetIsUsingExplicitTermsAndConditions

`func (o *CreateOfferRequest) GetIsUsingExplicitTermsAndConditions() bool`

GetIsUsingExplicitTermsAndConditions returns the IsUsingExplicitTermsAndConditions field if non-nil, zero value otherwise.

### GetIsUsingExplicitTermsAndConditionsOk

`func (o *CreateOfferRequest) GetIsUsingExplicitTermsAndConditionsOk() (*bool, bool)`

GetIsUsingExplicitTermsAndConditionsOk returns a tuple with the IsUsingExplicitTermsAndConditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsUsingExplicitTermsAndConditions

`func (o *CreateOfferRequest) SetIsUsingExplicitTermsAndConditions(v bool)`

SetIsUsingExplicitTermsAndConditions sets IsUsingExplicitTermsAndConditions field to given value.

### HasIsUsingExplicitTermsAndConditions

`func (o *CreateOfferRequest) HasIsUsingExplicitTermsAndConditions() bool`

HasIsUsingExplicitTermsAndConditions returns a boolean if a field has been set.

### GetTermsAndConditions

`func (o *CreateOfferRequest) GetTermsAndConditions() string`

GetTermsAndConditions returns the TermsAndConditions field if non-nil, zero value otherwise.

### GetTermsAndConditionsOk

`func (o *CreateOfferRequest) GetTermsAndConditionsOk() (*string, bool)`

GetTermsAndConditionsOk returns a tuple with the TermsAndConditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTermsAndConditions

`func (o *CreateOfferRequest) SetTermsAndConditions(v string)`

SetTermsAndConditions sets TermsAndConditions field to given value.

### HasTermsAndConditions

`func (o *CreateOfferRequest) HasTermsAndConditions() bool`

HasTermsAndConditions returns a boolean if a field has been set.

### GetIsForceTermsAndConditions

`func (o *CreateOfferRequest) GetIsForceTermsAndConditions() bool`

GetIsForceTermsAndConditions returns the IsForceTermsAndConditions field if non-nil, zero value otherwise.

### GetIsForceTermsAndConditionsOk

`func (o *CreateOfferRequest) GetIsForceTermsAndConditionsOk() (*bool, bool)`

GetIsForceTermsAndConditionsOk returns a tuple with the IsForceTermsAndConditions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsForceTermsAndConditions

`func (o *CreateOfferRequest) SetIsForceTermsAndConditions(v bool)`

SetIsForceTermsAndConditions sets IsForceTermsAndConditions field to given value.

### HasIsForceTermsAndConditions

`func (o *CreateOfferRequest) HasIsForceTermsAndConditions() bool`

HasIsForceTermsAndConditions returns a boolean if a field has been set.

### GetIsCapsEnabled

`func (o *CreateOfferRequest) GetIsCapsEnabled() bool`

GetIsCapsEnabled returns the IsCapsEnabled field if non-nil, zero value otherwise.

### GetIsCapsEnabledOk

`func (o *CreateOfferRequest) GetIsCapsEnabledOk() (*bool, bool)`

GetIsCapsEnabledOk returns a tuple with the IsCapsEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsCapsEnabled

`func (o *CreateOfferRequest) SetIsCapsEnabled(v bool)`

SetIsCapsEnabled sets IsCapsEnabled field to given value.

### HasIsCapsEnabled

`func (o *CreateOfferRequest) HasIsCapsEnabled() bool`

HasIsCapsEnabled returns a boolean if a field has been set.

### GetDailyConversionCap

`func (o *CreateOfferRequest) GetDailyConversionCap() int32`

GetDailyConversionCap returns the DailyConversionCap field if non-nil, zero value otherwise.

### GetDailyConversionCapOk

`func (o *CreateOfferRequest) GetDailyConversionCapOk() (*int32, bool)`

GetDailyConversionCapOk returns a tuple with the DailyConversionCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDailyConversionCap

`func (o *CreateOfferRequest) SetDailyConversionCap(v int32)`

SetDailyConversionCap sets DailyConversionCap field to given value.

### HasDailyConversionCap

`func (o *CreateOfferRequest) HasDailyConversionCap() bool`

HasDailyConversionCap returns a boolean if a field has been set.

### GetWeeklyConversionCap

`func (o *CreateOfferRequest) GetWeeklyConversionCap() int32`

GetWeeklyConversionCap returns the WeeklyConversionCap field if non-nil, zero value otherwise.

### GetWeeklyConversionCapOk

`func (o *CreateOfferRequest) GetWeeklyConversionCapOk() (*int32, bool)`

GetWeeklyConversionCapOk returns a tuple with the WeeklyConversionCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWeeklyConversionCap

`func (o *CreateOfferRequest) SetWeeklyConversionCap(v int32)`

SetWeeklyConversionCap sets WeeklyConversionCap field to given value.

### HasWeeklyConversionCap

`func (o *CreateOfferRequest) HasWeeklyConversionCap() bool`

HasWeeklyConversionCap returns a boolean if a field has been set.

### GetMonthlyConversionCap

`func (o *CreateOfferRequest) GetMonthlyConversionCap() int32`

GetMonthlyConversionCap returns the MonthlyConversionCap field if non-nil, zero value otherwise.

### GetMonthlyConversionCapOk

`func (o *CreateOfferRequest) GetMonthlyConversionCapOk() (*int32, bool)`

GetMonthlyConversionCapOk returns a tuple with the MonthlyConversionCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMonthlyConversionCap

`func (o *CreateOfferRequest) SetMonthlyConversionCap(v int32)`

SetMonthlyConversionCap sets MonthlyConversionCap field to given value.

### HasMonthlyConversionCap

`func (o *CreateOfferRequest) HasMonthlyConversionCap() bool`

HasMonthlyConversionCap returns a boolean if a field has been set.

### GetGlobalConversionCap

`func (o *CreateOfferRequest) GetGlobalConversionCap() int32`

GetGlobalConversionCap returns the GlobalConversionCap field if non-nil, zero value otherwise.

### GetGlobalConversionCapOk

`func (o *CreateOfferRequest) GetGlobalConversionCapOk() (*int32, bool)`

GetGlobalConversionCapOk returns a tuple with the GlobalConversionCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGlobalConversionCap

`func (o *CreateOfferRequest) SetGlobalConversionCap(v int32)`

SetGlobalConversionCap sets GlobalConversionCap field to given value.

### HasGlobalConversionCap

`func (o *CreateOfferRequest) HasGlobalConversionCap() bool`

HasGlobalConversionCap returns a boolean if a field has been set.

### GetDailyPayoutCap

`func (o *CreateOfferRequest) GetDailyPayoutCap() int32`

GetDailyPayoutCap returns the DailyPayoutCap field if non-nil, zero value otherwise.

### GetDailyPayoutCapOk

`func (o *CreateOfferRequest) GetDailyPayoutCapOk() (*int32, bool)`

GetDailyPayoutCapOk returns a tuple with the DailyPayoutCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDailyPayoutCap

`func (o *CreateOfferRequest) SetDailyPayoutCap(v int32)`

SetDailyPayoutCap sets DailyPayoutCap field to given value.

### HasDailyPayoutCap

`func (o *CreateOfferRequest) HasDailyPayoutCap() bool`

HasDailyPayoutCap returns a boolean if a field has been set.

### GetWeeklyPayoutCap

`func (o *CreateOfferRequest) GetWeeklyPayoutCap() int32`

GetWeeklyPayoutCap returns the WeeklyPayoutCap field if non-nil, zero value otherwise.

### GetWeeklyPayoutCapOk

`func (o *CreateOfferRequest) GetWeeklyPayoutCapOk() (*int32, bool)`

GetWeeklyPayoutCapOk returns a tuple with the WeeklyPayoutCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWeeklyPayoutCap

`func (o *CreateOfferRequest) SetWeeklyPayoutCap(v int32)`

SetWeeklyPayoutCap sets WeeklyPayoutCap field to given value.

### HasWeeklyPayoutCap

`func (o *CreateOfferRequest) HasWeeklyPayoutCap() bool`

HasWeeklyPayoutCap returns a boolean if a field has been set.

### GetMonthlyPayoutCap

`func (o *CreateOfferRequest) GetMonthlyPayoutCap() int32`

GetMonthlyPayoutCap returns the MonthlyPayoutCap field if non-nil, zero value otherwise.

### GetMonthlyPayoutCapOk

`func (o *CreateOfferRequest) GetMonthlyPayoutCapOk() (*int32, bool)`

GetMonthlyPayoutCapOk returns a tuple with the MonthlyPayoutCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMonthlyPayoutCap

`func (o *CreateOfferRequest) SetMonthlyPayoutCap(v int32)`

SetMonthlyPayoutCap sets MonthlyPayoutCap field to given value.

### HasMonthlyPayoutCap

`func (o *CreateOfferRequest) HasMonthlyPayoutCap() bool`

HasMonthlyPayoutCap returns a boolean if a field has been set.

### GetGlobalPayoutCap

`func (o *CreateOfferRequest) GetGlobalPayoutCap() int32`

GetGlobalPayoutCap returns the GlobalPayoutCap field if non-nil, zero value otherwise.

### GetGlobalPayoutCapOk

`func (o *CreateOfferRequest) GetGlobalPayoutCapOk() (*int32, bool)`

GetGlobalPayoutCapOk returns a tuple with the GlobalPayoutCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGlobalPayoutCap

`func (o *CreateOfferRequest) SetGlobalPayoutCap(v int32)`

SetGlobalPayoutCap sets GlobalPayoutCap field to given value.

### HasGlobalPayoutCap

`func (o *CreateOfferRequest) HasGlobalPayoutCap() bool`

HasGlobalPayoutCap returns a boolean if a field has been set.

### GetDailyRevenueCap

`func (o *CreateOfferRequest) GetDailyRevenueCap() int32`

GetDailyRevenueCap returns the DailyRevenueCap field if non-nil, zero value otherwise.

### GetDailyRevenueCapOk

`func (o *CreateOfferRequest) GetDailyRevenueCapOk() (*int32, bool)`

GetDailyRevenueCapOk returns a tuple with the DailyRevenueCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDailyRevenueCap

`func (o *CreateOfferRequest) SetDailyRevenueCap(v int32)`

SetDailyRevenueCap sets DailyRevenueCap field to given value.

### HasDailyRevenueCap

`func (o *CreateOfferRequest) HasDailyRevenueCap() bool`

HasDailyRevenueCap returns a boolean if a field has been set.

### GetWeeklyRevenueCap

`func (o *CreateOfferRequest) GetWeeklyRevenueCap() int32`

GetWeeklyRevenueCap returns the WeeklyRevenueCap field if non-nil, zero value otherwise.

### GetWeeklyRevenueCapOk

`func (o *CreateOfferRequest) GetWeeklyRevenueCapOk() (*int32, bool)`

GetWeeklyRevenueCapOk returns a tuple with the WeeklyRevenueCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWeeklyRevenueCap

`func (o *CreateOfferRequest) SetWeeklyRevenueCap(v int32)`

SetWeeklyRevenueCap sets WeeklyRevenueCap field to given value.

### HasWeeklyRevenueCap

`func (o *CreateOfferRequest) HasWeeklyRevenueCap() bool`

HasWeeklyRevenueCap returns a boolean if a field has been set.

### GetMonthlyRevenueCap

`func (o *CreateOfferRequest) GetMonthlyRevenueCap() int32`

GetMonthlyRevenueCap returns the MonthlyRevenueCap field if non-nil, zero value otherwise.

### GetMonthlyRevenueCapOk

`func (o *CreateOfferRequest) GetMonthlyRevenueCapOk() (*int32, bool)`

GetMonthlyRevenueCapOk returns a tuple with the MonthlyRevenueCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMonthlyRevenueCap

`func (o *CreateOfferRequest) SetMonthlyRevenueCap(v int32)`

SetMonthlyRevenueCap sets MonthlyRevenueCap field to given value.

### HasMonthlyRevenueCap

`func (o *CreateOfferRequest) HasMonthlyRevenueCap() bool`

HasMonthlyRevenueCap returns a boolean if a field has been set.

### GetGlobalRevenueCap

`func (o *CreateOfferRequest) GetGlobalRevenueCap() int32`

GetGlobalRevenueCap returns the GlobalRevenueCap field if non-nil, zero value otherwise.

### GetGlobalRevenueCapOk

`func (o *CreateOfferRequest) GetGlobalRevenueCapOk() (*int32, bool)`

GetGlobalRevenueCapOk returns a tuple with the GlobalRevenueCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGlobalRevenueCap

`func (o *CreateOfferRequest) SetGlobalRevenueCap(v int32)`

SetGlobalRevenueCap sets GlobalRevenueCap field to given value.

### HasGlobalRevenueCap

`func (o *CreateOfferRequest) HasGlobalRevenueCap() bool`

HasGlobalRevenueCap returns a boolean if a field has been set.

### GetDailyClickCap

`func (o *CreateOfferRequest) GetDailyClickCap() int32`

GetDailyClickCap returns the DailyClickCap field if non-nil, zero value otherwise.

### GetDailyClickCapOk

`func (o *CreateOfferRequest) GetDailyClickCapOk() (*int32, bool)`

GetDailyClickCapOk returns a tuple with the DailyClickCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDailyClickCap

`func (o *CreateOfferRequest) SetDailyClickCap(v int32)`

SetDailyClickCap sets DailyClickCap field to given value.

### HasDailyClickCap

`func (o *CreateOfferRequest) HasDailyClickCap() bool`

HasDailyClickCap returns a boolean if a field has been set.

### GetWeeklyClickCap

`func (o *CreateOfferRequest) GetWeeklyClickCap() int32`

GetWeeklyClickCap returns the WeeklyClickCap field if non-nil, zero value otherwise.

### GetWeeklyClickCapOk

`func (o *CreateOfferRequest) GetWeeklyClickCapOk() (*int32, bool)`

GetWeeklyClickCapOk returns a tuple with the WeeklyClickCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWeeklyClickCap

`func (o *CreateOfferRequest) SetWeeklyClickCap(v int32)`

SetWeeklyClickCap sets WeeklyClickCap field to given value.

### HasWeeklyClickCap

`func (o *CreateOfferRequest) HasWeeklyClickCap() bool`

HasWeeklyClickCap returns a boolean if a field has been set.

### GetMonthlyClickCap

`func (o *CreateOfferRequest) GetMonthlyClickCap() int32`

GetMonthlyClickCap returns the MonthlyClickCap field if non-nil, zero value otherwise.

### GetMonthlyClickCapOk

`func (o *CreateOfferRequest) GetMonthlyClickCapOk() (*int32, bool)`

GetMonthlyClickCapOk returns a tuple with the MonthlyClickCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMonthlyClickCap

`func (o *CreateOfferRequest) SetMonthlyClickCap(v int32)`

SetMonthlyClickCap sets MonthlyClickCap field to given value.

### HasMonthlyClickCap

`func (o *CreateOfferRequest) HasMonthlyClickCap() bool`

HasMonthlyClickCap returns a boolean if a field has been set.

### GetGlobalClickCap

`func (o *CreateOfferRequest) GetGlobalClickCap() int32`

GetGlobalClickCap returns the GlobalClickCap field if non-nil, zero value otherwise.

### GetGlobalClickCapOk

`func (o *CreateOfferRequest) GetGlobalClickCapOk() (*int32, bool)`

GetGlobalClickCapOk returns a tuple with the GlobalClickCap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGlobalClickCap

`func (o *CreateOfferRequest) SetGlobalClickCap(v int32)`

SetGlobalClickCap sets GlobalClickCap field to given value.

### HasGlobalClickCap

`func (o *CreateOfferRequest) HasGlobalClickCap() bool`

HasGlobalClickCap returns a boolean if a field has been set.

### GetRedirectMode

`func (o *CreateOfferRequest) GetRedirectMode() string`

GetRedirectMode returns the RedirectMode field if non-nil, zero value otherwise.

### GetRedirectModeOk

`func (o *CreateOfferRequest) GetRedirectModeOk() (*string, bool)`

GetRedirectModeOk returns a tuple with the RedirectMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectMode

`func (o *CreateOfferRequest) SetRedirectMode(v string)`

SetRedirectMode sets RedirectMode field to given value.

### HasRedirectMode

`func (o *CreateOfferRequest) HasRedirectMode() bool`

HasRedirectMode returns a boolean if a field has been set.

### GetIsUsingSuppressionList

`func (o *CreateOfferRequest) GetIsUsingSuppressionList() bool`

GetIsUsingSuppressionList returns the IsUsingSuppressionList field if non-nil, zero value otherwise.

### GetIsUsingSuppressionListOk

`func (o *CreateOfferRequest) GetIsUsingSuppressionListOk() (*bool, bool)`

GetIsUsingSuppressionListOk returns a tuple with the IsUsingSuppressionList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsUsingSuppressionList

`func (o *CreateOfferRequest) SetIsUsingSuppressionList(v bool)`

SetIsUsingSuppressionList sets IsUsingSuppressionList field to given value.

### HasIsUsingSuppressionList

`func (o *CreateOfferRequest) HasIsUsingSuppressionList() bool`

HasIsUsingSuppressionList returns a boolean if a field has been set.

### GetSuppressionListId

`func (o *CreateOfferRequest) GetSuppressionListId() int32`

GetSuppressionListId returns the SuppressionListId field if non-nil, zero value otherwise.

### GetSuppressionListIdOk

`func (o *CreateOfferRequest) GetSuppressionListIdOk() (*int32, bool)`

GetSuppressionListIdOk returns a tuple with the SuppressionListId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuppressionListId

`func (o *CreateOfferRequest) SetSuppressionListId(v int32)`

SetSuppressionListId sets SuppressionListId field to given value.

### HasSuppressionListId

`func (o *CreateOfferRequest) HasSuppressionListId() bool`

HasSuppressionListId returns a boolean if a field has been set.

### GetIsDuplicateFilterEnabled

`func (o *CreateOfferRequest) GetIsDuplicateFilterEnabled() bool`

GetIsDuplicateFilterEnabled returns the IsDuplicateFilterEnabled field if non-nil, zero value otherwise.

### GetIsDuplicateFilterEnabledOk

`func (o *CreateOfferRequest) GetIsDuplicateFilterEnabledOk() (*bool, bool)`

GetIsDuplicateFilterEnabledOk returns a tuple with the IsDuplicateFilterEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsDuplicateFilterEnabled

`func (o *CreateOfferRequest) SetIsDuplicateFilterEnabled(v bool)`

SetIsDuplicateFilterEnabled sets IsDuplicateFilterEnabled field to given value.

### HasIsDuplicateFilterEnabled

`func (o *CreateOfferRequest) HasIsDuplicateFilterEnabled() bool`

HasIsDuplicateFilterEnabled returns a boolean if a field has been set.

### GetDuplicateFilterTargetingAction

`func (o *CreateOfferRequest) GetDuplicateFilterTargetingAction() string`

GetDuplicateFilterTargetingAction returns the DuplicateFilterTargetingAction field if non-nil, zero value otherwise.

### GetDuplicateFilterTargetingActionOk

`func (o *CreateOfferRequest) GetDuplicateFilterTargetingActionOk() (*string, bool)`

GetDuplicateFilterTargetingActionOk returns a tuple with the DuplicateFilterTargetingAction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuplicateFilterTargetingAction

`func (o *CreateOfferRequest) SetDuplicateFilterTargetingAction(v string)`

SetDuplicateFilterTargetingAction sets DuplicateFilterTargetingAction field to given value.

### HasDuplicateFilterTargetingAction

`func (o *CreateOfferRequest) HasDuplicateFilterTargetingAction() bool`

HasDuplicateFilterTargetingAction returns a boolean if a field has been set.

### GetNetworkTrackingDomainId

`func (o *CreateOfferRequest) GetNetworkTrackingDomainId() int32`

GetNetworkTrackingDomainId returns the NetworkTrackingDomainId field if non-nil, zero value otherwise.

### GetNetworkTrackingDomainIdOk

`func (o *CreateOfferRequest) GetNetworkTrackingDomainIdOk() (*int32, bool)`

GetNetworkTrackingDomainIdOk returns a tuple with the NetworkTrackingDomainId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkTrackingDomainId

`func (o *CreateOfferRequest) SetNetworkTrackingDomainId(v int32)`

SetNetworkTrackingDomainId sets NetworkTrackingDomainId field to given value.

### HasNetworkTrackingDomainId

`func (o *CreateOfferRequest) HasNetworkTrackingDomainId() bool`

HasNetworkTrackingDomainId returns a boolean if a field has been set.

### GetIsUseSecureLink

`func (o *CreateOfferRequest) GetIsUseSecureLink() bool`

GetIsUseSecureLink returns the IsUseSecureLink field if non-nil, zero value otherwise.

### GetIsUseSecureLinkOk

`func (o *CreateOfferRequest) GetIsUseSecureLinkOk() (*bool, bool)`

GetIsUseSecureLinkOk returns a tuple with the IsUseSecureLink field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsUseSecureLink

`func (o *CreateOfferRequest) SetIsUseSecureLink(v bool)`

SetIsUseSecureLink sets IsUseSecureLink field to given value.

### HasIsUseSecureLink

`func (o *CreateOfferRequest) HasIsUseSecureLink() bool`

HasIsUseSecureLink returns a boolean if a field has been set.

### GetIsAllowDeepLink

`func (o *CreateOfferRequest) GetIsAllowDeepLink() bool`

GetIsAllowDeepLink returns the IsAllowDeepLink field if non-nil, zero value otherwise.

### GetIsAllowDeepLinkOk

`func (o *CreateOfferRequest) GetIsAllowDeepLinkOk() (*bool, bool)`

GetIsAllowDeepLinkOk returns a tuple with the IsAllowDeepLink field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsAllowDeepLink

`func (o *CreateOfferRequest) SetIsAllowDeepLink(v bool)`

SetIsAllowDeepLink sets IsAllowDeepLink field to given value.

### HasIsAllowDeepLink

`func (o *CreateOfferRequest) HasIsAllowDeepLink() bool`

HasIsAllowDeepLink returns a boolean if a field has been set.

### GetIsSessionTrackingEnabled

`func (o *CreateOfferRequest) GetIsSessionTrackingEnabled() bool`

GetIsSessionTrackingEnabled returns the IsSessionTrackingEnabled field if non-nil, zero value otherwise.

### GetIsSessionTrackingEnabledOk

`func (o *CreateOfferRequest) GetIsSessionTrackingEnabledOk() (*bool, bool)`

GetIsSessionTrackingEnabledOk returns a tuple with the IsSessionTrackingEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsSessionTrackingEnabled

`func (o *CreateOfferRequest) SetIsSessionTrackingEnabled(v bool)`

SetIsSessionTrackingEnabled sets IsSessionTrackingEnabled field to given value.

### HasIsSessionTrackingEnabled

`func (o *CreateOfferRequest) HasIsSessionTrackingEnabled() bool`

HasIsSessionTrackingEnabled returns a boolean if a field has been set.

### GetSessionTrackingLifespanHour

`func (o *CreateOfferRequest) GetSessionTrackingLifespanHour() int32`

GetSessionTrackingLifespanHour returns the SessionTrackingLifespanHour field if non-nil, zero value otherwise.

### GetSessionTrackingLifespanHourOk

`func (o *CreateOfferRequest) GetSessionTrackingLifespanHourOk() (*int32, bool)`

GetSessionTrackingLifespanHourOk returns a tuple with the SessionTrackingLifespanHour field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSessionTrackingLifespanHour

`func (o *CreateOfferRequest) SetSessionTrackingLifespanHour(v int32)`

SetSessionTrackingLifespanHour sets SessionTrackingLifespanHour field to given value.

### HasSessionTrackingLifespanHour

`func (o *CreateOfferRequest) HasSessionTrackingLifespanHour() bool`

HasSessionTrackingLifespanHour returns a boolean if a field has been set.

### GetSessionTrackingMinimumLifespanSecond

`func (o *CreateOfferRequest) GetSessionTrackingMinimumLifespanSecond() int32`

GetSessionTrackingMinimumLifespanSecond returns the SessionTrackingMinimumLifespanSecond field if non-nil, zero value otherwise.

### GetSessionTrackingMinimumLifespanSecondOk

`func (o *CreateOfferRequest) GetSessionTrackingMinimumLifespanSecondOk() (*int32, bool)`

GetSessionTrackingMinimumLifespanSecondOk returns a tuple with the SessionTrackingMinimumLifespanSecond field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSessionTrackingMinimumLifespanSecond

`func (o *CreateOfferRequest) SetSessionTrackingMinimumLifespanSecond(v int32)`

SetSessionTrackingMinimumLifespanSecond sets SessionTrackingMinimumLifespanSecond field to given value.

### HasSessionTrackingMinimumLifespanSecond

`func (o *CreateOfferRequest) HasSessionTrackingMinimumLifespanSecond() bool`

HasSessionTrackingMinimumLifespanSecond returns a boolean if a field has been set.

### GetIsViewThroughSessionTrackingEnabled

`func (o *CreateOfferRequest) GetIsViewThroughSessionTrackingEnabled() bool`

GetIsViewThroughSessionTrackingEnabled returns the IsViewThroughSessionTrackingEnabled field if non-nil, zero value otherwise.

### GetIsViewThroughSessionTrackingEnabledOk

`func (o *CreateOfferRequest) GetIsViewThroughSessionTrackingEnabledOk() (*bool, bool)`

GetIsViewThroughSessionTrackingEnabledOk returns a tuple with the IsViewThroughSessionTrackingEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsViewThroughSessionTrackingEnabled

`func (o *CreateOfferRequest) SetIsViewThroughSessionTrackingEnabled(v bool)`

SetIsViewThroughSessionTrackingEnabled sets IsViewThroughSessionTrackingEnabled field to given value.

### HasIsViewThroughSessionTrackingEnabled

`func (o *CreateOfferRequest) HasIsViewThroughSessionTrackingEnabled() bool`

HasIsViewThroughSessionTrackingEnabled returns a boolean if a field has been set.

### GetViewThroughSessionTrackingLifespanMinute

`func (o *CreateOfferRequest) GetViewThroughSessionTrackingLifespanMinute() int32`

GetViewThroughSessionTrackingLifespanMinute returns the ViewThroughSessionTrackingLifespanMinute field if non-nil, zero value otherwise.

### GetViewThroughSessionTrackingLifespanMinuteOk

`func (o *CreateOfferRequest) GetViewThroughSessionTrackingLifespanMinuteOk() (*int32, bool)`

GetViewThroughSessionTrackingLifespanMinuteOk returns a tuple with the ViewThroughSessionTrackingLifespanMinute field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetViewThroughSessionTrackingLifespanMinute

`func (o *CreateOfferRequest) SetViewThroughSessionTrackingLifespanMinute(v int32)`

SetViewThroughSessionTrackingLifespanMinute sets ViewThroughSessionTrackingLifespanMinute field to given value.

### HasViewThroughSessionTrackingLifespanMinute

`func (o *CreateOfferRequest) HasViewThroughSessionTrackingLifespanMinute() bool`

HasViewThroughSessionTrackingLifespanMinute returns a boolean if a field has been set.

### GetViewThroughSessionTrackingMinimalLifespanSecond

`func (o *CreateOfferRequest) GetViewThroughSessionTrackingMinimalLifespanSecond() int32`

GetViewThroughSessionTrackingMinimalLifespanSecond returns the ViewThroughSessionTrackingMinimalLifespanSecond field if non-nil, zero value otherwise.

### GetViewThroughSessionTrackingMinimalLifespanSecondOk

`func (o *CreateOfferRequest) GetViewThroughSessionTrackingMinimalLifespanSecondOk() (*int32, bool)`

GetViewThroughSessionTrackingMinimalLifespanSecondOk returns a tuple with the ViewThroughSessionTrackingMinimalLifespanSecond field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetViewThroughSessionTrackingMinimalLifespanSecond

`func (o *CreateOfferRequest) SetViewThroughSessionTrackingMinimalLifespanSecond(v int32)`

SetViewThroughSessionTrackingMinimalLifespanSecond sets ViewThroughSessionTrackingMinimalLifespanSecond field to given value.

### HasViewThroughSessionTrackingMinimalLifespanSecond

`func (o *CreateOfferRequest) HasViewThroughSessionTrackingMinimalLifespanSecond() bool`

HasViewThroughSessionTrackingMinimalLifespanSecond returns a boolean if a field has been set.

### GetIsBlockAlreadyConverted

`func (o *CreateOfferRequest) GetIsBlockAlreadyConverted() bool`

GetIsBlockAlreadyConverted returns the IsBlockAlreadyConverted field if non-nil, zero value otherwise.

### GetIsBlockAlreadyConvertedOk

`func (o *CreateOfferRequest) GetIsBlockAlreadyConvertedOk() (*bool, bool)`

GetIsBlockAlreadyConvertedOk returns a tuple with the IsBlockAlreadyConverted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsBlockAlreadyConverted

`func (o *CreateOfferRequest) SetIsBlockAlreadyConverted(v bool)`

SetIsBlockAlreadyConverted sets IsBlockAlreadyConverted field to given value.

### HasIsBlockAlreadyConverted

`func (o *CreateOfferRequest) HasIsBlockAlreadyConverted() bool`

HasIsBlockAlreadyConverted returns a boolean if a field has been set.

### GetAlreadyConvertedAction

`func (o *CreateOfferRequest) GetAlreadyConvertedAction() string`

GetAlreadyConvertedAction returns the AlreadyConvertedAction field if non-nil, zero value otherwise.

### GetAlreadyConvertedActionOk

`func (o *CreateOfferRequest) GetAlreadyConvertedActionOk() (*string, bool)`

GetAlreadyConvertedActionOk returns a tuple with the AlreadyConvertedAction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlreadyConvertedAction

`func (o *CreateOfferRequest) SetAlreadyConvertedAction(v string)`

SetAlreadyConvertedAction sets AlreadyConvertedAction field to given value.

### HasAlreadyConvertedAction

`func (o *CreateOfferRequest) HasAlreadyConvertedAction() bool`

HasAlreadyConvertedAction returns a boolean if a field has been set.

### GetVisibility

`func (o *CreateOfferRequest) GetVisibility() string`

GetVisibility returns the Visibility field if non-nil, zero value otherwise.

### GetVisibilityOk

`func (o *CreateOfferRequest) GetVisibilityOk() (*string, bool)`

GetVisibilityOk returns a tuple with the Visibility field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVisibility

`func (o *CreateOfferRequest) SetVisibility(v string)`

SetVisibility sets Visibility field to given value.

### HasVisibility

`func (o *CreateOfferRequest) HasVisibility() bool`

HasVisibility returns a boolean if a field has been set.

### GetConversionMethod

`func (o *CreateOfferRequest) GetConversionMethod() string`

GetConversionMethod returns the ConversionMethod field if non-nil, zero value otherwise.

### GetConversionMethodOk

`func (o *CreateOfferRequest) GetConversionMethodOk() (*string, bool)`

GetConversionMethodOk returns a tuple with the ConversionMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConversionMethod

`func (o *CreateOfferRequest) SetConversionMethod(v string)`

SetConversionMethod sets ConversionMethod field to given value.

### HasConversionMethod

`func (o *CreateOfferRequest) HasConversionMethod() bool`

HasConversionMethod returns a boolean if a field has been set.

### GetIsWhitelistCheckEnabled

`func (o *CreateOfferRequest) GetIsWhitelistCheckEnabled() bool`

GetIsWhitelistCheckEnabled returns the IsWhitelistCheckEnabled field if non-nil, zero value otherwise.

### GetIsWhitelistCheckEnabledOk

`func (o *CreateOfferRequest) GetIsWhitelistCheckEnabledOk() (*bool, bool)`

GetIsWhitelistCheckEnabledOk returns a tuple with the IsWhitelistCheckEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsWhitelistCheckEnabled

`func (o *CreateOfferRequest) SetIsWhitelistCheckEnabled(v bool)`

SetIsWhitelistCheckEnabled sets IsWhitelistCheckEnabled field to given value.

### HasIsWhitelistCheckEnabled

`func (o *CreateOfferRequest) HasIsWhitelistCheckEnabled() bool`

HasIsWhitelistCheckEnabled returns a boolean if a field has been set.

### GetIsUseScrubRate

`func (o *CreateOfferRequest) GetIsUseScrubRate() bool`

GetIsUseScrubRate returns the IsUseScrubRate field if non-nil, zero value otherwise.

### GetIsUseScrubRateOk

`func (o *CreateOfferRequest) GetIsUseScrubRateOk() (*bool, bool)`

GetIsUseScrubRateOk returns a tuple with the IsUseScrubRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsUseScrubRate

`func (o *CreateOfferRequest) SetIsUseScrubRate(v bool)`

SetIsUseScrubRate sets IsUseScrubRate field to given value.

### HasIsUseScrubRate

`func (o *CreateOfferRequest) HasIsUseScrubRate() bool`

HasIsUseScrubRate returns a boolean if a field has been set.

### GetScrubRateStatus

`func (o *CreateOfferRequest) GetScrubRateStatus() string`

GetScrubRateStatus returns the ScrubRateStatus field if non-nil, zero value otherwise.

### GetScrubRateStatusOk

`func (o *CreateOfferRequest) GetScrubRateStatusOk() (*string, bool)`

GetScrubRateStatusOk returns a tuple with the ScrubRateStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScrubRateStatus

`func (o *CreateOfferRequest) SetScrubRateStatus(v string)`

SetScrubRateStatus sets ScrubRateStatus field to given value.

### HasScrubRateStatus

`func (o *CreateOfferRequest) HasScrubRateStatus() bool`

HasScrubRateStatus returns a boolean if a field has been set.

### GetScrubRatePercentage

`func (o *CreateOfferRequest) GetScrubRatePercentage() int32`

GetScrubRatePercentage returns the ScrubRatePercentage field if non-nil, zero value otherwise.

### GetScrubRatePercentageOk

`func (o *CreateOfferRequest) GetScrubRatePercentageOk() (*int32, bool)`

GetScrubRatePercentageOk returns a tuple with the ScrubRatePercentage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScrubRatePercentage

`func (o *CreateOfferRequest) SetScrubRatePercentage(v int32)`

SetScrubRatePercentage sets ScrubRatePercentage field to given value.

### HasScrubRatePercentage

`func (o *CreateOfferRequest) HasScrubRatePercentage() bool`

HasScrubRatePercentage returns a boolean if a field has been set.

### GetSessionDefinition

`func (o *CreateOfferRequest) GetSessionDefinition() string`

GetSessionDefinition returns the SessionDefinition field if non-nil, zero value otherwise.

### GetSessionDefinitionOk

`func (o *CreateOfferRequest) GetSessionDefinitionOk() (*string, bool)`

GetSessionDefinitionOk returns a tuple with the SessionDefinition field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSessionDefinition

`func (o *CreateOfferRequest) SetSessionDefinition(v string)`

SetSessionDefinition sets SessionDefinition field to given value.

### HasSessionDefinition

`func (o *CreateOfferRequest) HasSessionDefinition() bool`

HasSessionDefinition returns a boolean if a field has been set.

### GetSessionDuration

`func (o *CreateOfferRequest) GetSessionDuration() int32`

GetSessionDuration returns the SessionDuration field if non-nil, zero value otherwise.

### GetSessionDurationOk

`func (o *CreateOfferRequest) GetSessionDurationOk() (*int32, bool)`

GetSessionDurationOk returns a tuple with the SessionDuration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSessionDuration

`func (o *CreateOfferRequest) SetSessionDuration(v int32)`

SetSessionDuration sets SessionDuration field to given value.

### HasSessionDuration

`func (o *CreateOfferRequest) HasSessionDuration() bool`

HasSessionDuration returns a boolean if a field has been set.

### GetAppIdentifier

`func (o *CreateOfferRequest) GetAppIdentifier() string`

GetAppIdentifier returns the AppIdentifier field if non-nil, zero value otherwise.

### GetAppIdentifierOk

`func (o *CreateOfferRequest) GetAppIdentifierOk() (*string, bool)`

GetAppIdentifierOk returns a tuple with the AppIdentifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAppIdentifier

`func (o *CreateOfferRequest) SetAppIdentifier(v string)`

SetAppIdentifier sets AppIdentifier field to given value.

### HasAppIdentifier

`func (o *CreateOfferRequest) HasAppIdentifier() bool`

HasAppIdentifier returns a boolean if a field has been set.

### GetIsDescriptionPlainText

`func (o *CreateOfferRequest) GetIsDescriptionPlainText() bool`

GetIsDescriptionPlainText returns the IsDescriptionPlainText field if non-nil, zero value otherwise.

### GetIsDescriptionPlainTextOk

`func (o *CreateOfferRequest) GetIsDescriptionPlainTextOk() (*bool, bool)`

GetIsDescriptionPlainTextOk returns a tuple with the IsDescriptionPlainText field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsDescriptionPlainText

`func (o *CreateOfferRequest) SetIsDescriptionPlainText(v bool)`

SetIsDescriptionPlainText sets IsDescriptionPlainText field to given value.

### HasIsDescriptionPlainText

`func (o *CreateOfferRequest) HasIsDescriptionPlainText() bool`

HasIsDescriptionPlainText returns a boolean if a field has been set.

### GetIsUseDirectLinking

`func (o *CreateOfferRequest) GetIsUseDirectLinking() bool`

GetIsUseDirectLinking returns the IsUseDirectLinking field if non-nil, zero value otherwise.

### GetIsUseDirectLinkingOk

`func (o *CreateOfferRequest) GetIsUseDirectLinkingOk() (*bool, bool)`

GetIsUseDirectLinkingOk returns a tuple with the IsUseDirectLinking field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsUseDirectLinking

`func (o *CreateOfferRequest) SetIsUseDirectLinking(v bool)`

SetIsUseDirectLinking sets IsUseDirectLinking field to given value.

### HasIsUseDirectLinking

`func (o *CreateOfferRequest) HasIsUseDirectLinking() bool`

HasIsUseDirectLinking returns a boolean if a field has been set.

### GetIsFailTrafficEnabled

`func (o *CreateOfferRequest) GetIsFailTrafficEnabled() bool`

GetIsFailTrafficEnabled returns the IsFailTrafficEnabled field if non-nil, zero value otherwise.

### GetIsFailTrafficEnabledOk

`func (o *CreateOfferRequest) GetIsFailTrafficEnabledOk() (*bool, bool)`

GetIsFailTrafficEnabledOk returns a tuple with the IsFailTrafficEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsFailTrafficEnabled

`func (o *CreateOfferRequest) SetIsFailTrafficEnabled(v bool)`

SetIsFailTrafficEnabled sets IsFailTrafficEnabled field to given value.

### HasIsFailTrafficEnabled

`func (o *CreateOfferRequest) HasIsFailTrafficEnabled() bool`

HasIsFailTrafficEnabled returns a boolean if a field has been set.

### GetRedirectRoutingMethod

`func (o *CreateOfferRequest) GetRedirectRoutingMethod() string`

GetRedirectRoutingMethod returns the RedirectRoutingMethod field if non-nil, zero value otherwise.

### GetRedirectRoutingMethodOk

`func (o *CreateOfferRequest) GetRedirectRoutingMethodOk() (*string, bool)`

GetRedirectRoutingMethodOk returns a tuple with the RedirectRoutingMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectRoutingMethod

`func (o *CreateOfferRequest) SetRedirectRoutingMethod(v string)`

SetRedirectRoutingMethod sets RedirectRoutingMethod field to given value.

### HasRedirectRoutingMethod

`func (o *CreateOfferRequest) HasRedirectRoutingMethod() bool`

HasRedirectRoutingMethod returns a boolean if a field has been set.

### GetRedirectInternalRoutingType

`func (o *CreateOfferRequest) GetRedirectInternalRoutingType() string`

GetRedirectInternalRoutingType returns the RedirectInternalRoutingType field if non-nil, zero value otherwise.

### GetRedirectInternalRoutingTypeOk

`func (o *CreateOfferRequest) GetRedirectInternalRoutingTypeOk() (*string, bool)`

GetRedirectInternalRoutingTypeOk returns a tuple with the RedirectInternalRoutingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectInternalRoutingType

`func (o *CreateOfferRequest) SetRedirectInternalRoutingType(v string)`

SetRedirectInternalRoutingType sets RedirectInternalRoutingType field to given value.

### HasRedirectInternalRoutingType

`func (o *CreateOfferRequest) HasRedirectInternalRoutingType() bool`

HasRedirectInternalRoutingType returns a boolean if a field has been set.

### GetMeta

`func (o *CreateOfferRequest) GetMeta() Meta`

GetMeta returns the Meta field if non-nil, zero value otherwise.

### GetMetaOk

`func (o *CreateOfferRequest) GetMetaOk() (*Meta, bool)`

GetMetaOk returns a tuple with the Meta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMeta

`func (o *CreateOfferRequest) SetMeta(v Meta)`

SetMeta sets Meta field to given value.

### HasMeta

`func (o *CreateOfferRequest) HasMeta() bool`

HasMeta returns a boolean if a field has been set.

### GetCreatives

`func (o *CreateOfferRequest) GetCreatives() []Creative`

GetCreatives returns the Creatives field if non-nil, zero value otherwise.

### GetCreativesOk

`func (o *CreateOfferRequest) GetCreativesOk() (*[]Creative, bool)`

GetCreativesOk returns a tuple with the Creatives field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatives

`func (o *CreateOfferRequest) SetCreatives(v []Creative)`

SetCreatives sets Creatives field to given value.

### HasCreatives

`func (o *CreateOfferRequest) HasCreatives() bool`

HasCreatives returns a boolean if a field has been set.

### GetInternalRedirects

`func (o *CreateOfferRequest) GetInternalRedirects() []InternalRedirect`

GetInternalRedirects returns the InternalRedirects field if non-nil, zero value otherwise.

### GetInternalRedirectsOk

`func (o *CreateOfferRequest) GetInternalRedirectsOk() (*[]InternalRedirect, bool)`

GetInternalRedirectsOk returns a tuple with the InternalRedirects field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternalRedirects

`func (o *CreateOfferRequest) SetInternalRedirects(v []InternalRedirect)`

SetInternalRedirects sets InternalRedirects field to given value.

### HasInternalRedirects

`func (o *CreateOfferRequest) HasInternalRedirects() bool`

HasInternalRedirects returns a boolean if a field has been set.

### GetRuleset

`func (o *CreateOfferRequest) GetRuleset() Ruleset`

GetRuleset returns the Ruleset field if non-nil, zero value otherwise.

### GetRulesetOk

`func (o *CreateOfferRequest) GetRulesetOk() (*Ruleset, bool)`

GetRulesetOk returns a tuple with the Ruleset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRuleset

`func (o *CreateOfferRequest) SetRuleset(v Ruleset)`

SetRuleset sets Ruleset field to given value.

### HasRuleset

`func (o *CreateOfferRequest) HasRuleset() bool`

HasRuleset returns a boolean if a field has been set.

### GetTrafficFilters

`func (o *CreateOfferRequest) GetTrafficFilters() []TrafficFilter`

GetTrafficFilters returns the TrafficFilters field if non-nil, zero value otherwise.

### GetTrafficFiltersOk

`func (o *CreateOfferRequest) GetTrafficFiltersOk() (*[]TrafficFilter, bool)`

GetTrafficFiltersOk returns a tuple with the TrafficFilters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrafficFilters

`func (o *CreateOfferRequest) SetTrafficFilters(v []TrafficFilter)`

SetTrafficFilters sets TrafficFilters field to given value.

### HasTrafficFilters

`func (o *CreateOfferRequest) HasTrafficFilters() bool`

HasTrafficFilters returns a boolean if a field has been set.

### GetEmail

`func (o *CreateOfferRequest) GetEmail() EmailSettings`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *CreateOfferRequest) GetEmailOk() (*EmailSettings, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *CreateOfferRequest) SetEmail(v EmailSettings)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *CreateOfferRequest) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetEmailOptout

`func (o *CreateOfferRequest) GetEmailOptout() EmailOptoutSettings`

GetEmailOptout returns the EmailOptout field if non-nil, zero value otherwise.

### GetEmailOptoutOk

`func (o *CreateOfferRequest) GetEmailOptoutOk() (*EmailOptoutSettings, bool)`

GetEmailOptoutOk returns a tuple with the EmailOptout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmailOptout

`func (o *CreateOfferRequest) SetEmailOptout(v EmailOptoutSettings)`

SetEmailOptout sets EmailOptout field to given value.

### HasEmailOptout

`func (o *CreateOfferRequest) HasEmailOptout() bool`

HasEmailOptout returns a boolean if a field has been set.

### GetLabels

`func (o *CreateOfferRequest) GetLabels() []string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *CreateOfferRequest) GetLabelsOk() (*[]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *CreateOfferRequest) SetLabels(v []string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *CreateOfferRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetSourceNames

`func (o *CreateOfferRequest) GetSourceNames() []string`

GetSourceNames returns the SourceNames field if non-nil, zero value otherwise.

### GetSourceNamesOk

`func (o *CreateOfferRequest) GetSourceNamesOk() (*[]string, bool)`

GetSourceNamesOk returns a tuple with the SourceNames field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceNames

`func (o *CreateOfferRequest) SetSourceNames(v []string)`

SetSourceNames sets SourceNames field to given value.

### HasSourceNames

`func (o *CreateOfferRequest) HasSourceNames() bool`

HasSourceNames returns a boolean if a field has been set.

### GetPayoutRevenue

`func (o *CreateOfferRequest) GetPayoutRevenue() []PayoutRevenue`

GetPayoutRevenue returns the PayoutRevenue field if non-nil, zero value otherwise.

### GetPayoutRevenueOk

`func (o *CreateOfferRequest) GetPayoutRevenueOk() (*[]PayoutRevenue, bool)`

GetPayoutRevenueOk returns a tuple with the PayoutRevenue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPayoutRevenue

`func (o *CreateOfferRequest) SetPayoutRevenue(v []PayoutRevenue)`

SetPayoutRevenue sets PayoutRevenue field to given value.


### GetThumbnailFile

`func (o *CreateOfferRequest) GetThumbnailFile() ThumbnailFile`

GetThumbnailFile returns the ThumbnailFile field if non-nil, zero value otherwise.

### GetThumbnailFileOk

`func (o *CreateOfferRequest) GetThumbnailFileOk() (*ThumbnailFile, bool)`

GetThumbnailFileOk returns a tuple with the ThumbnailFile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThumbnailFile

`func (o *CreateOfferRequest) SetThumbnailFile(v ThumbnailFile)`

SetThumbnailFile sets ThumbnailFile field to given value.

### HasThumbnailFile

`func (o *CreateOfferRequest) HasThumbnailFile() bool`

HasThumbnailFile returns a boolean if a field has been set.

### GetIntegrations

`func (o *CreateOfferRequest) GetIntegrations() Integrations`

GetIntegrations returns the Integrations field if non-nil, zero value otherwise.

### GetIntegrationsOk

`func (o *CreateOfferRequest) GetIntegrationsOk() (*Integrations, bool)`

GetIntegrationsOk returns a tuple with the Integrations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntegrations

`func (o *CreateOfferRequest) SetIntegrations(v Integrations)`

SetIntegrations sets Integrations field to given value.

### HasIntegrations

`func (o *CreateOfferRequest) HasIntegrations() bool`

HasIntegrations returns a boolean if a field has been set.

### GetChannels

`func (o *CreateOfferRequest) GetChannels() []Channel`

GetChannels returns the Channels field if non-nil, zero value otherwise.

### GetChannelsOk

`func (o *CreateOfferRequest) GetChannelsOk() (*[]Channel, bool)`

GetChannelsOk returns a tuple with the Channels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannels

`func (o *CreateOfferRequest) SetChannels(v []Channel)`

SetChannels sets Channels field to given value.

### HasChannels

`func (o *CreateOfferRequest) HasChannels() bool`

HasChannels returns a boolean if a field has been set.

### GetRequirementKpis

`func (o *CreateOfferRequest) GetRequirementKpis() []RequirementKPI`

GetRequirementKpis returns the RequirementKpis field if non-nil, zero value otherwise.

### GetRequirementKpisOk

`func (o *CreateOfferRequest) GetRequirementKpisOk() (*[]RequirementKPI, bool)`

GetRequirementKpisOk returns a tuple with the RequirementKpis field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequirementKpis

`func (o *CreateOfferRequest) SetRequirementKpis(v []RequirementKPI)`

SetRequirementKpis sets RequirementKpis field to given value.

### HasRequirementKpis

`func (o *CreateOfferRequest) HasRequirementKpis() bool`

HasRequirementKpis returns a boolean if a field has been set.

### GetRequirementTrackingParameters

`func (o *CreateOfferRequest) GetRequirementTrackingParameters() []RequirementTrackingParameter`

GetRequirementTrackingParameters returns the RequirementTrackingParameters field if non-nil, zero value otherwise.

### GetRequirementTrackingParametersOk

`func (o *CreateOfferRequest) GetRequirementTrackingParametersOk() (*[]RequirementTrackingParameter, bool)`

GetRequirementTrackingParametersOk returns a tuple with the RequirementTrackingParameters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequirementTrackingParameters

`func (o *CreateOfferRequest) SetRequirementTrackingParameters(v []RequirementTrackingParameter)`

SetRequirementTrackingParameters sets RequirementTrackingParameters field to given value.

### HasRequirementTrackingParameters

`func (o *CreateOfferRequest) HasRequirementTrackingParameters() bool`

HasRequirementTrackingParameters returns a boolean if a field has been set.

### GetEmailAttributionMethod

`func (o *CreateOfferRequest) GetEmailAttributionMethod() string`

GetEmailAttributionMethod returns the EmailAttributionMethod field if non-nil, zero value otherwise.

### GetEmailAttributionMethodOk

`func (o *CreateOfferRequest) GetEmailAttributionMethodOk() (*string, bool)`

GetEmailAttributionMethodOk returns a tuple with the EmailAttributionMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmailAttributionMethod

`func (o *CreateOfferRequest) SetEmailAttributionMethod(v string)`

SetEmailAttributionMethod sets EmailAttributionMethod field to given value.

### HasEmailAttributionMethod

`func (o *CreateOfferRequest) HasEmailAttributionMethod() bool`

HasEmailAttributionMethod returns a boolean if a field has been set.

### GetAttributionMethod

`func (o *CreateOfferRequest) GetAttributionMethod() string`

GetAttributionMethod returns the AttributionMethod field if non-nil, zero value otherwise.

### GetAttributionMethodOk

`func (o *CreateOfferRequest) GetAttributionMethodOk() (*string, bool)`

GetAttributionMethodOk returns a tuple with the AttributionMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributionMethod

`func (o *CreateOfferRequest) SetAttributionMethod(v string)`

SetAttributionMethod sets AttributionMethod field to given value.

### HasAttributionMethod

`func (o *CreateOfferRequest) HasAttributionMethod() bool`

HasAttributionMethod returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


