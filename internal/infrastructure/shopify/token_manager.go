package shopify

import (
	"context"
	"fmt"
	"time"

	"archie-core-shopify-layer/internal/ports"
	"github.com/rs/zerolog"
)

// TokenManager manages Shopify access tokens with refresh capabilities
type TokenManager struct {
	encryptionSvc ports.EncryptionService
	logger        zerolog.Logger
}

// NewTokenManager creates a new token manager
func NewTokenManager(encryptionSvc ports.EncryptionService, logger zerolog.Logger) *TokenManager {
	return &TokenManager{
		encryptionSvc: encryptionSvc,
		logger:        logger,
	}
}

// EncryptToken encrypts an access token before storage
func (tm *TokenManager) EncryptToken(token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("token cannot be empty")
	}
	return tm.encryptionSvc.Encrypt(token)
}

// DecryptToken decrypts an access token after retrieval
func (tm *TokenManager) DecryptToken(encryptedToken string) (string, error) {
	if encryptedToken == "" {
		return "", fmt.Errorf("encrypted token cannot be empty")
	}
	return tm.encryptionSvc.Decrypt(encryptedToken)
}

// TokenInfo represents token metadata
type TokenInfo struct {
	Token       string
	ExpiresAt   *time.Time
	Scopes      []string
	ShopDomain  string
	LastUsed    time.Time
}

// ValidateToken checks if a token is still valid
// Note: Shopify access tokens don't expire unless revoked, but we can check if they're still valid
func (tm *TokenManager) ValidateToken(ctx context.Context, token string, shopDomain string) (bool, error) {
	// In a real implementation, you would make a lightweight API call to Shopify
	// to verify the token is still valid. For now, we'll assume non-empty tokens are valid.
	if token == "" {
		return false, fmt.Errorf("token is empty")
	}

	// TODO: Make a lightweight API call to Shopify to verify token validity
	// Example: GET /admin/api/2024-01/shop.json with the token
	// If it returns 401, the token is invalid

	return true, nil
}

// ShouldRefresh checks if a token should be refreshed
// Shopify tokens don't expire, but they can be revoked
func (tm *TokenManager) ShouldRefresh(tokenInfo *TokenInfo) bool {
	// Shopify access tokens don't expire, but we can check:
	// 1. If token hasn't been used in a long time (stale)
	// 2. If we've received errors indicating token is invalid

	if tokenInfo.ExpiresAt != nil && time.Now().After(*tokenInfo.ExpiresAt) {
		return true
	}

	// Check if token is stale (not used in 30 days)
	staleThreshold := 30 * 24 * time.Hour
	if time.Since(tokenInfo.LastUsed) > staleThreshold {
		return true
	}

	return false
}

