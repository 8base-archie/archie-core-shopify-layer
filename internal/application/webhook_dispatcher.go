package application

import (
	"context"
	"fmt"
	"sync"

	"archie-core-shopify-layer/internal/domain"
	"github.com/rs/zerolog"
)

// WebhookDispatcher dispatches webhook events to registered handlers
type WebhookDispatcher struct {
	handlers []domain.WebhookHandler
	logger   zerolog.Logger
	mu       sync.RWMutex
}

// NewWebhookDispatcher creates a new webhook dispatcher
func NewWebhookDispatcher(logger zerolog.Logger) *WebhookDispatcher {
	return &WebhookDispatcher{
		handlers: make([]domain.WebhookHandler, 0),
		logger:   logger,
	}
}

// RegisterHandler registers a webhook handler
func (d *WebhookDispatcher) RegisterHandler(handler domain.WebhookHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers = append(d.handlers, handler)
	d.logger.Info().
		Msg("Webhook handler registered")
}

// Dispatch dispatches a webhook event to appropriate handlers
func (d *WebhookDispatcher) Dispatch(ctx context.Context, event *domain.WebhookEvent) error {
	d.mu.RLock()
	handlers := make([]domain.WebhookHandler, len(d.handlers))
	copy(handlers, d.handlers)
	d.mu.RUnlock()

	handled := false
	for _, handler := range handlers {
		if handler.CanHandle(event.Topic) {
			handled = true
			if err := handler.Handle(ctx, event); err != nil {
				d.logger.Error().
					Err(err).
					Str("topic", event.Topic).
					Str("handler", fmt.Sprintf("%T", handler)).
					Msg("Webhook handler failed")
				// Continue to other handlers even if one fails
				continue
			}
			d.logger.Info().
				Str("topic", event.Topic).
				Str("handler", fmt.Sprintf("%T", handler)).
				Msg("Webhook event handled successfully")
		}
	}

	if !handled {
		d.logger.Warn().
			Str("topic", event.Topic).
			Msg("No handler found for webhook topic")
	}

	return nil
}

