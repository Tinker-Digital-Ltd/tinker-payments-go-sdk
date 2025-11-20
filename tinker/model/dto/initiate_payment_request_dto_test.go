package dto

import (
	"testing"

	"github.com/tinker/tinker-payments-go-sdk/tinker/types"
)

func TestInitiatePaymentRequestDto_ToMap(t *testing.T) {
	customerPhone := "+254712345678"
	customerEmail := "test@example.com"
	transactionDesc := "Test transaction"
	metadata := map[string]interface{}{
		"order_id": "123",
	}

	dto := &InitiatePaymentRequestDto{
		Amount:           100.00,
		Currency:         "KES",
		Gateway:          types.MPESA,
		MerchantReference: "ORDER-123",
		ReturnURL:        "https://example.com/return",
		CustomerPhone:    &customerPhone,
		CustomerEmail:    &customerEmail,
		TransactionDesc:  &transactionDesc,
		Metadata:         metadata,
	}

	result := dto.ToMap()

	if result["amount"] != 100.00 {
		t.Errorf("amount = %v, want 100.00", result["amount"])
	}
	if result["currency"] != "KES" {
		t.Errorf("currency = %v, want KES", result["currency"])
	}
	if result["gateway"] != "mpesa" {
		t.Errorf("gateway = %v, want mpesa", result["gateway"])
	}
	if result["merchantReference"] != "ORDER-123" {
		t.Errorf("merchantReference = %v, want ORDER-123", result["merchantReference"])
	}
	if result["returnUrl"] != "https://example.com/return" {
		t.Errorf("returnUrl = %v, want https://example.com/return", result["returnUrl"])
	}
	if result["customerPhone"] != customerPhone {
		t.Errorf("customerPhone = %v, want %v", result["customerPhone"], customerPhone)
	}
	if result["customerEmail"] != customerEmail {
		t.Errorf("customerEmail = %v, want %v", result["customerEmail"], customerEmail)
	}
	if result["transactionDesc"] != transactionDesc {
		t.Errorf("transactionDesc = %v, want %v", result["transactionDesc"], transactionDesc)
	}
	if result["metadata"] == nil {
		t.Error("metadata should not be nil")
	}
}

func TestInitiatePaymentRequestDto_ToMap_OptionalFields(t *testing.T) {
	dto := &InitiatePaymentRequestDto{
		Amount:           100.00,
		Currency:         "KES",
		Gateway:          types.MPESA,
		MerchantReference: "ORDER-123",
		ReturnURL:        "https://example.com/return",
	}

	result := dto.ToMap()

	if result["customerPhone"] != nil {
		t.Error("customerPhone should not be in map when nil")
	}
	if result["customerEmail"] != nil {
		t.Error("customerEmail should not be in map when nil")
	}
	if result["transactionDesc"] != nil {
		t.Error("transactionDesc should not be in map when nil")
	}
	if result["metadata"] != nil {
		t.Error("metadata should not be in map when nil")
	}
}

