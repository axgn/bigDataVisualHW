package pbdbapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// WithBaseURL overrides the API base URL.
func WithBaseURL(raw string) Option {
	return func(c *Client) error {
		u, err := url.Parse(strings.TrimSpace(raw))
		if err != nil {
			return fmt.Errorf("parse base URL: %w", err)
		}
		if u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("invalid base URL %q", raw)
		}
		c.baseURL = u
		return nil
	}
}

// WithTimeout sets timeout when the underlying client is *http.Client.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		if timeout <= 0 {
			return fmt.Errorf("timeout must be > 0")
		}
		hc, ok := c.httpClient.(*http.Client)
		if !ok {
			return nil
		}
		hc.Timeout = timeout
		return nil
	}
}

// WithUserAgent overrides the user agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) error {
		ua = strings.TrimSpace(ua)
		if ua == "" {
			return fmt.Errorf("user agent cannot be empty")
		}
		c.userAgent = ua
		return nil
	}
}

// WithHTTPClient replaces the default HTTP client implementation.
func WithHTTPClient(hc HTTPClient) Option {
	return func(c *Client) error {
		if hc == nil {
			return fmt.Errorf("http client cannot be nil")
		}
		c.httpClient = hc
		return nil
	}
}

// WithRetry configures retry attempts and base backoff delay.
func WithRetry(max int, baseDelay time.Duration) Option {
	return func(c *Client) error {
		if max < 0 {
			return fmt.Errorf("max retries cannot be negative")
		}
		if baseDelay <= 0 {
			return fmt.Errorf("base delay must be > 0")
		}
		c.retry.Enabled = true
		c.retry.MaxRetries = max
		c.retry.BaseDelay = baseDelay
		return nil
	}
}

// WithRetryEnabled toggles retries for 429/5xx responses.
func WithRetryEnabled(enabled bool) Option {
	return func(c *Client) error {
		c.retry.Enabled = enabled
		return nil
	}
}
