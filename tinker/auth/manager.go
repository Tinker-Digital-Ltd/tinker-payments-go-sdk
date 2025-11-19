package auth

import (
	"encoding/base64"
	"net/url"
	"sync"
	"time"

	"github.com/tinker/tinker-payments-go-sdk/tinker/config"
	"github.com/tinker/tinker-payments-go-sdk/tinker/errors"
	"github.com/tinker/tinker-payments-go-sdk/tinker/http"
)

type Manager struct {
	config     *config.Configuration
	httpClient http.Client
	token      string
	expiresAt  int64
	mu         sync.RWMutex
}

func NewManager(cfg *config.Configuration, httpClient http.Client) *Manager {
	return &Manager{
		config:     cfg,
		httpClient: httpClient,
	}
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

	credentials := base64.StdEncoding.EncodeToString(
		[]byte(m.config.APIPublicKey + ":" + m.config.APISecretKey),
	)

	authURL := config.AUTH_TOKEN_URL
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Accept":       "application/json",
	}
	body := "credentials=" + url.QueryEscape(credentials)

	resp, err := m.httpClient.Post(authURL, headers, []byte(body))
	if err != nil {
		return "", errors.NewNetworkException("Failed to authenticate: "+err.Error(), errors.AUTHENTICATION_ERROR, err)
	}

	result, err := resp.JSON()
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		message := "Authentication failed"
		if msg, ok := result["message"].(string); ok {
			message = msg
		}
		return "", errors.NewAuthenticationException(message, 0, nil)
	}

	token, ok := result["token"].(string)
	if !ok || token == "" {
		return "", errors.NewNetworkException("Invalid authentication response: token missing", errors.AUTHENTICATION_ERROR, nil)
	}

	expiresIn := 3600
	if ei, ok := result["expires_in"].(float64); ok {
		expiresIn = int(ei)
	}

	m.token = token
	m.expiresAt = time.Now().Unix() + int64(expiresIn)

	return m.token, nil
}
