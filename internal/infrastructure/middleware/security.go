package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
)

// SecurityMiddleware provides security-related middleware functions

// InputValidationMiddleware validates and sanitizes input
func InputValidationMiddleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate shop domain format
			if shop := r.URL.Query().Get("shop"); shop != "" {
				if !isValidShopDomain(shop) {
					logger.Warn().
						Str("shop", shop).
						Str("ip", r.RemoteAddr).
						Msg("Invalid shop domain format")
					http.Error(w, "Invalid shop domain format", http.StatusBadRequest)
					return
				}
			}

			// Validate shop domain from header
			if shop := r.Header.Get("X-Shop-Domain"); shop != "" {
				if !isValidShopDomain(shop) {
					logger.Warn().
						Str("shop", shop).
						Str("ip", r.RemoteAddr).
						Msg("Invalid shop domain format in header")
					http.Error(w, "Invalid shop domain format", http.StatusBadRequest)
					return
				}
			}

			// Validate project ID format (alphanumeric, dashes, underscores)
			if projectID := r.Header.Get("X-Project-ID"); projectID != "" {
				if !isValidProjectID(projectID) {
					logger.Warn().
						Str("projectId", projectID).
						Str("ip", r.RemoteAddr).
						Msg("Invalid project ID format")
					http.Error(w, "Invalid project ID format", http.StatusBadRequest)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// isValidShopDomain validates Shopify shop domain format
func isValidShopDomain(domain string) bool {
	// Shopify domains: mystore.myshopify.com or custom domains
	// For simplicity, check basic format
	if domain == "" {
		return false
	}

	// Check for myshopify.com domain
	if strings.HasSuffix(domain, ".myshopify.com") {
		shopName := strings.TrimSuffix(domain, ".myshopify.com")
		// Shop name: 3-40 characters, lowercase letters, numbers, hyphens
		matched, _ := regexp.MatchString(`^[a-z0-9-]{3,40}$`, shopName)
		return matched
	}

	// Custom domain validation (basic)
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]?\.[a-zA-Z]{2,}$`, domain)
	return matched
}

// isValidProjectID validates project ID format
func isValidProjectID(projectID string) bool {
	if projectID == "" {
		return false
	}
	// Alphanumeric, dashes, underscores, 1-64 characters
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]{1,64}$`, projectID)
	return matched
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add security headers
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Don't expose server information
			w.Header().Set("Server", "")

			next.ServeHTTP(w, r)
		})
	}
}

// AuditLoggingMiddleware logs security-relevant events
func AuditLoggingMiddleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log authentication attempts
			projectID := r.Header.Get("X-Project-ID")
			shopDomain := r.URL.Query().Get("shop")
			if shopDomain == "" {
				shopDomain = r.Header.Get("X-Shop-Domain")
			}

			// Log sensitive operations
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
				logger.Info().
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("projectId", projectID).
					Str("shop", shopDomain).
					Str("ip", r.RemoteAddr).
					Str("userAgent", r.UserAgent()).
					Msg("Security audit: Sensitive operation")
			}

			next.ServeHTTP(w, r)
		})
	}
}

