package webhook

import (
	"os"
	"strings"
	"testing"
)

func TestGoSignedOnboardingWebhook(t *testing.T) {
	payload, err := os.ReadFile("testdata/setup-webhook.json")
	if err != nil {
		t.Fatal(err)
	}
	handler := NewHandler()
	event, err := handler.Handle(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !event.IsOnboardingEvent() || !event.IsSetupTest() || event.IsPaymentEvent() {
		t.Fatal("incorrect event classification")
	}
	if event.OnboardingData() == nil || !strings.Contains(event.OnboardingData().Message, "webhook verification") {
		t.Fatal("missing onboarding data")
	}
	transaction, err := handler.HandleAsTransaction(payload)
	if err != nil || transaction != nil || event.PaymentData() != nil {
		t.Fatal("setup test became a payment")
	}
	if !handler.VerifySignature(event, "webhook-secret") {
		t.Fatal("sender signature rejected")
	}
	if handler.VerifySignature(event, "wrong-secret") || handler.VerifySignature(event, "") {
		t.Fatal("incorrect secret accepted")
	}
	for _, pair := range [][2]string{{"verification", "forged"}, {"HMAC-SHA256", "none"}, {"sha256=", "invalid="}} {
		changed, err := handler.Handle([]byte(strings.ReplaceAll(string(payload), pair[0], pair[1])))
		if err != nil {
			t.Fatal(err)
		}
		if handler.VerifySignature(changed, "webhook-secret") {
			t.Fatal("tampered event accepted")
		}
	}
	future, err := handler.Handle([]byte(strings.ReplaceAll(string(payload), "setup.webhook_test", "onboarding.future")))
	if err != nil || future.IsSetupTest() {
		t.Fatal("unrecognized event treated as setup test")
	}
	if _, err := handler.Handle([]byte(strings.ReplaceAll(string(payload), `"source":"onboarding"`, `"source":"unknown"`))); err == nil {
		t.Fatal("unknown source accepted")
	}
}
