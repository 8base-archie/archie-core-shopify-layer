package domain

import (
	"fmt"
	"time"
)

// ShopifyConfig represents the domain entity for Shopify configuration
// This is stored within a Project document in MongoDB: projects.settings.shopify_configs[]
type ShopifyConfig struct {
	ID            string
	ProjectID     string // The project ID (from X-Project-ID header)
	Environment   string // The environment name (from environment header, e.g., "master")
	EncryptedKey  string // Encrypted API secret
	APIKey        string // API key (not encrypted, public)
	WebhookSecret string // Webhook secret for verification
	WebhookURL    string // Webhook URL
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewShopifyConfig creates a new Shopify configuration with validation
func NewShopifyConfig(projectID, environment, encryptedKey, apiKey, webhookSecret, webhookURL string) (*ShopifyConfig, error) {
	// Validate projectID
	if projectID == "" {
		return nil, fmt.Errorf("projectID cannot be empty")
	}

	// Validate environment
	if environment == "" {
		environment = DefaultEnvironment
	}

	// Validate encryptedKey
	if encryptedKey == "" {
		return nil, fmt.Errorf("encryptedKey cannot be empty")
	}

	// Validate apiKey
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey cannot be empty")
	}

	now := time.Now()
	return &ShopifyConfig{
		ProjectID:     projectID,
		Environment:   environment,
		EncryptedKey:  encryptedKey,
		APIKey:        apiKey,
		WebhookSecret: webhookSecret,
		WebhookURL:    webhookURL,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// Validate validates that the configuration is valid
func (c *ShopifyConfig) Validate() error {
	if c.ProjectID == "" {
		return fmt.Errorf("projectID cannot be empty")
	}
	if c.Environment == "" {
		return fmt.Errorf("environment cannot be empty")
	}
	if c.EncryptedKey == "" {
		return fmt.Errorf("encryptedKey cannot be empty")
	}
	if c.APIKey == "" {
		return fmt.Errorf("apiKey cannot be empty")
	}
	return nil
}

// Update updates the modifiable fields of the configuration
func (c *ShopifyConfig) Update(encryptedKey, apiKey, webhookSecret, webhookURL string) error {
	if encryptedKey != "" {
		c.EncryptedKey = encryptedKey
	}
	if apiKey != "" {
		c.APIKey = apiKey
	}
	if webhookSecret != "" {
		c.WebhookSecret = webhookSecret
	}
	if webhookURL != "" {
		c.WebhookURL = webhookURL
	}
	c.UpdatedAt = time.Now()
	return c.Validate()
}

// ShopifyCredentials is kept for backward compatibility but deprecated
// Use ShopifyConfig instead
type ShopifyCredentials struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	ProjectID   string    `json:"project_id" bson:"project_id"`
	Environment string    `json:"environment" bson:"environment"`
	APIKey      string    `json:"api_key" bson:"api_key"`
	APISecret   string    `json:"api_secret" bson:"api_secret"` // Encrypted
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}
