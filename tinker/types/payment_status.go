package types

type PaymentStatus string

const (
	PENDING   PaymentStatus = "pending"
	SUCCESS   PaymentStatus = "success"
	CANCELLED PaymentStatus = "cancelled"
	FAILED    PaymentStatus = "failed"
)
