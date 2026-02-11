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
	return result
}
