package tinker

import (
	"testing"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
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
}

func TestPayments_Managers(t *testing.T) {
	payments := NewPayments("public-key", "secret-key", nil)
	if payments.Transactions() == nil {
		t.Fatal("Transactions() returned nil")
	}
	if payments.Subscriptions() == nil {
		t.Fatal("Subscriptions() returned nil")
	}
	if payments.Webhooks() == nil {
		t.Fatal("Webhooks() returned nil")
	}
}

type mockHttpClient struct{}

func (m *mockHttpClient) Get(url string, headers map[string]string) (*http.Response, error) {
	return http.NewResponse(200, []byte(`{}`), nil), nil
}

func (m *mockHttpClient) Post(url string, headers map[string]string, body []byte) (*http.Response, error) {
	return http.NewResponse(200, []byte(`{}`), nil), nil
}
