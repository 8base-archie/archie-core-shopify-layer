package webhook_handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"archie-core-shopify-layer/internal/domain"
	"github.com/rs/zerolog"
)

// ProductHandler handles product-related webhook events
type ProductHandler struct {
	logger zerolog.Logger
}

// NewProductHandler creates a new product webhook handler
func NewProductHandler(logger zerolog.Logger) *ProductHandler {
	return &ProductHandler{
		logger: logger,
	}
}

// CanHandle returns true if this handler can process the given topic
func (h *ProductHandler) CanHandle(topic string) bool {
	return topic == "products/create" ||
		topic == "products/update" ||
		topic == "products/delete"
}

// Handle processes a product webhook event
func (h *ProductHandler) Handle(ctx context.Context, event *domain.WebhookEvent) error {
	// Parse product from payload
	var productData map[string]interface{}
	if err := json.Unmarshal(event.Payload, &productData); err != nil {
		return fmt.Errorf("failed to parse product webhook payload: %w", err)
	}

	h.logger.Info().
		Str("topic", event.Topic).
		Str("shop", event.Shop).
		Interface("product", productData).
		Msg("Processing product webhook event")

	// TODO: Implement product processing logic
	// - Update product cache
	// - Sync with external systems
	// - Update search index
	// - etc.

	return nil
}

