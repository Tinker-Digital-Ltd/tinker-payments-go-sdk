package dto

import (
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/types"
)

type InitiationDataDto struct {
	PaymentReference string
	Status           types.PaymentStatus
	AuthorizationURL *string
}

func NewInitiationDataDto(data map[string]interface{}) *InitiationDataDto {
	dto := &InitiationDataDto{}

	if pr, ok := data["payment_reference"].(string); ok {
		dto.PaymentReference = pr
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

	if au, ok := data["authorization_url"].(string); ok {
		dto.AuthorizationURL = &au
	}

	return dto
}

func (dto *InitiationDataDto) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"payment_reference": dto.PaymentReference,
		"status":            string(dto.Status),
	}
	if dto.AuthorizationURL != nil {
		result["authorization_url"] = *dto.AuthorizationURL
	}
	return result
}
