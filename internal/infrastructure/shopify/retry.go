package shopify

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// RetryConfig configures retry behavior
type RetryConfig struct {
	MaxRetries      int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffFactor   float64
	RetryableErrors []int // HTTP status codes to retry
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:      3,
		InitialDelay:    100 * time.Millisecond,
		MaxDelay:        5 * time.Second,
		BackoffFactor:   2.0,
		RetryableErrors: []int{429, 500, 502, 503, 504}, // Rate limit + server errors
	}
}

// RetryableError checks if an error is retryable
func isRetryableError(err error, statusCode int, retryableStatuses []int) bool {
	// Check if status code is retryable
	for _, code := range retryableStatuses {
		if statusCode == code {
			return true
		}
	}

	// Check for network errors
	if err != nil {
		if netErr, ok := err.(net.Error); ok {
			return netErr.Timeout() || netErr.Temporary()
		}
		// Check for context timeout
		if err == context.DeadlineExceeded || err == context.Canceled {
			return false // Don't retry context cancellations
		}
	}

	return false
}

// RetryableFunc is a function that can be retried
type RetryableFunc func() (interface{}, *http.Response, error)

// ExecuteWithRetry executes a function with retry logic
func ExecuteWithRetry(
	ctx context.Context,
	fn RetryableFunc,
	config RetryConfig,
	logger zerolog.Logger,
) (interface{}, *http.Response, error) {
	var lastErr error
	var lastResp *http.Response
	delay := config.InitialDelay

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			logger.Warn().
				Int("attempt", attempt).
				Dur("delay", delay).
				Err(lastErr).
				Msg("Retrying Shopify API call")

			select {
			case <-ctx.Done():
				return nil, nil, ctx.Err()
			case <-time.After(delay):
			}

			// Exponential backoff
			delay = time.Duration(float64(delay) * config.BackoffFactor)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		}

		// Execute function
		result, resp, err := fn()

		// Check if we should retry
		statusCode := 0
		if resp != nil {
			statusCode = resp.StatusCode
		}

		if err == nil && !isRetryableError(err, statusCode, config.RetryableErrors) {
			// Success or non-retryable error
			return result, resp, err
		}

		// Store error for retry
		lastErr = err
		lastResp = resp

		// Check if we've exhausted retries
		if attempt >= config.MaxRetries {
			break
		}
	}

	// All retries exhausted
	logger.Error().
		Int("attempts", config.MaxRetries+1).
		Err(lastErr).
		Msg("All retry attempts exhausted")

	return nil, lastResp, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// RetryableHTTPRequest wraps an HTTP request with retry logic
func RetryableHTTPRequest(
	ctx context.Context,
	client *http.Client,
	req *http.Request,
	config RetryConfig,
	logger zerolog.Logger,
) (*http.Response, error) {
	var lastResp *http.Response
	var lastErr error
	delay := config.InitialDelay

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if attempt > 0 {
			logger.Warn().
				Int("attempt", attempt).
				Dur("delay", delay).
				Str("url", req.URL.String()).
				Msg("Retrying HTTP request")

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}

			delay = time.Duration(float64(delay) * config.BackoffFactor)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		}

		// Create new request with context for each retry
		newReq := req.Clone(ctx)

		// Execute request
		resp, err := client.Do(newReq)

		statusCode := 0
		if resp != nil {
			statusCode = resp.StatusCode
		}

		if err == nil && !isRetryableError(err, statusCode, config.RetryableErrors) {
			return resp, err
		}

		// Close response body if we're retrying
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		lastResp = resp
		lastErr = err

		if attempt >= config.MaxRetries {
			break
		}
	}

	return lastResp, fmt.Errorf("max retries exceeded: %w", lastErr)
}

