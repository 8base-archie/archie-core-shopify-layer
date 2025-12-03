package domain

import (
	"time"
)

// Shop represents a Shopify store tenant
type Shop struct {
	ID          string    `json:"id" bson:"_id"`
	Domain      string    `json:"domain" bson:"domain"`
	AccessToken string    `json:"-" bson:"access_token"` // Encrypted
	Scopes      []string  `json:"scopes" bson:"scopes"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// Token represents an OAuth token
type Token struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

// WebhookEvent represents a received webhook
type WebhookEvent struct {
	ID        string    `json:"id" bson:"_id"`
	Topic     string    `json:"topic" bson:"topic"`
	Shop      string    `json:"shop" bson:"shop"`
	Payload   []byte    `json:"payload" bson:"payload"`
	Verified  bool      `json:"verified" bson:"verified"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
