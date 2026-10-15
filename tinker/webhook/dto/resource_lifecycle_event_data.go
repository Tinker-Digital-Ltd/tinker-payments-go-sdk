package dto

import (
	"encoding/json"
	"fmt"
)

// ResourceLifecycleEventData retains a financial request independently of payment outcomes.
type ResourceLifecycleEventData[T any] struct {
	Resource T `json:"resource"`
}

func NewResourceLifecycleEventData[T any](data map[string]interface{}) (*ResourceLifecycleEventData[T], error) {
	if _, ok := data["resource"].(map[string]interface{}); !ok {
		return nil, fmt.Errorf("missing resource in lifecycle webhook")
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result ResourceLifecycleEventData[T]
	if err = json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
