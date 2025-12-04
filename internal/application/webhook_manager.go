package application

import (
	"context"
	"fmt"

	"archie-core-shopify-layer/internal/domain"

	goshopify "github.com/bold-commerce/go-shopify/v4"
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
	shopifyService *ShopifyService
	logger         zerolog.Logger
	webhookURL     string
}

// NewWebhookManager creates a new webhook manager
func NewWebhookManager(
	shopifyService *ShopifyService,
	logger zerolog.Logger,
	webhookURL string,
) *WebhookManager {
	return &WebhookManager{
		shopifyService: shopifyService,
		logger:         logger,
		webhookURL:     webhookURL,
	}
}

// SubscribeToWebhooks subscribes to common webhook topics for a shop
func (m *WebhookManager) SubscribeToWebhooks(ctx context.Context, shopDomain string, accessToken string, topics []WebhookTopic) error {
	// Extract projectID and environment from context
	projectID := domain.GetProjectIDFromContext(ctx)
	environment := domain.GetEnvironmentFromContext(ctx)
	
	if environment == "" {
		environment = domain.DefaultEnvironment
	}

	// Build webhook URL with project/environment
	webhookAddress := fmt.Sprintf("%s/%s/%s", m.webhookURL, projectID, environment)

	for _, topic := range topics {
		if err := m.createWebhook(ctx, shopDomain, accessToken, string(topic), webhookAddress); err != nil {
			m.logger.Error().
				Err(err).
				Str("shop", shopDomain).
				Str("topic", string(topic)).
				Msg("Failed to create webhook")
			return fmt.Errorf("failed to create webhook for topic %s: %w", topic, err)
		}

		m.logger.Info().
			Str("shop", shopDomain).
			Str("topic", string(topic)).
			Msg("Webhook subscription created")
	}

	return nil
}

// createWebhook creates a single webhook subscription using Shopify Admin API
func (m *WebhookManager) createWebhook(ctx context.Context, shopDomain string, accessToken string, topic string, address string) error {
	// Get config to get API key/secret for creating Shopify client
	config, err := m.shopifyService.GetConfig(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	// Decrypt API secret
	apiSecret, err := m.shopifyService.encryptionSvc.Decrypt(config.EncryptedKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt API secret: %w", err)
	}

	// Create a client with the shop's accessToken (not the API key/secret)
	// For webhook operations, we need a client that uses the shop's accessToken
	// The go-shopify client can be created with App + shopDomain + accessToken
	goshopifyApp := goshopify.App{
		ApiKey:    config.APIKey,
		ApiSecret: apiSecret,
	}
	goshopifyClient, err := goshopify.NewClient(goshopifyApp, shopDomain, accessToken)
	if err != nil {
		return fmt.Errorf("failed to create Shopify client: %w", err)
	}

	// Create webhook
	webhook := goshopify.Webhook{
		Topic:   topic,
		Address: address,
		Format:  "json",
	}
	created, err := goshopifyClient.Webhook.Create(ctx, webhook)
	if err != nil {
		return fmt.Errorf("failed to create webhook via API: %w", err)
	}

	m.logger.Info().
		Str("shop", shopDomain).
		Str("topic", topic).
		Str("address", address).
		Uint64("webhookId", created.Id).
		Msg("Webhook subscription created successfully")

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

