package tinker

import (
	"testing"

	"github.com/tinker/tinker-payments-go-sdk/tinker/http"
)

func TestNewPayments(t *testing.T) {
	payments := NewPayments("public-key", "secret-key", nil)
	if payments == nil {
		t.Fatal("NewPayments returned nil")
	}
	if payments.config == nil {
		t.Error("config should not be nil")
	}
	if payments.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
	if payments.authManager == nil {
		t.Error("authManager should not be nil")
	}
}

func TestNewPaymentsWithCustomClient(t *testing.T) {
	mockClient := &mockHttpClient{}
	payments := NewPayments("public-key", "secret-key", mockClient)
	if payments == nil {
		t.Fatal("NewPayments returned nil")
	}
	if payments.httpClient != mockClient {
		t.Error("custom httpClient was not set")
	}
}

func TestPayments_Transactions(t *testing.T) {
	payments := NewPayments("public-key", "secret-key", nil)
	manager := payments.Transactions()
	if manager == nil {
		t.Fatal("Transactions() returned nil")
	}
	if payments.transactionManager == nil {
		t.Error("transactionManager should be initialized")
	}
	manager2 := payments.Transactions()
	if manager != manager2 {
		t.Error("Transactions() should return the same instance")
	}
}

func TestPayments_Webhooks(t *testing.T) {
	payments := NewPayments("public-key", "secret-key", nil)
	handler := payments.Webhooks()
	if handler == nil {
		t.Fatal("Webhooks() returned nil")
	}
	if payments.webhookHandler == nil {
		t.Error("webhookHandler should be initialized")
	}
	handler2 := payments.Webhooks()
	if handler != handler2 {
		t.Error("Webhooks() should return the same instance")
	}
}

type mockHttpClient struct{}

func (m *mockHttpClient) Post(url string, headers map[string]string, body []byte) (*http.Response, error) {
	return http.NewResponse(200, []byte(`{}`), nil), nil
}
