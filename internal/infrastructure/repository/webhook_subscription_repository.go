package repository

import (
	"context"
	"fmt"
	"time"

	"archie-core-shopify-layer/internal/domain"
	"archie-core-shopify-layer/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoWebhookSubscriptionRepository implements WebhookSubscriptionRepository using MongoDB
type MongoWebhookSubscriptionRepository struct {
	collection *mongo.Collection
}

// NewMongoWebhookSubscriptionRepository creates a new webhook subscription repository
func NewMongoWebhookSubscriptionRepository(db *mongo.Database) ports.WebhookSubscriptionRepository {
	return &MongoWebhookSubscriptionRepository{
		collection: db.Collection("webhook_subscriptions"),
	}
}

// SaveWebhookSubscription saves or updates a webhook subscription
func (r *MongoWebhookSubscriptionRepository) SaveWebhookSubscription(ctx context.Context, subscription *domain.WebhookSubscription) error {
	doc := bson.M{
		"projectId":   subscription.ProjectID,
		"environment": subscription.Environment,
		"shopDomain":  subscription.ShopDomain,
		"webhookId":   subscription.WebhookID,
		"topic":       subscription.Topic,
		"address":     subscription.Address,
		"updatedAt":   time.Now(),
	}

	if subscription.ID == "" {
		doc["_id"] = primitive.NewObjectID()
		doc["createdAt"] = time.Now()
		_, err := r.collection.InsertOne(ctx, doc)
		if err != nil {
			return fmt.Errorf("failed to save webhook subscription: %w", err)
		}
		subscription.ID = doc["_id"].(primitive.ObjectID).Hex()
	} else {
		objID, err := primitive.ObjectIDFromHex(subscription.ID)
		if err != nil {
			return fmt.Errorf("invalid subscription ID: %w", err)
		}
		filter := bson.M{"_id": objID}
		update := bson.M{"$set": doc}
		_, err = r.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return fmt.Errorf("failed to update webhook subscription: %w", err)
		}
	}

	return nil
}

// GetWebhookSubscription retrieves a webhook subscription
func (r *MongoWebhookSubscriptionRepository) GetWebhookSubscription(ctx context.Context, projectID string, environment string, shopDomain string, topic string) (*domain.WebhookSubscription, error) {
	var doc struct {
		ID          primitive.ObjectID `bson:"_id"`
		ProjectID   string             `bson:"projectId"`
		Environment string             `bson:"environment"`
		ShopDomain  string             `bson:"shopDomain"`
		WebhookID   int64              `bson:"webhookId"`
		Topic       string             `bson:"topic"`
		Address     string             `bson:"address"`
		CreatedAt   time.Time          `bson:"createdAt"`
		UpdatedAt   time.Time          `bson:"updatedAt"`
	}

	filter := bson.M{
		"projectId":   projectID,
		"environment": environment,
		"shopDomain":  shopDomain,
		"topic":       topic,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook subscription: %w", err)
	}

	return &domain.WebhookSubscription{
		ID:          doc.ID.Hex(),
		ProjectID:   doc.ProjectID,
		Environment: doc.Environment,
		ShopDomain:  doc.ShopDomain,
		WebhookID:   doc.WebhookID,
		Topic:       doc.Topic,
		Address:     doc.Address,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}, nil
}

// ListWebhookSubscriptions lists webhook subscriptions for a shop
func (r *MongoWebhookSubscriptionRepository) ListWebhookSubscriptions(ctx context.Context, projectID string, environment string, shopDomain string) ([]*domain.WebhookSubscription, error) {
	filter := bson.M{
		"projectId":   projectID,
		"environment": environment,
		"shopDomain":  shopDomain,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhook subscriptions: %w", err)
	}
	defer cursor.Close(ctx)

	var subscriptions []*domain.WebhookSubscription
	for cursor.Next(ctx) {
		var doc struct {
			ID          primitive.ObjectID `bson:"_id"`
			ProjectID   string             `bson:"projectId"`
			Environment string             `bson:"environment"`
			ShopDomain  string             `bson:"shopDomain"`
			WebhookID   int64              `bson:"webhookId"`
			Topic       string             `bson:"topic"`
			Address     string             `bson:"address"`
			CreatedAt   time.Time          `bson:"createdAt"`
			UpdatedAt   time.Time          `bson:"updatedAt"`
		}

		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		subscriptions = append(subscriptions, &domain.WebhookSubscription{
			ID:          doc.ID.Hex(),
			ProjectID:   doc.ProjectID,
			Environment: doc.Environment,
			ShopDomain:  doc.ShopDomain,
			WebhookID:   doc.WebhookID,
			Topic:       doc.Topic,
			Address:     doc.Address,
			CreatedAt:   doc.CreatedAt,
			UpdatedAt:   doc.UpdatedAt,
		})
	}

	return subscriptions, nil
}

// DeleteWebhookSubscription deletes a webhook subscription
func (r *MongoWebhookSubscriptionRepository) DeleteWebhookSubscription(ctx context.Context, subscriptionID string) error {
	objID, err := primitive.ObjectIDFromHex(subscriptionID)
	if err != nil {
		return fmt.Errorf("invalid subscription ID: %w", err)
	}

	filter := bson.M{"_id": objID}
	_, err = r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete webhook subscription: %w", err)
	}

	return nil
}

