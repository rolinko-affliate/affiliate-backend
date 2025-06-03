# ReportingData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Imp** | Pointer to **int32** | Impressions | [optional] 
**TotalClick** | Pointer to **int32** | Total clicks | [optional] 
**UniqueClick** | Pointer to **int32** | Unique clicks | [optional] 
**InvalidClick** | Pointer to **int32** | Invalid clicks | [optional] 
**DuplicateClick** | Pointer to **int32** | Duplicate clicks | [optional] 
**GrossClick** | Pointer to **int32** | Gross clicks | [optional] 
**Ctr** | Pointer to **float32** | Click-through rate | [optional] 
**Cv** | Pointer to **int32** | Conversions | [optional] 
**InvalidCvScrub** | Pointer to **int32** | Invalid conversions scrubbed | [optional] 
**ViewThroughCv** | Pointer to **int32** | View-through conversions | [optional] 
**TotalCv** | Pointer to **int32** | Total conversions | [optional] 
**Event** | Pointer to **int32** | Events | [optional] 
**Cvr** | Pointer to **float32** | Conversion rate | [optional] 
**Evr** | Pointer to **float32** | Event rate | [optional] 
**Cpc** | Pointer to **float32** | Cost per click | [optional] 
**Cpm** | Pointer to **float32** | Cost per mille | [optional] 
**Cpa** | Pointer to **float32** | Cost per acquisition | [optional] 
**Epc** | Pointer to **float32** | Earnings per click | [optional] 
**Rpc** | Pointer to **float32** | Revenue per click | [optional] 
**Rpa** | Pointer to **float32** | Revenue per acquisition | [optional] 
**Rpm** | Pointer to **float32** | Revenue per mille | [optional] 
**Payout** | Pointer to **float32** | Payout amount | [optional] 
**Revenue** | Pointer to **float32** | Revenue amount | [optional] 
**EventRevenue** | Pointer to **float32** | Event revenue | [optional] 
**GrossSales** | Pointer to **float32** | Gross sales | [optional] 
**Profit** | Pointer to **float32** | Profit | [optional] 
**Margin** | Pointer to **float32** | Margin | [optional] 
**Roas** | Pointer to **float32** | Return on ad spend | [optional] 
**AvgSaleValue** | Pointer to **float32** | Average sale value | [optional] 
**MediaBuyingCost** | Pointer to **float32** | Media buying cost | [optional] 

## Methods

### NewReportingData

`func NewReportingData() *ReportingData`

NewReportingData instantiates a new ReportingData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReportingDataWithDefaults

`func NewReportingDataWithDefaults() *ReportingData`

NewReportingDataWithDefaults instantiates a new ReportingData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetImp

`func (o *ReportingData) GetImp() int32`

GetImp returns the Imp field if non-nil, zero value otherwise.

### GetImpOk

`func (o *ReportingData) GetImpOk() (*int32, bool)`

GetImpOk returns a tuple with the Imp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImp

`func (o *ReportingData) SetImp(v int32)`

SetImp sets Imp field to given value.

### HasImp

`func (o *ReportingData) HasImp() bool`

HasImp returns a boolean if a field has been set.

### GetTotalClick

`func (o *ReportingData) GetTotalClick() int32`

GetTotalClick returns the TotalClick field if non-nil, zero value otherwise.

### GetTotalClickOk

`func (o *ReportingData) GetTotalClickOk() (*int32, bool)`

GetTotalClickOk returns a tuple with the TotalClick field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalClick

`func (o *ReportingData) SetTotalClick(v int32)`

SetTotalClick sets TotalClick field to given value.

### HasTotalClick

`func (o *ReportingData) HasTotalClick() bool`

HasTotalClick returns a boolean if a field has been set.

### GetUniqueClick

`func (o *ReportingData) GetUniqueClick() int32`

GetUniqueClick returns the UniqueClick field if non-nil, zero value otherwise.

### GetUniqueClickOk

`func (o *ReportingData) GetUniqueClickOk() (*int32, bool)`

GetUniqueClickOk returns a tuple with the UniqueClick field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUniqueClick

`func (o *ReportingData) SetUniqueClick(v int32)`

SetUniqueClick sets UniqueClick field to given value.

### HasUniqueClick

`func (o *ReportingData) HasUniqueClick() bool`

HasUniqueClick returns a boolean if a field has been set.

### GetInvalidClick

`func (o *ReportingData) GetInvalidClick() int32`

GetInvalidClick returns the InvalidClick field if non-nil, zero value otherwise.

### GetInvalidClickOk

`func (o *ReportingData) GetInvalidClickOk() (*int32, bool)`

GetInvalidClickOk returns a tuple with the InvalidClick field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInvalidClick

`func (o *ReportingData) SetInvalidClick(v int32)`

SetInvalidClick sets InvalidClick field to given value.

### HasInvalidClick

`func (o *ReportingData) HasInvalidClick() bool`

HasInvalidClick returns a boolean if a field has been set.

### GetDuplicateClick

`func (o *ReportingData) GetDuplicateClick() int32`

GetDuplicateClick returns the DuplicateClick field if non-nil, zero value otherwise.

### GetDuplicateClickOk

`func (o *ReportingData) GetDuplicateClickOk() (*int32, bool)`

GetDuplicateClickOk returns a tuple with the DuplicateClick field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuplicateClick

`func (o *ReportingData) SetDuplicateClick(v int32)`

SetDuplicateClick sets DuplicateClick field to given value.

### HasDuplicateClick

`func (o *ReportingData) HasDuplicateClick() bool`

HasDuplicateClick returns a boolean if a field has been set.

### GetGrossClick

`func (o *ReportingData) GetGrossClick() int32`

GetGrossClick returns the GrossClick field if non-nil, zero value otherwise.

### GetGrossClickOk

`func (o *ReportingData) GetGrossClickOk() (*int32, bool)`

GetGrossClickOk returns a tuple with the GrossClick field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrossClick

`func (o *ReportingData) SetGrossClick(v int32)`

SetGrossClick sets GrossClick field to given value.

### HasGrossClick

`func (o *ReportingData) HasGrossClick() bool`

HasGrossClick returns a boolean if a field has been set.

### GetCtr

`func (o *ReportingData) GetCtr() float32`

GetCtr returns the Ctr field if non-nil, zero value otherwise.

### GetCtrOk

`func (o *ReportingData) GetCtrOk() (*float32, bool)`

GetCtrOk returns a tuple with the Ctr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCtr

`func (o *ReportingData) SetCtr(v float32)`

SetCtr sets Ctr field to given value.

### HasCtr

`func (o *ReportingData) HasCtr() bool`

HasCtr returns a boolean if a field has been set.

### GetCv

`func (o *ReportingData) GetCv() int32`

GetCv returns the Cv field if non-nil, zero value otherwise.

### GetCvOk

`func (o *ReportingData) GetCvOk() (*int32, bool)`

GetCvOk returns a tuple with the Cv field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCv

`func (o *ReportingData) SetCv(v int32)`

SetCv sets Cv field to given value.

### HasCv

`func (o *ReportingData) HasCv() bool`

HasCv returns a boolean if a field has been set.

### GetInvalidCvScrub

`func (o *ReportingData) GetInvalidCvScrub() int32`

GetInvalidCvScrub returns the InvalidCvScrub field if non-nil, zero value otherwise.

### GetInvalidCvScrubOk

`func (o *ReportingData) GetInvalidCvScrubOk() (*int32, bool)`

GetInvalidCvScrubOk returns a tuple with the InvalidCvScrub field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInvalidCvScrub

`func (o *ReportingData) SetInvalidCvScrub(v int32)`

SetInvalidCvScrub sets InvalidCvScrub field to given value.

### HasInvalidCvScrub

`func (o *ReportingData) HasInvalidCvScrub() bool`

HasInvalidCvScrub returns a boolean if a field has been set.

### GetViewThroughCv

`func (o *ReportingData) GetViewThroughCv() int32`

GetViewThroughCv returns the ViewThroughCv field if non-nil, zero value otherwise.

### GetViewThroughCvOk

`func (o *ReportingData) GetViewThroughCvOk() (*int32, bool)`

GetViewThroughCvOk returns a tuple with the ViewThroughCv field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetViewThroughCv

`func (o *ReportingData) SetViewThroughCv(v int32)`

SetViewThroughCv sets ViewThroughCv field to given value.

### HasViewThroughCv

`func (o *ReportingData) HasViewThroughCv() bool`

HasViewThroughCv returns a boolean if a field has been set.

### GetTotalCv

`func (o *ReportingData) GetTotalCv() int32`

GetTotalCv returns the TotalCv field if non-nil, zero value otherwise.

### GetTotalCvOk

`func (o *ReportingData) GetTotalCvOk() (*int32, bool)`

GetTotalCvOk returns a tuple with the TotalCv field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalCv

`func (o *ReportingData) SetTotalCv(v int32)`

SetTotalCv sets TotalCv field to given value.

### HasTotalCv

`func (o *ReportingData) HasTotalCv() bool`

HasTotalCv returns a boolean if a field has been set.

### GetEvent

`func (o *ReportingData) GetEvent() int32`

GetEvent returns the Event field if non-nil, zero value otherwise.

### GetEventOk

`func (o *ReportingData) GetEventOk() (*int32, bool)`

GetEventOk returns a tuple with the Event field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvent

`func (o *ReportingData) SetEvent(v int32)`

SetEvent sets Event field to given value.

### HasEvent

`func (o *ReportingData) HasEvent() bool`

HasEvent returns a boolean if a field has been set.

### GetCvr

`func (o *ReportingData) GetCvr() float32`

GetCvr returns the Cvr field if non-nil, zero value otherwise.

### GetCvrOk

`func (o *ReportingData) GetCvrOk() (*float32, bool)`

GetCvrOk returns a tuple with the Cvr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCvr

`func (o *ReportingData) SetCvr(v float32)`

SetCvr sets Cvr field to given value.

### HasCvr

`func (o *ReportingData) HasCvr() bool`

HasCvr returns a boolean if a field has been set.

### GetEvr

`func (o *ReportingData) GetEvr() float32`

GetEvr returns the Evr field if non-nil, zero value otherwise.

### GetEvrOk

`func (o *ReportingData) GetEvrOk() (*float32, bool)`

GetEvrOk returns a tuple with the Evr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvr

`func (o *ReportingData) SetEvr(v float32)`

SetEvr sets Evr field to given value.

### HasEvr

`func (o *ReportingData) HasEvr() bool`

HasEvr returns a boolean if a field has been set.

### GetCpc

`func (o *ReportingData) GetCpc() float32`

GetCpc returns the Cpc field if non-nil, zero value otherwise.

### GetCpcOk

`func (o *ReportingData) GetCpcOk() (*float32, bool)`

GetCpcOk returns a tuple with the Cpc field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpc

`func (o *ReportingData) SetCpc(v float32)`

SetCpc sets Cpc field to given value.

### HasCpc

`func (o *ReportingData) HasCpc() bool`

HasCpc returns a boolean if a field has been set.

### GetCpm

`func (o *ReportingData) GetCpm() float32`

GetCpm returns the Cpm field if non-nil, zero value otherwise.

### GetCpmOk

`func (o *ReportingData) GetCpmOk() (*float32, bool)`

GetCpmOk returns a tuple with the Cpm field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpm

`func (o *ReportingData) SetCpm(v float32)`

SetCpm sets Cpm field to given value.

### HasCpm

`func (o *ReportingData) HasCpm() bool`

HasCpm returns a boolean if a field has been set.

### GetCpa

`func (o *ReportingData) GetCpa() float32`

GetCpa returns the Cpa field if non-nil, zero value otherwise.

### GetCpaOk

`func (o *ReportingData) GetCpaOk() (*float32, bool)`

GetCpaOk returns a tuple with the Cpa field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpa

`func (o *ReportingData) SetCpa(v float32)`

SetCpa sets Cpa field to given value.

### HasCpa

`func (o *ReportingData) HasCpa() bool`

HasCpa returns a boolean if a field has been set.

### GetEpc

`func (o *ReportingData) GetEpc() float32`

GetEpc returns the Epc field if non-nil, zero value otherwise.

### GetEpcOk

`func (o *ReportingData) GetEpcOk() (*float32, bool)`

GetEpcOk returns a tuple with the Epc field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEpc

`func (o *ReportingData) SetEpc(v float32)`

SetEpc sets Epc field to given value.

### HasEpc

`func (o *ReportingData) HasEpc() bool`

HasEpc returns a boolean if a field has been set.

### GetRpc

`func (o *ReportingData) GetRpc() float32`

GetRpc returns the Rpc field if non-nil, zero value otherwise.

### GetRpcOk

`func (o *ReportingData) GetRpcOk() (*float32, bool)`

GetRpcOk returns a tuple with the Rpc field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRpc

`func (o *ReportingData) SetRpc(v float32)`

SetRpc sets Rpc field to given value.

### HasRpc

`func (o *ReportingData) HasRpc() bool`

HasRpc returns a boolean if a field has been set.

### GetRpa

`func (o *ReportingData) GetRpa() float32`

GetRpa returns the Rpa field if non-nil, zero value otherwise.

### GetRpaOk

`func (o *ReportingData) GetRpaOk() (*float32, bool)`

GetRpaOk returns a tuple with the Rpa field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRpa

`func (o *ReportingData) SetRpa(v float32)`

SetRpa sets Rpa field to given value.

### HasRpa

`func (o *ReportingData) HasRpa() bool`

HasRpa returns a boolean if a field has been set.

### GetRpm

`func (o *ReportingData) GetRpm() float32`

GetRpm returns the Rpm field if non-nil, zero value otherwise.

### GetRpmOk

`func (o *ReportingData) GetRpmOk() (*float32, bool)`

GetRpmOk returns a tuple with the Rpm field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRpm

`func (o *ReportingData) SetRpm(v float32)`

SetRpm sets Rpm field to given value.

### HasRpm

`func (o *ReportingData) HasRpm() bool`

HasRpm returns a boolean if a field has been set.

### GetPayout

`func (o *ReportingData) GetPayout() float32`

GetPayout returns the Payout field if non-nil, zero value otherwise.

### GetPayoutOk

`func (o *ReportingData) GetPayoutOk() (*float32, bool)`

GetPayoutOk returns a tuple with the Payout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPayout

`func (o *ReportingData) SetPayout(v float32)`

SetPayout sets Payout field to given value.

### HasPayout

`func (o *ReportingData) HasPayout() bool`

HasPayout returns a boolean if a field has been set.

### GetRevenue

`func (o *ReportingData) GetRevenue() float32`

GetRevenue returns the Revenue field if non-nil, zero value otherwise.

### GetRevenueOk

`func (o *ReportingData) GetRevenueOk() (*float32, bool)`

GetRevenueOk returns a tuple with the Revenue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevenue

`func (o *ReportingData) SetRevenue(v float32)`

SetRevenue sets Revenue field to given value.

### HasRevenue

`func (o *ReportingData) HasRevenue() bool`

HasRevenue returns a boolean if a field has been set.

### GetEventRevenue

`func (o *ReportingData) GetEventRevenue() float32`

GetEventRevenue returns the EventRevenue field if non-nil, zero value otherwise.

### GetEventRevenueOk

`func (o *ReportingData) GetEventRevenueOk() (*float32, bool)`

GetEventRevenueOk returns a tuple with the EventRevenue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEventRevenue

`func (o *ReportingData) SetEventRevenue(v float32)`

SetEventRevenue sets EventRevenue field to given value.

### HasEventRevenue

`func (o *ReportingData) HasEventRevenue() bool`

HasEventRevenue returns a boolean if a field has been set.

### GetGrossSales

`func (o *ReportingData) GetGrossSales() float32`

GetGrossSales returns the GrossSales field if non-nil, zero value otherwise.

### GetGrossSalesOk

`func (o *ReportingData) GetGrossSalesOk() (*float32, bool)`

GetGrossSalesOk returns a tuple with the GrossSales field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrossSales

`func (o *ReportingData) SetGrossSales(v float32)`

SetGrossSales sets GrossSales field to given value.

### HasGrossSales

`func (o *ReportingData) HasGrossSales() bool`

HasGrossSales returns a boolean if a field has been set.

### GetProfit

`func (o *ReportingData) GetProfit() float32`

GetProfit returns the Profit field if non-nil, zero value otherwise.

### GetProfitOk

`func (o *ReportingData) GetProfitOk() (*float32, bool)`

GetProfitOk returns a tuple with the Profit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProfit

`func (o *ReportingData) SetProfit(v float32)`

SetProfit sets Profit field to given value.

### HasProfit

`func (o *ReportingData) HasProfit() bool`

HasProfit returns a boolean if a field has been set.

### GetMargin

`func (o *ReportingData) GetMargin() float32`

GetMargin returns the Margin field if non-nil, zero value otherwise.

### GetMarginOk

`func (o *ReportingData) GetMarginOk() (*float32, bool)`

GetMarginOk returns a tuple with the Margin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMargin

`func (o *ReportingData) SetMargin(v float32)`

SetMargin sets Margin field to given value.

### HasMargin

`func (o *ReportingData) HasMargin() bool`

HasMargin returns a boolean if a field has been set.

### GetRoas

`func (o *ReportingData) GetRoas() float32`

GetRoas returns the Roas field if non-nil, zero value otherwise.

### GetRoasOk

`func (o *ReportingData) GetRoasOk() (*float32, bool)`

GetRoasOk returns a tuple with the Roas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoas

`func (o *ReportingData) SetRoas(v float32)`

SetRoas sets Roas field to given value.

### HasRoas

`func (o *ReportingData) HasRoas() bool`

HasRoas returns a boolean if a field has been set.

### GetAvgSaleValue

`func (o *ReportingData) GetAvgSaleValue() float32`

GetAvgSaleValue returns the AvgSaleValue field if non-nil, zero value otherwise.

### GetAvgSaleValueOk

`func (o *ReportingData) GetAvgSaleValueOk() (*float32, bool)`

GetAvgSaleValueOk returns a tuple with the AvgSaleValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvgSaleValue

`func (o *ReportingData) SetAvgSaleValue(v float32)`

SetAvgSaleValue sets AvgSaleValue field to given value.

### HasAvgSaleValue

`func (o *ReportingData) HasAvgSaleValue() bool`

HasAvgSaleValue returns a boolean if a field has been set.

### GetMediaBuyingCost

`func (o *ReportingData) GetMediaBuyingCost() float32`

GetMediaBuyingCost returns the MediaBuyingCost field if non-nil, zero value otherwise.

### GetMediaBuyingCostOk

`func (o *ReportingData) GetMediaBuyingCostOk() (*float32, bool)`

GetMediaBuyingCostOk returns a tuple with the MediaBuyingCost field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMediaBuyingCost

`func (o *ReportingData) SetMediaBuyingCost(v float32)`

SetMediaBuyingCost sets MediaBuyingCost field to given value.

### HasMediaBuyingCost

`func (o *ReportingData) HasMediaBuyingCost() bool`

HasMediaBuyingCost returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


