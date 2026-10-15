# Tinker Payments Go SDK

Official Go SDK for [Tinker Payments API](https://payments.tinker.co.ke/docs).

## Installation

```bash
go get github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk
```

## Requirements

- Go 1.18 or higher
- No external runtime dependencies (uses standard library)

## Quick Start

```go
package main

import (
    "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker"
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
    "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker"
    "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/types"
    "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model/dto"
    "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/errors"
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
    "fmt"
    "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker"
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
    } else if event.IsSubscriptionEvent() {
        subscriptionData := event.SubscriptionData()
        // Handle subscription.created, subscription.cancelled, etc.
    } else if event.IsInvoiceEvent() {
        invoiceData := event.InvoiceData()
        // Handle invoice.paid, invoice.failed
    } else if event.IsSettlementEvent() {
        settlementData := event.SettlementData()
        // Handle settlement.processed
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
    Get(url string, headers map[string]string) (*http.Response, error)
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

## Documentation

For detailed API documentation, visit [Tinker Payments API Documentation](https://payments.tinker.co.ke/docs).

## Development

After checking out the repo, run `go mod download` to install dependencies. Then, run `go test ./...` to run the tests.

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk.

## License

MIT License


## Exact-price recurring billing (v0.2.0)

`payments.Subscriptions()` supports `Checkout(ManagedCheckout)`, `BillingResource(id)`,
`ChangeQuantity`, `ScheduleQuantity`, `ClearScheduledQuantity`, `SetCancelAtRenewal`,
and `BillingOperation(idempotencyKey)`. These methods require a Tinker backend with
managed recurring billing enabled. Checkout amounts are integer minor units; no
client-side price conversion or automatic discount is applied.

Use a stable, unique idempotency key per command and reuse it only with the same
payload. After a timeout, `BillingOperation` retrieves the durable result without
replaying a provider charge. A missing or pending result does not prove a charge
failed. Do not retry with a new key to work around an uncertain result.

API failures can be inspected with `errors.As(err, &apiError)` where `apiError` is
`*tinker/errors.ApiException`. It preserves `HTTPStatus`, `ErrorCode`, `ProviderCode`,
`RequestID`, `RetryAfter` and `Outcome`, while retaining existing `Code`/`GetCode()`
compatibility. Only an explicit `Outcome == "not_applied"` confirms no mutation.

### Update a managed subscription's payment method

`payments.Subscriptions().PaymentMethodPortal(subscriptionID, returnURL, operationKey)`
creates a short-lived hosted payment-method link. It requires the managed-payment
endpoint on the Tinker server. Reuse the operation key after an uncertain response;
use a new key when deliberately opening a fresh session after the old link expires.
The flow cannot change seat quantity or price. Paying an overdue invoice remains a
separate, explicit action through its hosted invoice URL.

The server verifies app ownership and restricts this flow to a managed subscription
with a dedicated customer. Existing standard merchant contracts are excluded.

### Verified captured-payment history (v0.4.0)

`payments.Subscriptions().BillingPaymentHistory(subscriptionID, cursor, snapshot)`
reads app-scoped captured cash payments, tax-exclusive eligible amounts, successful
refunds, and dispute principal withdrawals/reinstatements. Pass an empty cursor and
zero snapshot for the first page. Persist the returned snapshot and use it with
`next_cursor` until `has_more` is false. Amounts are integer minor units; timestamps
are UTC RFC3339. Reinstatements identify the withdrawal through `reverses_id`.

An item with `reconciliation_reason` has incomplete or ambiguous provider evidence
and must be quarantined rather than credited. Commit a scan only after every page
has been verified. The snapshot is an event cutoff, not an immutable provider
snapshot: invoices and adjustments can change while reading. Rescan from the first
page regularly to discover later payments, refunds and reinstatements, and apply
facts idempotently by payment/adjustment ID. A zero-cash invoice is omitted from
`items` while still advancing invoice pagination. Currency conversion and customer
credits require explicit reconciliation; Stripe processing and dispute fees are
excluded from eligible principal.

This method requires the payment-history endpoint on the Tinker payment server.
The SDK release does not deploy or enable that server endpoint.

### Setup webhook acknowledgement

The dashboard sends a signed `setup.webhook_test` event with `source: "onboarding"`, `data.message`, and setup metadata. Parse it with the webhook handler, verify its signature using the application webhook secret, then acknowledge it with HTTP 200, 202, or 204. Setup events have onboarding data and do not convert to payment transactions. Event classification alone does not authenticate a request. Invalid signatures must be rejected; no payment lookup or fulfillment is needed for a valid setup test.

The signature uses HMAC-SHA256 over the JSON envelope excluding `security`, preserving envelope field order and Go JSON map ordering and escaping.

## Refunds, wallet withdrawals and checkout cancellation

These APIs require the merchant lifecycle backend routes under `/v1/merchant` to be deployed. SDK package releases must follow backend deployment and verified sandbox parity. Existing apps need the owner to enable withdrawal requests in app settings.

`payouts` means a withdrawal from the merchant's wallet to its own active, verified payout method. Amounts use currency units: withdrawing JPY 100 with a fee of 2 reserves 100 and pays 98. Obtain a quote and pass its `fee_amount` as `expected_fee_amount`; a fee conflict includes a refreshed quote in the structured API error. No method accepts arbitrary recipient banking details.

Supply a stable `Idempotency-Key` for every refund, payout or cancellation request. Keys are trimmed and must be 1–255 UTF-8 bytes. Reuse the same key and semantic request after transport failures; a changed request conflicts. A replay returns the original resource and fee snapshot. Start a new key for a new intended operation.

A refund's `pending_review`/`approved` state and a payout's `pending_approval`/`approved` state mean review is outstanding. Only `refunded`, `paid` or provider-confirmed `cancelled` establishes completion. Preserve cancellation `outcome_unknown` and poll by ID; do not blindly submit another request. Unsupported cancellation raises HTTP 422 `CHECKOUT_CANCELLATION_UNSUPPORTED` and retains the persisted `request_id`/`status` in error details. Lifecycle conflicts use HTTP 409; replaying a persisted rejection retains its 422/409 result.

Refund history includes refundable amounts and `refund_capability`. Cancellation history includes `capability`, the newest 100 requests and `has_more`; retrieve older known requests by ID. Payout listing accepts `page`, `per_page` (maximum 100), `status` and `currency`. Reads are scoped to the app; payout reads show only requests attributed to that app.

The signed events `refund.updated.v1`, `withdrawal.updated.v1` and `checkout.cancellation.updated.v1` have source `payment` and typed `data.resource`. Verify its signature before using its typed cancellation data. It never converts to a payment transaction. Events can arrive late or more than once; polling remains available.

```go
// Models are in tinker/model. Handle each error before using its result.
review, err := client.Transactions().Refund(paymentID, model.RefundRequest{
    Reason: model.RefundReasonRequestedByCustomer,
}, "refund-order-42")
if err != nil { return err }
fmt.Println(review.Refund.Status)
methods, err := client.Payouts().Methods()
if err != nil { return err }
if len(methods) == 0 { return fmt.Errorf("configure a verified payout method in the dashboard") }
destination := methods[0]
currency := destination.Currency
quote, err := client.Payouts().Quote(model.PayoutQuoteRequest{
    PayoutMethodID: destination.ID, Amount: 100, Currency: &currency,
})
if err != nil { return err }
requested, err := client.Payouts().Request(model.PayoutRequest{
    PayoutMethodID: destination.ID, Amount: 100, Currency: currency, ExpectedFeeAmount: quote.FeeAmount,
}, "withdrawal-42")
if err != nil { return err }
payout, err := client.Payouts().Get(requested.Withdrawal.ID)
if err != nil { return err }
fmt.Println(payout.Status)
result, err := client.Transactions().CancelCheckout(unpaidPaymentID, model.CancellationRequest{
    Reason: "Customer left checkout",
}, "cancel-order-42")
if err != nil {
    // *errors.ApiException retains HTTPStatus, ErrorCode and full Details,
    // including error.details.request_id for persisted rejection history.
    return err
}
current, err := client.Transactions().GetCancellation(unpaidPaymentID, result.Cancellation.ID)
if err != nil { return err }
fmt.Println(current.Status)
```

`Transactions()` also exposes `Refunds`, `GetRefund` and `Cancellations`; `Payouts().List(model.PayoutFilters{})` provides paginated tracking. Typed lifecycle resources, reasons, statuses and result envelopes are in `tinker/model`. Cancellation webhook data is available from `Event.CheckoutCancellationData()`.
