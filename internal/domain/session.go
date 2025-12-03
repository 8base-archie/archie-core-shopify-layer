package domain

import "time"

// Session represents an OAuth session
type Session struct {
	ID        string    `json:"id" bson:"_id"`
	Shop      string    `json:"shop" bson:"shop"`
	State     string    `json:"state" bson:"state"`
	Scopes    []string  `json:"scopes" bson:"scopes"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
