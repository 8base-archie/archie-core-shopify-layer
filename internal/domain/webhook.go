package domain

import "time"

// Webhook represents a Shopify webhook subscription
type Webhook struct {
	ID        int64     `json:"id" bson:"id"`
	Address   string    `json:"address" bson:"address"`
	Topic     string    `json:"topic" bson:"topic"`
	Format    string    `json:"format" bson:"format"` // json or xml
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// WebhookSubscription represents a webhook subscription stored in our system
type WebhookSubscription struct {
	ID          string    `json:"id" bson:"_id"`
	ProjectID   string    `json:"project_id" bson:"project_id"`
	Environment string    `json:"environment" bson:"environment"`
	ShopDomain  string    `json:"shop_domain" bson:"shop_domain"`
	WebhookID   int64     `json:"webhook_id" bson:"webhook_id"` // Shopify webhook ID
	Topic       string    `json:"topic" bson:"topic"`
	Address     string    `json:"address" bson:"address"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

