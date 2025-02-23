package chatgpt

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultTimeout = 30 * time.Second
	defaultBaseURL = "https://api.openai.com/v1/chat/completions"
)

// HTTPClient defines an interface for making HTTP POST requests.
type HTTPClient interface {
	// Post sends an HTTP POST request to the specified URL with the provided data and headers.
	Post(ctx context.Context, url string, data string, headers map[string]string) (string, error)
}

// Client represents a ChatGPT API client.
type Client struct {
	httpClient HTTPClient    // HTTP client to perform requests
	apiKey     string        // API key for authentication
	baseURL    string        // Base URL for the ChatGPT API
	timeout    time.Duration // Timeout for API requests
}

// Option defines a functional option for configuring the Client.
type Option func(*Client)

// MustClient creates and returns a new ChatGPT API client instance.
// It requires a valid HTTPClient implementation and an API key.
func MustClient(httpClient HTTPClient, apiKey string, opts ...Option) *Client {
	if httpClient == nil {
		panic("http client cannot be nil")
	}

	c := &Client{
		httpClient: httpClient,
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		timeout:    defaultTimeout,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithTimeout returns an Option that sets a custom timeout for API requests.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithBaseURL returns an Option that sets a custom base URL for the API.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// SendMessage sends a message to the ChatGPT API using the provided Request.
// It validates the request, sets up a context with timeout, marshals the request into JSON,
// sends it via the HTTPClient, and then unmarshals the response.
// Returns a pointer to the Response struct or an error if any step fails.
func (c *Client) SendMessage(ctx context.Context, req *Request) (*Response, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.apiKey,
	}
	payload, err := req.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request")
	}

	body, err := c.httpClient.Post(ctx, c.baseURL, string(payload), headers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}
	var resp Response
	if err = resp.Unmarshal([]byte(body)); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	if len(resp.Choices) == 0 {
		return nil, ErrEmptyResponse
	}
	return &resp, nil
}
