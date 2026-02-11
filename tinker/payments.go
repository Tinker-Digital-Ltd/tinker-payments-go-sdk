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
	subscriptionManager *api.SubscriptionManager
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

func (p *Payments) Webhooks() *webhook.Handler {
	if p.webhookHandler == nil {
		p.webhookHandler = webhook.NewHandler()
	}
	return p.webhookHandler
}

func (p *Payments) LastAuthMeta() *model.ApiMeta {
	return p.authManager.LastMeta()
}
