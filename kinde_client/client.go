package kinde_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// Client represents the Kinde API client
type Client struct {
	HTTPClient *http.Client
	Token      string
	Auth       authStruct
	IssuerUrl  string
	limiter    *rate.Limiter
	logger     *log.Logger
}

type authStruct struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithTimeout sets the client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.HTTPClient.Timeout = timeout
	}
}

// WithRateLimit sets the rate limit (requests per second)
func WithRateLimit(rps int) ClientOption {
	return func(c *Client) {
		c.limiter = rate.NewLimiter(rate.Limit(rps), rps)
	}
}

// WithLogger sets a custom logger
func WithLogger(logger *log.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// NewClient creates a new Kinde API client
func NewClient(ctx context.Context, issuerUrl, clientId, clientSecret string, opts ...ClientOption) (*Client, error) {
	c := &Client{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		IssuerUrl: issuerUrl,
		limiter:   rate.NewLimiter(rate.Limit(100), 100), // Default: 100 requests per second
		logger:    log.New(io.Discard, "", 0),            // Default: no logging
	}

	c.Auth = authStruct{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	tok, err := c.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	c.Token = *tok
	return c, nil
}

// doRequest handles the common request logic with retries and rate limiting
func (c *Client) doRequest(ctx context.Context, method, url string, body interface{}, maxRetries int) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Wait for rate limiter
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit wait failed: %w", err)
		}

		// Prepare request
		var reqBody io.Reader
		if body != nil {
			jsonBody, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
			reqBody = bytes.NewBuffer(jsonBody)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
		if body != nil {
			req.Header.Add("Content-Type", "application/json")
		}

		// Log request
		c.logger.Printf("Making %s request to %s", method, url)

		// Execute request
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return nil, lastErr
		}

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return nil, lastErr
		}

		// Log response
		c.logger.Printf("Received response with status %d", resp.StatusCode)

		// Handle response
		if !isHttpStatusCodeSuccess(resp.StatusCode) {
			lastErr = fmt.Errorf("received status code %d: %s", resp.StatusCode, string(respBody))
			if isRetryableError(resp.StatusCode) && attempt < maxRetries {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return nil, lastErr
		}

		return respBody, nil
	}

	return nil, lastErr
}

// doGetRequest performs a GET request with retries and rate limiting
func doGetRequest[T any](c *Client, ctx context.Context, url string) (T, error) {
	var output T
	respBody, err := c.doRequest(ctx, "GET", url, nil, 3)
	if err != nil {
		return output, err
	}

	if err := json.Unmarshal(respBody, &output); err != nil {
		return output, fmt.Errorf("failed to decode response: %w", err)
	}

	return output, nil
}

// doPostRequest performs a POST request with retries and rate limiting
func doPostRequest[T any](c *Client, ctx context.Context, url string, body interface{}) (T, error) {
	var output T
	respBody, err := c.doRequest(ctx, "POST", url, body, 3)
	if err != nil {
		return output, err
	}

	if err := json.Unmarshal(respBody, &output); err != nil {
		return output, fmt.Errorf("failed to decode response: %w", err)
	}

	return output, nil
}

// isHttpStatusCodeSuccess checks if the status code indicates success
func isHttpStatusCodeSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// isRetryableError determines if an error should trigger a retry
func isRetryableError(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusInternalServerError ||
		statusCode == http.StatusBadGateway ||
		statusCode == http.StatusServiceUnavailable ||
		statusCode == http.StatusGatewayTimeout
}
