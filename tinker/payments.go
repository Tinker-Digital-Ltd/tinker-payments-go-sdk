package tinker

import (
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/api"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/auth"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/config"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/http"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/webhook"
)

type Payments struct {
	config              *config.Configuration
	httpClient          http.Client
	authManager         *auth.Manager
	transactionManager  *api.TransactionManager
	payoutManager       *api.PayoutManager
	subscriptionManager *api.SubscriptionManager
	customerManager     *api.CustomerManager
	entitlementManager  *api.EntitlementManager
	productManager      *api.ProductManager
	offeringManager     *api.OfferingManager
	metricsManager      *api.MetricsManager
	experimentManager   *api.ExperimentManager
	integrationManager  *api.IntegrationManager
	webhookHandler      *webhook.Handler
}

func NewPayments(apiPublicKey, apiSecretKey string, httpClient http.Client, baseURL ...string) *Payments {
	cfg := config.NewConfiguration(apiPublicKey, apiSecretKey, baseURL...)

	var client http.Client
	if httpClient != nil {
		client = httpClient
	} else {
		client = http.NewHttpClient()
	}

	authMgr := auth.NewManager(cfg, client)
	return &Payments{config: cfg, httpClient: client, authManager: authMgr}
}

func (p *Payments) Transactions() *api.TransactionManager {
	if p.transactionManager == nil {
		p.transactionManager = api.NewTransactionManager(p.config, p.httpClient, p.authManager)
	}
	return p.transactionManager
}

func (p *Payments) Subscriptions() *api.SubscriptionManager {
	if p.subscriptionManager == nil {
		p.subscriptionManager = api.NewSubscriptionManager(p.config, p.httpClient, p.authManager)
	}
	return p.subscriptionManager
}

func (p *Payments) Payouts() *api.PayoutManager {
	if p.payoutManager == nil {
		p.payoutManager = api.NewPayoutManager(p.config, p.httpClient, p.authManager)
	}
	return p.payoutManager
}

func (p *Payments) Customers() *api.CustomerManager {
	if p.customerManager == nil {
		p.customerManager = api.NewCustomerManager(p.config, p.httpClient, p.authManager)
	}
	return p.customerManager
}

func (p *Payments) Entitlements() *api.EntitlementManager {
	if p.entitlementManager == nil {
		p.entitlementManager = api.NewEntitlementManager(p.config, p.httpClient, p.authManager)
	}
	return p.entitlementManager
}

func (p *Payments) Products() *api.ProductManager {
	if p.productManager == nil {
		p.productManager = api.NewProductManager(p.config, p.httpClient, p.authManager)
	}
	return p.productManager
}

func (p *Payments) Offerings() *api.OfferingManager {
	if p.offeringManager == nil {
		p.offeringManager = api.NewOfferingManager(p.config, p.httpClient, p.authManager)
	}
	return p.offeringManager
}

func (p *Payments) Metrics() *api.MetricsManager {
	if p.metricsManager == nil {
		p.metricsManager = api.NewMetricsManager(p.config, p.httpClient, p.authManager)
	}
	return p.metricsManager
}

func (p *Payments) Experiments() *api.ExperimentManager {
	if p.experimentManager == nil {
		p.experimentManager = api.NewExperimentManager(p.config, p.httpClient, p.authManager)
	}
	return p.experimentManager
}

func (p *Payments) Integrations() *api.IntegrationManager {
	if p.integrationManager == nil {
		p.integrationManager = api.NewIntegrationManager(p.config, p.httpClient, p.authManager)
	}
	return p.integrationManager
}

func (p *Payments) Webhooks() *webhook.Handler {
	if p.webhookHandler == nil {
		p.webhookHandler = webhook.NewHandler()
	}
	return p.webhookHandler
}

func (p *Payments) LastAuthMeta() *model.ApiMeta {
	return p.authManager.LastMeta()
}
