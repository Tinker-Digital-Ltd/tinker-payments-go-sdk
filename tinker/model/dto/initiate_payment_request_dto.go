package dto

import (
	"github.com/tinker/tinker-payments-go-sdk/tinker/types"
)

type InitiatePaymentRequestDto struct {
	Amount            float64
	Currency          string
	Gateway           types.Gateway
	MerchantReference string
	ReturnURL         string
	CustomerPhone     *string
	CustomerEmail     *string
	TransactionDesc   *string
	Metadata          map[string]interface{}
}

func (dto *InitiatePaymentRequestDto) ToMap() map[string]interface{} {
	payload := map[string]interface{}{
		"amount":            dto.Amount,
		"currency":          dto.Currency,
		"gateway":           string(dto.Gateway),
		"merchantReference": dto.MerchantReference,
		"returnUrl":         dto.ReturnURL,
	}

	if dto.CustomerPhone != nil {
		payload["customerPhone"] = *dto.CustomerPhone
	}
	if dto.CustomerEmail != nil {
		payload["customerEmail"] = *dto.CustomerEmail
	}
	if dto.TransactionDesc != nil {
		payload["transactionDesc"] = *dto.TransactionDesc
	}
	if dto.Metadata != nil {
		payload["metadata"] = dto.Metadata
	}

	return payload
}
