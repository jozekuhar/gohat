package sqlc

type ChannelProvider string

func (p ChannelProvider) String() string {
	return string(p)
}

const (
	ChannelProviderShopify     = "shopify"
	ChannelProviderWooCommerce = "woocommerce"
	ChannelProviderManual      = "manual"
)

type ChannelStatus string

func (s ChannelStatus) String() string {
	return string(s)
}

const (
	ChannelStatusActive   ChannelStatus = "active"
	ChannelStatusInactive ChannelStatus = "inactive"
)
