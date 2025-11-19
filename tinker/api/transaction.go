package api

import (
	"github.com/tinker/tinker-payments-go-sdk/tinker/auth"
	"github.com/tinker/tinker-payments-go-sdk/tinker/config"
	"github.com/tinker/tinker-payments-go-sdk/tinker/http"
	"github.com/tinker/tinker-payments-go-sdk/tinker/model"
	"github.com/tinker/tinker-payments-go-sdk/tinker/model/dto"
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
