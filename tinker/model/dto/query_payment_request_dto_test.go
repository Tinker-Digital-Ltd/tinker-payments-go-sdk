package dto

import (
	"testing"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/types"
)

func TestQueryPaymentRequestDto_ToMap(t *testing.T) {
	dto := &QueryPaymentRequestDto{
		PaymentReference: "TXN-123",
		Gateway:          types.MPESA,
	}

	result := dto.ToMap()

	if result["payment_reference"] != "TXN-123" {
		t.Errorf("payment_reference = %v, want TXN-123", result["payment_reference"])
	}
	if result["gateway"] != "mpesa" {
		t.Errorf("gateway = %v, want mpesa", result["gateway"])
	}
}
