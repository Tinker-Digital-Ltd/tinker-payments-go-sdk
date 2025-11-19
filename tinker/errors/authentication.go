package errors

type AuthenticationException struct {
	Message string
	Code    int
	Cause   error
}

func (e *AuthenticationException) Error() string {
	return e.Message
}

func (e *AuthenticationException) Unwrap() error {
	return e.Cause
}

func (e *AuthenticationException) GetCode() int {
	return e.Code
}

func NewAuthenticationException(message string, code int, cause error) *AuthenticationException {
	if code == 0 {
		code = AUTHENTICATION_ERROR
	}
	return &AuthenticationException{
		Message: message,
		Code:    code,
		Cause:   cause,
	}
}
