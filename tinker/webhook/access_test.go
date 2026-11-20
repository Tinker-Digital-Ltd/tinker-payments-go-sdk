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
	changed, err := handler.Handle([]byte(`{"id":"e","type":"subscription.product_changed","source":"subscription","data":{"subscription_id":"s","product_id":"p2","previous_product_id":"p1","previous_expiration_at":"x","cancel_reason":"UNKNOWN","environment":"SANDBOX","presented_offering_id":"o","transaction_id":"t1","original_transaction_id":"t0"}}`))
	if err != nil {
		t.Fatal(err)
	}
	d := changed.SubscriptionData()
	if *d.ProductID != "p2" || *d.PreviousProductID != "p1" || *d.PreviousExpirationAt != "x" || *d.CancelReason != "UNKNOWN" || d.ExpirationReason != nil || *d.Environment != "SANDBOX" || *d.PresentedOfferingID != "o" || *d.TransactionID != "t1" || *d.OriginalTransactionID != "t0" {
		t.Fatalf("subscription.product_changed misparsed: %+v", d)
	}
	if purchase.EntitlementData().Environment != nil {
		t.Fatal("environment must be nil when absent")
	}
}
