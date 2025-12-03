package entity

import (
	"time"

	"archie-core-shopify-layer/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoCredentialsDoc represents credentials in MongoDB
type MongoCredentialsDoc struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ProjectID   string             `bson:"projectId"`
	Environment string             `bson:"environment"`
	APIKey      string             `bson:"apiKey"`
	APISecret   string             `bson:"apiSecret"` // Encrypted
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

// ToDomain converts MongoDB document to domain model
func (d *MongoCredentialsDoc) ToDomain() *domain.ShopifyCredentials {
	return &domain.ShopifyCredentials{
		ID:          d.ID.Hex(),
		ProjectID:   d.ProjectID,
		Environment: d.Environment,
		APIKey:      d.APIKey,
		APISecret:   d.APISecret,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

// MongoCredentialsDocFromDomain converts domain model to MongoDB document
func MongoCredentialsDocFromDomain(creds *domain.ShopifyCredentials) *MongoCredentialsDoc {
	doc := &MongoCredentialsDoc{
		ProjectID:   creds.ProjectID,
		Environment: creds.Environment,
		APIKey:      creds.APIKey,
		APISecret:   creds.APISecret,
		CreatedAt:   creds.CreatedAt,
		UpdatedAt:   creds.UpdatedAt,
	}

	if creds.ID != "" {
		if objID, err := primitive.ObjectIDFromHex(creds.ID); err == nil {
			doc.ID = objID
		}
	}

	return doc
}
