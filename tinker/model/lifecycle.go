package model

type RefundReason string

const (
	RefundReasonRequestedByCustomer RefundReason = "requested_by_customer"
	RefundReasonDuplicate           RefundReason = "duplicate"
	RefundReasonFraudulent          RefundReason = "fraudulent"
	RefundReasonOther               RefundReason = "other"
)

type RefundStatus string

const (
	RefundStatusPendingReview RefundStatus = "pending_review"
	RefundStatusProcessing    RefundStatus = "processing"
	RefundStatusApproved      RefundStatus = "approved"
	RefundStatusFailed        RefundStatus = "failed"
	RefundStatusRefunded      RefundStatus = "refunded"
	RefundStatusRejected      RefundStatus = "rejected"
)

type PayoutStatus string

const (
	PayoutStatusPendingApproval PayoutStatus = "pending_approval"
	PayoutStatusApproved        PayoutStatus = "approved"
	PayoutStatusProcessing      PayoutStatus = "processing"
	PayoutStatusPaid            PayoutStatus = "paid"
	PayoutStatusRejected        PayoutStatus = "rejected"
)

type CancellationStatus string

const (
	CancellationStatusRequested      CancellationStatus = "requested"
	CancellationStatusProcessing     CancellationStatus = "processing"
	CancellationStatusCancelled      CancellationStatus = "cancelled"
	CancellationStatusRejected       CancellationStatus = "rejected"
	CancellationStatusOutcomeUnknown CancellationStatus = "outcome_unknown"
)

type RefundRequest struct {
	Amount *float64     `json:"amount,omitempty"`
	Reason RefundReason `json:"reason"`
	Note   *string      `json:"note,omitempty"`
}

type Refund struct {
	RequestedByAppID *string      `json:"requested_by_app_id,omitempty"`
	ID               string       `json:"id"`
	PaymentID        string       `json:"payment_id"`
	Amount           float64      `json:"amount"`
	Currency         string       `json:"currency"`
	Reason           RefundReason `json:"reason"`
	Note             *string      `json:"note,omitempty"`
	Source           string       `json:"source"`
	Status           RefundStatus `json:"status"`
	Execution        *string      `json:"execution,omitempty"`
	ProviderRefundID *string      `json:"provider_refund_id,omitempty"`
	ManualReference  *string      `json:"manual_reference,omitempty"`
	RejectionReason  *string      `json:"rejection_reason,omitempty"`
	FailureMessage   *string      `json:"failure_message,omitempty"`
	ReviewedAt       *string      `json:"reviewed_at,omitempty"`
	CompletedAt      *string      `json:"completed_at,omitempty"`
	CreatedAt        string       `json:"created_at"`
	UpdatedAt        string       `json:"updated_at"`
}

type RefundCapability struct {
	Supported         bool     `json:"supported"`
	Partial           bool     `json:"partial"`
	Reasons           []string `json:"reasons"`
	MaxAmount         float64  `json:"max_amount"`
	ReviewRequired    bool     `json:"review_required"`
	Execution         string   `json:"execution"`
	UnavailableReason *string  `json:"unavailable_reason,omitempty"`
}

type RefundHistory struct {
	RefundedAmount   float64          `json:"refunded_amount"`
	RefundableAmount float64          `json:"refundable_amount"`
	RefundCapability RefundCapability `json:"refund_capability"`
	Refunds          []Refund         `json:"refunds"`
}

type RefundResult struct {
	Refund   Refund        `json:"refund"`
	Payment  PaymentDetail `json:"payment"`
	Replayed bool          `json:"replayed"`
}

type PaymentDetail struct {
	ID               string            `json:"id"`
	Status           string            `json:"status"`
	Amount           float64           `json:"amount"`
	Currency         string            `json:"currency"`
	RefundableAmount *float64          `json:"refundable_amount,omitempty"`
	RefundCapability *RefundCapability `json:"refund_capability,omitempty"`
	Refunds          []Refund          `json:"refunds,omitempty"`
}

type PayoutMethod struct {
	ID                string  `json:"id"`
	MethodType        string  `json:"method_type"`
	Provider          *string `json:"provider,omitempty"`
	Currency          string  `json:"currency"`
	Status            string  `json:"status"`
	IsVerified        bool    `json:"is_verified"`
	IsDefault         bool    `json:"is_default"`
	MaskedDestination string  `json:"masked_destination"`
}

type WithdrawalFeeRule struct {
	Fixed   float64 `json:"fixed"`
	Percent float64 `json:"percent"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
}

type WithdrawalArrival struct {
	Min  int    `json:"min"`
	Max  int    `json:"max"`
	Unit string `json:"unit"`
}

type PayoutQuote struct {
	MethodType  string             `json:"method_type"`
	Currency    string             `json:"currency"`
	Amount      *float64           `json:"amount"`
	FeeAmount   float64            `json:"fee_amount"`
	NetAmount   *float64           `json:"net_amount"`
	Fee         WithdrawalFeeRule  `json:"fee"`
	ReviewTime  string             `json:"review_time"`
	Arrival     *WithdrawalArrival `json:"arrival"`
	ArrivalText string             `json:"arrival_text"`
	Configured  bool               `json:"configured"`
	IsEstimate  bool               `json:"is_estimate"`
}

type PayoutRequest struct {
	PayoutMethodID    string  `json:"payout_method_id"`
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"`
	ExpectedFeeAmount float64 `json:"expected_fee_amount"`
}

type Payout struct {
	ID                string       `json:"id"`
	AccountID         string       `json:"account_id"`
	PayoutMethodID    string       `json:"payout_method_id"`
	Amount            float64      `json:"amount"`
	Currency          string       `json:"currency"`
	FeeAmount         float64      `json:"fee_amount"`
	NetAmount         float64      `json:"net_amount"`
	Status            PayoutStatus `json:"status"`
	Source            *string      `json:"source,omitempty"`
	RequestedByAppID  *string      `json:"requested_by_app_id,omitempty"`
	RequestedByUserID *string      `json:"requested_by_user_id,omitempty"`
	RejectionReason   *string      `json:"rejection_reason,omitempty"`
	RequestedAt       string       `json:"requested_at"`
	ApprovedAt        *string      `json:"approved_at,omitempty"`
	ProcessingAt      *string      `json:"processing_at,omitempty"`
	PaidAt            *string      `json:"paid_at,omitempty"`
	RejectedAt        *string      `json:"rejected_at,omitempty"`
	CreatedAt         string       `json:"created_at"`
	UpdatedAt         string       `json:"updated_at"`
}

type PayoutResult struct {
	Withdrawal Payout `json:"withdrawal"`
	Replayed   bool   `json:"replayed"`
}

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PayoutPage struct {
	Payouts    []Payout   `json:"payouts"`
	Pagination Pagination `json:"pagination"`
}

type CancellationRequest struct {
	Reason string `json:"reason"`
}

type Cancellation struct {
	ID               string             `json:"id"`
	PaymentID        string             `json:"payment_id"`
	Status           CancellationStatus `json:"status"`
	Reason           string             `json:"reason"`
	ReasonCode       string             `json:"reason_code"`
	Source           string             `json:"source"`
	RequestedByAppID *string            `json:"requested_by_app_id,omitempty"`
	CreatedAt        string             `json:"created_at"`
	UpdatedAt        string             `json:"updated_at"`
	CompletedAt      *string            `json:"completed_at,omitempty"`
	FailureMessage   *string            `json:"failure_message,omitempty"`
}

type CancellationCapability struct {
	Supported  bool   `json:"supported"`
	Eligible   bool   `json:"eligible"`
	ReasonCode string `json:"reason_code"`
}

type CancellationHistory struct {
	HasMore       bool                   `json:"has_more"`
	Cancellations []Cancellation         `json:"cancellations"`
	Capability    CancellationCapability `json:"capability"`
}

type CancellationResult struct {
	Cancellation Cancellation `json:"cancellation"`
	Replayed     bool         `json:"replayed"`
}

type PayoutFilters struct {
	Page     *int          `json:"page,omitempty"`
	PerPage  *int          `json:"per_page,omitempty"`
	Status   *PayoutStatus `json:"status,omitempty"`
	Currency *string       `json:"currency,omitempty"`
}

type PayoutQuoteRequest struct {
	PayoutMethodID string  `json:"payout_method_id"`
	Amount         float64 `json:"amount"`
	Currency       *string `json:"currency,omitempty"`
}
