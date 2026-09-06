package errors

type ApiException struct {
	HTTPStatus   int
	ErrorCode    string
	ProviderCode string
	RequestID    string
	RetryAfter   string
	// Outcome is not_applied only when the server confirms no mutation occurred.
	Outcome string
	Message string
	Code    int
}

func (e *ApiException) Error() string {
	return e.Message
}

func (e *ApiException) GetCode() int {
	return e.Code
}

func NewApiException(message string, code int) *ApiException {
	if code == 0 {
		code = API_ERROR
	}
	return &ApiException{
		Message: message,
		Code:    code,
	}
}
