package http

import (
	"testing"
	"time"
)

func TestNewHttpClient(t *testing.T) {
	client := NewHttpClient()
	if client == nil {
		t.Fatal("NewHttpClient returned nil")
	}
	if client.timeout != 30*time.Second {
		t.Errorf("timeout = %v, want 30s", client.timeout)
	}
	if client.client == nil {
		t.Error("http.Client should not be nil")
	}
}

func TestHttpClient_Post(t *testing.T) {
	client := NewHttpClient()
	url := "https://httpbin.org/post"
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	body := []byte(`{"test": "data"}`)

	resp, err := client.Post(url, headers, body)
	if err != nil {
		t.Skipf("Post() error = %v (skipping integration test)", err)
	}
	if resp == nil {
		t.Fatal("Post() returned nil response")
	}
	if resp.StatusCode == 0 {
		t.Error("StatusCode should not be zero")
	}
}

func TestHttpClient_Post_NilBody(t *testing.T) {
	client := NewHttpClient()
	url := "https://httpbin.org/post"
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := client.Post(url, headers, nil)
	if err != nil {
		t.Skipf("Post() error = %v (skipping integration test)", err)
	}
	if resp == nil {
		t.Fatal("Post() returned nil response")
	}
}

