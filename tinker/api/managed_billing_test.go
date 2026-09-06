package api

import (
	"encoding/json"
	"errors"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/auth"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	sdkerrors "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/errors"
	sdkhttp "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
	"strings"
	"testing"
)

type managedClient struct {
	mockHttpClient
	endpoint string
	payload  map[string]interface{}
}

func (c *managedClient) Post(endpoint string, headers map[string]string, body []byte) (*sdkhttp.Response, error) {
	if !strings.Contains(endpoint, "/auth/token") {
		c.endpoint = endpoint
		_ = json.Unmarshal(body, &c.payload)
	}
	return c.mockHttpClient.Post(endpoint, headers, body)
}
func TestManagedBillingExactAnnualAndSeatCommands(t *testing.T) {
	cfg := config.NewConfiguration("public", "secret")
	c := &managedClient{}
	a := auth.NewManager(cfg, c)
	s := NewSubscriptionManager(cfg, c, a)
	_, e := s.Checkout(ManagedCheckout{ExternalCustomerID: "workspace", UnitAmount: 15000, Quantity: 3, Interval: "year", Currency: "EUR", IdempotencyKey: "same-operation-key"})
	if e != nil {
		t.Fatal(e)
	}
	if !strings.HasSuffix(c.endpoint, "/merchant/billing/checkout") || c.payload["unit_amount"] != float64(15000) || c.payload["quantity"] != float64(3) || c.payload["interval"] != "year" {
		t.Fatalf("incorrect annual checkout: %#v", c.payload)
	}
	_, e = s.ChangeQuantity("sub_test", 4, "increase-operation")
	if e != nil {
		t.Fatal(e)
	}
	if !strings.HasSuffix(c.endpoint, "/sub_test/quantity") || c.payload["quantity"] != float64(4) || c.payload["idempotency_key"] != "increase-operation" {
		t.Fatalf("incorrect increase: %#v", c.payload)
	}
	_, e = s.SetCancelAtRenewal("sub_test", false, "resume-operation")
	if e != nil {
		t.Fatal(e)
	}
	if c.payload["cancel"] != false {
		t.Fatal("resume must preserve explicit false")
	}
}

func TestManagedBillingErrorsPreserveRecoveryInformation(t *testing.T) {
	bm := &BaseManager{}
	resp := sdkhttp.NewResponse(429, []byte(`{}`), map[string][]string{"Retry-After": {"30"}, "X-Request-Id": {"req_api"}})
	var failure *sdkerrors.ApiException
	err := bm.apiError(resp, map[string]interface{}{"error": map[string]interface{}{"message": "Try later", "code": "RATE_LIMITED", "provider_code": "rate_limit", "outcome": "not_applied", "request_id": "req_provider"}})
	if !errors.As(err, &failure) || failure.HTTPStatus != 429 || failure.ErrorCode != "RATE_LIMITED" || failure.ProviderCode != "rate_limit" || failure.Outcome != "not_applied" || failure.RetryAfter != "30" || failure.RequestID != "req_provider" {
		t.Fatalf("missing metadata: %#v", failure)
	}
	err = bm.apiError(sdkhttp.NewResponse(502, nil, nil), nil)
	if !errors.As(err, &failure) || failure.HTTPStatus != 502 || failure.Outcome != "" {
		t.Fatal("unstructured failures must remain uncertain")
	}
}

func TestPaymentMethodPortalHasNoSeatOrPriceMutation(t *testing.T) {
	cfg := config.NewConfiguration("public", "secret")
	c := &managedClient{}
	s := NewSubscriptionManager(cfg, c, auth.NewManager(cfg, c))
	_, e := s.PaymentMethodPortal("sub_test", "https://talon.example/dashboard/billing", "payment-method-key")
	if e != nil {
		t.Fatal(e)
	}
	if !strings.HasSuffix(c.endpoint, "/sub_test/payment-method") || c.payload["return_url"] != "https://talon.example/dashboard/billing" || len(c.payload) != 2 {
		t.Fatalf("unexpected payment-method payload: %#v", c.payload)
	}
}
