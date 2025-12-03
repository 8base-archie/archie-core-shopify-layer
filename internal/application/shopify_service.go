package application

import (
	"context"
	"fmt"

	"archie-core-shopify-layer/internal/domain"
	"archie-core-shopify-layer/internal/ports"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/rs/zerolog"
)

// ShopifyService implements the application business logic
type ShopifyService struct {
	repository ports.Repository
	client     ports.ShopifyClient
	logger     zerolog.Logger
	apiKey     string
	apiSecret  string
}

// NewShopifyService creates a new Shopify service
func NewShopifyService(
	repository ports.Repository,
	client ports.ShopifyClient,
	logger zerolog.Logger,
	apiKey string,
	apiSecret string,
) *ShopifyService {
	return &ShopifyService{
		repository: repository,
		client:     client,
		logger:     logger,
		apiKey:     apiKey,
		apiSecret:  apiSecret,
	}
}

// GenerateAuthURL generates the OAuth authorization URL
func (s *ShopifyService) GenerateAuthURL(ctx context.Context, shop string, scopes []string) (string, error) {
	authURL, err := s.client.GenerateAuthURL(shop, scopes)
	if err != nil {
		s.logger.Error().Err(err).Str("shop", shop).Msg("Failed to generate auth URL")
		return "", fmt.Errorf("failed to generate auth URL: %w", err)
	}

	return authURL, nil
}

// ExchangeToken exchanges the authorization code for an access token
func (s *ShopifyService) ExchangeToken(ctx context.Context, shop string, code string) (*domain.Shop, error) {
	// Exchange code for access token
	accessToken, err := s.client.ExchangeToken(ctx, shop, code)
	if err != nil {
		s.logger.Error().Err(err).Str("shop", shop).Msg("Failed to exchange token")
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get shop information
	shopInfo, err := s.client.GetShop(ctx, shop, accessToken)
	if err != nil {
		s.logger.Error().Err(err).Str("shop", shop).Msg("Failed to get shop info")
		return nil, fmt.Errorf("failed to get shop info: %w", err)
	}

	// Create domain shop entity
	domainShop := &domain.Shop{
		Domain:      shopInfo.Domain,
		AccessToken: accessToken,
		Scopes:      []string{}, // TODO: Extract scopes from response
	}

	// Save shop to repository
	if err := s.repository.SaveShop(ctx, domainShop); err != nil {
		s.logger.Error().Err(err).Str("shop", shop).Msg("Failed to save shop")
		return nil, fmt.Errorf("failed to save shop: %w", err)
	}

	return domainShop, nil
}

// GetShop retrieves shop information
func (s *ShopifyService) GetShop(ctx context.Context, domain string) (*domain.Shop, error) {
	shop, err := s.repository.GetShop(ctx, domain)
	if err != nil {
		s.logger.Error().Err(err).Str("domain", domain).Msg("Failed to get shop")
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	return shop, nil
}

// GetProducts retrieves products for a shop
func (s *ShopifyService) GetProducts(ctx context.Context, domain string) ([]goshopify.Product, error) {
	// Get shop from repository
	shop, err := s.repository.GetShop(ctx, domain)
	if err != nil {
		s.logger.Error().Err(err).Str("domain", domain).Msg("Failed to get shop")
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	if shop == nil {
		return nil, fmt.Errorf("shop not found: %s", domain)
	}

	// Get products from Shopify API
	products, err := s.client.GetProducts(ctx, domain, shop.AccessToken, nil)
	if err != nil {
		s.logger.Error().Err(err).Str("domain", domain).Msg("Failed to get products")
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}

// ProcessWebhook processes a Shopify webhook event
func (s *ShopifyService) ProcessWebhook(ctx context.Context, topic string, shop string, payload []byte, verified bool) error {
	// Create webhook event
	event := &domain.WebhookEvent{
		Topic:    topic,
		Shop:     shop,
		Payload:  payload,
		Verified: verified,
	}

	// Log webhook to repository
	if err := s.repository.LogWebhook(ctx, event); err != nil {
		s.logger.Error().Err(err).Str("topic", topic).Str("shop", shop).Msg("Failed to log webhook")
		return fmt.Errorf("failed to log webhook: %w", err)
	}

	s.logger.Info().Str("topic", topic).Str("shop", shop).Bool("verified", verified).Msg("Webhook processed")
	return nil
}
