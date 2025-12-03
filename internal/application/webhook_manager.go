package application

import (
	"context"
	"fmt"

	"archie-core-shopify-layer/internal/ports"

	"github.com/rs/zerolog"
)

// WebhookTopic represents common Shopify webhook topics
type WebhookTopic string

const (
	TopicOrdersCreate    WebhookTopic = "orders/create"
	TopicOrdersUpdated   WebhookTopic = "orders/updated"
	TopicOrdersCancelled WebhookTopic = "orders/cancelled"
	TopicProductsCreate  WebhookTopic = "products/create"
	TopicProductsUpdate  WebhookTopic = "products/update"
	TopicProductsDelete  WebhookTopic = "products/delete"
	TopicCustomersCreate WebhookTopic = "customers/create"
	TopicCustomersUpdate WebhookTopic = "customers/update"
	TopicCustomersDelete WebhookTopic = "customers/delete"
	TopicAppUninstalled  WebhookTopic = "app/uninstalled"
)

// WebhookManager manages webhook subscriptions
type WebhookManager struct {
	client     ports.ShopifyClient
	logger     zerolog.Logger
	webhookURL string
}

// NewWebhookManager creates a new webhook manager
func NewWebhookManager(
	client ports.ShopifyClient,
	logger zerolog.Logger,
	webhookURL string,
) *WebhookManager {
	return &WebhookManager{
		client:     client,
		logger:     logger,
		webhookURL: webhookURL,
	}
}

// SubscribeToWebhooks subscribes to common webhook topics
func (m *WebhookManager) SubscribeToWebhooks(ctx context.Context, shop string, accessToken string, topics []WebhookTopic) error {
	for _, topic := range topics {
		if err := m.createWebhook(ctx, shop, accessToken, string(topic)); err != nil {
			m.logger.Error().
				Err(err).
				Str("shop", shop).
				Str("topic", string(topic)).
				Msg("Failed to create webhook")
			return fmt.Errorf("failed to create webhook for topic %s: %w", topic, err)
		}

		m.logger.Info().
			Str("shop", shop).
			Str("topic", string(topic)).
			Msg("Webhook subscription created")
	}

	return nil
}

// createWebhook creates a single webhook subscription
func (m *WebhookManager) createWebhook(ctx context.Context, shop string, accessToken string, topic string) error {
	// Note: This is a placeholder implementation
	// The actual implementation would use the Shopify Admin API to create webhooks
	// For now, we'll log the intent
	m.logger.Info().
		Str("shop", shop).
		Str("topic", topic).
		Str("address", m.webhookURL).
		Msg("Creating webhook subscription")

	// TODO: Implement actual webhook creation using Shopify Admin API
	// This would require extending the ShopifyClient interface or using the goshopify client directly

	return nil
}

// GetDefaultTopics returns the default webhook topics to subscribe to
func (m *WebhookManager) GetDefaultTopics() []WebhookTopic {
	return []WebhookTopic{
		TopicOrdersCreate,
		TopicOrdersUpdated,
		TopicProductsCreate,
		TopicProductsUpdate,
		TopicAppUninstalled,
	}
}
