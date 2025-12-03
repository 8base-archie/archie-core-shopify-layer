package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"archie-core-shopify-layer/graph"
	"archie-core-shopify-layer/graph/generated"
	"archie-core-shopify-layer/internal/application"
	"archie-core-shopify-layer/internal/domain"
	"archie-core-shopify-layer/internal/infrastructure/encryption"
	"archie-core-shopify-layer/internal/infrastructure/repository"
	shopifyinfra "archie-core-shopify-layer/internal/infrastructure/shopify"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Initialize logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if err := godotenv.Load(); err != nil {
		_ = fmt.Sprint("⚠️  Warning: .env file not found")
	}

	// Get configuration from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	shopifyAPIKey := os.Getenv("SHOPIFY_API_KEY")
	if shopifyAPIKey == "" {
		logger.Fatal().Msg("SHOPIFY_API_KEY environment variable is required")
	}

	shopifyAPISecret := os.Getenv("SHOPIFY_API_SECRET")
	if shopifyAPISecret == "" {
		logger.Fatal().Msg("SHOPIFY_API_SECRET environment variable is required")
	}

	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:8080"
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}
	defer client.Disconnect(context.Background())

	db := client.Database(os.Getenv("MONGODB_DATABASE"))

	// Get encryption key
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		logger.Fatal().Msg("ENCRYPTION_KEY environment variable is required")
	}

	// Initialize infrastructure
	repo := repository.NewMongoRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	shopifyClient := shopifyinfra.NewClient(shopifyAPIKey, shopifyAPISecret)
	webhookVerifier := shopifyinfra.NewWebhookVerifier(shopifyAPISecret)

	// Initialize encryption service
	encryptionService, err := encryption.NewService(encryptionKey)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize encryption service")
	}

	// Initialize application services
	shopifyService := application.NewShopifyService(
		repo,
		shopifyClient,
		logger,
		shopifyAPIKey,
		shopifyAPISecret,
	)

	credentialsService := application.NewCredentialsService(
		repo,
		encryptionService,
		logger,
	)

	webhookManager := application.NewWebhookManager(
		shopifyClient,
		logger,
		appURL+"/webhooks/shopify",
	)

	// Create GraphQL resolver
	resolver := graph.NewResolver(shopifyService, credentialsService)

	// Create GraphQL executable schema
	execSchema := generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	})

	// Create GraphQL handler
	srv := handler.NewDefaultServer(execSchema)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	// Routes
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	// OAuth routes
	r.Get("/auth/shopify", oauthInitHandler(sessionRepo, shopifyAPIKey, appURL, logger))
	r.Get("/auth/callback", oauthCallbackHandler(sessionRepo, shopifyService, webhookManager, shopifyAPIKey, shopifyAPISecret, logger))

	// Webhook endpoint: POST /webhooks/shopify/{shop}
	r.Post("/webhooks/shopify/{shop}", webhookHandler(shopifyService, webhookVerifier, shopifyAPISecret, logger))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info().Str("port", port).Msg("Starting API server")
	logger.Info().Msg("GraphQL Playground available at http://localhost:" + port + "/")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}

// oauthInitHandler initiates the OAuth flow
func oauthInitHandler(sessionRepo *repository.SessionRepository, apiKey string, appURL string, logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shop := r.URL.Query().Get("shop")
		if shop == "" {
			http.Error(w, "shop parameter is required", http.StatusBadRequest)
			return
		}

		// Generate random state for CSRF protection
		stateBytes := make([]byte, 16)
		if _, err := rand.Read(stateBytes); err != nil {
			logger.Error().Err(err).Msg("Failed to generate state")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		state := hex.EncodeToString(stateBytes)

		// Save session
		session := &domain.Session{
			Shop:      shop,
			State:     state,
			Scopes:    []string{"read_products", "write_products", "read_orders", "write_orders"},
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}

		if err := sessionRepo.CreateSession(r.Context(), session); err != nil {
			logger.Error().Err(err).Msg("Failed to create session")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Build authorization URL
		scopes := "read_products,write_products,read_orders,write_orders"
		redirectURI := appURL + "/auth/callback"
		authURL := fmt.Sprintf(
			"https://%s/admin/oauth/authorize?client_id=%s&scope=%s&redirect_uri=%s&state=%s",
			shop,
			apiKey,
			url.QueryEscape(scopes),
			url.QueryEscape(redirectURI),
			state,
		)

		http.Redirect(w, r, authURL, http.StatusFound)
	}
}

// oauthCallbackHandler handles the OAuth callback
func oauthCallbackHandler(
	sessionRepo *repository.SessionRepository,
	shopifyService *application.ShopifyService,
	webhookManager *application.WebhookManager,
	apiKey string,
	apiSecret string,
	logger zerolog.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Verify HMAC
		if err := shopifyinfra.VerifyHMAC(r.URL.Query(), apiSecret); err != nil {
			logger.Warn().Err(err).Msg("HMAC verification failed")
			http.Error(w, "Invalid request", http.StatusUnauthorized)
			return
		}

		// Get parameters
		shop := r.URL.Query().Get("shop")
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if shop == "" || code == "" || state == "" {
			http.Error(w, "Missing required parameters", http.StatusBadRequest)
			return
		}

		// Verify state
		session, err := sessionRepo.GetSession(ctx, state)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get session")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if session == nil || session.Shop != shop {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		// Delete session
		sessionRepo.DeleteSession(ctx, state)

		// Exchange token
		shopDomain, err := shopifyService.ExchangeToken(ctx, shop, code)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to exchange token")
			http.Error(w, "Failed to complete installation", http.StatusInternalServerError)
			return
		}

		// Subscribe to webhooks
		topics := webhookManager.GetDefaultTopics()
		// Note: We would need the access token here to subscribe to webhooks
		// This is a placeholder for now
		logger.Info().
			Str("shop", shop).
			Interface("topics", topics).
			Msg("Would subscribe to webhooks")

		// Redirect to success page
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<html>
				<head><title>Installation Complete</title></head>
				<body>
					<h1>Installation Complete!</h1>
					<p>Shop: %s</p>
					<p>You can now close this window.</p>
				</body>
			</html>
		`, shopDomain.Domain)
	}
}

// webhookHandler handles Shopify webhook requests
func webhookHandler(
	shopifyService *application.ShopifyService,
	webhookVerifier *shopifyinfra.WebhookVerifier,
	apiSecret string,
	logger zerolog.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract shop from URL
		shop := chi.URLParam(r, "shop")
		if shop == "" {
			http.Error(w, "shop is required", http.StatusBadRequest)
			return
		}

		// Get webhook topic from header
		topic := r.Header.Get("X-Shopify-Topic")
		if topic == "" {
			logger.Warn().Msg("Missing X-Shopify-Topic header")
			http.Error(w, "Missing X-Shopify-Topic header", http.StatusBadRequest)
			return
		}

		// Read request body
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to read webhook payload")
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Verify webhook signature
		hmacHeader := r.Header.Get("X-Shopify-Hmac-SHA256")
		if err := webhookVerifier.Verify(payload, hmacHeader); err != nil {
			logger.Warn().Err(err).Str("shop", shop).Msg("Webhook signature verification failed")
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		// Process webhook event
		if err := shopifyService.ProcessWebhook(ctx, topic, shop, payload, true); err != nil {
			logger.Error().
				Err(err).
				Str("topic", topic).
				Str("shop", shop).
				Msg("Failed to process webhook event")

			// Return 500 to trigger Shopify retry
			http.Error(w, "Failed to process webhook event", http.StatusInternalServerError)
			return
		}

		// Return success
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"received": "true",
		})
	}
}
