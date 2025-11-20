package webhook

import (
	"encoding/json"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/errors"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(payload []byte) (*Event, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, errors.NewInvalidPayloadException("Invalid JSON payload: "+err.Error(), 0, err)
	}

	return NewEvent(data)
}

func (h *Handler) HandleAsTransaction(payload []byte) (*model.Transaction, error) {
	event, err := h.Handle(payload)
	if err != nil {
		return nil, err
	}

	return event.ToTransaction(), nil
}
