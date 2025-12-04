package webhook_handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"archie-core-shopify-layer/internal/domain"
	"github.com/rs/zerolog"
)

// OrderHandler handles order-related webhook events
type OrderHandler struct {
	logger zerolog.Logger
}

// NewOrderHandler creates a new order webhook handler
func NewOrderHandler(logger zerolog.Logger) *OrderHandler {
	return &OrderHandler{
		logger: logger,
	}
}

// CanHandle returns true if this handler can process the given topic
func (h *OrderHandler) CanHandle(topic string) bool {
	return topic == "orders/create" ||
		topic == "orders/updated" ||
		topic == "orders/cancelled" ||
		topic == "orders/paid" ||
		topic == "orders/fulfilled" ||
		topic == "orders/partially_fulfilled"
}

// Handle processes an order webhook event
func (h *OrderHandler) Handle(ctx context.Context, event *domain.WebhookEvent) error {
	// Parse order from payload
	var orderData map[string]interface{}
	if err := json.Unmarshal(event.Payload, &orderData); err != nil {
		return fmt.Errorf("failed to parse order webhook payload: %w", err)
	}

	h.logger.Info().
		Str("topic", event.Topic).
		Str("shop", event.Shop).
		Interface("order", orderData).
		Msg("Processing order webhook event")

	// TODO: Implement order processing logic
	// - Update order status in database
	// - Trigger business logic
	// - Send notifications
	// - etc.

	return nil
}

