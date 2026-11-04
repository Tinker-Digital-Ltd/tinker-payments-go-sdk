package api

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/auth"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	apierrors "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/errors"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model/dto"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/types"
)

func TestGatewaySharedContract(t *testing.T) {
	for _, row := range sharedLifecycleFixture(t).Gateway.Cases {
		row.Headers = nil
		client := &lifecycleHTTP{t: t, row: row}
		cfg := config.NewConfiguration("public", "secret", "https://example.test")
		tx := NewTransactionManager(cfg, client, auth.NewManager(cfg, client))
		data, _ := row.Response["data"].(map[string]interface{})
		switch row.Operation {
		case "listPaymentCapabilities":
			providers, err := tx.Capabilities("KES", "KE")
			if err != nil {
				t.Fatal(err)
			}
			raw, _ := json.Marshal(providers)
			var actual interface{}
			_ = json.Unmarshal(raw, &actual)
			if !reflect.DeepEqual(actual, data["providers"]) {
				t.Errorf("unexpected providers %s", raw)
			}
		case "initiatePayment":
			phone, method := "+254712345678", "apple_pay"
			request := &dto.InitiatePaymentRequestDto{Amount: 100, Currency: "KES", Gateway: types.Gateway(row.Body["gateway"].(string)), MerchantReference: row.Body["merchantReference"].(string), ReturnURL: "https://example.com/return"}
			if row.Status == 422 {
				request.CustomerPhone, request.PaymentMethod = &phone, &method
			}
			transaction, err := tx.Initiate(request)
			if row.Status == 422 {
				var failure *apierrors.FeatureNotSupportedException
				if !errors.As(err, &failure) || failure.Gateway != "mpesa" || failure.Feature != "wallet" || !reflect.DeepEqual(failure.Alternatives, []string{"stripe"}) || failure.HTTPStatus != 422 {
					t.Fatalf("unexpected error %#v", err)
				}
				var api *apierrors.ApiException
				if !errors.As(err, &api) || api.ErrorCode != "FEATURE_NOT_SUPPORTED" {
					t.Fatal("not an ApiException")
				}
			} else if err != nil || transaction.InitiationData.PaymentReference != data["paymentReference"] {
				t.Fatalf("unexpected result %#v %v", transaction, err)
			}
		case "queryPayment":
			transaction, err := tx.Query(&dto.QueryPaymentRequestDto{PaymentReference: "PAY-future1", Gateway: "futurepay"})
			if err != nil || transaction.Status != types.REVERSED {
				t.Fatalf("unexpected result %#v %v", transaction, err)
			}
		}
		if client.calls != 1 {
			t.Errorf("unexpected calls %d", client.calls)
		}
	}
}

func TestUnknownStatusIsKept(t *testing.T) {
	if types.ParsePaymentStatus("future_status") != "future_status" || types.ParsePaymentStatus("completed") != types.SUCCESS || types.ParsePaymentStatus("") != types.PENDING {
		t.Fatal("unexpected status parsing")
	}
}
