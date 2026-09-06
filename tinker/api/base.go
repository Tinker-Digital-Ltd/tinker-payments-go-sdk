package api

import (
	"encoding/json"
	nethttp "net/http"
	"strings"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/auth"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/errors"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model"
)

type BaseManager struct {
	config      *config.Configuration
	httpClient  http.Client
	authManager *auth.Manager
	lastMeta    *model.ApiMeta
}

func NewBaseManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *BaseManager {
	return &BaseManager{config: cfg, httpClient: httpClient, authManager: authManager}
}

func (bm *BaseManager) LastMeta() *model.ApiMeta { return bm.lastMeta }

func (bm *BaseManager) request(method, endpoint string, data map[string]interface{}) (map[string]interface{}, error) {
	baseURL := strings.TrimSuffix(bm.config.BaseURL, "/")
	endpoint = strings.TrimPrefix(endpoint, "/")
	url := baseURL + "/" + endpoint

	token, err := bm.authManager.Token()
	if err != nil {
		return nil, err
	}

	headers := map[string]string{"Authorization": "Bearer " + token, "Accept": "application/json", "Content-Type": "application/json"}
	var body []byte
	if len(data) > 0 {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, errors.NewNetworkException("Failed to serialize request: "+err.Error(), 0, err)
		}
	}

	var resp *http.Response
	if strings.EqualFold(method, "GET") {
		resp, err = bm.httpClient.Get(url, headers)
	} else {
		resp, err = bm.httpClient.Post(url, headers, body)
	}
	if err != nil {
		return nil, err
	}

	result, err := resp.JSON()
	if err != nil {
		if resp.StatusCode >= 400 {
			return nil, bm.apiError(resp, nil)
		}
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, bm.apiError(resp, result)
	}

	if result != nil {
		if metaMap, ok := result["meta"].(map[string]interface{}); ok {
			bm.lastMeta = model.NewApiMeta(metaMap)
		}
		if success, ok := result["success"].(bool); ok {
			if !success {
				return nil, bm.apiError(resp, result)
			}
			if dataMap, ok := result["data"].(map[string]interface{}); ok {
				return dataMap, nil
			}
			return map[string]interface{}{"value": result["data"]}, nil
		}
	}

	if result == nil {
		return map[string]interface{}{}, nil
	}
	return result, nil
}

func (bm *BaseManager) extractErrorMessage(result map[string]interface{}) string {
	if result != nil {
		if errorMap, ok := result["error"].(map[string]interface{}); ok {
			if msg, ok := errorMap["message"].(string); ok && msg != "" {
				return msg
			}
			if code, ok := errorMap["code"].(string); ok && code != "" {
				return code
			}
		}
		if msg, ok := result["message"].(string); ok && msg != "" {
			return msg
		}
	}
	return "Unknown error"
}

func (bm *BaseManager) apiError(resp *http.Response, result map[string]interface{}) error {
	e := errors.NewApiException(bm.extractErrorMessage(result), 0)
	e.HTTPStatus = resp.StatusCode
	e.RequestID = nethttp.Header(resp.Headers).Get("X-Request-ID")
	e.RetryAfter = nethttp.Header(resp.Headers).Get("Retry-After")
	if detail, ok := result["error"].(map[string]interface{}); ok {
		e.ErrorCode, _ = detail["code"].(string)
		e.ProviderCode, _ = detail["provider_code"].(string)
		e.Outcome, _ = detail["outcome"].(string)
		if id, ok := detail["request_id"].(string); ok && id != "" {
			e.RequestID = id
		}
	}
	return e
}
