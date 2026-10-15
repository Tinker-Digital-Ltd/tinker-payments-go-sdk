# v0.6.0 — Refunds, merchant withdrawals and checkout cancellation

- Request refunds and read refund history and status through the transaction manager.
- List verified payout methods, obtain fee quotes, request merchant wallet withdrawals, and track payouts.
- Request checkout cancellation and inspect supported, rejected and uncertain outcomes.
- Use stable idempotency keys for mutations and retain structured errors, including request IDs and refreshed fee quotes.
- Parse signed refund, withdrawal and checkout-cancellation events as typed resources without creating payment transactions.
- Expose conflicts between checkout cancellation evidence and later successful collection.

Existing payment initiation, payment history, subscription and webhook APIs remain available. Withdrawal permission is controlled by the app owner. See the README for examples, status handling and retry behaviour.
