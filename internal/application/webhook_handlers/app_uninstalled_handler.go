package webhook_handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"archie-core-shopify-layer/internal/domain"
	"github.com/rs/zerolog"
)

// AppUninstalledHandler handles app uninstalled webhook events
type AppUninstalledHandler struct {
	logger zerolog.Logger
}

// NewAppUninstalledHandler creates a new app uninstalled webhook handler
func NewAppUninstalledHandler(logger zerolog.Logger) *AppUninstalledHandler {
	return &AppUninstalledHandler{
		logger: logger,
	}
}

// CanHandle returns true if this handler can process the given topic
func (h *AppUninstalledHandler) CanHandle(topic string) bool {
	return topic == "app/uninstalled"
}

// Handle processes an app uninstalled webhook event
func (h *AppUninstalledHandler) Handle(ctx context.Context, event *domain.WebhookEvent) error {
	// Parse shop data from payload
	var shopData map[string]interface{}
	if err := json.Unmarshal(event.Payload, &shopData); err != nil {
		return fmt.Errorf("failed to parse app uninstalled webhook payload: %w", err)
	}

	h.logger.Info().
		Str("topic", event.Topic).
		Str("shop", event.Shop).
		Interface("shop", shopData).
		Msg("Processing app uninstalled webhook event")

	// TODO: Implement app uninstalled logic
	// - Clean up shop data
	// - Delete webhook subscriptions
	// - Revoke access tokens
	// - Send notification
	// - etc.

	return nil
}

