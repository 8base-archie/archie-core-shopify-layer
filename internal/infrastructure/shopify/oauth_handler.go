package shopify

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// VerifyHMAC verifies the HMAC signature from Shopify OAuth requests
func VerifyHMAC(queryParams url.Values, apiSecret string) error {
	// Extract the HMAC from query parameters
	receivedHMAC := queryParams.Get("hmac")
	if receivedHMAC == "" {
		return fmt.Errorf("missing HMAC parameter")
	}

	// Create a copy of query parameters and remove HMAC
	params := make(url.Values)
	for key, values := range queryParams {
		if key != "hmac" {
			params[key] = values
		}
	}

	// Sort parameters and build query string
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var pairs []string
	for _, key := range keys {
		for _, value := range params[key] {
			pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
		}
	}
	message := strings.Join(pairs, "&")

	// Calculate HMAC
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(message))
	expectedHMAC := hex.EncodeToString(mac.Sum(nil))

	// Compare HMACs
	if !hmac.Equal([]byte(receivedHMAC), []byte(expectedHMAC)) {
		return fmt.Errorf("HMAC verification failed")
	}

	return nil
}

// VerifyWebhookHMAC verifies the HMAC signature from Shopify webhooks
func VerifyWebhookHMAC(body []byte, hmacHeader string, apiSecret string) error {
	if hmacHeader == "" {
		return fmt.Errorf("missing HMAC header")
	}

	// Calculate HMAC
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write(body)
	expectedHMAC := mac.Sum(nil)

	// Decode received HMAC (base64 encoded)
	receivedHMAC, err := base64.StdEncoding.DecodeString(hmacHeader)
	if err != nil {
		// Try hex decoding as fallback
		receivedHMAC, err = hex.DecodeString(hmacHeader)
		if err != nil {
			return fmt.Errorf("failed to decode HMAC: %w", err)
		}
	}

	// Compare HMACs
	if !hmac.Equal(receivedHMAC, expectedHMAC) {
		return fmt.Errorf("webhook HMAC verification failed")
	}

	return nil
}
