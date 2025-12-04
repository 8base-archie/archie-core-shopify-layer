package shopify

import (
	"context"
	"fmt"
	"sync"

	"archie-core-shopify-layer/internal/ports"
	"github.com/rs/zerolog"
)

// ClientPool implements ports.ShopifyClientPool using in-memory caching
type ClientPool struct {
	clients     sync.Map // map[string]ports.ShopifyClient
	mu          sync.Mutex
	onceMap     sync.Map // map[string]*sync.Once - to guarantee unique creation per tenant
	logger      zerolog.Logger
	rateLimiter *RateLimiter
	retryConfig RetryConfig
}

// NewClientPool creates a new Shopify client pool
func NewClientPool(logger zerolog.Logger) ports.ShopifyClientPool {
	return NewClientPoolWithOptions(logger, nil, DefaultRetryConfig())
}

// NewClientPoolWithOptions creates a client pool with rate limiting and retry options
func NewClientPoolWithOptions(logger zerolog.Logger, rateLimiter *RateLimiter, retryConfig RetryConfig) ports.ShopifyClientPool {
	return &ClientPool{
		logger:      logger,
		rateLimiter: rateLimiter,
		retryConfig: retryConfig,
	}
}

// GetClient retrieves or creates a Shopify client for a tenant
// Thread-safe: guarantees that only one client is created per tenantID
func (p *ClientPool) GetClient(ctx context.Context, tenantID, apiKey, apiSecret string) (ports.ShopifyClient, error) {
	// Check cache first (fast path)
	if cached, ok := p.clients.Load(tenantID); ok {
		if client, ok := cached.(ports.ShopifyClient); ok {
			return client, nil
		}
	}

	// Get or create sync.Once for this tenantID
	onceInterface, _ := p.onceMap.LoadOrStore(tenantID, &sync.Once{})
	once := onceInterface.(*sync.Once)

	var client ports.ShopifyClient
	var err error

	// Guarantee that only one goroutine creates the client for this tenantID
	once.Do(func() {
		// Double-check: another goroutine may have created the client while we were waiting
		if cached, ok := p.clients.Load(tenantID); ok {
			if existingClient, ok := cached.(ports.ShopifyClient); ok {
				client = existingClient
				return
			}
		}

	// Create new client with rate limiting and retry options
	newClient := NewClientWithOptions(apiKey, apiSecret, p.rateLimiter, p.retryConfig, p.logger)
		if newClient == nil {
			err = fmt.Errorf("failed to create Shopify client")
			// Clean up the once if it fails to allow retry
			p.onceMap.Delete(tenantID)
			return
		}

		// Cache the client
		p.clients.Store(tenantID, newClient)
		client = newClient

		p.logger.Info().
			Str("tenantId", tenantID).
			Msg("Created new Shopify client for tenant")
	})

	if err != nil {
		return nil, err
	}

	// If client is nil, it means another goroutine created it, load it from cache
	if client == nil {
		if cached, ok := p.clients.Load(tenantID); ok {
			if existingClient, ok := cached.(ports.ShopifyClient); ok {
				return existingClient, nil
			}
		}
		return nil, fmt.Errorf("failed to get or create client for tenant %s", tenantID)
	}

	return client, nil
}

// InvalidateClient removes a client from the cache
// Thread-safe: cleans both the client and the associated sync.Once
func (p *ClientPool) InvalidateClient(tenantID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.clients.Delete(tenantID)
	p.onceMap.Delete(tenantID) // Allow recreation if needed

	p.logger.Info().
		Str("tenantId", tenantID).
		Msg("Invalidated Shopify client for tenant")
}

