package auth

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
)

func TestNewManager(t *testing.T) {
	cfg := config.NewConfiguration("public-key", "secret-key")
	mockClient := &mockHttpClient{}

	manager := NewManager(cfg, mockClient)
	if manager == nil {
		t.Fatal("NewManager returned nil")
	}
	if manager.config == nil {
		t.Error("config should not be nil")
	}
	if manager.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
}

func TestManager_Token(t *testing.T) {
	cfg := config.NewConfiguration("public-key", "secret-key")
	mockClient := &mockHttpClient{
		response: map[string]interface{}{
			"token":      "test-token",
			"expires_in": 3600,
		},
	}

	manager := NewManager(cfg, mockClient)
	token, err := manager.Token()
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	if token == "" {
		t.Error("Token() should return a non-empty token")
	}
	if token != "test-token" {
		t.Errorf("Token() = %v, want test-token", token)
	}
}

func TestManager_Token_Cached(t *testing.T) {
	cfg := config.NewConfiguration("public-key", "secret-key")
	mockClient := &mockHttpClient{
		response: map[string]interface{}{
			"token":      "test-token",
			"expires_in": 3600,
		},
	}

	manager := NewManager(cfg, mockClient)
	token1, err := manager.Token()
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}

	token2, err := manager.Token()
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}

	if token1 != token2 {
		t.Error("Token() should return cached token on second call")
	}
}

func TestManager_tokenValid(t *testing.T) {
	manager := &Manager{
		token:     "",
		expiresAt: 0,
	}
	if manager.tokenValid() {
		t.Error("tokenValid() should return false when token is empty")
	}

	manager.token = "test-token"
	manager.expiresAt = time.Now().Unix() + 100
	if !manager.tokenValid() {
		t.Error("tokenValid() should return true when token is valid")
	}

	manager.expiresAt = time.Now().Unix() - 100
	if manager.tokenValid() {
		t.Error("tokenValid() should return false when token is expired")
	}
}

type mockHttpClient struct {
	response map[string]interface{}
	err      error
}

func (m *mockHttpClient) Post(url string, headers map[string]string, body []byte) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}

	responseBody := `{}`
	if m.response != nil {
		bodyBytes, _ := json.Marshal(m.response)
		responseBody = string(bodyBytes)
	}

	return http.NewResponse(200, []byte(responseBody), nil), nil
}
