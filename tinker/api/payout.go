package api

import (
	"fmt"
	"math"
	"net/url"
	"strconv"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/auth"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model"
)

// PayoutManager withdraws the merchant wallet to its own verified destinations.
type PayoutManager struct{ *BaseManager }

func NewPayoutManager(cfg *config.Configuration, client http.Client, authentication *auth.Manager) *PayoutManager {
	return &PayoutManager{NewBaseManager(cfg, client, authentication)}
}
func (pm *PayoutManager) Methods() ([]model.PayoutMethod, error) {
	result, err := lifecycleRequest[[]model.PayoutMethod](pm.BaseManager, "GET", "merchant/payout-methods", nil, "")
	if err != nil {
		return nil, err
	}
	return *result, nil
}
func (pm *PayoutManager) Quote(request model.PayoutQuoteRequest) (*model.PayoutQuote, error) {
	if err := positiveAmount(request.Amount); err != nil {
		return nil, err
	}
	query := url.Values{"payout_method_id": {request.PayoutMethodID}, "amount": {strconv.FormatFloat(request.Amount, 'f', -1, 64)}}
	if request.Currency != nil {
		query.Set("currency", *request.Currency)
	}
	return lifecycleRequest[model.PayoutQuote](pm.BaseManager, "GET", "merchant/payouts/quote?"+query.Encode(), nil, "")
}
func (pm *PayoutManager) Request(request model.PayoutRequest, key string) (*model.PayoutResult, error) {
	if err := positiveAmount(request.Amount); err != nil {
		return nil, err
	}
	if math.IsNaN(request.ExpectedFeeAmount) || math.IsInf(request.ExpectedFeeAmount, 0) || request.ExpectedFeeAmount < 0 {
		return nil, fmt.Errorf("expected_fee_amount must be finite and non-negative")
	}
	return lifecycleRequest[model.PayoutResult](pm.BaseManager, "POST", "merchant/payouts", request, key)
}
func (pm *PayoutManager) List(filters model.PayoutFilters) (*model.PayoutPage, error) {
	query := url.Values{}
	if filters.Page != nil {
		query.Set("page", strconv.Itoa(*filters.Page))
	}
	if filters.PerPage != nil {
		query.Set("per_page", strconv.Itoa(*filters.PerPage))
	}
	if filters.Status != nil {
		query.Set("status", string(*filters.Status))
	}
	if filters.Currency != nil {
		query.Set("currency", *filters.Currency)
	}
	path := "merchant/payouts"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	return lifecycleRequest[model.PayoutPage](pm.BaseManager, "GET", path, nil, "")
}
func (pm *PayoutManager) Get(id string) (*model.Payout, error) {
	return lifecycleRequest[model.Payout](pm.BaseManager, "GET", "merchant/payouts/"+url.PathEscape(id), nil, "")
}
