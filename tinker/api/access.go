package api

import (
	"net/url"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/auth"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
)

type Body = map[string]interface{}

func listOf(res Body, err error) ([]Body, error) {
	if err != nil {
		return nil, err
	}
	out := []Body{}
	if value, ok := res["value"].([]interface{}); ok {
		for _, v := range value {
			if m, ok := v.(Body); ok {
				out = append(out, m)
			}
		}
	}
	return out, nil
}

type CustomerManager struct{ *BaseManager }

func NewCustomerManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *CustomerManager {
	return &CustomerManager{BaseManager: NewBaseManager(cfg, httpClient, authManager)}
}

func customerPath(id, suffix string) string {
	return "merchant/customers/" + url.PathEscape(id) + suffix
}

func (m *CustomerManager) Get(id string) (Body, error) {
	return m.request("GET", customerPath(id, ""), nil)
}

func (m *CustomerManager) GetEntitlement(id, identifier string) (Body, error) {
	return m.request("GET", customerPath(id, "/entitlements/"+url.PathEscape(identifier)), nil)
}

func (m *CustomerManager) GrantPromotionalEntitlement(id, identifier string, body Body) (Body, error) {
	return m.request("POST", customerPath(id, "/entitlements/"+url.PathEscape(identifier)+"/promotional"), body)
}

func (m *CustomerManager) RevokePromotionalEntitlements(id, identifier string) (Body, error) {
	return m.request("POST", customerPath(id, "/entitlements/"+url.PathEscape(identifier)+"/revoke_promotionals"), nil)
}

func (m *CustomerManager) Alias(id string, body Body) (Body, error) {
	return m.request("POST", customerPath(id, "/alias"), body)
}

func (m *CustomerManager) Transfer(id string, body Body) (Body, error) {
	return m.request("POST", customerPath(id, "/transfer"), body)
}

func (m *CustomerManager) SetAttributes(id string, body Body) (Body, error) {
	return m.request("POST", customerPath(id, "/attributes"), body)
}

func (m *CustomerManager) GetOfferings(id string) (Body, error) {
	return m.request("GET", customerPath(id, "/offerings"), nil)
}

func (m *CustomerManager) CreateToken(id string) (Body, error) {
	return m.request("POST", customerPath(id, "/token"), nil)
}

type EntitlementManager struct{ *BaseManager }

func NewEntitlementManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *EntitlementManager {
	return &EntitlementManager{BaseManager: NewBaseManager(cfg, httpClient, authManager)}
}

func entitlementPath(id, suffix string) string {
	return "merchant/entitlements/" + url.PathEscape(id) + suffix
}

func (m *EntitlementManager) List() ([]Body, error) {
	return listOf(m.request("GET", "merchant/entitlements", nil))
}

func (m *EntitlementManager) Create(body Body) (Body, error) {
	return m.request("POST", "merchant/entitlements", body)
}

func (m *EntitlementManager) Get(id string) (Body, error) {
	return m.request("GET", entitlementPath(id, ""), nil)
}

func (m *EntitlementManager) Update(id string, body Body) (Body, error) {
	return m.request("PATCH", entitlementPath(id, ""), body)
}

func (m *EntitlementManager) Delete(id string) (Body, error) {
	return m.request("DELETE", entitlementPath(id, ""), nil)
}

func (m *EntitlementManager) AttachProducts(id string, body Body) (Body, error) {
	return m.request("POST", entitlementPath(id, "/attach_products"), body)
}

func (m *EntitlementManager) DetachProducts(id string, body Body) (Body, error) {
	return m.request("POST", entitlementPath(id, "/detach_products"), body)
}

type ProductManager struct{ *BaseManager }

func NewProductManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *ProductManager {
	return &ProductManager{BaseManager: NewBaseManager(cfg, httpClient, authManager)}
}

func (m *ProductManager) List() ([]Body, error) {
	return listOf(m.request("GET", "merchant/products", nil))
}

func (m *ProductManager) Create(body Body) (Body, error) {
	return m.request("POST", "merchant/products", body)
}

func (m *ProductManager) Get(id string) (Body, error) {
	return m.request("GET", "merchant/products/"+url.PathEscape(id), nil)
}

func (m *ProductManager) Update(id string, body Body) (Body, error) {
	return m.request("PATCH", "merchant/products/"+url.PathEscape(id), body)
}

type OfferingManager struct{ *BaseManager }

func NewOfferingManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *OfferingManager {
	return &OfferingManager{BaseManager: NewBaseManager(cfg, httpClient, authManager)}
}

func (m *OfferingManager) List() ([]Body, error) {
	return listOf(m.request("GET", "merchant/offerings", nil))
}

func (m *OfferingManager) Create(body Body) (Body, error) {
	return m.request("POST", "merchant/offerings", body)
}

func (m *OfferingManager) Get(id string) (Body, error) {
	return m.request("GET", "merchant/offerings/"+url.PathEscape(id), nil)
}

func (m *OfferingManager) Update(id string, body Body) (Body, error) {
	return m.request("PATCH", "merchant/offerings/"+url.PathEscape(id), body)
}

func (m *OfferingManager) Delete(id string) (Body, error) {
	return m.request("DELETE", "merchant/offerings/"+url.PathEscape(id), nil)
}

type MetricsManager struct{ *BaseManager }

func NewMetricsManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *MetricsManager {
	return &MetricsManager{BaseManager: NewBaseManager(cfg, httpClient, authManager)}
}

func (m *MetricsManager) Get(from, to string) (Body, error) {
	q := url.Values{}
	if from != "" {
		q.Set("from", from)
	}
	if to != "" {
		q.Set("to", to)
	}
	endpoint := "merchant/metrics"
	if len(q) > 0 {
		endpoint += "?" + q.Encode()
	}
	return m.request("GET", endpoint, nil)
}

type ExperimentManager struct{ *BaseManager }

func NewExperimentManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *ExperimentManager {
	return &ExperimentManager{BaseManager: NewBaseManager(cfg, httpClient, authManager)}
}

func (m *ExperimentManager) List() ([]Body, error) {
	return listOf(m.request("GET", "merchant/experiments", nil))
}

func (m *ExperimentManager) Create(body Body) (Body, error) {
	return m.request("POST", "merchant/experiments", body)
}

func (m *ExperimentManager) Get(id string) (Body, error) {
	return m.request("GET", "merchant/experiments/"+url.PathEscape(id), nil)
}

func (m *ExperimentManager) Stop(id string) (Body, error) {
	return m.request("POST", "merchant/experiments/"+url.PathEscape(id)+"/stop", nil)
}

type IntegrationManager struct{ *BaseManager }

func NewIntegrationManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *IntegrationManager {
	return &IntegrationManager{BaseManager: NewBaseManager(cfg, httpClient, authManager)}
}

func (m *IntegrationManager) List() ([]Body, error) {
	return listOf(m.request("GET", "merchant/integrations", nil))
}

func (m *IntegrationManager) Upsert(kind string, body Body) (Body, error) {
	return m.request("PUT", "merchant/integrations/"+url.PathEscape(kind), body)
}

func (m *IntegrationManager) Delete(kind string) (Body, error) {
	return m.request("DELETE", "merchant/integrations/"+url.PathEscape(kind), nil)
}

func (m *IntegrationManager) Test(kind string) (Body, error) {
	return m.request("POST", "merchant/integrations/"+url.PathEscape(kind)+"/test", nil)
}

func (sm *SubscriptionManager) ListRenewalPrompts(subscriptionID string) ([]Body, error) {
	q := url.Values{"subscription_id": {subscriptionID}}
	return listOf(sm.request("GET", config.SUBSCRIPTION_BASE_PATH+"/renewal_prompts?"+q.Encode(), nil))
}

func (sm *SubscriptionManager) SendRenewalPrompt(subscriptionID string) (Body, error) {
	return sm.request("POST", config.SUBSCRIPTION_BASE_PATH+"/"+url.PathEscape(subscriptionID)+"/renewal_prompt", nil)
}
