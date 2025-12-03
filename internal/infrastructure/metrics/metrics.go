package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shopify_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shopify_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Shopify API metrics
	ShopifyAPICallsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shopify_api_calls_total",
			Help: "Total number of Shopify API calls",
		},
		[]string{"operation", "status"},
	)

	ShopifyAPICallDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shopify_api_call_duration_seconds",
			Help:    "Shopify API call duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	ShopifyAPIErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shopify_api_errors_total",
			Help: "Total number of Shopify API errors",
		},
		[]string{"operation", "error_type"},
	)

	// Cache metrics
	CacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "shopify_cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "shopify_cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	// Webhook metrics
	WebhooksReceived = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shopify_webhooks_received_total",
			Help: "Total number of webhooks received",
		},
		[]string{"topic", "verified"},
	)

	WebhookProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shopify_webhook_processing_duration_seconds",
			Help:    "Webhook processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"topic"},
	)

	// Database metrics
	DatabaseQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shopify_database_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "status"},
	)

	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shopify_database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)
