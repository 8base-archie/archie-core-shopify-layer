package entity

import (
	"time"

	"archie-core-shopify-layer/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoShopDoc represents a Shopify store in MongoDB
type MongoShopDoc struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Domain      string             `bson:"domain"`
	AccessToken string             `bson:"accessToken"` // Encrypted
	Scopes      []string           `bson:"scopes"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

// ToDomain converts the MongoDB document to a domain entity
func (d *MongoShopDoc) ToDomain() *domain.Shop {
	return &domain.Shop{
		ID:          d.ID.Hex(),
		Domain:      d.Domain,
		AccessToken: d.AccessToken,
		Scopes:      d.Scopes,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

// MongoShopDocFromDomain converts a domain entity to a MongoDB document
func MongoShopDocFromDomain(shop *domain.Shop) *MongoShopDoc {
	doc := &MongoShopDoc{
		Domain:      shop.Domain,
		AccessToken: shop.AccessToken,
		Scopes:      shop.Scopes,
		CreatedAt:   shop.CreatedAt,
		UpdatedAt:   shop.UpdatedAt,
	}

	if shop.ID != "" {
		if objID, err := primitive.ObjectIDFromHex(shop.ID); err == nil {
			doc.ID = objID
		}
	}

	return doc
}
