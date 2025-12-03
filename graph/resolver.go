package graph

import "archie-core-shopify-layer/internal/application"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	shopifyService     *application.ShopifyService
	credentialsService *application.CredentialsService
}

// NewResolver creates a new GraphQL resolver
func NewResolver(shopifyService *application.ShopifyService, credentialsService *application.CredentialsService) *Resolver {
	return &Resolver{
		shopifyService:     shopifyService,
		credentialsService: credentialsService,
	}
}
