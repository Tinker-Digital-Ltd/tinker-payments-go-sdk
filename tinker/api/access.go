package api

import (
	"net/url"
	"strconv"

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
	return config.CUSTOMERS_PATH + "/" + url.PathEscape(id) + suffix
}

func (m *CustomerManager) GetCustomer(id string) (Body, error) {
	return m.request("GET", customerPath(id, ""), nil)
}

func (m *CustomerManager) GetEntitlement(id, identifier string) (Body, error) {
	return m.request("GET", customerPath(id, "/entitlements/"+url.PathEscape(identifier)), nil)
}

// PromotionalGrant sets Duration or EndTime, not both; both empty grants access until revoked. EndTime is RFC3339.
type PromotionalGrant struct {
	Duration string
	EndTime  string
}

func (m *CustomerManager) GrantPromotionalEntitlement(id, identifier string, grant PromotionalGrant) (Body, error) {
	body := Body{}
	if grant.Duration != "" {
		body["duration"] = grant.Duration
	}
	if grant.EndTime != "" {
		body["end_time"] = grant.EndTime
	}
	return m.request("POST", customerPath(id, "/entitlements/"+url.PathEscape(identifier)+"/promotional"), body)
}

func (m *CustomerManager) RevokePromotionalEntitlements(id, identifier string) (Body, error) {
	return m.request("POST", customerPath(id, "/entitlements/"+url.PathEscape(identifier)+"/revoke_promotionals"), nil)
}

func (m *CustomerManager) Alias(id, alias string) (Body, error) {
	return m.request("POST", customerPath(id, "/alias"), Body{"alias": alias})
}

func (m *CustomerManager) Transfer(id, to string) (Body, error) {
	return m.request("POST", customerPath(id, "/transfer"), Body{"to": to})
}

func (m *CustomerManager) SetAttributes(id string, attributes map[string]*string) (Body, error) {
	return m.request("POST", customerPath(id, "/attributes"), Body{"attributes": attributes})
}

func (m *CustomerManager) GetOfferings(id string) (Body, error) {
	return m.request("GET", customerPath(id, "/offerings"), nil)
}

func (m *CustomerManager) CreateToken(id string, expiresIn int) (Body, error) {
	var body Body
	if expiresIn != 0 {
		body = Body{"expires_in": expiresIn}
	}
	return m.request("POST", customerPath(id, "/token"), body)
}

type EntitlementManager struct{ *BaseManager }

func NewEntitlementManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *EntitlementManager {
	return &EntitlementManager{BaseManager: NewBaseManager(cfg, httpClient, authManager)}
}

func entitlementPath(id, suffix string) string {
	return config.ENTITLEMENTS_PATH + "/" + url.PathEscape(id) + suffix
}

func (m *EntitlementManager) List() ([]Body, error) {
	return listOf(m.request("GET", config.ENTITLEMENTS_PATH, nil))
}

func (m *EntitlementManager) Create(body Body) (Body, error) {
	return m.request("POST", config.ENTITLEMENTS_PATH, body)
}

func (m *EntitlementManager) Get(id string) (Body, error) {
	return m.request("GET", entitlementPath(id, ""), nil)
}

func (m *EntitlementManager) Update(id, displayName string) (Body, error) {
	return m.request("PATCH", entitlementPath(id, ""), Body{"display_name": displayName})
}

func (m *EntitlementManager) Delete(id string) (Body, error) {
	return m.request("DELETE", entitlementPath(id, ""), nil)
}

func (m *EntitlementManager) AttachProducts(id string, productIDs []string) (Body, error) {
	return m.request("POST", entitlementPath(id, "/attach_products"), Body{"product_ids": productIDs})
}

func (m *EntitlementManager) DetachProducts(id string, productIDs []string) (Body, error) {
	return m.request("POST", entitlementPath(id, "/detach_products"), Body{"product_ids": productIDs})
}

type ProductManager struct{ *BaseManager }

func NewProductManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *ProductManager {
	return &ProductManager{BaseManager: NewBaseManager(cfg, httpClient, authManager)}
}

func (m *ProductManager) List() ([]Body, error) {
	return listOf(m.request("GET", config.PRODUCTS_PATH, nil))
}

func (m *ProductManager) Create(body Body) (Body, error) {
	return m.request("POST", config.PRODUCTS_PATH, body)
}

func (m *ProductManager) Get(id string) (Body, error) {
	return m.request("GET", config.PRODUCTS_PATH+"/"+url.PathEscape(id), nil)
}

func (m *ProductManager) Update(id string, body Body) (Body, error) {
	return m.request("PATCH", config.PRODUCTS_PATH+"/"+url.PathEscape(id), body)
}

type OfferingManager struct{ *BaseManager }

func NewOfferingManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *OfferingManager {
	return &OfferingManager{BaseManager: NewBaseManager(cfg, httpClient, authManager)}
}

func (m *OfferingManager) List() ([]Body, error) {
	return listOf(m.request("GET", config.OFFERINGS_PATH, nil))
}

func (m *OfferingManager) Create(body Body) (Body, error) {
	return m.request("POST", config.OFFERINGS_PATH, body)
}

func (m *OfferingManager) Get(id string) (Body, error) {
	return m.request("GET", config.OFFERINGS_PATH+"/"+url.PathEscape(id), nil)
}

func (m *OfferingManager) Update(id string, body Body) (Body, error) {
	return m.request("PATCH", config.OFFERINGS_PATH+"/"+url.PathEscape(id), body)
}

func (m *OfferingManager) Delete(id string) (Body, error) {
	return m.request("DELETE", config.OFFERINGS_PATH+"/"+url.PathEscape(id), nil)
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
	endpoint := config.METRICS_PATH
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
	return listOf(m.request("GET", config.EXPERIMENTS_PATH, nil))
}

func (m *ExperimentManager) Create(body Body) (Body, error) {
	return m.request("POST", config.EXPERIMENTS_PATH, body)
}

func (m *ExperimentManager) Get(id string) (Body, error) {
	return m.request("GET", config.EXPERIMENTS_PATH+"/"+url.PathEscape(id), nil)
}

func (m *ExperimentManager) Stop(id string) (Body, error) {
	return m.request("POST", config.EXPERIMENTS_PATH+"/"+url.PathEscape(id)+"/stop", nil)
}

type IntegrationManager struct{ *BaseManager }

func NewIntegrationManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *IntegrationManager {
	return &IntegrationManager{BaseManager: NewBaseManager(cfg, httpClient, authManager)}
}

func (m *IntegrationManager) List() ([]Body, error) {
	return listOf(m.request("GET", config.INTEGRATIONS_PATH, nil))
}

func (m *IntegrationManager) Put(kind string, body Body) (Body, error) {
	return m.request("PUT", config.INTEGRATIONS_PATH+"/"+url.PathEscape(kind), body)
}

func (m *IntegrationManager) Delete(kind string) (Body, error) {
	return m.request("DELETE", config.INTEGRATIONS_PATH+"/"+url.PathEscape(kind), nil)
}

func (m *IntegrationManager) Test(kind string) (Body, error) {
	return m.request("POST", config.INTEGRATIONS_PATH+"/"+url.PathEscape(kind)+"/test", nil)
}

func (sm *SubscriptionManager) ListRenewalPrompts(subscriptionID string, page, pageSize int) (Body, error) {
	q := url.Values{}
	if subscriptionID != "" {
		q.Set("subscription_id", subscriptionID)
	}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		q.Set("page_size", strconv.Itoa(pageSize))
	}
	endpoint := config.SUBSCRIPTION_RENEWAL_PROMPTS_PATH
	if len(q) > 0 {
		endpoint += "?" + q.Encode()
	}
	return sm.request("GET", endpoint, nil)
}

func (sm *SubscriptionManager) SendRenewalPrompt(subscriptionID string) (Body, error) {
	return sm.request("POST", config.SUBSCRIPTION_BASE_PATH+"/"+url.PathEscape(subscriptionID)+"/renewal_prompt", nil)
}
