package model

type PaymentCapabilities struct {
	Providers []ProviderCapability `json:"providers"`
}

type ProviderCapability struct {
	Gateway      string            `json:"gateway"`
	Name         string            `json:"name"`
	Methods      []string          `json:"methods"`
	Currencies   []string          `json:"currencies"`
	Countries    []string          `json:"countries"`
	Capabilities map[string]string `json:"capabilities"`
}
