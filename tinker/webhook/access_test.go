package webhook

import "testing"

func TestAccessEvents(t *testing.T) {
	handler := NewHandler()
	purchase, err := handler.Handle([]byte(`{"id":"e","type":"purchase.completed","source":"payment","data":{"external_customer_id":"u1","entitlement_ids":["pro"],"store":"MPESA","period_type":"NORMAL","product_id":"p1","amount":10}}`))
	if err != nil || purchase.EntitlementData() == nil || *purchase.EntitlementData().ProductID != "p1" || purchase.EntitlementData().EntitlementIDs[0] != "pro" || purchase.PaymentData() != nil {
		t.Fatalf("purchase.completed misparsed: %v", err)
	}
	granted, err := handler.Handle([]byte(`{"id":"e","type":"entitlement.granted","source":"entitlement","data":{"id":"g"}}`))
	if err != nil || *granted.EntitlementData().ID != "g" {
		t.Fatalf("entitlement.granted misparsed: %v", err)
	}
	moved, err := handler.Handle([]byte(`{"id":"e","type":"customer.transferred","source":"customer","data":{"transferred_from":["a"],"transferred_to":["b"]}}`))
	if err != nil || moved.CustomerData().TransferredTo[0] != "b" {
		t.Fatalf("customer.transferred misparsed: %v", err)
	}
	sub, err := handler.Handle([]byte(`{"id":"e","type":"subscription.billing_issue","source":"subscription","data":{"subscription_id":"s","external_customer_id":"u1","product_id":"p","entitlement_ids":["pro"],"grace_period_expiration_at":"g"}}`))
	if err != nil || *sub.SubscriptionData().GracePeriodExpirationAt != "g" || sub.SubscriptionData().EntitlementIDs[0] != "pro" || *sub.SubscriptionData().ExternalCustomerID != "u1" {
		t.Fatalf("subscription.billing_issue misparsed: %v", err)
	}
}
