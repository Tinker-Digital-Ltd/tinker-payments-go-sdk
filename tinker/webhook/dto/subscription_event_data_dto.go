package dto

type SubscriptionEventDataDto struct {
	ID                 string
	Status             string
	PlanID             string
	AccountID          string
	CurrentPeriodStart *string
	CurrentPeriodEnd   *string
	CreatedAt          string
	CancelledAt        *string
	PausedAt           *string
	ReactivatedAt      *string

	ProductID               *string
	EntitlementIDs          []string
	Store                   *string
	PeriodType              *string
	ExpirationAt            *string
	GracePeriodExpirationAt *string
	ExternalCustomerID      *string
	CancelReason            *string
	ExpirationReason        *string
	Environment             *string
	PresentedOfferingID     *string
	TransactionID           *string
	OriginalTransactionID   *string
	PreviousProductID       *string
	PreviousExpirationAt    *string
}

func NewSubscriptionEventDataDto(data map[string]interface{}) *SubscriptionEventDataDto {
	dto := &SubscriptionEventDataDto{}
	if id, ok := data["subscription_id"].(string); ok {
		dto.ID = id
	} else if id, ok := data["id"].(string); ok {
		dto.ID = id
	}
	if status, ok := data["status"].(string); ok {
		dto.Status = status
	}
	if planID, ok := data["plan_id"].(string); ok {
		dto.PlanID = planID
	}
	if accountID, ok := data["account_id"].(string); ok {
		dto.AccountID = accountID
	} else if customerID, ok := data["customer_id"].(string); ok {
		dto.AccountID = customerID
	}
	if cps, ok := data["current_period_start"].(string); ok {
		dto.CurrentPeriodStart = &cps
	}
	if cpe, ok := data["current_period_end"].(string); ok {
		dto.CurrentPeriodEnd = &cpe
	}
	if createdAt, ok := data["created_at"].(string); ok {
		dto.CreatedAt = createdAt
	} else if dto.CurrentPeriodStart != nil {
		dto.CreatedAt = *dto.CurrentPeriodStart
	}
	if cancelledAt, ok := data["cancelled_at"].(string); ok {
		dto.CancelledAt = &cancelledAt
	}
	if pausedAt, ok := data["paused_at"].(string); ok {
		dto.PausedAt = &pausedAt
	}
	if reactivatedAt, ok := data["reactivated_at"].(string); ok {
		dto.ReactivatedAt = &reactivatedAt
	}
	dto.ProductID = optionalString(data, "product_id")
	dto.EntitlementIDs = stringList(data["entitlement_ids"])
	dto.Store = optionalString(data, "store")
	dto.PeriodType = optionalString(data, "period_type")
	dto.ExpirationAt = optionalString(data, "expiration_at")
	dto.GracePeriodExpirationAt = optionalString(data, "grace_period_expiration_at")
	dto.ExternalCustomerID = optionalString(data, "external_customer_id")
	dto.CancelReason = optionalString(data, "cancel_reason")
	dto.ExpirationReason = optionalString(data, "expiration_reason")
	dto.Environment = optionalString(data, "environment")
	dto.PresentedOfferingID = optionalString(data, "presented_offering_id")
	dto.TransactionID = optionalString(data, "transaction_id")
	dto.OriginalTransactionID = optionalString(data, "original_transaction_id")
	dto.PreviousProductID = optionalString(data, "previous_product_id")
	dto.PreviousExpirationAt = optionalString(data, "previous_expiration_at")
	return dto
}

func (dto *SubscriptionEventDataDto) ToMap() map[string]interface{} {
	result := map[string]interface{}{"id": dto.ID, "subscription_id": dto.ID, "status": dto.Status, "plan_id": dto.PlanID, "account_id": dto.AccountID, "created_at": dto.CreatedAt}
	if dto.CurrentPeriodStart != nil {
		result["current_period_start"] = *dto.CurrentPeriodStart
	}
	if dto.CurrentPeriodEnd != nil {
		result["current_period_end"] = *dto.CurrentPeriodEnd
	}
	if dto.CancelledAt != nil {
		result["cancelled_at"] = *dto.CancelledAt
	}
	if dto.PausedAt != nil {
		result["paused_at"] = *dto.PausedAt
	}
	if dto.ReactivatedAt != nil {
		result["reactivated_at"] = *dto.ReactivatedAt
	}
	for key, value := range map[string]*string{"product_id": dto.ProductID, "store": dto.Store, "period_type": dto.PeriodType, "expiration_at": dto.ExpirationAt, "grace_period_expiration_at": dto.GracePeriodExpirationAt, "external_customer_id": dto.ExternalCustomerID, "cancel_reason": dto.CancelReason, "expiration_reason": dto.ExpirationReason, "environment": dto.Environment, "presented_offering_id": dto.PresentedOfferingID, "transaction_id": dto.TransactionID, "original_transaction_id": dto.OriginalTransactionID, "previous_product_id": dto.PreviousProductID, "previous_expiration_at": dto.PreviousExpirationAt} {
		if value != nil {
			result[key] = *value
		}
	}
	if dto.EntitlementIDs != nil {
		result["entitlement_ids"] = dto.EntitlementIDs
	}
	return result
}
