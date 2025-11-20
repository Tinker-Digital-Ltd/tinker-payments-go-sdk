package dto

import (
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/types"
)

type QueryPaymentRequestDto struct {
	PaymentReference string
	Gateway          types.Gateway
}

func (dto *QueryPaymentRequestDto) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"payment_reference": dto.PaymentReference,
		"gateway":           string(dto.Gateway),
	}
}
