package webhook

type Security struct {
	Signature string
	Algorithm string
}

func NewSecurity(security map[string]interface{}) *Security {
	s := &Security{
		Signature: "",
		Algorithm: "HMAC-SHA256",
	}

	if signature, ok := security["signature"].(string); ok {
		s.Signature = signature
	}

	if algorithm, ok := security["algorithm"].(string); ok {
		s.Algorithm = algorithm
	}

	return s
}
