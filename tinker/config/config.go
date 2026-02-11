package config

import "strings"

type Configuration struct {
	APIPublicKey string
	APISecretKey string
	BaseURL      string
	AuthURL      string
}

func NewConfiguration(apiPublicKey, apiSecretKey string, baseURL ...string) *Configuration {
	resolvedBaseURL := ""
	if len(baseURL) > 0 {
		resolvedBaseURL = baseURL[0]
	}

	if strings.TrimSpace(resolvedBaseURL) == "" {
		if strings.HasPrefix(apiPublicKey, "pk_test_") || strings.HasPrefix(apiSecretKey, "sk_test_") {
			resolvedBaseURL = SANDBOX_BASE_URL
		} else {
			resolvedBaseURL = PRODUCTION_BASE_URL
		}
	}

	resolvedBaseURL = strings.TrimSuffix(resolvedBaseURL, "/")
	if !strings.HasSuffix(resolvedBaseURL, API_VERSION_PATH) {
		resolvedBaseURL += API_VERSION_PATH
	}

	return &Configuration{
		APIPublicKey: apiPublicKey,
		APISecretKey: apiSecretKey,
		BaseURL:      resolvedBaseURL + "/",
		AuthURL:      resolvedBaseURL + AUTH_TOKEN_PATH,
	}
}

func (c *Configuration) APIKey() string {
	return c.APISecretKey
}
