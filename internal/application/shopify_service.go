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
// It depends on ports (interfaces) not concrete implementations
type ShopifyService struct {
	repository     ports.Repository
	configRepo     ports.ShopifyConfigRepository
	encryptionSvc  ports.EncryptionService
	clientPool     ports.ShopifyClientPool
	logger         zerolog.Logger
	webhookBaseURL string
}

// NewShopifyService creates a new Shopify application service
func NewShopifyService(
	repository ports.Repository,
	configRepo ports.ShopifyConfigRepository,
	encryptionSvc ports.EncryptionService,
	clientPool ports.ShopifyClientPool,
	logger zerolog.Logger,
	webhookBaseURL string,
) *ShopifyService {
	return &ShopifyService{
		repository:     repository,
		configRepo:     configRepo,
		encryptionSvc:  encryptionSvc,
		clientPool:     clientPool,
		logger:         logger,
		webhookBaseURL: webhookBaseURL,
	}
}

// GetClientForTenant retrieves a Shopify client for a project and environment
// tenantID is actually projectID in this context
func (s *ShopifyService) GetClientForTenant(ctx context.Context, tenantID string) (ports.ShopifyClient, error) {
	// Extract projectID and environment from context (type-safe)
	projectID := domain.GetProjectIDFromContext(ctx)
	environment := domain.GetEnvironmentFromContext(ctx)

	if projectID == "" {
		projectID = tenantID // Fallback
	}
	if environment == "" {
		environment = domain.DefaultEnvironment // Default
	}

	config, err := s.configRepo.GetByTenantID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, fmt.Errorf("shopify not configured for project %s and environment %s", projectID, environment)
	}

	// Decrypt API secret
	apiSecret, err := s.encryptionSvc.Decrypt(config.EncryptedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt API secret: %w", err)
	}

	// Get client from pool using projectID-environment as key
	return s.clientPool.GetClient(ctx, projectID+"-"+environment, config.APIKey, apiSecret)
}

// GetConfig retrieves the Shopify configuration for a project and environment
func (s *ShopifyService) GetConfig(ctx context.Context, tenantID string) (*domain.ShopifyConfig, error) {
	// Extract projectID and environment from context (type-safe)
	projectID := domain.GetProjectIDFromContext(ctx)
	environment := domain.GetEnvironmentFromContext(ctx)

	if projectID == "" {
		projectID = tenantID // Fallback
	}
	if environment == "" {
		environment = domain.DefaultEnvironment // Default
	}

	config, err := s.configRepo.GetByTenantID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, fmt.Errorf("shopify not configured for project %s and environment %s", projectID, environment)
	}

	return config, nil
}

// GenerateAuthURL generates the OAuth authorization URL
func (s *ShopifyService) GenerateAuthURL(ctx context.Context, shop string, scopes []string) (string, error) {
	// Get client for tenant
	client, err := s.GetClientForTenant(ctx, "")
	if err != nil {
		return "", fmt.Errorf("failed to get client: %w", err)
	}

	authURL, err := client.GenerateAuthURL(shop, scopes)
	if err != nil {
		s.logger.Error().Err(err).Str("shop", shop).Msg("Failed to generate auth URL")
		return "", fmt.Errorf("failed to generate auth URL: %w", err)
	}

	return authURL, nil
}

// ExchangeToken exchanges the authorization code for an access token
func (s *ShopifyService) ExchangeToken(ctx context.Context, shop string, code string) (*domain.Shop, error) {
	// Get client for tenant
	client, err := s.GetClientForTenant(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	// Exchange code for access token
	accessToken, err := client.ExchangeToken(ctx, shop, code)
	if err != nil {
		s.logger.Error().Err(err).Str("shop", shop).Msg("Failed to exchange token")
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get shop information
	shopInfo, err := client.GetShop(ctx, shop, accessToken)
	if err != nil {
		s.logger.Error().Err(err).Str("shop", shop).Msg("Failed to get shop info")
		return nil, fmt.Errorf("failed to get shop info: %w", err)
	}

	// Encrypt access token before storage
	encryptedToken, err := s.encryptionSvc.Encrypt(accessToken)
	if err != nil {
		s.logger.Error().Err(err).Str("shop", shop).Msg("Failed to encrypt access token")
		return nil, fmt.Errorf("failed to encrypt access token: %w", err)
	}

	// Create domain shop entity
	domainShop := &domain.Shop{
		Domain:      shopInfo.Domain,
		AccessToken: encryptedToken, // Store encrypted token
		Scopes:      []string{},     // TODO: Extract scopes from response
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

	// Decrypt access token
	if shop.AccessToken != "" {
		decryptedToken, err := s.encryptionSvc.Decrypt(shop.AccessToken)
		if err != nil {
			s.logger.Error().Err(err).Str("domain", domain).Msg("Failed to decrypt access token")
			return nil, fmt.Errorf("failed to decrypt access token: %w", err)
		}
		// Return shop with decrypted token (for internal use only)
		// Note: This is a copy, the domain entity still has encrypted token
		shop.AccessToken = decryptedToken
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

	// Get client for tenant
	client, err := s.GetClientForTenant(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	// Get products from Shopify API
	products, err := client.GetProducts(ctx, domain, shop.AccessToken, nil)
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
