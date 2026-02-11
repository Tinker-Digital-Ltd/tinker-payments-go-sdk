package api

import (
	"net/url"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/auth"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model/dto"
)

type SubscriptionManager struct{ *BaseManager }

func NewSubscriptionManager(cfg *config.Configuration, httpClient http.Client, authManager *auth.Manager) *SubscriptionManager {
	return &SubscriptionManager{BaseManager: NewBaseManager(cfg, httpClient, authManager)}
}

func (sm *SubscriptionManager) CreatePlan(request *dto.CreateSubscriptionPlanRequestDto) (map[string]interface{}, error) {
	return sm.request("POST", config.SUBSCRIPTION_PLANS_PATH, request.ToMap())
}

func (sm *SubscriptionManager) ListPlans() ([]map[string]interface{}, error) {
	res, err := sm.request("GET", config.SUBSCRIPTION_PLANS_PATH, nil)
	if err != nil {
		return nil, err
	}
	out := []map[string]interface{}{}
	if value, ok := res["value"].([]interface{}); ok {
		for _, v := range value {
			if m, ok := v.(map[string]interface{}); ok {
				out = append(out, m)
			}
		}
	}
	return out, nil
}

func (sm *SubscriptionManager) Create(request *dto.CreateSubscriptionRequestDto) (map[string]interface{}, error) {
	return sm.request("POST", config.SUBSCRIPTION_BASE_PATH, request.ToMap())
}

func (sm *SubscriptionManager) List(planID, externalCustomerID string) ([]map[string]interface{}, error) {
	endpoint := config.SUBSCRIPTION_BASE_PATH
	q := url.Values{}
	if planID != "" {
		q.Add("plan_id", planID)
	}
	if externalCustomerID != "" {
		q.Add("external_customer_id", externalCustomerID)
	}
	if len(q) > 0 {
		endpoint += "?" + q.Encode()
	}

	res, err := sm.request("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	out := []map[string]interface{}{}
	if value, ok := res["value"].([]interface{}); ok {
		for _, v := range value {
			if m, ok := v.(map[string]interface{}); ok {
				out = append(out, m)
			}
		}
	}
	return out, nil
}

func (sm *SubscriptionManager) Cancel(subscriptionID string) (map[string]interface{}, error) {
	return sm.request("POST", config.SUBSCRIPTION_BASE_PATH+"/"+subscriptionID+"/cancel", nil)
}
