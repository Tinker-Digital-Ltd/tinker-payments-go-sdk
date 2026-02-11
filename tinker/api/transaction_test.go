package api

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/auth"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model/dto"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/types"
)

func TestTransactionManager_Initiate(t *testing.T) {
	cfg := config.NewConfiguration("public-key", "secret-key")
	mockClient := &mockHttpClient{response: map[string]interface{}{"paymentReference": "TXN-123", "authorizationUrl": "https://example.com/auth", "status": "pending"}}
	authMgr := auth.NewManager(cfg, mockClient)
	manager := NewTransactionManager(cfg, mockClient, authMgr)
	request := &dto.InitiatePaymentRequestDto{Amount: 100.00, Currency: "KES", Gateway: types.MPESA, MerchantReference: "ORDER-123", ReturnURL: "https://example.com/return"}
	transaction, err := manager.Initiate(request)
	if err != nil {
		t.Fatalf("Initiate() error = %v", err)
	}
	if transaction == nil || transaction.InitiationData == nil {
		t.Fatal("InitiationData should not be nil")
	}
}

func TestTransactionManager_Query(t *testing.T) {
	cfg := config.NewConfiguration("public-key", "secret-key")
	mockClient := &mockHttpClient{response: map[string]interface{}{"id": "123", "reference": "TXN-123", "amount": 100.00, "currency": "KES", "status": "success"}}
	authMgr := auth.NewManager(cfg, mockClient)
	manager := NewTransactionManager(cfg, mockClient, authMgr)
	request := &dto.QueryPaymentRequestDto{PaymentReference: "TXN-123", Gateway: types.MPESA}
	transaction, err := manager.Query(request)
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if transaction == nil || transaction.QueryData == nil {
		t.Fatal("QueryData should not be nil")
	}
}

type mockHttpClient struct {
	response map[string]interface{}
	err      error
}

func (m *mockHttpClient) Get(url string, headers map[string]string) (*http.Response, error) {
	return http.NewResponse(200, []byte(`{}`), nil), nil
}

func (m *mockHttpClient) Post(url string, headers map[string]string, body []byte) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	if strings.Contains(url, "/auth/token") {
		authResponse := map[string]interface{}{"token": "test-token", "expires_in": 3600}
		bodyBytes, _ := json.Marshal(authResponse)
		return http.NewResponse(200, bodyBytes, nil), nil
	}
	responseBody := `{}`
	if m.response != nil {
		bodyBytes, _ := json.Marshal(m.response)
		responseBody = string(bodyBytes)
	}
	return http.NewResponse(200, []byte(responseBody), nil), nil
}
