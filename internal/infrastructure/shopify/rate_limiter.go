package shopify

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// RateLimiter handles Shopify API rate limiting
// Shopify allows 40 requests per app per store per minute (leaky bucket)
type RateLimiter struct {
	mu          sync.RWMutex
	buckets     map[string]*rateBucket // key: shopDomain
	logger      zerolog.Logger
	maxRequests int
	window      time.Duration
}

type rateBucket struct {
	mu           sync.Mutex
	tokens       int
	lastRefill   time.Time
	maxTokens    int
	refillRate   int // tokens per second
	window       time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(logger zerolog.Logger) *RateLimiter {
	// Shopify allows 40 requests per minute per store
	// We'll use a conservative 35 requests per minute to leave buffer
	maxRequests := 35
	window := time.Minute

	return &RateLimiter{
		buckets:     make(map[string]*rateBucket),
		logger:      logger,
		maxRequests: maxRequests,
		window:      window,
	}
}

// Wait waits until a request can be made for the given shop
func (rl *RateLimiter) Wait(ctx context.Context, shopDomain string) error {
	bucket := rl.getBucket(shopDomain)
	return bucket.wait(ctx)
}

// UpdateFromResponse updates rate limit state from HTTP response headers
func (rl *RateLimiter) UpdateFromResponse(shopDomain string, resp *http.Response) {
	bucket := rl.getBucket(shopDomain)
	bucket.updateFromHeaders(resp.Header)
}

// getBucket gets or creates a rate bucket for a shop
func (rl *RateLimiter) getBucket(shopDomain string) *rateBucket {
	rl.mu.RLock()
	bucket, exists := rl.buckets[shopDomain]
	rl.mu.RUnlock()

	if exists {
		return bucket
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double-check after acquiring write lock
	if bucket, exists := rl.buckets[shopDomain]; exists {
		return bucket
	}

	// Create new bucket
	bucket = &rateBucket{
		tokens:     rl.maxRequests,
		lastRefill: time.Now(),
		maxTokens:  rl.maxRequests,
		refillRate: rl.maxRequests / int(rl.window.Seconds()),
		window:     rl.window,
	}
	rl.buckets[shopDomain] = bucket

	return bucket
}

// wait waits until a token is available
func (rb *rateBucket) wait(ctx context.Context) error {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	// Refill tokens based on elapsed time
	rb.refill()

	// If tokens available, consume one
	if rb.tokens > 0 {
		rb.tokens--
		return nil
	}

	// Calculate wait time until next token is available
	waitTime := time.Duration(float64(time.Second) / float64(rb.refillRate))
	if waitTime < 100*time.Millisecond {
		waitTime = 100 * time.Millisecond // Minimum wait
	}

	// Wait with context cancellation support
	timer := time.NewTimer(waitTime)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		rb.refill()
		if rb.tokens > 0 {
			rb.tokens--
			return nil
		}
		return fmt.Errorf("rate limit exceeded, please retry")
	}
}

// refill refills tokens based on elapsed time
func (rb *rateBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(rb.lastRefill)

	if elapsed <= 0 {
		return
	}

	// Calculate tokens to add
	tokensToAdd := int(elapsed.Seconds()) * rb.refillRate
	if tokensToAdd > 0 {
		rb.tokens = min(rb.tokens+tokensToAdd, rb.maxTokens)
		rb.lastRefill = now
	}
}

// updateFromHeaders updates rate limit state from Shopify response headers
func (rb *rateBucket) updateFromHeaders(headers http.Header) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	// Shopify provides rate limit info in headers:
	// X-Shopify-Shop-Api-Call-Limit: 40/40
	// X-Shopify-API-Version: 2024-01

	if limitHeader := headers.Get("X-Shopify-Shop-Api-Call-Limit"); limitHeader != "" {
		// Parse "40/40" format
		var used, limit int
		if _, err := fmt.Sscanf(limitHeader, "%d/%d", &used, &limit); err == nil {
			rb.tokens = max(0, limit-used)
			rb.maxTokens = limit
		}
	}

	// Update refill time if we're at limit
	if rb.tokens == 0 {
		rb.lastRefill = time.Now()
	}
}

// GetRemainingTokens returns remaining tokens for a shop (for monitoring)
func (rl *RateLimiter) GetRemainingTokens(shopDomain string) int {
	bucket := rl.getBucket(shopDomain)
	bucket.mu.Lock()
	defer bucket.mu.Unlock()
	bucket.refill()
	return bucket.tokens
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

