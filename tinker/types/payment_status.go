package types

type PaymentStatus string

const (
	PENDING            PaymentStatus = "pending"
	AUTHORISED         PaymentStatus = "authorised"
	BLOCKED            PaymentStatus = "blocked"
	EXPIRED            PaymentStatus = "expired"
	SUCCESS            PaymentStatus = "success"
	CANCELLED          PaymentStatus = "cancelled"
	FAILED             PaymentStatus = "failed"
	PARTIALLY_REFUNDED PaymentStatus = "partially_refunded"
	REFUNDED           PaymentStatus = "refunded"
	REVERSED           PaymentStatus = "reversed"
)

func ParsePaymentStatus(value string) PaymentStatus {
	switch value {
	case "completed":
		return SUCCESS
	case "":
		return PENDING
	}
	return PaymentStatus(value)
}
