package dto

import "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model"

// CheckoutCancellationEventData carries a request outcome, never a payment.
type CheckoutCancellationEventData = ResourceLifecycleEventData[model.Cancellation]

func NewCheckoutCancellationEventData(data map[string]interface{}) (*CheckoutCancellationEventData, error) {
	return NewResourceLifecycleEventData[model.Cancellation](data)
}
