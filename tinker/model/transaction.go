package model

import (
	"github.com/tinker/tinker-payments-go-sdk/tinker/types"
	"github.com/tinker/tinker-payments-go-sdk/tinker/model/dto"
)

type Transaction struct {
	Status         types.PaymentStatus
	InitiationData *dto.InitiationDataDto
	QueryData      *dto.QueryDataDto
	CallbackData   *dto.CallbackDataDto
}

func NewTransaction(data map[string]interface{}) *Transaction {
	txn := &Transaction{}

	if _, hasPaymentRef := data["payment_reference"]; hasPaymentRef {
		if _, hasID := data["id"]; !hasID {
			txn.InitiationData = dto.NewInitiationDataDto(data)
			txn.Status = txn.InitiationData.Status
			return txn
		}
	}

	if _, hasID := data["id"]; hasID {
		if _, hasRef := data["reference"]; hasRef {
			txn.QueryData = dto.NewQueryDataDto(data)
			txn.CallbackData = dto.NewCallbackDataDto(data)
			txn.Status = txn.QueryData.Status
			return txn
		}
	}

	statusValue := "pending"
	if s, ok := data["status"].(string); ok {
		statusValue = s
	}
	switch statusValue {
	case "success":
		txn.Status = types.SUCCESS
	case "cancelled":
		txn.Status = types.CANCELLED
	case "failed":
		txn.Status = types.FAILED
	default:
		txn.Status = types.PENDING
	}

	return txn
}

func (t *Transaction) IsSuccessful() bool {
	return t.Status == types.SUCCESS
}

func (t *Transaction) IsPending() bool {
	return t.Status == types.PENDING
}

func (t *Transaction) IsCancelled() bool {
	return t.Status == types.CANCELLED
}

func (t *Transaction) IsFailed() bool {
	return t.Status == types.FAILED
}
