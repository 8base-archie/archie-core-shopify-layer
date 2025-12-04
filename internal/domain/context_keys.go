package domain

import "context"

// ContextKey is a type-safe key for context values
type ContextKey string

const (
	// ProjectIDKey is the key for the project ID in the context
	ProjectIDKey ContextKey = "projectId"
	// EnvironmentKey is the key for the environment in the context
	EnvironmentKey ContextKey = "environment"
	// TenantIDKey is the key for the tenant ID in the context (backward compatibility)
	TenantIDKey ContextKey = "tenantId"
)

// GetProjectIDFromContext extracts the project ID from context in a type-safe way
func GetProjectIDFromContext(ctx context.Context) string {
	if projectID, ok := ctx.Value(ProjectIDKey).(string); ok {
		return projectID
	}
	return ""
}

// GetEnvironmentFromContext extracts the environment from context in a type-safe way
func GetEnvironmentFromContext(ctx context.Context) string {
	if environment, ok := ctx.Value(EnvironmentKey).(string); ok {
		return environment
	}
	return ""
}

// GetTenantIDFromContext extracts the tenant ID from context in a type-safe way (backward compatibility)
func GetTenantIDFromContext(ctx context.Context) string {
	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
		return tenantID
	}
	return ""
}

// WithProjectID adds the project ID to the context in a type-safe way
func WithProjectID(ctx context.Context, projectID string) context.Context {
	return context.WithValue(ctx, ProjectIDKey, projectID)
}

// WithEnvironment adds the environment to the context in a type-safe way
func WithEnvironment(ctx context.Context, environment string) context.Context {
	return context.WithValue(ctx, EnvironmentKey, environment)
}

// WithTenantID adds the tenant ID to the context in a type-safe way (backward compatibility)
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

