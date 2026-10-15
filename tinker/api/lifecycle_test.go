package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/auth"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	apierrors "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/errors"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model"
)

type lifecycleCase struct {
	Operation, Method, Path string
	Headers                 map[string]string
	Body                    map[string]interface{}
	Status                  int
	Response                map[string]interface{}
}
type lifecycleFixture struct {
	Cases  []lifecycleCase
	Errors []lifecycleCase
}

func sharedLifecycleFixture(t *testing.T) lifecycleFixture {
	t.Helper()
	snapshot, err := os.ReadFile("testdata/merchant-lifecycle.json")
	if err != nil {
		t.Fatal(err)
	}
	canonical := os.Getenv("TINKER_LIFECYCLE_FIXTURE")
	if canonical == "" {
		canonical = "../../../tinker-payments-openapi/fixtures/merchant-lifecycle.json"
	}
	if raw, err := os.ReadFile(canonical); err == nil && string(raw) != string(snapshot) {
		t.Fatal("Shared fixture snapshot drifted")
	}
	var fixture lifecycleFixture
	if err = json.Unmarshal(snapshot, &fixture); err != nil {
		t.Fatal(err)
	}
	return fixture
}

type lifecycleHTTP struct {
	t     *testing.T
	row   lifecycleCase
	calls int
}

func (client *lifecycleHTTP) Get(path string, headers map[string]string) (*http.Response, error) {
	return client.send("GET", path, headers, nil)
}
func (client *lifecycleHTTP) Post(path string, headers map[string]string, body []byte) (*http.Response, error) {
	if strings.HasSuffix(path, "/auth/token") {
		return http.NewResponse(200, []byte(`{"token":"token","expires_in":3600}`), nil), nil
	}
	return client.send("POST", path, headers, body)
}
func (client *lifecycleHTTP) send(method, path string, headers map[string]string, body []byte) (*http.Response, error) {
	client.calls++
	if client.row.Path != "" {
		actual, _ := url.Parse(path)
		expected, _ := url.Parse("https://example.test/v1/" + client.row.Path)
		if method != client.row.Method || actual.EscapedPath() != expected.EscapedPath() || !reflect.DeepEqual(actual.Query(), expected.Query()) {
			client.t.Errorf("unexpected %s %s", method, path)
		}
		var payload map[string]interface{}
		if len(body) > 0 {
			if err := json.Unmarshal(body, &payload); err != nil {
				client.t.Fatal(err)
			}
		}
		if !reflect.DeepEqual(payload, client.row.Body) {
			client.t.Errorf("unexpected payload %#v", payload)
		}
		if headers["Authorization"] != "Bearer token" || headers["Idempotency-Key"] != client.row.Headers["Idempotency-Key"] {
			client.t.Errorf("unexpected headers %#v", headers)
		}
	}
	raw, err := json.Marshal(client.row.Response)
	if err != nil {
		return nil, err
	}
	return http.NewResponse(client.row.Status, raw, nil), nil
}
func invokeLifecycle(tx *TransactionManager, pm *PayoutManager, row lifecycleCase) (interface{}, error) {
	switch row.Operation {
	case "refund":
		amount := 100.0
		return tx.Refund("payment/a b", model.RefundRequest{Amount: &amount, Reason: model.RefundReasonRequestedByCustomer}, " stable-key ")
	case "refunds":
		return tx.Refunds("payment/a b")
	case "get_refund":
		return tx.GetRefund("payment/a b", "refund/a b")
	case "payout_methods":
		return pm.Methods()
	case "payout_quote":
		currency := "JPY"
		return pm.Quote(model.PayoutQuoteRequest{PayoutMethodID: "method/a &b", Amount: 100, Currency: &currency})
	case "request_payout":
		return pm.Request(model.PayoutRequest{PayoutMethodID: "method/a &b", Amount: 100, Currency: "JPY", ExpectedFeeAmount: 2}, "stable-key")
	case "payouts":
		page, perPage := 2, 10
		status, currency := model.PayoutStatusPendingApproval, "JPY"
		return pm.List(model.PayoutFilters{Page: &page, PerPage: &perPage, Status: &status, Currency: &currency})
	case "get_payout":
		return pm.Get("payout/a b")
	case "cancel_checkout":
		return tx.CancelCheckout("payment/a b", model.CancellationRequest{Reason: "Customer left checkout"}, "stable-key")
	case "cancellations":
		return tx.Cancellations("payment/a b")
	case "get_cancellation":
		return tx.GetCancellation("payment/a b", "cancel/a b")
	default:
		return nil, fmt.Errorf("unhandled fixture operation %s", row.Operation)
	}
}
func TestLifecycleSharedContract(t *testing.T) {
	for _, row := range sharedLifecycleFixture(t).Cases {
		t.Run(row.Operation, func(t *testing.T) {
			client := &lifecycleHTTP{t: t, row: row}
			cfg := config.NewConfiguration("public", "secret", "https://example.test")
			authentication := auth.NewManager(cfg, client)
			result, err := invokeLifecycle(NewTransactionManager(cfg, client, authentication), NewPayoutManager(cfg, client, authentication), row)
			if err != nil {
				t.Fatal(err)
			}
			raw, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}
			var actual interface{}
			if err = json.Unmarshal(raw, &actual); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(actual, row.Response["data"]) {
				t.Errorf("unexpected result %s", raw)
			}
			if client.calls != 1 {
				t.Errorf("unexpected calls %d", client.calls)
			}
		})
	}
}
func TestLifecycleStructuredErrors(t *testing.T) {
	for _, row := range sharedLifecycleFixture(t).Errors {
		client := &lifecycleHTTP{t: t, row: row}
		cfg := config.NewConfiguration("public", "secret")
		tx := NewTransactionManager(cfg, client, auth.NewManager(cfg, client))
		_, err := tx.CancelCheckout("payment", model.CancellationRequest{Reason: "stop"}, "key")
		failure, ok := err.(*apierrors.ApiException)
		if !ok || failure.HTTPStatus != row.Status || !reflect.DeepEqual(failure.Details, row.Response) {
			t.Fatalf("lost error details: %#v", err)
		}
	}
}
func TestLifecycleRejectsInvalidValues(t *testing.T) {
	client := &lifecycleHTTP{t: t}
	cfg := config.NewConfiguration("public", "secret")
	authentication := auth.NewManager(cfg, client)
	tx := NewTransactionManager(cfg, client, authentication)
	pm := NewPayoutManager(cfg, client, authentication)
	for _, key := range []string{"", " ", strings.Repeat("é", 128)} {
		if _, err := tx.CancelCheckout("payment", model.CancellationRequest{Reason: "stop"}, key); err == nil {
			t.Fatal("accepted invalid key")
		}
	}
	for _, amount := range []float64{math.NaN(), math.Inf(1), 0, -1} {
		if _, err := pm.Quote(model.PayoutQuoteRequest{Amount: amount}); err == nil {
			t.Fatal("accepted invalid amount")
		}
	}
	if client.calls != 0 {
		t.Fatal("invalid values reached transport")
	}
}
