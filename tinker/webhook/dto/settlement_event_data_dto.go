package dto

type SettlementEventDataDto struct {
	ID             string
	Status         string
	Amount         float64
	Currency       string
	SettlementDate string
	CreatedAt      string
	ProcessedAt    *string
}

func NewSettlementEventDataDto(data map[string]interface{}) *SettlementEventDataDto {
	dto := &SettlementEventDataDto{}

	if id, ok := data["id"].(string); ok {
		dto.ID = id
	}

	if status, ok := data["status"].(string); ok {
		dto.Status = status
	}

	if amt, ok := data["amount"].(float64); ok {
		dto.Amount = amt
	}

	if curr, ok := data["currency"].(string); ok {
		dto.Currency = curr
	}

	if settlementDate, ok := data["settlement_date"].(string); ok {
		dto.SettlementDate = settlementDate
	}

	if createdAt, ok := data["created_at"].(string); ok {
		dto.CreatedAt = createdAt
	}

	if processedAt, ok := data["processed_at"].(string); ok {
		dto.ProcessedAt = &processedAt
	}

	return dto
}

func (dto *SettlementEventDataDto) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"id":              dto.ID,
		"status":          dto.Status,
		"amount":          dto.Amount,
		"currency":        dto.Currency,
		"settlement_date": dto.SettlementDate,
		"created_at":      dto.CreatedAt,
	}
	if dto.ProcessedAt != nil {
		result["processed_at"] = *dto.ProcessedAt
	}
	return result
}
