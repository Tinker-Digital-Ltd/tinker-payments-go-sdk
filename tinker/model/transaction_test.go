package model

import (
	"testing"

	"github.com/tinker/tinker-payments-go-sdk/tinker/types"
)

func TestNewTransaction_WithInitiationData(t *testing.T) {
	data := map[string]interface{}{
		"payment_reference": "TXN-123",
		"authorization_url": "https://example.com/auth",
		"status":            "pending",
	}

	transaction := NewTransaction(data)
	if transaction == nil {
		t.Fatal("NewTransaction returned nil")
	}
	if transaction.InitiationData == nil {
		t.Error("InitiationData should not be nil")
	}
	if transaction.Status != types.PENDING {
		t.Errorf("Status = %v, want PENDING", transaction.Status)
	}
}

func TestNewTransaction_WithQueryData(t *testing.T) {
	data := map[string]interface{}{
		"id":        "123",
		"reference": "TXN-123",
		"amount":    100.00,
		"currency":  "KES",
		"status":    "success",
	}

	transaction := NewTransaction(data)
	if transaction == nil {
		t.Fatal("NewTransaction returned nil")
	}
	if transaction.QueryData == nil {
		t.Error("QueryData should not be nil")
	}
	if transaction.CallbackData == nil {
		t.Error("CallbackData should not be nil")
	}
	if transaction.Status != types.SUCCESS {
		t.Errorf("Status = %v, want SUCCESS", transaction.Status)
	}
}

func TestNewTransaction_WithStatusOnly(t *testing.T) {
	data := map[string]interface{}{
		"status": "failed",
	}

	transaction := NewTransaction(data)
	if transaction == nil {
		t.Fatal("NewTransaction returned nil")
	}
	if transaction.Status != types.FAILED {
		t.Errorf("Status = %v, want FAILED", transaction.Status)
	}
}

func TestTransaction_IsSuccessful(t *testing.T) {
	transaction := &Transaction{Status: types.SUCCESS}
	if !transaction.IsSuccessful() {
		t.Error("IsSuccessful() should return true for SUCCESS status")
	}

	transaction.Status = types.PENDING
	if transaction.IsSuccessful() {
		t.Error("IsSuccessful() should return false for non-SUCCESS status")
	}
}

func TestTransaction_IsPending(t *testing.T) {
	transaction := &Transaction{Status: types.PENDING}
	if !transaction.IsPending() {
		t.Error("IsPending() should return true for PENDING status")
	}

	transaction.Status = types.SUCCESS
	if transaction.IsPending() {
		t.Error("IsPending() should return false for non-PENDING status")
	}
}

func TestTransaction_IsCancelled(t *testing.T) {
	transaction := &Transaction{Status: types.CANCELLED}
	if !transaction.IsCancelled() {
		t.Error("IsCancelled() should return true for CANCELLED status")
	}

	transaction.Status = types.SUCCESS
	if transaction.IsCancelled() {
		t.Error("IsCancelled() should return false for non-CANCELLED status")
	}
}

func TestTransaction_IsFailed(t *testing.T) {
	transaction := &Transaction{Status: types.FAILED}
	if !transaction.IsFailed() {
		t.Error("IsFailed() should return true for FAILED status")
	}

	transaction.Status = types.SUCCESS
	if transaction.IsFailed() {
		t.Error("IsFailed() should return false for non-FAILED status")
	}
}
