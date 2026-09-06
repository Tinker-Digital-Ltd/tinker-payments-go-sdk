package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestSignatureMatchesTinkerServerEnvelope(t *testing.T) {
	// This field order is the Tinker server's marshalWithoutSecurity contract.
	payload := `{"id":"event-1","type":"payment.completed","source":"payment","timestamp":"2026-09-06T00:00:00Z","data":{"amount":15,"currency":"EUR"},"meta":{"app_id":"app-1"}}`
	mac := hmac.New(sha256.New, []byte("test-secret"))
	mac.Write([]byte(payload))
	signed := payload[:len(payload)-1] + `,"security":{"signature":"sha256=` + hex.EncodeToString(mac.Sum(nil)) + `"}}`
	h := NewHandler()
	event, err := h.Handle([]byte(signed))
	if err != nil {
		t.Fatal(err)
	}
	if !h.VerifySignature(event, "test-secret") {
		t.Fatal("rejected server signature")
	}
	if h.VerifySignature(event, "other-secret") {
		t.Fatal("accepted wrong secret")
	}
	event.RawData["amount"] = float64(150)
	if h.VerifySignature(event, "test-secret") {
		t.Fatal("accepted changed amount")
	}
	if h.VerifySignature(nil, "test-secret") || h.VerifySignature(event, "") {
		t.Fatal("accepted missing security")
	}
}
