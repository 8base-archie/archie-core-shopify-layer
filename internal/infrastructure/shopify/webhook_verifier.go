package shopify

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// WebhookVerifier handles webhook signature verification
type WebhookVerifier struct {
	apiSecret string
}

// NewWebhookVerifier creates a new webhook verifier
func NewWebhookVerifier(apiSecret string) *WebhookVerifier {
	return &WebhookVerifier{
		apiSecret: apiSecret,
	}
}

// Verify verifies the webhook signature
func (v *WebhookVerifier) Verify(body []byte, hmacHeader string) error {
	if hmacHeader == "" {
		return fmt.Errorf("missing X-Shopify-Hmac-SHA256 header")
	}

	// Calculate HMAC
	mac := hmac.New(sha256.New, []byte(v.apiSecret))
	mac.Write(body)
	expectedHMAC := mac.Sum(nil)

	// Decode received HMAC (base64 encoded)
	receivedHMAC, err := base64.StdEncoding.DecodeString(hmacHeader)
	if err != nil {
		return fmt.Errorf("failed to decode HMAC: %w", err)
	}

	// Compare HMACs
	if !hmac.Equal(receivedHMAC, expectedHMAC) {
		return fmt.Errorf("webhook signature verification failed")
	}

	return nil
}
