package webhook

import (
	"encoding/json"
	"testing"
)

func TestNewHandler(t *testing.T) {
	handler := NewHandler()
	if handler == nil {
		t.Fatal("NewHandler returned nil")
	}
}

func TestHandler_Handle(t *testing.T) {
	handler := NewHandler()

	payload := map[string]interface{}{
		"id":       "evt_123",
		"type":     "payment.completed",
		"source":   "payment",
		"data":     map[string]interface{}{"reference": "TXN-123"},
		"meta":     map[string]interface{}{"app_id": "app_123"},
		"security": map[string]interface{}{"signature": "sig_123"},
	}

	payloadBytes, _ := json.Marshal(payload)

	event, err := handler.Handle(payloadBytes)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if event == nil {
		t.Fatal("Handle() returned nil event")
	}
	if event.ID != "evt_123" {
		t.Errorf("Event ID = %v, want evt_123", event.ID)
	}
	if event.Type != "payment.completed" {
		t.Errorf("Event Type = %v, want payment.completed", event.Type)
	}
	if event.Source != "payment" {
		t.Errorf("Event Source = %v, want payment", event.Source)
	}
}

func TestHandler_Handle_InvalidJSON(t *testing.T) {
	handler := NewHandler()
	invalidPayload := []byte("invalid json")

	event, err := handler.Handle(invalidPayload)
	if err == nil {
		t.Error("Handle() should return error for invalid JSON")
	}
	if event != nil {
		t.Error("Handle() should return nil event on error")
	}
}

func TestHandler_HandleAsTransaction(t *testing.T) {
	handler := NewHandler()

	payload := map[string]interface{}{
		"id":     "evt_123",
		"type":   "payment.completed",
		"source": "payment",
		"data": map[string]interface{}{
			"id":        "123",
			"reference": "TXN-123",
			"amount":    100.00,
			"currency":  "KES",
			"status":    "success",
		},
		"meta":     map[string]interface{}{"app_id": "app_123"},
		"security": map[string]interface{}{"signature": "sig_123"},
	}

	payloadBytes, _ := json.Marshal(payload)

	transaction, err := handler.HandleAsTransaction(payloadBytes)
	if err != nil {
		t.Fatalf("HandleAsTransaction() error = %v", err)
	}
	if transaction == nil {
		t.Fatal("HandleAsTransaction() returned nil transaction")
	}
}

func TestHandler_HandleAsTransaction_NonPaymentEvent(t *testing.T) {
	handler := NewHandler()

	payload := map[string]interface{}{
		"id":       "evt_123",
		"type":     "subscription.created",
		"source":   "subscription",
		"data":     map[string]interface{}{"id": "sub_123"},
		"meta":     map[string]interface{}{"app_id": "app_123"},
		"security": map[string]interface{}{"signature": "sig_123"},
	}

	payloadBytes, _ := json.Marshal(payload)

	transaction, err := handler.HandleAsTransaction(payloadBytes)
	if err != nil {
		t.Fatalf("HandleAsTransaction() error = %v", err)
	}
	if transaction != nil {
		t.Error("HandleAsTransaction() should return nil for non-payment events")
	}
}
