package webhook

import (
	"testing"
)

func TestNewEvent(t *testing.T) {
	payload := map[string]interface{}{
		"id":        "evt_123",
		"type":      "payment.completed",
		"source":    "payment",
		"timestamp": "2024-01-01T00:00:00Z",
		"data":      map[string]interface{}{"reference": "TXN-123"},
		"meta":      map[string]interface{}{"app_id": "app_123"},
		"security":  map[string]interface{}{"signature": "sig_123"},
	}

	event, err := NewEvent(payload)
	if err != nil {
		t.Fatalf("NewEvent() error = %v", err)
	}
	if event == nil {
		t.Fatal("NewEvent() returned nil")
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

func TestNewEvent_MissingData(t *testing.T) {
	payload := map[string]interface{}{
		"id":   "evt_123",
		"type": "payment.completed",
	}

	event, err := NewEvent(payload)
	if err == nil {
		t.Error("NewEvent() should return error when data is missing")
	}
	if event != nil {
		t.Error("NewEvent() should return nil event on error")
	}
}

func TestEvent_IsPaymentEvent(t *testing.T) {
	event := &Event{Source: "payment"}
	if !event.IsPaymentEvent() {
		t.Error("IsPaymentEvent() should return true for payment source")
	}

	event.Source = "subscription"
	if event.IsPaymentEvent() {
		t.Error("IsPaymentEvent() should return false for non-payment source")
	}
}

func TestEvent_IsSubscriptionEvent(t *testing.T) {
	event := &Event{Source: "subscription"}
	if !event.IsSubscriptionEvent() {
		t.Error("IsSubscriptionEvent() should return true for subscription source")
	}

	event.Source = "payment"
	if event.IsSubscriptionEvent() {
		t.Error("IsSubscriptionEvent() should return false for non-subscription source")
	}
}

func TestEvent_IsInvoiceEvent(t *testing.T) {
	event := &Event{Source: "invoice"}
	if !event.IsInvoiceEvent() {
		t.Error("IsInvoiceEvent() should return true for invoice source")
	}

	event.Source = "payment"
	if event.IsInvoiceEvent() {
		t.Error("IsInvoiceEvent() should return false for non-invoice source")
	}
}

func TestEvent_IsSettlementEvent(t *testing.T) {
	event := &Event{Source: "settlement"}
	if !event.IsSettlementEvent() {
		t.Error("IsSettlementEvent() should return true for settlement source")
	}

	event.Source = "payment"
	if event.IsSettlementEvent() {
		t.Error("IsSettlementEvent() should return false for non-settlement source")
	}
}
