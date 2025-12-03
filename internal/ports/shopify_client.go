package ports

import (
	"context"

	shopify "github.com/bold-commerce/go-shopify/v4"
)

// ShopifyClient defines the interface for Shopify API operations
type ShopifyClient interface {
	// Authentication
	GenerateAuthURL(shop string, scopes []string) (string, error)
	ExchangeToken(ctx context.Context, shop string, code string) (string, error)

	// Shop API
	GetShop(ctx context.Context, shop string, accessToken string) (*shopify.Shop, error)

	// Product API
	GetProducts(ctx context.Context, shop string, accessToken string, options interface{}) ([]shopify.Product, error)
	GetProduct(ctx context.Context, shop string, accessToken string, productID int64) (*shopify.Product, error)
	CreateProduct(ctx context.Context, shop string, accessToken string, product *shopify.Product) (*shopify.Product, error)
	UpdateProduct(ctx context.Context, shop string, accessToken string, product *shopify.Product) (*shopify.Product, error)
	DeleteProduct(ctx context.Context, shop string, accessToken string, productID int64) error

	// Order API
	GetOrders(ctx context.Context, shop string, accessToken string, options interface{}) ([]shopify.Order, error)
	GetOrder(ctx context.Context, shop string, accessToken string, orderID int64) (*shopify.Order, error)
	CreateOrder(ctx context.Context, shop string, accessToken string, order *shopify.Order) (*shopify.Order, error)
	UpdateOrder(ctx context.Context, shop string, accessToken string, order *shopify.Order) (*shopify.Order, error)
	CancelOrder(ctx context.Context, shop string, accessToken string, orderID int64) (*shopify.Order, error)

	// Customer API
	GetCustomers(ctx context.Context, shop string, accessToken string, options interface{}) ([]shopify.Customer, error)
	GetCustomer(ctx context.Context, shop string, accessToken string, customerID int64) (*shopify.Customer, error)
	CreateCustomer(ctx context.Context, shop string, accessToken string, customer *shopify.Customer) (*shopify.Customer, error)
	UpdateCustomer(ctx context.Context, shop string, accessToken string, customer *shopify.Customer) (*shopify.Customer, error)
	DeleteCustomer(ctx context.Context, shop string, accessToken string, customerID int64) error
	SearchCustomers(ctx context.Context, shop string, accessToken string, query string) ([]shopify.Customer, error)

	// Inventory API
	GetInventoryLevels(ctx context.Context, shop string, accessToken string, options interface{}) ([]shopify.InventoryLevel, error)
	UpdateInventoryLevel(ctx context.Context, shop string, accessToken string, inventoryLevel *shopify.InventoryLevel) (*shopify.InventoryLevel, error)
}
