package errors

const FeatureNotSupportedCode = "FEATURE_NOT_SUPPORTED"

type FeatureNotSupportedException struct {
	*ApiException
	Gateway      string
	Feature      string
	Alternatives []string
}

func (e *FeatureNotSupportedException) Unwrap() error { return e.ApiException }
