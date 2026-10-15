package webhook

import (
	"os"
	"strings"
	"testing"
)

func TestSignedCheckoutCancellation(t *testing.T) {
	raw, err := os.ReadFile("testdata/checkout-cancellation-webhook.json")
	if err != nil {
		t.Fatal(err)
	}
	if canonical, err := os.ReadFile("../../../tinker-payments-openapi/fixtures/checkout-cancellation-webhook.json"); err == nil && string(canonical) != string(raw) {
		t.Fatal("Signed fixture snapshot drifted")
	}
	handler := NewHandler()
	event, err := handler.Handle(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !event.IsCheckoutCancellation() || event.CheckoutCancellationData() == nil || event.CheckoutCancellationData().Resource.Status != "outcome_unknown" {
		t.Fatal("lost cancellation state")
	}
	if event.PaymentData() != nil || event.ToTransaction() != nil {
		t.Fatal("cancellation became payment")
	}
	if !handler.VerifySignature(event, "webhook-secret") || handler.VerifySignature(event, "wrong-secret") {
		t.Fatal("incorrect signature handling")
	}
	changed, err := handler.Handle([]byte(strings.Replace(string(raw), "outcome_unknown", "cancelled", 1)))
	if err != nil {
		t.Fatal(err)
	}
	if handler.VerifySignature(changed, "webhook-secret") {
		t.Fatal("accepted tampered outcome")
	}
}

func TestSignedRefundAndWithdrawalLifecycle(t *testing.T) {
	for _, name := range []string{"refund", "withdrawal"} {
		t.Run(name, func(t *testing.T) {
			raw, err := os.ReadFile("testdata/" + name + "-webhook.json")
			if err != nil {
				t.Fatal(err)
			}
			if canonical, err := os.ReadFile("../../../tinker-payments-openapi/fixtures/" + name + "-webhook.json"); err == nil && string(canonical) != string(raw) {
				t.Fatal("Signed fixture snapshot drifted")
			}
			handler := NewHandler()
			event, err := handler.Handle(raw)
			if err != nil {
				t.Fatal(err)
			}
			future, err := handler.Handle([]byte(strings.Replace(string(raw), ".updated.v1", ".updated.v2", 1)))
			if err != nil || future.ToTransaction() != nil {
				t.Fatal("future lifecycle event became payment")
			}
			if event.ToTransaction() != nil || event.PaymentData() != nil {
				t.Fatal("lifecycle update became payment")
			}
			if name == "refund" && (event.RefundData() == nil || event.RefundData().Resource.Status != "pending_review") {
				t.Fatal("lost refund review state")
			}
			if name == "withdrawal" && (event.WithdrawalData() == nil || event.WithdrawalData().Resource.Status != "pending_approval") {
				t.Fatal("lost withdrawal approval state")
			}
			if !handler.VerifySignature(event, "webhook-secret") || handler.VerifySignature(event, "wrong-secret") {
				t.Fatal("incorrect signature handling")
			}
			changed, err := handler.Handle([]byte(strings.Replace(string(raw), "pending_", "forged_", 1)))
			if err != nil {
				t.Fatal(err)
			}
			if handler.VerifySignature(changed, "webhook-secret") {
				t.Fatal("accepted tampered status")
			}
		})
	}
}
