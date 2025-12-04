package graph

import (
	"context"

	"archie-core-shopify-layer/internal/domain"
)

// getTenantID extracts tenant ID (projectID) from context
func getTenantID(ctx context.Context) string {
	// Try to get projectID first (new approach)
	if projectID, ok := ctx.Value(domain.ProjectIDKey).(string); ok && projectID != "" {
		return projectID
	}
	// Fallback to tenantId for backward compatibility
	if tenantID, ok := ctx.Value(domain.TenantIDKey).(string); ok {
		return tenantID
	}
	return ""
}

// getEnvironment extracts environment from context
func getEnvironment(ctx context.Context) string {
	if environment, ok := ctx.Value(domain.EnvironmentKey).(string); ok {
		return environment
	}
	return domain.DefaultEnvironment // Default environment
}
