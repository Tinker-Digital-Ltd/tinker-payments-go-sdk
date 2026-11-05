package api

import (
	"encoding/json"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/auth"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model/dto"
	"net/url"
)

type TransactionManager struct {
	*BaseManager
}

func NewTransactionManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *TransactionManager {
	return &TransactionManager{
		BaseManager: NewBaseManager(cfg, httpClient, authManager),
	}
}

func (tm *TransactionManager) Initiate(request *dto.InitiatePaymentRequestDto) (*model.Transaction, error) {
	payload := request.ToMap()
	response, err := tm.request("POST", config.PAYMENT_INITIATE_PATH, payload)
	if err != nil {
		return nil, err
	}

	return model.NewTransaction(response), nil
}

func (tm *TransactionManager) Query(request *dto.QueryPaymentRequestDto) (*model.Transaction, error) {
	payload := request.ToMap()
	response, err := tm.request("POST", config.PAYMENT_QUERY_PATH, payload)
	if err != nil {
		return nil, err
	}

	return model.NewTransaction(response), nil
}

func (tm *TransactionManager) Capabilities(currency, country string) ([]model.ProviderCapability, error) {
	query := url.Values{}
	if currency != "" {
		query.Set("currency", currency)
	}
	if country != "" {
		query.Set("country", country)
	}
	path := config.PAYMENT_CAPABILITIES_PATH
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	response, err := tm.request("GET", path, nil)
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(response["providers"])
	if err != nil {
		return nil, err
	}
	var providers []model.ProviderCapability
	return providers, json.Unmarshal(raw, &providers)
}
