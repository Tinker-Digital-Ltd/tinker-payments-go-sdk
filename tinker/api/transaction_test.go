package api

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/tinker/tinker-payments-go-sdk/tinker/auth"
	"github.com/tinker/tinker-payments-go-sdk/tinker/config"
	"github.com/tinker/tinker-payments-go-sdk/tinker/http"
	"github.com/tinker/tinker-payments-go-sdk/tinker/model/dto"
	"github.com/tinker/tinker-payments-go-sdk/tinker/types"
)

func TestNewTransactionManager(t *testing.T) {
	cfg := config.NewConfiguration("public-key", "secret-key")
	mockClient := &mockHttpClient{}
	authMgr := auth.NewManager(cfg, mockClient)

	manager := NewTransactionManager(cfg, mockClient, authMgr)
	if manager == nil {
		t.Fatal("NewTransactionManager returned nil")
	}
	if manager.BaseManager == nil {
		t.Error("BaseManager should not be nil")
	}
}

func TestTransactionManager_Initiate(t *testing.T) {
	cfg := config.NewConfiguration("public-key", "secret-key")
	mockClient := &mockHttpClient{
		response: map[string]interface{}{
			"payment_reference": "TXN-123",
			"authorization_url":  "https://example.com/auth",
			"status":             "pending",
		},
	}
	authMgr := auth.NewManager(cfg, mockClient)
	manager := NewTransactionManager(cfg, mockClient, authMgr)

	customerPhone := "+254712345678"
	request := &dto.InitiatePaymentRequestDto{
		Amount:           100.00,
		Currency:         "KES",
		Gateway:          types.MPESA,
		MerchantReference: "ORDER-123",
		ReturnURL:        "https://example.com/return",
		CustomerPhone:    &customerPhone,
	}

	transaction, err := manager.Initiate(request)
	if err != nil {
		t.Fatalf("Initiate() error = %v", err)
	}
	if transaction == nil {
		t.Fatal("Initiate() returned nil transaction")
	}
	if transaction.InitiationData == nil {
		t.Error("InitiationData should not be nil")
	}
}

func TestTransactionManager_Query(t *testing.T) {
	cfg := config.NewConfiguration("public-key", "secret-key")
	mockClient := &mockHttpClient{
		response: map[string]interface{}{
			"id":       "123",
			"reference": "TXN-123",
			"amount":   100.00,
			"currency": "KES",
			"status":   "success",
		},
	}
	authMgr := auth.NewManager(cfg, mockClient)
	manager := NewTransactionManager(cfg, mockClient, authMgr)

	request := &dto.QueryPaymentRequestDto{
		PaymentReference: "TXN-123",
		Gateway:          types.MPESA,
	}

	transaction, err := manager.Query(request)
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if transaction == nil {
		t.Fatal("Query() returned nil transaction")
	}
	if transaction.QueryData == nil {
		t.Error("QueryData should not be nil")
	}
}

type mockHttpClient struct {
	response map[string]interface{}
	err     error
}

func (m *mockHttpClient) Post(url string, headers map[string]string, body []byte) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}

	if strings.Contains(url, "/auth/token") {
		authResponse := map[string]interface{}{
			"token":      "test-token",
			"expires_in": 3600,
		}
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

