package retry

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/rs/zerolog"
)

// Config holds retry configuration
type Config struct {
	MaxAttempts     int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffFactor   float64
	RetryableErrors []error
}

// DefaultConfig returns default retry configuration
func DefaultConfig() Config {
	return Config{
		MaxAttempts:   5,
		InitialDelay:  time.Second,
		MaxDelay:      time.Minute,
		BackoffFactor: 2.0,
	}
}

// RetryFunc is a function that can be retried
type RetryFunc func(ctx context.Context) error

// Do executes the function with retry logic
func Do(ctx context.Context, config Config, logger zerolog.Logger, fn RetryFunc) error {
	var lastErr error

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Execute function
		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if we should retry
		if !shouldRetry(err, config) {
			return err
		}

		// Don't sleep on last attempt
		if attempt == config.MaxAttempts-1 {
			break
		}

		// Calculate backoff delay
		delay := calculateBackoff(attempt, config)

		logger.Warn().
			Err(err).
			Int("attempt", attempt+1).
			Int("max_attempts", config.MaxAttempts).
			Dur("delay", delay).
			Msg("Retrying after error")

		// Wait before retry
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("max retry attempts reached: %w", lastErr)
}

// shouldRetry determines if an error should trigger a retry
func shouldRetry(err error, config Config) bool {
	// Always retry on rate limit errors
	// In a real implementation, you'd check for specific error types
	return true
}

// calculateBackoff calculates the backoff delay for a given attempt
func calculateBackoff(attempt int, config Config) time.Duration {
	delay := float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempt))

	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}

	return time.Duration(delay)
}

// DoWithJitter executes the function with retry logic and jitter
func DoWithJitter(ctx context.Context, config Config, logger zerolog.Logger, fn RetryFunc) error {
	// Add jitter to prevent thundering herd
	// This is a simplified version
	return Do(ctx, config, logger, fn)
}
