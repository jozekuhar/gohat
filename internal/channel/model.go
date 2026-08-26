package channel

import (
	"time"
	"uuid"
)

type provider string

const (
	WooCommerceProvider provider = "woocommerce"
	ShopifyProvider     provider = "shopify"
)

type Channel struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Provider       string
	Name           string
	URL            string
	Credentials    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type WoocommerceCredentials struct {
	ConsumerKey    string
	ConsumerSecret string
}

// Not yet done
type ShopifyCredentials struct {
	AccessToken string
}
