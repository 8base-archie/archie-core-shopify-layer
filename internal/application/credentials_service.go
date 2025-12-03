package application

import (
	"context"
	"fmt"

	"archie-core-shopify-layer/internal/domain"
	"archie-core-shopify-layer/internal/infrastructure/encryption"
	"archie-core-shopify-layer/internal/ports"

	"github.com/rs/zerolog"
)

// CredentialsService handles Shopify credentials management
type CredentialsService struct {
	repository ports.Repository
	encryption *encryption.Service
	logger     zerolog.Logger
}

// NewCredentialsService creates a new credentials service
func NewCredentialsService(
	repository ports.Repository,
	encryptionService *encryption.Service,
	logger zerolog.Logger,
) *CredentialsService {
	return &CredentialsService{
		repository: repository,
		encryption: encryptionService,
		logger:     logger,
	}
}

// SaveCredentials saves Shopify API credentials
func (s *CredentialsService) SaveCredentials(ctx context.Context, projectID string, environment string, apiKey string, apiSecret string) (*domain.ShopifyCredentials, error) {
	// Encrypt API secret
	encryptedSecret, err := s.encryption.Encrypt(apiSecret)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to encrypt API secret")
		return nil, fmt.Errorf("failed to encrypt API secret: %w", err)
	}

	creds := &domain.ShopifyCredentials{
		ProjectID:   projectID,
		Environment: environment,
		APIKey:      apiKey,
		APISecret:   encryptedSecret,
	}

	if err := s.repository.SaveCredentials(ctx, creds); err != nil {
		s.logger.Error().Err(err).Str("projectId", projectID).Str("environment", environment).Msg("Failed to save credentials")
		return nil, fmt.Errorf("failed to save credentials: %w", err)
	}

	s.logger.Info().Str("projectId", projectID).Str("environment", environment).Msg("Credentials saved successfully")
	return creds, nil
}

// GetCredentials retrieves credentials and decrypts the API secret
func (s *CredentialsService) GetCredentials(ctx context.Context, projectID string, environment string) (*domain.ShopifyCredentials, error) {
	creds, err := s.repository.GetCredentials(ctx, projectID, environment)
	if err != nil {
		s.logger.Error().Err(err).Str("projectId", projectID).Str("environment", environment).Msg("Failed to get credentials")
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	if creds == nil {
		return nil, nil
	}

	// Decrypt API secret
	decryptedSecret, err := s.encryption.Decrypt(creds.APISecret)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to decrypt API secret")
		return nil, fmt.Errorf("failed to decrypt API secret: %w", err)
	}

	creds.APISecret = decryptedSecret
	return creds, nil
}

// DeleteCredentials deletes credentials
func (s *CredentialsService) DeleteCredentials(ctx context.Context, projectID string, environment string) error {
	if err := s.repository.DeleteCredentials(ctx, projectID, environment); err != nil {
		s.logger.Error().Err(err).Str("projectId", projectID).Str("environment", environment).Msg("Failed to delete credentials")
		return fmt.Errorf("failed to delete credentials: %w", err)
	}

	s.logger.Info().Str("projectId", projectID).Str("environment", environment).Msg("Credentials deleted successfully")
	return nil
}
