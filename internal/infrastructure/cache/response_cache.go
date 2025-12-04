package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// CacheEntry represents a cached response
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
	Key       string
}

// ResponseCache provides caching for API responses
type ResponseCache struct {
	mu        sync.RWMutex
	entries   map[string]*CacheEntry
	logger    zerolog.Logger
	defaultTTL time.Duration
	maxSize   int
}

// NewResponseCache creates a new response cache
func NewResponseCache(logger zerolog.Logger, defaultTTL time.Duration, maxSize int) *ResponseCache {
	cache := &ResponseCache{
		entries:    make(map[string]*CacheEntry),
		logger:     logger,
		defaultTTL: defaultTTL,
		maxSize:    maxSize,
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a cached value
func (c *ResponseCache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.entries, key)
		c.mu.Unlock()
		c.mu.RLock()
		return nil, false
	}

	return entry.Data, true
}

// Set stores a value in the cache
func (c *ResponseCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = c.defaultTTL
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Check max size
	if len(c.entries) >= c.maxSize {
		// Evict oldest entry (simple FIFO)
		c.evictOldest()
	}

	entry := &CacheEntry{
		Data:      value,
		ExpiresAt: time.Now().Add(ttl),
		Key:       key,
	}

	c.entries[key] = entry
	return nil
}

// Delete removes a key from cache
func (c *ResponseCache) Delete(ctx context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Clear removes all entries
func (c *ResponseCache) Clear(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*CacheEntry)
}

// GenerateKey generates a cache key from components
func GenerateKey(components ...string) string {
	combined := ""
	for _, comp := range components {
		combined += comp + ":"
	}
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

// evictOldest removes the oldest entry (simple implementation)
func (c *ResponseCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, entry := range c.entries {
		if first || entry.ExpiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.ExpiresAt
			first = false
		}
	}

	if oldestKey != "" {
		delete(c.entries, oldestKey)
		c.logger.Debug().
			Str("key", oldestKey).
			Msg("Evicted cache entry")
	}
}

// cleanup periodically removes expired entries
func (c *ResponseCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		expired := make([]string, 0)

		for key, entry := range c.entries {
			if now.After(entry.ExpiresAt) {
				expired = append(expired, key)
			}
		}

		for _, key := range expired {
			delete(c.entries, key)
		}
		c.mu.Unlock()

		if len(expired) > 0 {
			c.logger.Debug().
				Int("count", len(expired)).
				Msg("Cleaned up expired cache entries")
		}
	}
}

// GetStats returns cache statistics
func (c *ResponseCache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"size":      len(c.entries),
		"max_size":  c.maxSize,
		"default_ttl": c.defaultTTL.String(),
	}
}

