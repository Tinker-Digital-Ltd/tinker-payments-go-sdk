package errors

type ApiException struct {
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
