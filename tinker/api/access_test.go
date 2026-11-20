package api

import (
	"strings"
	"testing"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/auth"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
)

func TestAccessManagers(t *testing.T) {
	cfg := config.NewConfiguration("public-key", "secret-key")
	client := &planHTTP{}
	a := auth.NewManager(cfg, client)
	customers := NewCustomerManager(cfg, client, a)
	entitlements := NewEntitlementManager(cfg, client, a)
	products := NewProductManager(cfg, client, a)
	offerings := NewOfferingManager(cfg, client, a)
	metrics := NewMetricsManager(cfg, client, a)
	experiments := NewExperimentManager(cfg, client, a)
	integrations := NewIntegrationManager(cfg, client, a)
	subs := NewSubscriptionManager(cfg, client, a)
	body := Body{"k": "v"}

	cases := []struct {
		name         string
		call         func() error
		method, path string
		body         bool
	}{
		{"customers list", func() error { _, e := customers.ListCustomers("", 0, 0); return e }, "GET", "merchant/customers", false},
		{"customers list paged", func() error { _, e := customers.ListCustomers("a b", 2, 50); return e }, "GET", "merchant/customers?page=2&page_size=50&search=a+b", false},
		{"customer get", func() error { _, e := customers.GetCustomer("u/1"); return e }, "GET", "merchant/customers/u%2F1", false},
		{"get entitlement", func() error { _, e := customers.GetEntitlement("u1", "pro"); return e }, "GET", "merchant/customers/u1/entitlements/pro", false},
		{"grant", func() error {
			_, e := customers.GrantPromotionalEntitlement("u1", "pro", PromotionalGrant{Duration: "P7D"})
			return e
		}, "POST", "merchant/customers/u1/entitlements/pro/promotional", true},
		{"revoke", func() error { _, e := customers.RevokePromotionalEntitlements("u1", "pro"); return e }, "POST", "merchant/customers/u1/entitlements/pro/revoke_promotionals", false},
		{"alias", func() error { _, e := customers.Alias("u1", "a"); return e }, "POST", "merchant/customers/u1/alias", true},
		{"transfer", func() error { _, e := customers.Transfer("u1", "u2"); return e }, "POST", "merchant/customers/u1/transfer", true},
		{"attributes", func() error { _, e := customers.SetAttributes("u1", map[string]*string{"k": nil}); return e }, "POST", "merchant/customers/u1/attributes", true},
		{"offerings", func() error { _, e := customers.GetOfferings("u1"); return e }, "GET", "merchant/customers/u1/offerings", false},
		{"token", func() error { _, e := customers.CreateToken("u1", 0); return e }, "POST", "merchant/customers/u1/token", false},
		{"token expiry", func() error { _, e := customers.CreateToken("u1", 300); return e }, "POST", "merchant/customers/u1/token", true},
		{"entitlements list", func() error { _, e := entitlements.List(); return e }, "GET", "merchant/entitlements", false},
		{"entitlements create", func() error { _, e := entitlements.Create(body); return e }, "POST", "merchant/entitlements", true},
		{"entitlements get", func() error { _, e := entitlements.Get("e1"); return e }, "GET", "merchant/entitlements/e1", false},
		{"entitlements update", func() error { _, e := entitlements.Update("e1", "n"); return e }, "PATCH", "merchant/entitlements/e1", true},
		{"entitlements delete", func() error { _, e := entitlements.Delete("e1"); return e }, "DELETE", "merchant/entitlements/e1", false},
		{"attach", func() error { _, e := entitlements.AttachProducts("e1", []string{"p"}); return e }, "POST", "merchant/entitlements/e1/attach_products", true},
		{"detach", func() error { _, e := entitlements.DetachProducts("e1", []string{"p"}); return e }, "POST", "merchant/entitlements/e1/detach_products", true},
		{"products list", func() error { _, e := products.List(); return e }, "GET", "merchant/products", false},
		{"products create", func() error { _, e := products.Create(body); return e }, "POST", "merchant/products", true},
		{"products get", func() error { _, e := products.Get("p1"); return e }, "GET", "merchant/products/p1", false},
		{"products update", func() error { _, e := products.Update("p1", body); return e }, "PATCH", "merchant/products/p1", true},
		{"offerings list", func() error { _, e := offerings.List(); return e }, "GET", "merchant/offerings", false},
		{"offerings create", func() error { _, e := offerings.Create(body); return e }, "POST", "merchant/offerings", true},
		{"offerings get", func() error { _, e := offerings.Get("o1"); return e }, "GET", "merchant/offerings/o1", false},
		{"offerings update", func() error { _, e := offerings.Update("o1", body); return e }, "PATCH", "merchant/offerings/o1", true},
		{"offerings delete", func() error { _, e := offerings.Delete("o1"); return e }, "DELETE", "merchant/offerings/o1", false},
		{"metrics", func() error { _, e := metrics.Get("", ""); return e }, "GET", "merchant/metrics", false},
		{"metrics range", func() error { _, e := metrics.Get("2026-01-01", "2026-01-31"); return e }, "GET", "merchant/metrics?from=2026-01-01&to=2026-01-31", false},
		{"experiments list", func() error { _, e := experiments.List(); return e }, "GET", "merchant/experiments", false},
		{"experiments create", func() error { _, e := experiments.Create(body); return e }, "POST", "merchant/experiments", true},
		{"experiments get", func() error { _, e := experiments.Get("x1"); return e }, "GET", "merchant/experiments/x1", false},
		{"experiments stop", func() error { _, e := experiments.Stop("x1"); return e }, "POST", "merchant/experiments/x1/stop", false},
		{"integrations list", func() error { _, e := integrations.List(); return e }, "GET", "merchant/integrations", false},
		{"integrations put", func() error { _, e := integrations.Put("slack", body); return e }, "PUT", "merchant/integrations/slack", true},
		{"integrations delete", func() error { _, e := integrations.Delete("slack"); return e }, "DELETE", "merchant/integrations/slack", false},
		{"integrations test", func() error { _, e := integrations.Test("slack"); return e }, "POST", "merchant/integrations/slack/test", false},
		{"renewal prompts", func() error { _, e := subs.ListRenewalPrompts("sub_1", 0, 0); return e }, "GET", "merchant/subscriptions/renewal_prompts?subscription_id=sub_1", false},
		{"renewal prompts paged", func() error { _, e := subs.ListRenewalPrompts("sub_1", 2, 50); return e }, "GET", "merchant/subscriptions/renewal_prompts?page=2&page_size=50&subscription_id=sub_1", false},
		{"renewal prompts all", func() error { _, e := subs.ListRenewalPrompts("", 0, 0); return e }, "GET", "merchant/subscriptions/renewal_prompts", false},
		{"send renewal prompt", func() error { _, e := subs.SendRenewalPrompt("sub_1"); return e }, "POST", "merchant/subscriptions/sub_1/renewal_prompt", false},
	}
	for _, c := range cases {
		if err := c.call(); err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if client.method != c.method || !strings.HasSuffix(client.url, "/v1/"+c.path) {
			t.Errorf("%s: got %s %s, want %s %s", c.name, client.method, client.url, c.method, c.path)
		}
		if (client.body != nil) != c.body {
			t.Errorf("%s: body = %v", c.name, client.body)
		}
	}
}

func TestPromotionalGrantBody(t *testing.T) {
	cfg := config.NewConfiguration("public-key", "secret-key")
	client := &planHTTP{}
	customers := NewCustomerManager(cfg, client, auth.NewManager(cfg, client))
	for _, c := range []struct {
		grant PromotionalGrant
		key   string
	}{{PromotionalGrant{Duration: "P7D"}, "duration"}, {PromotionalGrant{EndTime: "2030-01-01T00:00:00Z"}, "end_time"}} {
		if _, err := customers.GrantPromotionalEntitlement("u1", "pro", c.grant); err != nil {
			t.Fatal(err)
		}
		if len(client.body) != 1 || client.body[c.key] == nil {
			t.Errorf("body = %v, want only %s", client.body, c.key)
		}
	}
}
