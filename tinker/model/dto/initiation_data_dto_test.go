package dto

import "testing"

func TestNewInitiationDataDto_WalletFields(t *testing.T) {
	dto := NewInitiationDataDto(map[string]interface{}{
		"paymentReference": "pi_123",
		"status":           "pending",
		"clientSecret":     "pi_123_secret_abc",
		"publishableKey":   "pk_test_abc",
	})

	if dto.ClientSecret == nil || *dto.ClientSecret != "pi_123_secret_abc" {
		t.Errorf("ClientSecret = %v", dto.ClientSecret)
	}
	if dto.PublishableKey == nil || *dto.PublishableKey != "pk_test_abc" {
		t.Errorf("PublishableKey = %v", dto.PublishableKey)
	}
	if dto.AuthorizationURL != nil {
		t.Error("AuthorizationURL should be nil")
	}
}
