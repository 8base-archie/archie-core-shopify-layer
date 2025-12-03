package shopify

import (
	"context"
	"fmt"
	"strings"

	"archie-core-shopify-layer/internal/ports"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)

type client struct {
	apiKey    string
	apiSecret string
	app       goshopify.App
}

// NewClient creates a new Shopify client adapter
func NewClient(apiKey, apiSecret string) ports.ShopifyClient {
	app := goshopify.App{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
	}
	return &client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		app:       app,
	}
}

// createClient is a helper to create a goshopify client
func (c *client) createClient(shopDomain string, accessToken string) (*goshopify.Client, error) {
	client, err := goshopify.NewClient(c.app, shopDomain, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return client, nil
}

// Authentication methods

func (c *client) GenerateAuthURL(shop string, scopes []string) (string, error) {
	return c.app.AuthorizeUrl(shop, strings.Join(scopes, ","))
}

func (c *client) ExchangeToken(ctx context.Context, shop string, code string) (string, error) {
	token, err := c.app.GetAccessToken(ctx, shop, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange token: %w", err)
	}
	return token, nil
}

// Shop API

func (c *client) GetShop(ctx context.Context, shopDomain string, accessToken string) (*goshopify.Shop, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	shop, err := client.Shop.Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}
	return shop, nil
}

// Product API

func (c *client) GetProducts(ctx context.Context, shopDomain string, accessToken string, options interface{}) ([]goshopify.Product, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	products, err := client.Product.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}

func (c *client) GetProduct(ctx context.Context, shopDomain string, accessToken string, productID int64) (*goshopify.Product, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	product, err := client.Product.Get(ctx, uint64(productID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (c *client) CreateProduct(ctx context.Context, shopDomain string, accessToken string, product *goshopify.Product) (*goshopify.Product, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	created, err := client.Product.Create(ctx, *product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return created, nil
}

func (c *client) UpdateProduct(ctx context.Context, shopDomain string, accessToken string, product *goshopify.Product) (*goshopify.Product, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	updated, err := client.Product.Update(ctx, *product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	return updated, nil
}

func (c *client) DeleteProduct(ctx context.Context, shopDomain string, accessToken string, productID int64) error {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return err
	}
	err = client.Product.Delete(ctx, uint64(productID))
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

// Order API

func (c *client) GetOrders(ctx context.Context, shopDomain string, accessToken string, options interface{}) ([]goshopify.Order, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	orders, err := client.Order.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	return orders, nil
}

func (c *client) GetOrder(ctx context.Context, shopDomain string, accessToken string, orderID int64) (*goshopify.Order, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	order, err := client.Order.Get(ctx, uint64(orderID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

func (c *client) CreateOrder(ctx context.Context, shopDomain string, accessToken string, order *goshopify.Order) (*goshopify.Order, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	created, err := client.Order.Create(ctx, *order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return created, nil
}

func (c *client) UpdateOrder(ctx context.Context, shopDomain string, accessToken string, order *goshopify.Order) (*goshopify.Order, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	updated, err := client.Order.Update(ctx, *order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return updated, nil
}

func (c *client) CancelOrder(ctx context.Context, shopDomain string, accessToken string, orderID int64) (*goshopify.Order, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	cancelled, err := client.Order.Cancel(ctx, uint64(orderID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}
	return cancelled, nil
}

// Customer API

func (c *client) GetCustomers(ctx context.Context, shopDomain string, accessToken string, options interface{}) ([]goshopify.Customer, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	customers, err := client.Customer.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	return customers, nil
}

func (c *client) GetCustomer(ctx context.Context, shopDomain string, accessToken string, customerID int64) (*goshopify.Customer, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	customer, err := client.Customer.Get(ctx, uint64(customerID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	return customer, nil
}

func (c *client) CreateCustomer(ctx context.Context, shopDomain string, accessToken string, customer *goshopify.Customer) (*goshopify.Customer, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	created, err := client.Customer.Create(ctx, *customer)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}
	return created, nil
}

func (c *client) UpdateCustomer(ctx context.Context, shopDomain string, accessToken string, customer *goshopify.Customer) (*goshopify.Customer, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	updated, err := client.Customer.Update(ctx, *customer)
	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}
	return updated, nil
}

func (c *client) DeleteCustomer(ctx context.Context, shopDomain string, accessToken string, customerID int64) error {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return err
	}
	err = client.Customer.Delete(ctx, uint64(customerID))
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	return nil
}

func (c *client) SearchCustomers(ctx context.Context, shopDomain string, accessToken string, query string) ([]goshopify.Customer, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	customers, err := client.Customer.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}
	return customers, nil
}

// Inventory API

func (c *client) GetInventoryLevels(ctx context.Context, shopDomain string, accessToken string, options interface{}) ([]goshopify.InventoryLevel, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	levels, err := client.InventoryLevel.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory levels: %w", err)
	}
	return levels, nil
}

func (c *client) UpdateInventoryLevel(ctx context.Context, shopDomain string, accessToken string, inventoryLevel *goshopify.InventoryLevel) (*goshopify.InventoryLevel, error) {
	client, err := c.createClient(shopDomain, accessToken)
	if err != nil {
		return nil, err
	}
	// Note: The actual update method may vary based on the goshopify library version
	// This is a placeholder implementation
	updated, err := client.InventoryLevel.Set(ctx, *inventoryLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to update inventory level: %w", err)
	}
	return updated, nil
}
