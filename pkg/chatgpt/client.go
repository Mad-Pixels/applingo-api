package chatgpt

import (
	"context"

	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/pkg/errors"
)

// HTTPClient defines an interface for making HTTP POST requests.
type HTTPClient interface {
	Post(ctx context.Context, url string, data string, headers map[string]string) (string, error)
}

// Client represents a ChatGPT API client.
type Client struct {
	httpClient HTTPClient // HTTP client to perform requests
	apiKey     string     // API key for authentication
	baseURL    string     // Base URL for the ChatGPT API
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
		baseURL:    defaultOpenAIURL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithBaseURL returns an Option that sets a custom base URL for the API.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// IsAPIError checks if the provided error is an OpenAI API error of the specified type.
func IsAPIError(err error, errorType APIErrorType) bool {
	if apiErr, ok := errors.Cause(err).(APIError); ok {
		return apiErr.Type == errorType
	}
	return false
}

// SendMessage sends a message to the ChatGPT API using the provided Request.
// It validates the request, marshals the request into JSON,
// sends it via the HTTPClient, and then unmarshals the response.
// Returns a pointer to the Response struct or an error if any step fails.
func (c *Client) SendMessage(ctx context.Context, req *Request) (*Response, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

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
		if httpErr, ok := err.(interface{ Body() string }); ok {
			var errResp ErrorResponse
			if unmarshalErr := serializer.UnmarshalJSON([]byte(httpErr.Body()), &errResp); unmarshalErr == nil && errResp.Error.Message != "" {
				return nil, errResp.Error
			}
		}
		return nil, errors.Wrap(err, "failed to send request")
	}

	var resp Response
	if err = resp.Unmarshal([]byte(body)); err != nil {
		if _, ok := err.(APIError); ok {
			return nil, err
		}
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}
	if len(resp.Choices) == 0 {
		return nil, ErrEmptyResponse
	}
	return &resp, nil
}

// SendMessageString is a convenience method that creates a simple request with a single user message
// and returns just the response text.
func (c *Client) SendMessageString(ctx context.Context, model OpenAIModel, message string) (string, error) {
	req := NewRequest(model, []Message{NewUserMessage(message)})
	resp, err := c.SendMessage(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.GetResponseText(), nil
}
