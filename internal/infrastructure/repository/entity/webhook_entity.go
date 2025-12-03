package entity

import (
	"time"

	"archie-core-shopify-layer/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoWebhookDoc represents a webhook event in MongoDB
type MongoWebhookDoc struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Topic     string             `bson:"topic"`
	Shop      string             `bson:"shop"`
	Payload   []byte             `bson:"payload"`
	Verified  bool               `bson:"verified"`
	CreatedAt time.Time          `bson:"createdAt"`
}

// ToDomain converts the MongoDB document to a domain entity
func (d *MongoWebhookDoc) ToDomain() *domain.WebhookEvent {
	return &domain.WebhookEvent{
		ID:        d.ID.Hex(),
		Topic:     d.Topic,
		Shop:      d.Shop,
		Payload:   d.Payload,
		Verified:  d.Verified,
		CreatedAt: d.CreatedAt,
	}
}

// MongoWebhookDocFromDomain converts a domain entity to a MongoDB document
func MongoWebhookDocFromDomain(event *domain.WebhookEvent) *MongoWebhookDoc {
	doc := &MongoWebhookDoc{
		Topic:     event.Topic,
		Shop:      event.Shop,
		Payload:   event.Payload,
		Verified:  event.Verified,
		CreatedAt: event.CreatedAt,
	}

	if event.ID != "" {
		if objID, err := primitive.ObjectIDFromHex(event.ID); err == nil {
			doc.ID = objID
		}
	}

	return doc
}
