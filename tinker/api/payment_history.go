package api

import (
	"net/url"
	"strconv"
)

// BillingPaymentHistory retrieves one app-scoped page of verified captured
// payments and itemized reversals. Retain snapshot across pagination; rescan
// completed history to discover later refunds and dispute reinstatements.
func (sm *SubscriptionManager) BillingPaymentHistory(id, cursor string, snapshot int64) (map[string]interface{}, error) {
	q := url.Values{}
	if cursor != "" {
		q.Set("starting_after", cursor)
	}
	if snapshot > 0 {
		q.Set("snapshot", strconv.FormatInt(snapshot, 10))
	}
	endpoint := "merchant/billing/resources/" + url.PathEscape(id) + "/payment-history"
	if len(q) > 0 {
		endpoint += "?" + q.Encode()
	}
	return sm.request("GET", endpoint, nil)
}
