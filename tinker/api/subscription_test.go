package api

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/auth"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model/dto"
)

type planHTTP struct {
	method, url string
	body        map[string]interface{}
}

func (c *planHTTP) Get(url string, headers map[string]string) (*http.Response, error) {
	return c.Do("GET", url, headers, nil)
}

func (c *planHTTP) Post(url string, headers map[string]string, body []byte) (*http.Response, error) {
	return c.Do("POST", url, headers, body)
}

func (c *planHTTP) Do(method, url string, headers map[string]string, body []byte) (*http.Response, error) {
	if strings.HasSuffix(url, "/auth/token") {
		return http.NewResponse(200, []byte(`{"token":"t","expires_in":3600}`), nil), nil
	}
	c.method, c.url, c.body = method, url, nil
	if len(body) > 0 {
		_ = json.Unmarshal(body, &c.body)
	}
	return http.NewResponse(200, []byte(`{"success":true,"data":{"ok":true}}`), nil), nil
}

func newPlanManager() (*SubscriptionManager, *planHTTP) {
	cfg := config.NewConfiguration("public-key", "secret-key")
	client := &planHTTP{}
	return NewSubscriptionManager(cfg, client, auth.NewManager(cfg, client)), client
}

func TestSubscriptionManager_UpdatePlan(t *testing.T) {
	manager, client := newPlanManager()
	plan := &dto.CreateSubscriptionPlanRequestDto{Name: "Pro", Amount: 10, Currency: "EUR", Intervals: []string{"monthly"}, IsActive: true}
	if _, err := manager.UpdatePlan("plan-1", plan); err != nil {
		t.Fatal(err)
	}
	if client.method != "PUT" || !strings.HasSuffix(client.url, "/merchant/subscriptions/plans/plan-1") || client.body["name"] != "Pro" {
		t.Fatalf("unexpected request %s %s %v", client.method, client.url, client.body)
	}
}

func TestCreatePlanRequestDto_FeaturesJSON(t *testing.T) {
	plan := &dto.CreateSubscriptionPlanRequestDto{Name: "Pro"}
	if _, ok := plan.ToMap()["features_json"]; ok {
		t.Fatal("features_json must be omitted when nil")
	}
	plan.FeaturesJSON = map[string]interface{}{"seats": 5}
	if plan.ToMap()["features_json"].(map[string]interface{})["seats"] != 5 {
		t.Fatal("features_json missing")
	}
}

func TestSubscriptionManager_DeletePlan(t *testing.T) {
	manager, client := newPlanManager()
	if _, err := manager.DeletePlan("plan-1"); err != nil {
		t.Fatal(err)
	}
	if client.method != "DELETE" || !strings.HasSuffix(client.url, "/merchant/subscriptions/plans/plan-1") || client.body != nil {
		t.Fatalf("unexpected request %s %s %v", client.method, client.url, client.body)
	}
}

func TestSubscriptionManager_ListInvoices(t *testing.T) {
	manager, client := newPlanManager()
	if _, err := manager.ListInvoices("sub-1", 2, 50); err != nil {
		t.Fatal(err)
	}
	if client.method != "GET" || !strings.HasSuffix(client.url, "/merchant/subscriptions/invoices?page=2&page_size=50&subscription_id=sub-1") {
		t.Fatalf("unexpected request %s %s", client.method, client.url)
	}
	if _, err := manager.ListInvoices("", 0, 0); err != nil || !strings.HasSuffix(client.url, "/merchant/subscriptions/invoices") {
		t.Fatalf("unexpected request %s %v", client.url, err)
	}
}
