# v0.9.0 — Expire abandoned checkouts

- `ExpireCheckout(id, key)` expires an open managed checkout so an abandoned link cannot create a second subscription. An already expired checkout is returned unchanged; a completed one returns HTTP 409 `CHECKOUT_COMPLETED`.

Requires a Tinker Payments server with checkout expiry support.

# v0.8.0 — Complimentary trials

- `ManagedCheckout.TrialEnd` sets an exact trial end (48 hours to 730 days ahead); `TrialDays` now allows up to 730 days.
- `SetTrialEnd` moves the next charge of a trialing subscription, or extends a paid one; `EndTrial` ends a trial and charges now.

Requires a Tinker Payments server with complimentary trial support. See the README for the rules on paid subscriptions and scheduled seat changes.

# v0.7.0 — Managed billing trials

- `ManagedCheckout.TrialDays` starts a 1 to 30 day trial that requires a payment method; the first charge happens at trial end.
- `trial_days` is sent only when above 0, so existing checkouts are unchanged.
- Checkout and billing resource responses include the provider-neutral `managed_state`; webhooks add `billing_status`, `trial_ends_at` and `subscription.trial_will_end`.

Requires a Tinker Payments server with trial support. See the README for access rules.

# v0.6.0 — Refunds, merchant withdrawals and checkout cancellation

- Request refunds and read refund history and status through the transaction manager.
- List verified payout methods, obtain fee quotes, request merchant wallet withdrawals, and track payouts.
- Request checkout cancellation and inspect supported, rejected and uncertain outcomes.
- Use stable idempotency keys for mutations and retain structured errors, including request IDs and refreshed fee quotes.
- Parse signed refund, withdrawal and checkout-cancellation events as typed resources without creating payment transactions.
- Expose conflicts between checkout cancellation evidence and later successful collection.

Existing payment initiation, payment history, subscription and webhook APIs remain available. Withdrawal permission is controlled by the app owner. See the README for examples, status handling and retry behaviour.
