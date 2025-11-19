package errors

type InvalidPayloadException struct {
	Message string
	Code    int
	Cause   error
}

func (e *InvalidPayloadException) Error() string {
	return e.Message
}

func (e *InvalidPayloadException) Unwrap() error {
	return e.Cause
}

func (e *InvalidPayloadException) GetCode() int {
	return e.Code
}

func NewInvalidPayloadException(message string, code int, cause error) *InvalidPayloadException {
	if code == 0 {
		code = INVALID_PAYLOAD
	}
	return &InvalidPayloadException{
		Message: message,
		Code:    code,
		Cause:   cause,
	}
}
