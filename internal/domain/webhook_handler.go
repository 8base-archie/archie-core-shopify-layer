package domain

import "context"

// WebhookHandler defines the interface for webhook event handlers
type WebhookHandler interface {
	// Handle processes a webhook event
	Handle(ctx context.Context, event *WebhookEvent) error
	
	// CanHandle returns true if this handler can process the given topic
	CanHandle(topic string) bool
}

// WebhookEvent represents a received webhook event
// This is already defined in entities.go but we'll keep it here for reference
// type WebhookEvent struct {
// 	ID        string    `json:"id" bson:"_id"`
// 	Topic     string    `json:"topic" bson:"topic"`
// 	Shop      string    `json:"shop" bson:"shop"`
// 	Payload   []byte    `json:"payload" bson:"payload"`
// 	Verified  bool      `json:"verified" bson:"verified"`
// 	CreatedAt time.Time `json:"created_at" bson:"created_at"`
// }

