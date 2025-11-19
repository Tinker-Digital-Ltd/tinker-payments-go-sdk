package errors

type ClientException struct {
	Message string
	Code    int
	Cause   error
}

func (e *ClientException) Error() string {
	return e.Message
}

func (e *ClientException) Unwrap() error {
	return e.Cause
}

func (e *ClientException) GetCode() int {
	return e.Code
}

func NewClientException(message string, code int, cause error) *ClientException {
	if code == 0 {
		code = CLIENT_ERROR
	}
	return &ClientException{
		Message: message,
		Code:    code,
		Cause:   cause,
	}
}
