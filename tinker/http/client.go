package http

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/tinker/tinker-payments-go-sdk/tinker/errors"
)

type Client interface {
	Post(url string, headers map[string]string, body []byte) (*Response, error)
}

type HttpClient struct {
	timeout time.Duration
	client  *http.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		timeout: 30 * time.Second,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *HttpClient) Post(url string, headers map[string]string, body []byte) (*Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return nil, errors.NewNetworkException("Network error: "+err.Error(), 0, err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.NewNetworkException("Network error: "+err.Error(), 0, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewNetworkException("Network error: "+err.Error(), 0, err)
	}

	return NewResponse(resp.StatusCode, respBody, resp.Header), nil
}
