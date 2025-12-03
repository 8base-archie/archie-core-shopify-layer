package domain

import "time"

// ShopifyCredentials represents Shopify API credentials
type ShopifyCredentials struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	ProjectID   string    `json:"project_id" bson:"project_id"`
	Environment string    `json:"environment" bson:"environment"`
	APIKey      string    `json:"api_key" bson:"api_key"`
	APISecret   string    `json:"api_secret" bson:"api_secret"` // Encrypted
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}
