package dto

type InvoiceEventDataDto struct {
	ID             string
	Status         string
	InvoiceNumber  string
	Amount         float64
	Currency       string
	SubscriptionID string
	CreatedAt      string
	PaidAt         *string
}

func NewInvoiceEventDataDto(data map[string]interface{}) *InvoiceEventDataDto {
	dto := &InvoiceEventDataDto{}

	if id, ok := data["id"].(string); ok {
		dto.ID = id
	}

	if status, ok := data["status"].(string); ok {
		dto.Status = status
	}

	if invoiceNumber, ok := data["invoice_number"].(string); ok {
		dto.InvoiceNumber = invoiceNumber
	}

	if amt, ok := data["amount"].(float64); ok {
		dto.Amount = amt
	}

	if curr, ok := data["currency"].(string); ok {
		dto.Currency = curr
	}

	if subscriptionID, ok := data["subscription_id"].(string); ok {
		dto.SubscriptionID = subscriptionID
	}

	if createdAt, ok := data["created_at"].(string); ok {
		dto.CreatedAt = createdAt
	}

	if paidAt, ok := data["paid_at"].(string); ok {
		dto.PaidAt = &paidAt
	}

	return dto
}

func (dto *InvoiceEventDataDto) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"id":              dto.ID,
		"status":          dto.Status,
		"invoice_number":  dto.InvoiceNumber,
		"amount":          dto.Amount,
		"currency":        dto.Currency,
		"subscription_id": dto.SubscriptionID,
		"created_at":      dto.CreatedAt,
	}
	if dto.PaidAt != nil {
		result["paid_at"] = *dto.PaidAt
	}
	return result
}
