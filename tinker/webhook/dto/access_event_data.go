package dto

func optionalString(data map[string]interface{}, key string) *string {
	if value, ok := data[key].(string); ok {
		return &value
	}
	return nil
}

func stringList(value interface{}) []string {
	items, _ := value.([]interface{})
	out := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

type EntitlementEventData struct {
	ID                    *string
	TransactionID         *string
	ProductID             *string
	ExternalCustomerID    string
	EntitlementIDs        []string
	Store                 string
	PeriodType            string
	Amount                *float64
	Currency              *string
	PurchasedAt           *string
	ExpirationAt          *string
	PresentedOfferingID   *string
	RevokedAt             *string
	Environment           *string
	OriginalTransactionID *string
}

func NewEntitlementEventData(data map[string]interface{}) *EntitlementEventData {
	out := &EntitlementEventData{
		ID:                    optionalString(data, "id"),
		TransactionID:         optionalString(data, "transaction_id"),
		ProductID:             optionalString(data, "product_id"),
		EntitlementIDs:        stringList(data["entitlement_ids"]),
		Currency:              optionalString(data, "currency"),
		PurchasedAt:           optionalString(data, "purchased_at"),
		ExpirationAt:          optionalString(data, "expiration_at"),
		PresentedOfferingID:   optionalString(data, "presented_offering_id"),
		RevokedAt:             optionalString(data, "revoked_at"),
		Environment:           optionalString(data, "environment"),
		OriginalTransactionID: optionalString(data, "original_transaction_id"),
	}
	out.ExternalCustomerID, _ = data["external_customer_id"].(string)
	out.Store, _ = data["store"].(string)
	out.PeriodType, _ = data["period_type"].(string)
	if amount, ok := data["amount"].(float64); ok {
		out.Amount = &amount
	}
	return out
}

type CustomerEventData struct {
	TransferredFrom []string
	TransferredTo   []string
}

func NewCustomerEventData(data map[string]interface{}) *CustomerEventData {
	return &CustomerEventData{TransferredFrom: stringList(data["transferred_from"]), TransferredTo: stringList(data["transferred_to"])}
}
