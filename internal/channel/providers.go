package channel

type WooCommerceCredentials struct {
	URL            string `json:"url"`
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
}

type ShopifyCredentials struct{}
