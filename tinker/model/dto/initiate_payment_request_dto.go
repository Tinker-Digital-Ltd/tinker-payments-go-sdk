package dto

import (
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/types"
)

type InitiatePaymentRequestDto struct {
	Amount            float64
	Currency          string
	Gateway           types.Gateway
	MerchantReference string
	ReturnURL         string
	CustomerPhone     *string
	CustomerEmail     *string
	CustomerName      *string
	Country           *string
	TransactionDesc   *string
	Metadata          map[string]interface{}
	PaymentMethod     *string
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
	if dto.CustomerName != nil {
		payload["customerName"] = *dto.CustomerName
	}
	if dto.Country != nil {
		payload["country"] = *dto.Country
	}
	if dto.TransactionDesc != nil {
		payload["transactionDesc"] = *dto.TransactionDesc
	}
	if dto.Metadata != nil {
		payload["metadata"] = dto.Metadata
	}
	if dto.PaymentMethod != nil {
		payload["paymentMethod"] = *dto.PaymentMethod
	}

	return payload
}
