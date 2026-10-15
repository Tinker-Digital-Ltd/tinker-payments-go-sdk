package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"strings"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model"
)

func paymentPath(id string) string { return "merchant/payment/" + url.PathEscape(id) }

func lifecycleRequest[T any](manager *BaseManager, method, path string, payload interface{}, key string) (*T, error) {
	var headers map[string]string
	if method == "POST" {
		key = strings.TrimSpace(key)
		if len(key) < 1 || len(key) > 255 {
			return nil, fmt.Errorf("Idempotency-Key must contain 1 to 255 UTF-8 bytes")
		}
		headers = map[string]string{"Idempotency-Key": key}
	}
	var data map[string]interface{}
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(raw, &data); err != nil {
			return nil, err
		}
	}
	response, err := manager.requestWithHeaders(method, path, data, headers)
	if err != nil {
		return nil, err
	}
	var decoded T
	var value interface{} = response
	if wrapped, ok := response["value"]; ok {
		value = wrapped
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(raw, &decoded); err != nil {
		return nil, err
	}
	return &decoded, nil
}

func positiveAmount(amount float64) error {
	if math.IsNaN(amount) || math.IsInf(amount, 0) || amount <= 0 {
		return fmt.Errorf("amount must be finite and positive")
	}
	return nil
}

// Refund requests review. Only status refunded confirms completion.
func (tm *TransactionManager) Refund(id string, request model.RefundRequest, key string) (*model.RefundResult, error) {
	if request.Amount != nil {
		if err := positiveAmount(*request.Amount); err != nil {
			return nil, err
		}
	}
	return lifecycleRequest[model.RefundResult](tm.BaseManager, "POST", paymentPath(id)+"/refund", request, key)
}
func (tm *TransactionManager) Refunds(id string) (*model.RefundHistory, error) {
	return lifecycleRequest[model.RefundHistory](tm.BaseManager, "GET", paymentPath(id)+"/refunds", nil, "")
}
func (tm *TransactionManager) GetRefund(id, refundID string) (*model.Refund, error) {
	return lifecycleRequest[model.Refund](tm.BaseManager, "GET", paymentPath(id)+"/refunds/"+url.PathEscape(refundID), nil, "")
}
func (tm *TransactionManager) CancelCheckout(id string, request model.CancellationRequest, key string) (*model.CancellationResult, error) {
	return lifecycleRequest[model.CancellationResult](tm.BaseManager, "POST", paymentPath(id)+"/cancellations", request, key)
}
func (tm *TransactionManager) Cancellations(id string) (*model.CancellationHistory, error) {
	return lifecycleRequest[model.CancellationHistory](tm.BaseManager, "GET", paymentPath(id)+"/cancellations", nil, "")
}
func (tm *TransactionManager) GetCancellation(id, requestID string) (*model.Cancellation, error) {
	return lifecycleRequest[model.Cancellation](tm.BaseManager, "GET", paymentPath(id)+"/cancellations/"+url.PathEscape(requestID), nil, "")
}
