package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/errors"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

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

func (h *Handler) VerifySignature(event *Event, webhookSecret string) bool {
	if webhookSecret == "" || event == nil || event.Security == nil {
		return false
	}
	signature := event.Security.Signature
	if len(signature) < 7 || signature[:7] != "sha256=" {
		return false
	}

	// Match the server's stable envelope order. Map encoding sorts keys and
	// produces a different HMAC even when the decoded event is identical.
	payloadWithoutSecurity := struct {
		ID        string                 `json:"id"`
		Type      string                 `json:"type"`
		Source    string                 `json:"source"`
		Timestamp *string                `json:"timestamp"`
		Data      map[string]interface{} `json:"data"`
		Meta      map[string]interface{} `json:"meta"`
	}{event.ID, event.Type, event.Source, event.Timestamp, event.RawData, event.RawMeta}

	payload, err := json.Marshal(payloadWithoutSecurity)
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write(payload)
	computed := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(computed))
}
