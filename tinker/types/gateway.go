package types

type Gateway string

const (
	MPESA    Gateway = "mpesa"
	PAYSTACK Gateway = "paystack"
	STRIPE   Gateway = "stripe"
)
