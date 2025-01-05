package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// ClientWrapper provides a interface for making HTTP requests.
type ClientWrapper struct {
	client *http.Client
}

// New creates and returns a new instance of ClientWrapper.
func New() *ClientWrapper {
	return &ClientWrapper{
		client: &http.Client{},
	}
}

// WithTimeout sets the timeout duration for all requests made through this client.
func (c *ClientWrapper) WithTimeout(timeout time.Duration) *ClientWrapper {
	c.client.Timeout = timeout
	return c
}

// Post performs an HTTP POST request with the provided context, URL, data, and headers.
// Returns the response body as a string and any error that occurred during the request.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - url: Target URL for the POST request
//   - data: Request body data as string
//   - headers: Map of request headers to be added
//
// Returns:
//   - string: Response body
//   - error: Any error that occurred during the request
func (c *ClientWrapper) Post(ctx context.Context, url string, data string, headers map[string]string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(data))
	if err != nil {
		return "", errors.Wrap(err, "failed to create HTTP request")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute request")
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}
	return string(body), nil
}
