package entity

import (
	"time"

	"archie-core-shopify-layer/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoProjectDoc represents the MongoDB project document structure
type MongoProjectDoc struct {
	ID        primitive.ObjectID      `bson:"_id,omitempty"`
	ProjectID string                  `bson:"projectId"`
	Settings  MongoProjectSettings    `bson:"settings"`
	UpdatedAt time.Time               `bson:"updatedAt"`
}

// MongoProjectSettings represents the settings section of a project
type MongoProjectSettings struct {
	ShopifyConfigs []MongoShopifyConfigDoc `bson:"shopify_configs"`
}

// MongoShopifyConfigDoc represents a Shopify config within settings.shopify_configs[]
type MongoShopifyConfigDoc struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Env           string             `bson:"env"` // Environment name (e.g., "master")
	EncryptedKey  string             `bson:"encryptedKey"`
	APIKey        string             `bson:"apiKey"`
	WebhookSecret string             `bson:"webhookSecret,omitempty"`
	WebhookURL    string             `bson:"webhookURL"`
	CreatedAt     time.Time          `bson:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt"`
}

// ToDomain converts the MongoDB document to a domain entity
func (d *MongoShopifyConfigDoc) ToDomain(projectID, environment string) *domain.ShopifyConfig {
	return &domain.ShopifyConfig{
		ID:            d.ID.Hex(),
		ProjectID:     projectID,
		Environment:   environment,
		EncryptedKey:  d.EncryptedKey,
		APIKey:        d.APIKey,
		WebhookSecret: d.WebhookSecret,
		WebhookURL:    d.WebhookURL,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}

// MongoShopifyConfigDocFromDomain converts a domain entity to a MongoDB document
func MongoShopifyConfigDocFromDomain(config *domain.ShopifyConfig) *MongoShopifyConfigDoc {
	doc := &MongoShopifyConfigDoc{
		Env:           config.Environment,
		EncryptedKey:  config.EncryptedKey,
		APIKey:        config.APIKey,
		WebhookSecret: config.WebhookSecret,
		WebhookURL:    config.WebhookURL,
		CreatedAt:     config.CreatedAt,
		UpdatedAt:     config.UpdatedAt,
	}

	if config.ID != "" {
		if objID, err := primitive.ObjectIDFromHex(config.ID); err == nil {
			doc.ID = objID
		}
	}

	return doc
}

