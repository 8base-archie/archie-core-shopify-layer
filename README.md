# Archie Core Shopify Layer

A service layer that wraps Shopify API capabilities, allowing users to configure credentials and use Shopify APIs through this service instead of calling Shopify directly.

## Architecture Overview

This project follows **Hexagonal Architecture** (Ports & Adapters) and **SOLID** principles:

```
┌─────────────────────────────────────────────────────────────┐
│                      Presentation Layer                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   GraphQL    │  │  REST API    │  │   Webhooks   │      │
│  │   Resolvers  │  │   Handlers   │  │   Handlers   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Shopify     │  │ Credentials  │  │   Webhook    │      │
│  │  Service     │  │   Service    │  │   Manager    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Entities    │  │    Errors    │  │   Session   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Repository  │  │   Shopify    │  │  Encryption  │      │
│  │  (MongoDB)   │  │   Client    │  │   Service    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## Current Features

### ✅ Implemented
- **OAuth Flow**: Shopify app installation and token exchange
- **Credentials Management**: Per-project/environment credential storage with encryption
- **Basic Shopify API**: Products, Orders, Customers, Inventory operations
- **Webhook Verification**: HMAC signature verification for webhooks
- **Webhook Logging**: Webhook events are logged to database
- **GraphQL API**: Basic GraphQL schema for Shopify operations
- **MongoDB Persistence**: Shop, credentials, and webhook storage

### ⚠️ Partially Implemented
- **Webhook Management**: WebhookManager exists but doesn't actually create webhooks via Shopify API
- **GraphQL Schema**: Basic schema exists but missing many operations
- **Error Handling**: Error types defined but not fully integrated

### ❌ Missing Features
- **REST API Gateway**: No REST endpoints for proxying Shopify API calls
- **Multi-tenancy Context**: No request context to identify project/environment
- **Extended Shopify APIs**: Missing Collections, Discounts, Fulfillments, Metafields, etc.
- **Storefront API**: Not implemented
- **GraphQL Admin API**: Not implemented
- **Rate Limiting**: No Shopify rate limit handling
- **Token Refresh**: No handling for expired tokens
- **Webhook Dispatcher**: No proper webhook event handler/dispatcher system
- **Access Token Encryption**: Access tokens stored but not encrypted
- **Request Middleware**: No authentication/authorization middleware
- **Retry Logic**: Infrastructure exists but not used

## Key Issues Identified

### 1. Multi-tenancy Architecture Gap
**Problem**: The service uses global API key/secret from environment variables, but credentials are stored per project/environment. There's no way to identify which project/environment is making a request.

**Impact**: Cannot support multiple projects/environments properly.

**Solution**: Implement request context middleware to extract project/environment from headers or query parameters.

### 2. Incomplete Webhook Management
**Problem**: `WebhookManager.createWebhook()` is a placeholder - doesn't actually create webhooks via Shopify Admin API.

**Impact**: Webhooks must be configured manually in Shopify admin.

**Solution**: Implement webhook CRUD operations using Shopify Admin API.

### 3. Limited API Coverage
**Problem**: Only Products, Orders, Customers, and Inventory APIs are implemented. Missing many Shopify APIs.

**Impact**: Users cannot access full Shopify functionality through the service.

**Solution**: Extend `ShopifyClient` interface and implementation to cover more APIs.

### 4. No REST API Gateway
**Problem**: Only GraphQL API exists. Users might want REST endpoints to proxy Shopify API calls directly.

**Impact**: Limited integration options for users.

**Solution**: Implement REST API gateway that proxies requests to Shopify.

### 5. Access Token Security
**Problem**: Access tokens are stored in plaintext (only API secrets are encrypted).

**Impact**: Security risk if database is compromised.

**Solution**: Encrypt access tokens before storage.

## Technology Stack

- **Language**: Go 1.25.1
- **Framework**: Chi Router
- **GraphQL**: gqlgen
- **Database**: MongoDB
- **Shopify SDK**: github.com/bold-commerce/go-shopify/v4
- **Logging**: zerolog
- **Encryption**: Custom encryption service

## Project Structure

```
archie-core-shopify-layer/
├── cmd/
│   └── api/              # Application entry point
├── graph/                # GraphQL layer
│   ├── generated/        # Generated GraphQL code
│   ├── model/            # GraphQL models
│   ├── schema/           # GraphQL schema definitions
│   └── schema.resolvers.go
├── internal/
│   ├── application/      # Application services (use cases)
│   ├── domain/           # Domain entities and business logic
│   ├── infrastructure/   # Infrastructure implementations
│   │   ├── cache/        # Caching layer
│   │   ├── encryption/   # Encryption service
│   │   ├── metrics/      # Metrics collection
│   │   ├── repository/   # Data persistence
│   │   ├── retry/        # Retry logic
│   │   └── shopify/      # Shopify client adapter
│   └── ports/            # Ports (interfaces)
└── k8s/                  # Kubernetes configurations
```

## Environment Variables

- `MONGODB_URI`: MongoDB connection string
- `MONGODB_DATABASE`: MongoDB database name
- `SHOPIFY_API_KEY`: Shopify API key (global, fallback)
- `SHOPIFY_API_SECRET`: Shopify API secret (global, fallback)
- `ENCRYPTION_KEY`: Encryption key for sensitive data
- `APP_URL`: Application URL for OAuth callbacks
- `PORT`: Server port (default: 8080)

## API Endpoints

### GraphQL
- `POST /query`: GraphQL endpoint
- `GET /`: GraphQL Playground

### REST
- `GET /auth/shopify`: Initiate OAuth flow
- `GET /auth/callback`: OAuth callback handler
- `POST /webhooks/shopify/{shop}`: Webhook receiver
- `GET /health`: Health check

## Development

### Prerequisites
- Go 1.25.1+
- MongoDB
- Shopify Partner account

### Running Locally
```bash
# Load environment variables
export MONGODB_URI="mongodb://localhost:27017"
export SHOPIFY_API_KEY="your_api_key"
export SHOPIFY_API_SECRET="your_api_secret"
export ENCRYPTION_KEY="your_encryption_key"
export APP_URL="http://localhost:8080"

# Run the application
go run cmd/api/main.go
```

## Future Enhancements

See `ENHANCEMENT_PLAN.md` for detailed enhancement roadmap.

