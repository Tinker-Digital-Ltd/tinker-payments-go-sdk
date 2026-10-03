package api

import (
	"errors"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/auth"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	sdkerrors "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/errors"
	sdkhttp "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
	"net/url"
	"strings"
	"testing"
)

type historyHTTPClient struct {
	mockHttpClient
	endpoint string
	headers  map[string]string
	status   int
	body     string
}

func (c *historyHTTPClient) Get(endpoint string, headers map[string]string) (*sdkhttp.Response, error) {
	c.endpoint = endpoint
	c.headers = headers
	return sdkhttp.NewResponse(c.status, []byte(c.body), nil), nil
}
func TestBillingPaymentHistoryPaginationAndEnvelope(t *testing.T) {
	c := &historyHTTPClient{status: 200, body: `{"success":true,"data":{"subscription_id":"sub_test","snapshot":123,"items":[],"has_more":true,"next_cursor":"in_next"}}`}
	cfg := config.NewConfiguration("public", "secret")
	manager := NewSubscriptionManager(cfg, c, auth.NewManager(cfg, c))
	page, err := manager.BillingPaymentHistory("sub_test", "in_cursor", 123)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := url.Parse(c.endpoint)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(parsed.Path, "/merchant/billing/resources/sub_test/payment-history") || parsed.Query().Get("snapshot") != "123" || parsed.Query().Get("starting_after") != "in_cursor" || c.headers["Authorization"] != "Bearer test-token" || page["next_cursor"] != "in_next" {
		t.Fatalf("incorrect request or response: %s %#v", c.endpoint, page)
	}
	_, err = manager.BillingPaymentHistory("sub_test", "", 0)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(c.endpoint, "?") {
		t.Fatal("initial page must let server establish snapshot")
	}
}
func TestBillingPaymentHistoryPreservesVerificationFailures(t *testing.T) {
	c := &historyHTTPClient{status: 502, body: `{"success":false,"error":{"code":"PAYMENT_HISTORY_UNAVAILABLE","message":"Could not verify complete payment history"}}`}
	cfg := config.NewConfiguration("public", "secret")
	manager := NewSubscriptionManager(cfg, c, auth.NewManager(cfg, c))
	_, err := manager.BillingPaymentHistory("sub_test", "", 0)
	var apiError *sdkerrors.ApiException
	if !errors.As(err, &apiError) || apiError.HTTPStatus != 502 || apiError.ErrorCode != "PAYMENT_HISTORY_UNAVAILABLE" {
		t.Fatalf("verification failure lost: %v", err)
	}
}
