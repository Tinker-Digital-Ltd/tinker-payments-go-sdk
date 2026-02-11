package dto

type SubscriptionCustomerDto struct {
	ExternalCustomerID string
	Name               string
	Email              *string
	Phone              *string
	Metadata           map[string]interface{}
}

func (dto *SubscriptionCustomerDto) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"external_customer_id": dto.ExternalCustomerID,
		"name":                 dto.Name,
	}
	if dto.Email != nil {
		m["email"] = *dto.Email
	}
	if dto.Phone != nil {
		m["phone"] = *dto.Phone
	}
	if dto.Metadata != nil {
		m["metadata"] = dto.Metadata
	}
	return m
}
