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
	if !event.IsCheckoutCancellation() || event.CheckoutCancellationData() == nil || event.CheckoutCancellationData().Cancellation.Status != "outcome_unknown" {
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
