package dto

type CreateSubscriptionPlanRequestDto struct {
	Name         string
	Amount       float64
	Currency     string
	Intervals    []string
	Description  string
	IsActive     bool
	FeaturesJSON map[string]interface{}
}

func (dto *CreateSubscriptionPlanRequestDto) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"name":        dto.Name,
		"description": dto.Description,
		"amount":      dto.Amount,
		"currency":    dto.Currency,
		"intervals":   dto.Intervals,
		"is_active":   dto.IsActive,
	}
	if dto.FeaturesJSON != nil {
		m["features_json"] = dto.FeaturesJSON
	}
	return m
}
