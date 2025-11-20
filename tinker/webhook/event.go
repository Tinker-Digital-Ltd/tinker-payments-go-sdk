package webhook

import (
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/errors"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model"
	webhookDto "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/webhook/dto"
)

type Event struct {
	ID        string
	Type      string
	Source    string
	Timestamp *string
	Data      interface{}
	Meta      *Meta
	Security  *Security
}

func NewEvent(payload map[string]interface{}) (*Event, error) {
	event := &Event{}

	if id, ok := payload["id"].(string); ok {
		event.ID = id
	}

	if eventType, ok := payload["type"].(string); ok {
		event.Type = eventType
	}

	if source, ok := payload["source"].(string); ok {
		event.Source = source
	}

	if timestamp, ok := payload["timestamp"].(string); ok {
		event.Timestamp = &timestamp
	}

	data, ok := payload["data"].(map[string]interface{})
	if !ok {
		return nil, errors.NewInvalidPayloadException("Webhook payload data must be a map", 0, nil)
	}

	eventData, err := createEventData(data, event.Source)
	if err != nil {
		return nil, err
	}
	event.Data = eventData

	metaData, ok := payload["meta"].(map[string]interface{})
	if !ok {
		metaData = make(map[string]interface{})
	}
	event.Meta = NewMeta(metaData)

	securityData, ok := payload["security"].(map[string]interface{})
	if !ok {
		securityData = make(map[string]interface{})
	}
	event.Security = NewSecurity(securityData)

	return event, nil
}

func (e *Event) IsPaymentEvent() bool {
	return e.Source == "payment"
}

func (e *Event) IsSubscriptionEvent() bool {
	return e.Source == "subscription"
}

func (e *Event) IsInvoiceEvent() bool {
	return e.Source == "invoice"
}

func (e *Event) IsSettlementEvent() bool {
	return e.Source == "settlement"
}

func (e *Event) PaymentData() *webhookDto.PaymentEventDataDto {
	if data, ok := e.Data.(*webhookDto.PaymentEventDataDto); ok {
		return data
	}
	return nil
}

func (e *Event) SubscriptionData() *webhookDto.SubscriptionEventDataDto {
	if data, ok := e.Data.(*webhookDto.SubscriptionEventDataDto); ok {
		return data
	}
	return nil
}

func (e *Event) InvoiceData() *webhookDto.InvoiceEventDataDto {
	if data, ok := e.Data.(*webhookDto.InvoiceEventDataDto); ok {
		return data
	}
	return nil
}

func (e *Event) SettlementData() *webhookDto.SettlementEventDataDto {
	if data, ok := e.Data.(*webhookDto.SettlementEventDataDto); ok {
		return data
	}
	return nil
}

func (e *Event) ToTransaction() *model.Transaction {
	if !e.IsPaymentEvent() {
		return nil
	}

	paymentData := e.PaymentData()
	if paymentData == nil {
		return nil
	}

	data := paymentData.ToMap()
	return model.NewTransaction(data)
}

func createEventData(data map[string]interface{}, source string) (interface{}, error) {
	switch source {
	case "payment":
		return webhookDto.NewPaymentEventDataDto(data), nil
	case "subscription":
		return webhookDto.NewSubscriptionEventDataDto(data), nil
	case "invoice":
		return webhookDto.NewInvoiceEventDataDto(data), nil
	case "settlement":
		return webhookDto.NewSettlementEventDataDto(data), nil
	default:
		return nil, errors.NewInvalidPayloadException("Unknown webhook source: "+source, 0, nil)
	}
}
