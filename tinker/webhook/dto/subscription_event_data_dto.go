package dto

type SubscriptionEventDataDto struct {
	ID            string
	Status        string
	PlanID        string
	CustomerID    string
	CreatedAt     string
	CancelledAt   *string
	PausedAt      *string
	ReactivatedAt *string
}

func NewSubscriptionEventDataDto(data map[string]interface{}) *SubscriptionEventDataDto {
	dto := &SubscriptionEventDataDto{}

	if id, ok := data["id"].(string); ok {
		dto.ID = id
	}

	if status, ok := data["status"].(string); ok {
		dto.Status = status
	}

	if planID, ok := data["plan_id"].(string); ok {
		dto.PlanID = planID
	}

	if customerID, ok := data["customer_id"].(string); ok {
		dto.CustomerID = customerID
	}

	if createdAt, ok := data["created_at"].(string); ok {
		dto.CreatedAt = createdAt
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
	result := map[string]interface{}{
		"id":          dto.ID,
		"status":      dto.Status,
		"plan_id":     dto.PlanID,
		"customer_id": dto.CustomerID,
		"created_at":  dto.CreatedAt,
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
