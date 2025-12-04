package ports

import "context"

// ShopifyClientPool defines the interface for managing Shopify client instances
type ShopifyClientPool interface {
	GetClient(ctx context.Context, tenantID, apiKey, apiSecret string) (ShopifyClient, error)
	InvalidateClient(tenantID string)
}

