package dto

import (
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/types"
)

type CallbackDataDto struct {
	ID        string
	Status    types.PaymentStatus
	Reference string
	Amount    float64
	Currency  string
	PaidAt    *string
	CreatedAt *string
	Channel   *string
}

func NewCallbackDataDto(data map[string]interface{}) *CallbackDataDto {
	dto := &CallbackDataDto{}

	if id, ok := data["id"].(string); ok {
		dto.ID = id
	}

	statusValue := "pending"
	if s, ok := data["status"].(string); ok {
		statusValue = s
	}
	switch statusValue {
	case "success":
		dto.Status = types.SUCCESS
	case "cancelled":
		dto.Status = types.CANCELLED
	case "failed":
		dto.Status = types.FAILED
	default:
		dto.Status = types.PENDING
	}

	if ref, ok := data["reference"].(string); ok {
		dto.Reference = ref
	}

	if amt, ok := data["amount"].(float64); ok {
		dto.Amount = amt
	}

	if curr, ok := data["currency"].(string); ok {
		dto.Currency = curr
	}

	if paidAt, ok := data["paid_at"].(string); ok {
		dto.PaidAt = &paidAt
	}

	if createdAt, ok := data["created_at"].(string); ok {
		dto.CreatedAt = &createdAt
	}

	if channel, ok := data["channel"].(string); ok {
		dto.Channel = &channel
	}

	return dto
}

func (dto *CallbackDataDto) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"id":        dto.ID,
		"status":    string(dto.Status),
		"reference": dto.Reference,
		"amount":    dto.Amount,
		"currency":  dto.Currency,
	}
	if dto.PaidAt != nil {
		result["paid_at"] = *dto.PaidAt
	}
	if dto.CreatedAt != nil {
		result["created_at"] = *dto.CreatedAt
	}
	if dto.Channel != nil {
		result["channel"] = *dto.Channel
	}
	return result
}
