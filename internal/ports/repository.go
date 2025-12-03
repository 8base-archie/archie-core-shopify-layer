package ports

import (
	"context"

	"archie-core-shopify-layer/internal/domain"
)

// Repository defines the interface for persistence
type Repository interface {
	// Shop operations
	SaveShop(ctx context.Context, shop *domain.Shop) error
	GetShop(ctx context.Context, domain string) (*domain.Shop, error)

	// Webhook operations
	LogWebhook(ctx context.Context, event *domain.WebhookEvent) error

	// Credentials operations
	SaveCredentials(ctx context.Context, creds *domain.ShopifyCredentials) error
	GetCredentials(ctx context.Context, projectID string, environment string) (*domain.ShopifyCredentials, error)
	DeleteCredentials(ctx context.Context, projectID string, environment string) error
}
