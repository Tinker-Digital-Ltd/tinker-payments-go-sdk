package api

import "net/url"

// ManagedCheckout uses exact minor-unit prices with an explicit recurring interval.
// Requires the merchant managed-billing capability on the Tinker server.
type ManagedCheckout struct {
	ExternalCustomerID string `json:"external_customer_id"`
	Email              string `json:"email"`
	Name               string `json:"name"`
	Currency           string `json:"currency"`
	UnitAmount         int64  `json:"unit_amount"`
	Quantity           int    `json:"quantity"`
	Interval           string `json:"interval"`
	ReturnURL          string `json:"return_url"`
	IdempotencyKey     string `json:"idempotency_key"`
}

func (sm *SubscriptionManager) Checkout(in ManagedCheckout) (map[string]interface{}, error) {
	return sm.request("POST", "merchant/billing/checkout", map[string]interface{}{"external_customer_id": in.ExternalCustomerID, "email": in.Email, "name": in.Name, "currency": in.Currency, "unit_amount": in.UnitAmount, "quantity": in.Quantity, "interval": in.Interval, "return_url": in.ReturnURL, "idempotency_key": in.IdempotencyKey})
}
func (sm *SubscriptionManager) BillingResource(id string) (map[string]interface{}, error) {
	return sm.request("GET", "merchant/billing/resources/"+url.PathEscape(id), nil)
}
func (sm *SubscriptionManager) ChangeQuantity(id string, quantity int, key string) (map[string]interface{}, error) {
	return sm.request("POST", "merchant/billing/resources/"+url.PathEscape(id)+"/quantity", map[string]interface{}{"quantity": quantity, "idempotency_key": key})
}
func (sm *SubscriptionManager) ScheduleQuantity(id string, quantity int, key string) (map[string]interface{}, error) {
	return sm.request("POST", "merchant/billing/resources/"+url.PathEscape(id)+"/schedule", map[string]interface{}{"quantity": quantity, "idempotency_key": key})
}
func (sm *SubscriptionManager) ClearScheduledQuantity(id, key string) (map[string]interface{}, error) {
	return sm.request("POST", "merchant/billing/resources/"+url.PathEscape(id)+"/unschedule", map[string]interface{}{"idempotency_key": key})
}
func (sm *SubscriptionManager) SetCancelAtRenewal(id string, cancel bool, key string) (map[string]interface{}, error) {
	return sm.request("POST", "merchant/billing/resources/"+url.PathEscape(id)+"/cancel", map[string]interface{}{"cancel": cancel, "idempotency_key": key})
}

// BillingOperation recovers the durable result without replaying a provider charge.
func (sm *SubscriptionManager) BillingOperation(key string) (map[string]interface{}, error) {
	return sm.request("GET", "merchant/billing/operations/"+url.PathEscape(key), nil)
}

// PaymentMethodPortal opens a hosted payment-method-only flow. It does not
// purchase seats, change a subscription price, or pay an outstanding invoice.
func (sm *SubscriptionManager) PaymentMethodPortal(id, returnURL, key string) (map[string]interface{}, error) {
	return sm.request("POST", "merchant/billing/resources/"+url.PathEscape(id)+"/payment-method", map[string]interface{}{"return_url": returnURL, "idempotency_key": key})
}
