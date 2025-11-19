package tinker

import (
	"github.com/tinker/tinker-payments-go-sdk/tinker/api"
	"github.com/tinker/tinker-payments-go-sdk/tinker/auth"
	"github.com/tinker/tinker-payments-go-sdk/tinker/config"
	"github.com/tinker/tinker-payments-go-sdk/tinker/http"
	"github.com/tinker/tinker-payments-go-sdk/tinker/webhook"
)

type Payments struct {
	config             *config.Configuration
	httpClient         http.Client
	authManager        *auth.Manager
	transactionManager *api.TransactionManager
	webhookHandler     *webhook.Handler
}

func NewPayments(apiPublicKey, apiSecretKey string, httpClient http.Client) *Payments {
	cfg := config.NewConfiguration(apiPublicKey, apiSecretKey)

	var client http.Client
	if httpClient != nil {
		client = httpClient
	} else {
		client = http.NewHttpClient()
	}

	authMgr := auth.NewManager(cfg, client)

	return &Payments{
		config:      cfg,
		httpClient:  client,
		authManager: authMgr,
	}
}

func (p *Payments) Transactions() *api.TransactionManager {
	if p.transactionManager == nil {
		p.transactionManager = api.NewTransactionManager(p.config, p.httpClient, p.authManager)
	}
	return p.transactionManager
}

func (p *Payments) Webhooks() *webhook.Handler {
	if p.webhookHandler == nil {
		p.webhookHandler = webhook.NewHandler()
	}
	return p.webhookHandler
}
