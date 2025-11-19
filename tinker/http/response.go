package http

import (
	"encoding/json"

	"github.com/tinker/tinker-payments-go-sdk/tinker/errors"
)

type Response struct {
	StatusCode int
	Body       []byte
	Headers    map[string][]string
}

func NewResponse(statusCode int, body []byte, headers map[string][]string) *Response {
	return &Response{
		StatusCode: statusCode,
		Body:       body,
		Headers:    headers,
	}
}

func (r *Response) JSON() (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return nil, errors.NewInvalidPayloadException("Invalid JSON response: "+err.Error(), 0, err)
	}
	return result, nil
}
