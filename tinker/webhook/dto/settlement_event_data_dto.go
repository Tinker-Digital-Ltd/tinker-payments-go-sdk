package dto

type SettlementEventDataDto struct {
	ID             string
	Status         string
	Amount         float64
	NetAmount      *float64
	Currency       string
	SettlementDate string
	CreatedAt      string
	ProcessedAt    *string
}

func NewSettlementEventDataDto(data map[string]interface{}) *SettlementEventDataDto {
	dto := &SettlementEventDataDto{}
	if id, ok := data["settlement_id"].(string); ok {
		dto.ID = id
	} else if id, ok := data["id"].(string); ok {
		dto.ID = id
	}
	if status, ok := data["status"].(string); ok {
		dto.Status = status
	}
	if amt, ok := data["amount"].(float64); ok {
		dto.Amount = amt
	}
	if netAmt, ok := data["net_amount"].(float64); ok {
		dto.NetAmount = &netAmt
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
		"id": dto.ID, "settlement_id": dto.ID, "status": dto.Status, "amount": dto.Amount,
		"currency": dto.Currency, "settlement_date": dto.SettlementDate, "created_at": dto.CreatedAt,
	}
	if dto.NetAmount != nil {
		result["net_amount"] = *dto.NetAmount
	}
	if dto.ProcessedAt != nil {
		result["processed_at"] = *dto.ProcessedAt
	}
	return result
}
