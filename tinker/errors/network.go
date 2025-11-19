package errors

type NetworkException struct {
	Message string
	Code    int
	Cause   error
}

func (e *NetworkException) Error() string {
	return e.Message
}

func (e *NetworkException) Unwrap() error {
	return e.Cause
}

func (e *NetworkException) GetCode() int {
	return e.Code
}

func NewNetworkException(message string, code int, cause error) *NetworkException {
	if code == 0 {
		code = NETWORK_ERROR
	}
	return &NetworkException{
		Message: message,
		Code:    code,
		Cause:   cause,
	}
}
