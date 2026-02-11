package auth

import (
	"encoding/base64"
	"net/url"
	"sync"
	"time"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/errors"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model"
)

type Manager struct {
	config     *config.Configuration
	httpClient http.Client
	token      string
	expiresAt  int64
	lastMeta   *model.ApiMeta
	mu         sync.RWMutex
}

func NewManager(cfg *config.Configuration, httpClient http.Client) *Manager {
	return &Manager{config: cfg, httpClient: httpClient}
}

func (m *Manager) Token() (string, error) {
	m.mu.RLock()
	if m.tokenValid() {
		token := m.token
		m.mu.RUnlock()
		return token, nil
	}
	m.mu.RUnlock()
	return m.fetchToken()
}

func (m *Manager) LastMeta() *model.ApiMeta {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastMeta
}

func (m *Manager) tokenValid() bool {
	if m.token == "" || m.expiresAt == 0 {
		return false
	}
	return time.Now().Unix() < (m.expiresAt - 60)
}

func (m *Manager) fetchToken() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.tokenValid() {
		return m.token, nil
	}

	credentials := base64.StdEncoding.EncodeToString([]byte(m.config.APIPublicKey + ":" + m.config.APISecretKey))
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded", "Accept": "application/json"}
	body := "credentials=" + url.QueryEscape(credentials)

	resp, err := m.httpClient.Post(m.config.AuthURL, headers, []byte(body))
	if err != nil {
		return "", errors.NewNetworkException("Failed to authenticate: "+err.Error(), errors.AUTHENTICATION_ERROR, err)
	}

	result, err := resp.JSON()
	if err != nil {
		return "", err
	}
	authData, err := m.extractAuthData(result)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", errors.NewAuthenticationException(m.extractErrorMessage(result), 0, nil)
	}

	token, ok := authData["token"].(string)
	if !ok || token == "" {
		return "", errors.NewNetworkException("Invalid authentication response: token missing", errors.AUTHENTICATION_ERROR, nil)
	}

	expiresIn := 3600
	if ei, ok := authData["expires_in"].(float64); ok {
		expiresIn = int(ei)
	}

	m.token = token
	m.expiresAt = time.Now().Unix() + int64(expiresIn)
	return m.token, nil
}

func (m *Manager) extractAuthData(result map[string]interface{}) (map[string]interface{}, error) {
	if result != nil {
		if meta, ok := result["meta"].(map[string]interface{}); ok {
			m.lastMeta = model.NewApiMeta(meta)
		}
		if success, ok := result["success"].(bool); ok {
			if !success {
				return nil, errors.NewAuthenticationException(m.extractErrorMessage(result), 0, nil)
			}
			if data, ok := result["data"].(map[string]interface{}); ok {
				return data, nil
			}
			return map[string]interface{}{}, nil
		}
	}
	if result == nil {
		return map[string]interface{}{}, nil
	}
	return result, nil
}

func (m *Manager) extractErrorMessage(result map[string]interface{}) string {
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
	return "Authentication failed"
}
