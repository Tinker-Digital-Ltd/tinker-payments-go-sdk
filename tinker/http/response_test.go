package http

import (
	"encoding/json"
	"testing"
)

func TestNewResponse(t *testing.T) {
	body := []byte(`{"test": "data"}`)
	headers := map[string][]string{
		"Content-Type": {"application/json"},
	}

	resp := NewResponse(200, body, headers)
	if resp == nil {
		t.Fatal("NewResponse returned nil")
	}
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %v, want 200", resp.StatusCode)
	}
	if len(resp.Body) == 0 {
		t.Error("Body should not be empty")
	}
	if resp.Headers == nil {
		t.Error("Headers should not be nil")
	}
}

func TestResponse_JSON(t *testing.T) {
	data := map[string]interface{}{
		"test": "data",
		"number": 123,
	}
	body, _ := json.Marshal(data)

	resp := NewResponse(200, body, nil)
	result, err := resp.JSON()
	if err != nil {
		t.Fatalf("JSON() error = %v", err)
	}
	if result == nil {
		t.Fatal("JSON() returned nil")
	}
	if result["test"] != "data" {
		t.Errorf("result[\"test\"] = %v, want data", result["test"])
	}
}

func TestResponse_JSON_InvalidJSON(t *testing.T) {
	body := []byte("invalid json")

	resp := NewResponse(200, body, nil)
	result, err := resp.JSON()
	if err == nil {
		t.Error("JSON() should return error for invalid JSON")
	}
	if result != nil {
		t.Error("JSON() should return nil result on error")
	}
}

