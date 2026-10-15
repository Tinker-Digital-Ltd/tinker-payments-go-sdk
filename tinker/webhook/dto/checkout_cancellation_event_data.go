package dto

import (
	"encoding/json"
	"fmt"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model"
)

// CheckoutCancellationEventData carries a request outcome, never a payment.
type CheckoutCancellationEventData struct {
	Cancellation model.Cancellation `json:"cancellation"`
}

func NewCheckoutCancellationEventData(data map[string]interface{}) (*CheckoutCancellationEventData, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result CheckoutCancellationEventData
	if err = json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	if result.Cancellation.ID == "" {
		return nil, fmt.Errorf("missing cancellation in webhook")
	}
	return &result, nil
}
