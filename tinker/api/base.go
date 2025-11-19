package api

import (
	"encoding/json"
	"strings"

	"github.com/tinker/tinker-payments-go-sdk/tinker/auth"
	"github.com/tinker/tinker-payments-go-sdk/tinker/config"
	"github.com/tinker/tinker-payments-go-sdk/tinker/errors"
	"github.com/tinker/tinker-payments-go-sdk/tinker/http"
)

type BaseManager struct {
	config      *config.Configuration
	httpClient  http.Client
	authManager *auth.Manager
}

func NewBaseManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *BaseManager {
	return &BaseManager{
		config:      cfg,
		httpClient:  httpClient,
		authManager: authManager,
	}
}

func (bm *BaseManager) request(method, endpoint string, data map[string]interface{}) (map[string]interface{}, error) {
	baseURL := strings.TrimSuffix(bm.config.BaseURL, "/")
	endpoint = strings.TrimPrefix(endpoint, "/")
	url := baseURL + "/" + endpoint

	token, err := bm.authManager.Token()
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
		"Content-Type":  "application/json",
	}

	var body []byte
	if len(data) > 0 {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, errors.NewNetworkException("Failed to serialize request: "+err.Error(), 0, err)
		}
	}

	resp, err := bm.httpClient.Post(url, headers, body)
	if err != nil {
		return nil, err
	}

	result, err := resp.JSON()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		message := "Unknown error"
		if msg, ok := result["message"].(string); ok {
			message = msg
		} else if errMsg, ok := result["error"].(string); ok {
			message = errMsg
		}
		return nil, errors.NewApiException(message, 0)
	}

	if result == nil {
		return make(map[string]interface{}), nil
	}

	return result, nil
}
