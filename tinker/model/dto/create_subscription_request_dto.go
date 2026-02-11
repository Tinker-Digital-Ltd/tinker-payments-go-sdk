package dto

type CreateSubscriptionRequestDto struct {
	PlanID          string
	Gateway         string
	Customer        *SubscriptionCustomerDto
	PaymentMethodID *string
	BillingPeriod   *string
}

func (dto *CreateSubscriptionRequestDto) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"plan_id":  dto.PlanID,
		"gateway":  dto.Gateway,
		"customer": dto.Customer.ToMap(),
	}
	if dto.PaymentMethodID != nil {
		m["payment_method_id"] = *dto.PaymentMethodID
	}
	if dto.BillingPeriod != nil {
		m["billing_period"] = *dto.BillingPeriod
	}
	return m
}
