package config

const (
	PRODUCTION_BASE_URL               = "https://api.tinkerpayments.com"
	SANDBOX_BASE_URL                  = "https://sandbox-api.tinkerpayments.com"
	API_VERSION_PATH                  = "/v1"
	AUTH_TOKEN_PATH                   = "/auth/token"
	PAYMENT_INITIATE_PATH             = "/merchant/payment/initiate"
	PAYMENT_QUERY_PATH                = "/merchant/payment/query"
	PAYMENT_CAPABILITIES_PATH         = "/merchant/payment/capabilities"
	SUBSCRIPTION_BASE_PATH            = "/merchant/subscriptions"
	SUBSCRIPTION_PLANS_PATH           = "/merchant/subscriptions/plans"
	SUBSCRIPTION_INVOICES_PATH        = "/merchant/subscriptions/invoices"
	SUBSCRIPTION_RENEWAL_PROMPTS_PATH = "/merchant/subscriptions/renewal_prompts"
	CUSTOMERS_PATH                    = "/merchant/customers"
	ENTITLEMENTS_PATH                 = "/merchant/entitlements"
	PRODUCTS_PATH                     = "/merchant/products"
	OFFERINGS_PATH                    = "/merchant/offerings"
	METRICS_PATH                      = "/merchant/metrics"
	EXPERIMENTS_PATH                  = "/merchant/experiments"
	INTEGRATIONS_PATH                 = "/merchant/integrations"
)
