package dto

import "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/types"

type InitiationDataDto struct {
	PaymentReference string
	Status           types.PaymentStatus
	AuthorizationURL *string
	ClientSecret     *string
	PublishableKey   *string
}

func NewInitiationDataDto(data map[string]interface{}) *InitiationDataDto {
	dto := &InitiationDataDto{}
	if pr, ok := data["paymentReference"].(string); ok {
		dto.PaymentReference = pr
	} else if pr, ok := data["payment_reference"].(string); ok {
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

	if au, ok := data["authorizationUrl"].(string); ok {
		dto.AuthorizationURL = &au
	} else if au, ok := data["authorization_url"].(string); ok {
		dto.AuthorizationURL = &au
	}
	if cs, ok := data["clientSecret"].(string); ok {
		dto.ClientSecret = &cs
	} else if cs, ok := data["client_secret"].(string); ok {
		dto.ClientSecret = &cs
	}
	if pk, ok := data["publishableKey"].(string); ok {
		dto.PublishableKey = &pk
	} else if pk, ok := data["publishable_key"].(string); ok {
		dto.PublishableKey = &pk
	}
	return dto
}

func (dto *InitiationDataDto) ToMap() map[string]interface{} {
	result := map[string]interface{}{"paymentReference": dto.PaymentReference, "status": string(dto.Status)}
	if dto.AuthorizationURL != nil {
		result["authorizationUrl"] = *dto.AuthorizationURL
	}
	if dto.ClientSecret != nil {
		result["clientSecret"] = *dto.ClientSecret
	}
	if dto.PublishableKey != nil {
		result["publishableKey"] = *dto.PublishableKey
	}
	return result
}
