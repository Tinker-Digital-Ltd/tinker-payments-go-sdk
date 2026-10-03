package http

import (
	"bytes"
	"io"
	httpstd "net/http"
	"time"

	"github.com/Tinker-Digital-Ltd/tinker-payments-go-sdk/tinker/errors"
)

type Client interface {
	Get(url string, headers map[string]string) (*Response, error)
	Post(url string, headers map[string]string, body []byte) (*Response, error)
}

type HttpClient struct {
	timeout time.Duration
	client  *httpstd.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		timeout: 30 * time.Second,
		client:  &httpstd.Client{Timeout: 30 * time.Second},
	}
}

func (c *HttpClient) Get(url string, headers map[string]string) (*Response, error) {
	req, err := httpstd.NewRequest("GET", url, nil)
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
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		_ = resp.Body.Close()
		return nil, errors.NewNetworkException("Network error: "+err.Error(), 0, err)
	}
	if err := resp.Body.Close(); err != nil {
		return nil, errors.NewNetworkException("Network error: "+err.Error(), 0, err)
	}
	return NewResponse(resp.StatusCode, respBody, resp.Header), nil
}

func (c *HttpClient) Post(url string, headers map[string]string, body []byte) (*Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := httpstd.NewRequest("POST", url, bodyReader)
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
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		_ = resp.Body.Close()
		return nil, errors.NewNetworkException("Network error: "+err.Error(), 0, err)
	}
	if err := resp.Body.Close(); err != nil {
		return nil, errors.NewNetworkException("Network error: "+err.Error(), 0, err)
	}

	return NewResponse(resp.StatusCode, respBody, resp.Header), nil
}
