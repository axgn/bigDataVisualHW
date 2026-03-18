package pbdbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://paleobiodb.org/data1.2"

// HTTPClient allows custom transport injection for tests and advanced usage.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RetryConfig controls retry behavior for transient API failures.
type RetryConfig struct {
	Enabled    bool
	MaxRetries int
	BaseDelay  time.Duration
}

// Client is the main entry point for the PBDB Data Service API.
type Client struct {
	httpClient HTTPClient
	baseURL    *url.URL
	userAgent  string
	retry      RetryConfig
}

// Option configures a Client.
type Option func(*Client) error

// NewClient creates a new PBDB API client with sensible defaults.
func NewClient(opts ...Option) (*Client, error) {
	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse default base URL: %w", err)
	}

	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
		userAgent:  "pbdbapi-go/0.1",
		retry: RetryConfig{
			Enabled:    true,
			MaxRetries: 3,
			BaseDelay:  500 * time.Millisecond,
		},
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) buildURL(path string, query url.Values) (string, error) {
	rel, err := url.Parse(strings.TrimSpace(path))
	if err != nil {
		return "", fmt.Errorf("parse relative path %q: %w", path, err)
	}

	u := c.baseURL.ResolveReference(rel)
	u.RawQuery = query.Encode()
	return u.String(), nil
}

func (c *Client) doJSON(ctx context.Context, path string, query url.Values, out any) error {
	fullURL, err := c.buildURL(path, query)
	if err != nil {
		return err
	}

	attempt := 0
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
		if err != nil {
			return fmt.Errorf("build request: %w", err)
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.userAgent)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		body, readErr := io.ReadAll(resp.Body)
		closeErr := resp.Body.Close()
		if readErr != nil {
			return fmt.Errorf("read response body: %w", readErr)
		}
		if closeErr != nil {
			return fmt.Errorf("close response body: %w", closeErr)
		}

		if shouldRetry(c.retry, resp.StatusCode, attempt) {
			if err := sleepWithContext(ctx, backoff(c.retry.BaseDelay, attempt)); err != nil {
				return err
			}
			attempt++
			continue
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			return &StatusError{StatusCode: resp.StatusCode, Body: string(body)}
		}

		if apiErr := detectAPIError(body); apiErr != nil {
			return apiErr
		}

		if out == nil {
			return nil
		}

		if err := json.Unmarshal(body, out); err != nil {
			return &DecodeError{Err: err}
		}

		return nil
	}
}

func shouldRetry(cfg RetryConfig, statusCode, attempt int) bool {
	if !cfg.Enabled || attempt >= cfg.MaxRetries {
		return false
	}

	return statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError
}

func backoff(base time.Duration, attempt int) time.Duration {
	delay := base
	for i := 0; i < attempt; i++ {
		delay *= 2
	}
	return delay
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
