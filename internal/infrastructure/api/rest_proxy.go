package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"archie-core-shopify-layer/internal/application"
	"archie-core-shopify-layer/internal/domain"
	"github.com/rs/zerolog"
)

// RESTProxy handles REST API proxy requests to Shopify
type RESTProxy struct {
	shopifyService *application.ShopifyService
	logger         zerolog.Logger
}

// NewRESTProxy creates a new REST proxy
func NewRESTProxy(shopifyService *application.ShopifyService, logger zerolog.Logger) *RESTProxy {
	return &RESTProxy{
		shopifyService: shopifyService,
		logger:         logger,
	}
}

// HandleProxyRequest handles proxied REST API requests
func (p *RESTProxy) HandleProxyRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract project ID and environment from context (set by middleware)
	projectID := domain.GetProjectIDFromContext(ctx)

	if projectID == "" {
		http.Error(w, "X-Project-ID header is required", http.StatusBadRequest)
		return
	}

	// Extract shop domain from query parameter or header
	shopDomain := r.URL.Query().Get("shop")
	if shopDomain == "" {
		shopDomain = r.Header.Get("X-Shop-Domain")
	}
	if shopDomain == "" {
		http.Error(w, "shop parameter or X-Shop-Domain header is required", http.StatusBadRequest)
		return
	}

	// Get shop to retrieve access token
	shop, err := p.shopifyService.GetShop(ctx, shopDomain)
	if err != nil {
		p.logger.Error().Err(err).Str("shop", shopDomain).Msg("Failed to get shop")
		http.Error(w, "Shop not found or not authenticated", http.StatusNotFound)
		return
	}

	// Extract Shopify API path from request
	// Format: /api/v1/{project}/{environment}/shopify/{resource}
	pathParts := strings.Split(r.URL.Path, "/")
	shopifyPathIndex := -1
	for i, part := range pathParts {
		if part == "shopify" {
			shopifyPathIndex = i + 1
			break
		}
	}

	if shopifyPathIndex == -1 || shopifyPathIndex >= len(pathParts) {
		http.Error(w, "Invalid API path", http.StatusBadRequest)
		return
	}

	shopifyPath := "/admin/api/2024-01/" + strings.Join(pathParts[shopifyPathIndex:], "/")
	if r.URL.RawQuery != "" {
		shopifyPath += "?" + r.URL.RawQuery
	}

	// Build Shopify API URL
	shopifyURL := fmt.Sprintf("https://%s%s", shopDomain, shopifyPath)

	// Create request to Shopify
	req, err := http.NewRequestWithContext(ctx, r.Method, shopifyURL, r.Body)
	if err != nil {
		p.logger.Error().Err(err).Msg("Failed to create Shopify request")
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers (excluding host and connection)
	for key, values := range r.Header {
		if key != "Host" && key != "Connection" && !strings.HasPrefix(key, "X-") {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	// Set Shopify-specific headers
	req.Header.Set("X-Shopify-Access-Token", shop.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	// Make request to Shopify
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		p.logger.Error().Err(err).Str("url", shopifyURL).Msg("Failed to make Shopify request")
		http.Error(w, "Failed to connect to Shopify", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.logger.Error().Err(err).Msg("Failed to read Shopify response")
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Transform response if needed
	transformedBody, err := p.transformResponse(body, r.URL.Path)
	if err != nil {
		p.logger.Warn().Err(err).Msg("Failed to transform response, using original")
		transformedBody = body
	}

	// Write response
	if _, err := w.Write(transformedBody); err != nil {
		p.logger.Error().Err(err).Msg("Failed to write response")
	}
}

// transformResponse transforms Shopify API response to standardized format
func (p *RESTProxy) transformResponse(body []byte, path string) ([]byte, error) {
	// Parse JSON response
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		// Not JSON, return as-is
		return body, nil
	}

	// Transform to standardized format
	transformed := map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"path": path,
		},
	}

	return json.Marshal(transformed)
}

