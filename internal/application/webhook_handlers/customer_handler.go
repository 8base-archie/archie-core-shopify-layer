package webhook_handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"archie-core-shopify-layer/internal/domain"
	"github.com/rs/zerolog"
)

// CustomerHandler handles customer-related webhook events
type CustomerHandler struct {
	logger zerolog.Logger
}

// NewCustomerHandler creates a new customer webhook handler
func NewCustomerHandler(logger zerolog.Logger) *CustomerHandler {
	return &CustomerHandler{
		logger: logger,
	}
}

// CanHandle returns true if this handler can process the given topic
func (h *CustomerHandler) CanHandle(topic string) bool {
	return topic == "customers/create" ||
		topic == "customers/update" ||
		topic == "customers/delete" ||
		topic == "customers/enable" ||
		topic == "customers/disable"
}

// Handle processes a customer webhook event
func (h *CustomerHandler) Handle(ctx context.Context, event *domain.WebhookEvent) error {
	// Parse customer from payload
	var customerData map[string]interface{}
	if err := json.Unmarshal(event.Payload, &customerData); err != nil {
		return fmt.Errorf("failed to parse customer webhook payload: %w", err)
	}

	h.logger.Info().
		Str("topic", event.Topic).
		Str("shop", event.Shop).
		Interface("customer", customerData).
		Msg("Processing customer webhook event")

	// TODO: Implement customer processing logic
	// - Update customer database
	// - Sync with CRM
	// - Send welcome emails
	// - etc.

	return nil
}

