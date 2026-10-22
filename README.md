# Tinker Payments Go SDK

Go 1.25.3+ SDK using the standard library, with automatic authentication, payments, subscriptions and signed webhooks. Keep secret keys on the server.

## Install and configure

```sh
go get github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk
```

Test key prefixes select `https://sandbox-api.tinkerpayments.com/v1/`; other keys select `https://api.tinkerpayments.com/v1/`. Pass an optional fourth `baseURL` argument to `tinker.NewPayments` to override it. The third argument accepts the SDK's `tinker/http.Client` transport interface; `nil` uses its default transport.

## Payments

```go
package main

import (
	"fmt"
	"log"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/model/dto"
	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/types"
)

func main() {
	client := tinker.NewPayments("pk_test_xxx", "sk_test_xxx", nil)
	phone := "+254712345678"
	payment, err := client.Transactions().Initiate(&dto.InitiatePaymentRequestDto{
		Amount: 100, Currency: "KES", Gateway: types.MPESA,
		MerchantReference: "ORDER-12345",
		ReturnURL:         "https://merchant.example/payment/return", CustomerPhone: &phone,
	})
	if err != nil {
		log.Fatal(err)
	}
	if payment.InitiationData != nil && payment.InitiationData.AuthorizationURL != nil {
		fmt.Println(*payment.InitiationData.AuthorizationURL)
	}
	result, err := client.Transactions().Query(&dto.QueryPaymentRequestDto{
		PaymentReference: "TXN-abc123xyz", Gateway: types.MPESA,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.IsSuccessful())
}
```

Payments use currency units. Mpesa requires `CustomerPhone`; `ReturnURL` is the browser destination, while webhooks are configured separately in the dashboard.

Plan subscriptions use `Subscriptions().CreatePlan`, `ListPlans`, `Create`, `List` and `Cancel` with request DTOs in `tinker/model/dto`.

## Managed recurring billing

Checkout uses exact integer minor-unit prices (`15000` = EUR 150.00), a `month` or `year` interval, and a quantity. The SDK does not convert prices or apply discounts. Responses are unwrapped data maps.

The following fragments assume a `client` and run inside functions returning `error`; handle errors before using results.

```go
// Import "github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/api".
billing := client.Subscriptions()
checkout, err := billing.Checkout(api.ManagedCheckout{
	ExternalCustomerID: "workspace-123", Email: "buyer@example.com", Name: "Buyer",
	Currency: "EUR", UnitAmount: 15000, Quantity: 3, Interval: "year",
	ReturnURL: "https://example.com/billing", IdempotencyKey: "checkout-123",
})
if err != nil {
	return err
}
fmt.Println(checkout)
operation, err := billing.BillingOperation("checkout-123")
if err != nil {
	return err
}
fmt.Println(operation)
```

### Trials with a required payment method

Set `TrialDays` (1 to 730) or an exact `TrialEnd` (`*time.Time`, 48 hours to 730 days ahead; cannot be combined with `TrialDays`) to collect a card at checkout and charge first when the trial ends. The SDK sends `trial_days` only when it is above 0 and `trial_end` only when set, so checkouts without a trial keep working against servers released before trial support. A gateway without trial support returns HTTP 422 `TRIAL_UNSUPPORTED` with `outcome: not_applied`. Limiting a customer to one trial is the app's responsibility.

```go
checkout, err := billing.Checkout(api.ManagedCheckout{
	ExternalCustomerID: "workspace-123", Email: "buyer@example.com", Name: "Buyer",
	Currency: "EUR", UnitAmount: 1500, Quantity: 3, Interval: "month",
	ReturnURL: "https://example.com/billing", IdempotencyKey: "checkout-124",
	TrialDays: 14,
})
```

#### Changing the trial end

`SetTrialEnd(id, end, key)` moves the next charge to `end` (1 hour to 730 days ahead). A trialing subscription can move earlier or later; a paid (`active`) subscription only to a time on or after `current_period_end`, and it then becomes `trialing` with nothing charged. Clear scheduled seat changes first. `EndTrial(id, key)` ends a trial now and charges the full period immediately; if the charge fails the subscription becomes `past_due`. State-dependent rejections return HTTP 502 `BILLING_GATEWAY_ERROR` with `outcome: not_applied`.

```go
_, err := billing.SetTrialEnd("resource-id", time.Now().AddDate(0, 0, 30), "trial-end-1")
_, err = billing.EndTrial("resource-id", "end-trial-1")
```

Ending a trial early sends `invoice.paid` and `payment.updated` (not `subscription.renewed`). Extending a paid subscription produces a €0 `invoice.paid`.

Checkout and billing resource responses include a provider-neutral `managed_state` map: `status` (`pending`, `trialing`, `active`, `past_due`, `paused` or `canceled`), `trial_ends_at`, `current_period_end` and `cancel_at_period_end`. Read it instead of raw provider fields; grant access while `status` is `trialing` or `active` and `current_period_end` is in the future. Managed subscription webhooks carry `billing_status` and `trial_ends_at`, and `subscription.trial_will_end` arrives up to three days before the trial ends. The trial's €0 invoice sends `invoice.paid` with amount 0 and is not a payment.

Other methods are `BillingResource`, `ChangeQuantity`, `ScheduleQuantity`, `ClearScheduledQuantity`, `SetCancelAtRenewal`, `SetTrialEnd`, `EndTrial`, `PaymentMethodPortal` and `BillingPaymentHistory(subscriptionID, cursor, snapshot)`. Start history with an empty cursor and zero snapshot. `api.ManagedCheckout` uses `int64` minor-unit amounts and `int` quantities.

Use a stable key per mutation and reuse it with the same payload after a lost response. Recover through the operation lookup; missing or pending results remain uncertain. Only an explicit `outcome` of `not_applied` confirms no mutation.

The payment-method portal changes payment details only; quantity, price and overdue invoice payment are separate actions. It requires an app-owned managed subscription with a dedicated customer. Reuse the operation key after a lost response and use a new key for a deliberately new session after expiry.

For payment history, keep the returned `snapshot` across `next_cursor` pages until `has_more` is false. Amounts are integer minor units and timestamps are UTC RFC3339. History covers captured cash, tax-exclusive eligible amounts, refunds, and dispute principal withdrawals/reinstatements (`reverses_id` identifies the reversed withdrawal).

Quarantine items with `reconciliation_reason` and verify every page before committing a scan. The snapshot is an event cutoff, not an immutable provider snapshot: rescan from the first page to discover later adjustments, then apply facts idempotently by payment/adjustment ID. Zero-cash invoices advance pagination without an item. Processing/dispute fees are excluded from eligible principal; currency conversion and customer credits need explicit reconciliation.

## Refunds, wallet withdrawals and checkout cancellation

The app owner must enable withdrawal requests in app settings. A payout is a withdrawal from the merchant's wallet to its own active, verified payout method. Amounts use currency units: withdrawing JPY 100 with a fee of 2 reserves 100 and pays 98. Obtain a quote and pass its `fee_amount` as `expected_fee_amount`; a fee conflict includes a refreshed quote in the structured API error. No method accepts arbitrary recipient banking details.

Supply a stable `Idempotency-Key` for every refund, payout or cancellation request. Keys are trimmed and must be 1–255 UTF-8 bytes. Reuse the same key and semantic request after transport failures; a changed request conflicts. A replay returns the original resource and fee snapshot. Start a new key for a new intended operation.

A refund's `pending_review`/`approved` state and a payout's `pending_approval`/`approved` state are not completion. Only `refunded`, `paid` or provider-confirmed `cancelled` establishes completion. Preserve cancellation `outcome_unknown` and poll by ID; do not blindly submit another request. Unsupported cancellation raises HTTP 422 `CHECKOUT_CANCELLATION_UNSUPPORTED` and retains the persisted `request_id`/`status` in error details. Lifecycle conflicts use HTTP 409; replaying a persisted rejection retains its 422/409 result.

Refund history includes refundable amounts and `refund_capability`. Cancellation history includes `capability`, the newest 100 requests and `has_more`; retrieve older known requests by ID. Payout listing accepts `page`, `per_page` (maximum 100), `status` and `currency`. Reads are scoped to the app; payout reads show only requests attributed to that app.

```go
// Models are in tinker/model. Handle each error before using its result.
review, err := client.Transactions().Refund(paymentID, model.RefundRequest{
	Reason: model.RefundReasonRequestedByCustomer,
}, "refund-order-42")
if err != nil {
	return err
}
fmt.Println(review.Refund.Status)
methods, err := client.Payouts().Methods()
if err != nil {
	return err
}
var destination model.PayoutMethod
for _, method := range methods {
	if method.IsVerified && method.Status == "active" {
		destination = method
		break
	}
}
if destination.ID == "" {
	return fmt.Errorf("configure an active, verified payout method in the dashboard")
}
currency := destination.Currency
quote, err := client.Payouts().Quote(model.PayoutQuoteRequest{
	PayoutMethodID: destination.ID, Amount: 100, Currency: &currency,
})
if err != nil {
	return err
}
requested, err := client.Payouts().Request(model.PayoutRequest{
	PayoutMethodID: destination.ID, Amount: 100, Currency: currency, ExpectedFeeAmount: quote.FeeAmount,
}, "withdrawal-42")
if err != nil {
	return err
}
payout, err := client.Payouts().Get(requested.Withdrawal.ID)
if err != nil {
	return err
}
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
if err != nil {
	return err
}
fmt.Println(current.Status)
```

`Transactions()` also exposes `Refunds`, `GetRefund` and `Cancellations`; `Payouts().List(model.PayoutFilters{})` provides paginated tracking. Typed lifecycle resources, reasons, statuses and result envelopes are in `tinker/model`. Cancellation webhook data is available from `Event.CheckoutCancellationData()`.

Cancellation resources expose `payment_conflict` when later successful payment evidence conflicts with expiry. Treat successful payment as authoritative and use the refund workflow where appropriate.

## Webhooks

```go
event, err := client.Webhooks().Handle(body)
if err != nil {
	return err
}
if !client.Webhooks().VerifySignature(event, webhookSecret) {
	return fmt.Errorf("invalid webhook signature")
}
transaction := event.ToTransaction() // nil for setup/lifecycle events
if transaction != nil {
	fmt.Println(transaction.IsSuccessful())
}
```

Verify the signature before consuming fields or converting the event to a transaction; reuse the parsed event rather than reading the request body twice. Deduplicate by event ID. The SDK verifies HMAC-SHA256 over the envelope excluding `security`, using the server's field order, map ordering and escaping.

After verification, acknowledge `setup.webhook_test` (`source: onboarding`) with HTTP 200, 202 or 204; it needs no payment lookup or fulfillment. The signed lifecycle events `refund.updated.v1`, `withdrawal.updated.v1` and `checkout.cancellation.updated.v1` carry `data.resource` with `source: payment`. Their typed resources never convert into payment transactions. Poll request IDs if delivery is delayed.

## Errors and response metadata

Use standard-library `errors.As` with `*tinker/errors.ApiException` to inspect `HTTPStatus`, `ErrorCode`, `ProviderCode`, `RequestID`, `RetryAfter`, `Outcome` and the full envelope in `Details`. Request IDs prefer the response body over `X-Request-ID`; `Retry-After` remains its original string. Non-JSON HTTP failures retain status/headers and an uncertain outcome. Only an explicit `Outcome == "not_applied"` confirms no mutation.

Managers expose `LastMeta()`; `client.LastAuthMeta()` exposes authentication response metadata.

## Development

```sh
go test ./...
go vet ./...
```

[Release notes](RELEASE_NOTES.md) record additions. MIT License.
