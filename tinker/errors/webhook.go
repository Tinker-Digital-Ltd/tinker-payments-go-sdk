package errors

type WebhookException struct {
	Message string
	Code    int
	Cause   error
}

func (e *WebhookException) Error() string {
	return e.Message
}

func (e *WebhookException) Unwrap() error {
	return e.Cause
}

func (e *WebhookException) GetCode() int {
	return e.Code
}

func NewWebhookException(message string, code int, cause error) *WebhookException {
	if code == 0 {
		code = WEBHOOK_ERROR
	}
	return &WebhookException{
		Message: message,
		Code:    code,
		Cause:   cause,
	}
}
