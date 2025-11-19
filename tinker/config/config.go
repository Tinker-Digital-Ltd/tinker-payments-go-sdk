package config

type Configuration struct {
	APIPublicKey string
	APISecretKey string
	BaseURL      string
}

func NewConfiguration(apiPublicKey, apiSecretKey string) *Configuration {
	return &Configuration{
		APIPublicKey: apiPublicKey,
		APISecretKey: apiSecretKey,
		BaseURL:      API_BASE_URL + "/",
	}
}

func (c *Configuration) APIKey() string {
	return c.APISecretKey
}
