# Tinker Payments Go SDK

Official Go SDK for [Tinker Payments API](https://payments.tinker.co.ke/docs).

## Installation

```bash
go get github.com/tinker/tinker-payments-go-sdk
```

## Requirements

- Go 1.18 or higher
- No external runtime dependencies (uses standard library)

## Quick Start

```go
package main

import (
    "github.com/tinker/tinker-payments-go-sdk/tinker"
)

func main() {
    payments := tinker.NewPayments(
        "your-public-key",
        "your-secret-key",
        nil, // nil uses default HTTP client
    )
}
```

## Usage

### Initiate a Payment

```go
package main

import (
    "fmt"
    "github.com/tinker/tinker-payments-go-sdk/tinker"
    "github.com/tinker/tinker-payments-go-sdk/tinker/types"
    "github.com/tinker/tinker-payments-go-sdk/tinker/model/dto"
    "github.com/tinker/tinker-payments-go-sdk/tinker/errors"
)

func main() {
    payments := tinker.NewPayments("your-public-key", "your-secret-key", nil)
    
    customerPhone := "+254712345678"
    transactionDesc := "Payment for order #12345"
    metadata := map[string]interface{}{
        "order_id": "12345",
    }
    
    initiateRequest := &dto.InitiatePaymentRequestDto{
        Amount:           100.00,
        Currency:         "KES",
        Gateway:          types.MPESA,
        MerchantReference: "ORDER-12345",
        ReturnURL:        "https://your-app.com/payment/return",
        CustomerPhone:    &customerPhone,
        TransactionDesc:  &transactionDesc,
        Metadata:         metadata,
    }
    
    transaction, err := payments.Transactions().Initiate(initiateRequest)
    if err != nil {
        if apiErr, ok := err.(*errors.ApiException); ok {
            fmt.Printf("API Error: %s\n", apiErr.Error())
        } else if netErr, ok := err.(*errors.NetworkException); ok {
            fmt.Printf("Network Error: %s\n", netErr.Error())
        } else {
            fmt.Printf("Error: %s\n", err.Error())
        }
        return
    }
    
    if transaction.InitiationData != nil && transaction.InitiationData.AuthorizationURL != nil {
        // Redirect user to authorization URL (Paystack, Stripe, etc.)
        fmt.Printf("Authorization URL: %s\n", *transaction.InitiationData.AuthorizationURL)
    }
}
```

**Note:** The `ReturnURL` is where users are redirected after payment completion. Webhooks are configured separately in your dashboard.

### Query a Transaction

```go
queryRequest := &dto.QueryPaymentRequestDto{
    PaymentReference: "TXN-abc123xyz",
    Gateway:          types.MPESA,
}

transaction, err := payments.Transactions().Query(queryRequest)
if err != nil {
    // Handle error
    return
}

if transaction.IsSuccessful() && transaction.QueryData != nil {
    queryData := transaction.QueryData
    fmt.Printf("Amount: %.2f %s\n", queryData.Amount, queryData.Currency)
}
```

### Handle Webhooks

Webhooks support multiple event types: payment, subscription, invoice, and settlement. Check the event type and handle accordingly:

```go
package main

import (
    "io"
    "net/http"
    "github.com/tinker/tinker-payments-go-sdk/tinker"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    payments := tinker.NewPayments("your-public-key", "your-secret-key", nil)
    
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Unable to read request body", http.StatusBadRequest)
        return
    }
    
    event, err := payments.Webhooks().Handle(body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Check event type
    if event.IsPaymentEvent() {
        paymentData := event.PaymentData()
        // Handle payment.completed, payment.failed, etc.
        fmt.Printf("Payment event: %s\n", event.Type)
    } else if event.IsSubscriptionEvent() {
        subscriptionData := event.SubscriptionData()
        // Handle subscription.created, subscription.cancelled, etc.
        fmt.Printf("Subscription event: %s\n", event.Type)
    } else if event.IsInvoiceEvent() {
        invoiceData := event.InvoiceData()
        // Handle invoice.paid, invoice.failed
        fmt.Printf("Invoice event: %s\n", event.Type)
    } else if event.IsSettlementEvent() {
        settlementData := event.SettlementData()
        // Handle settlement.processed
        fmt.Printf("Settlement event: %s\n", event.Type)
    }
    
    // Access event details
    fmt.Printf("Event type: %s\n", event.Type)        // e.g., "payment.completed"
    fmt.Printf("Event source: %s\n", event.Source)    // e.g., "payment"
    fmt.Printf("App ID: %s\n", event.Meta.AppID)
    fmt.Printf("Signature: %s\n", event.Security.Signature)
}
```

For payment events only, you can convert to a `Transaction` object:

```go
transaction, err := payments.Webhooks().HandleAsTransaction(body)
if err != nil {
    // Handle error
    return
}

if transaction != nil && transaction.IsSuccessful() {
    if transaction.CallbackData != nil {
        fmt.Printf("Payment successful: %s\n", transaction.CallbackData.Reference)
    }
}
```

## Custom HTTP Client

You can use your own HTTP client by implementing the `http.Client` interface:

```go
type Client interface {
    Post(url string, headers map[string]string, body []byte) (*http.Response, error)
}
```

Then pass it to the constructor:

```go
customClient := MyCustomHttpClient{}

payments := tinker.NewPayments(
    "your-public-key",
    "your-secret-key",
    customClient,
)
```

## Error Handling

```go
transaction, err := payments.Transactions().Initiate(request)
if err != nil {
    switch e := err.(type) {
    case *errors.ApiException:
        fmt.Printf("API Error (code %d): %s\n", e.GetCode(), e.Error())
    case *errors.NetworkException:
        fmt.Printf("Network Error (code %d): %s\n", e.GetCode(), e.Error())
    case *errors.AuthenticationException:
        fmt.Printf("Auth Error (code %d): %s\n", e.GetCode(), e.Error())
    default:
        fmt.Printf("Error: %s\n", err.Error())
    }
}
```

## Documentation

For detailed API documentation, visit [Tinker Payments API Documentation](https://payments.tinker.co.ke/docs).

## Development

After checking out the repo, run `go mod download` to install dependencies. Then, run `go test ./...` to run the tests.

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/tinker/payments-go-sdk.

## License

MIT License

