// Package httpclient provides a simple HTTP client wrapper with retry logic,
// customizable retry conditions, and request timeout support.
package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// RetryConditionFunc defines a function that determines whether a request should be retried
// based on the HTTP status code and response body.
type RetryConditionFunc func(statusCode int, responseBody string) bool

// DefaultRetryCondition returns true if the status code is between 500 and 599.
func DefaultRetryCondition(statusCode int, _ string) bool {
	return statusCode >= 500 && statusCode < 600
}

// HTTPError represents an HTTP error including details about the request and response.
type HTTPError struct {
	StatusCode int
	Err        error
	body       string
}

// Error implements the error interface for HTTPError.
func (e HTTPError) Error() string {
	return fmt.Sprintf("failed: %d - %s", e.StatusCode, e.Err)
}

// Body returns the response body associated with this HTTP error.
func (e HTTPError) Body() string {
	return e.body
}

// ClientWrapper provides an interface for making HTTP requests with configurable
// timeout, retry logic, and retry conditions.
type ClientWrapper struct {
	client *http.Client

	maxRetries     int
	retryDelay     time.Duration
	retryCondition RetryConditionFunc
}

// New creates and returns a new instance of ClientWrapper with default settings.
func New() *ClientWrapper {
	return &ClientWrapper{
		client: &http.Client{},

		maxRetries:     0,
		retryDelay:     time.Second,
		retryCondition: DefaultRetryCondition,
	}
}

// WithTimeout sets the timeout duration for all requests made through this client.
func (c *ClientWrapper) WithTimeout(timeout time.Duration) *ClientWrapper {
	c.client.Timeout = timeout
	return c
}

// WithMaxRetries sets the maximum number of retries for failed requests, as well as the delay between retries.
func (c *ClientWrapper) WithMaxRetries(maxRetries int, retryDelay time.Duration) *ClientWrapper {
	c.maxRetries = maxRetries
	c.retryDelay = retryDelay
	return c
}

// WithRetryCondition sets a custom condition function that determines whether to retry a request.
func (c *ClientWrapper) WithRetryCondition(retryCondition RetryConditionFunc) *ClientWrapper {
	c.retryCondition = retryCondition
	return c
}

// request performs an HTTP request with the specified method, URL, data, and headers.
// It supports retries according to the configured retry condition and maximum attempts.
func (c *ClientWrapper) request(ctx context.Context, method, url, data string, headers map[string]string) (string, error) {
	var (
		resp *http.Response
		err  error
		body string
	)

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(c.retryDelay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}

		var bodyReader io.Reader
		if data != "" {
			bodyReader = bytes.NewBufferString(data)
		}
		req, reqErr := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if reqErr != nil {
			err = errors.Wrap(reqErr, "failed to create HTTP request")
			continue
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err = c.client.Do(req)
		if err != nil {
			err = errors.Wrap(err, "failed to execute request")
			continue
		}
		bodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			err = errors.Wrap(readErr, "failed to read response body")
			continue
		}
		body = string(bodyBytes)
		if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
			return body, nil
		}

		httpErr := HTTPError{
			Err:        fmt.Errorf("request failed with status code %d", resp.StatusCode),
			StatusCode: resp.StatusCode,
			body:       body,
		}
		if !c.retryCondition(resp.StatusCode, body) || attempt == c.maxRetries {
			return body, httpErr
		}
		err = httpErr
	}
	if err != nil {
		return "", err
	}
	return "", errors.New("unknown error occurred")
}

// Post performs an HTTP POST request with the provided context, URL, data, and headers.
// It returns the response body as a string and any error that occurred during the request.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - url: Target URL for the POST request.
//   - data: Request body as a string.
//   - headers: Map of request headers to be added.
//
// Returns:
//   - string: The response body.
//   - error: Any error that occurred during the request.
func (c *ClientWrapper) Post(ctx context.Context, url, data string, headers map[string]string) (string, error) {
	return c.request(ctx, http.MethodPost, url, data, headers)
}

// Get performs an HTTP GET request with the provided context, URL, and headers.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - url: Target URL for the GET request.
//   - headers: Map of request headers to be added.
//
// Returns:
//   - string: The response body.
//   - error: Any error that occurred during the request.
func (c *ClientWrapper) Get(ctx context.Context, url string, headers map[string]string) (string, error) {
	return c.request(ctx, http.MethodGet, url, "", headers)
}
