package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"uuid"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/identity"
	"mimokocke/internal/shared/permissions"
)

type Service struct {
	channelRepo *repository
}

func NewService(channelRepo *repository) *Service {
	return &Service{
		channelRepo: channelRepo,
	}
}

func (s *Service) GetChannelsOverview(
	ctx context.Context,
	ident identity.Identity,
) ([]db.Channel, error) {
	if !ident.HasPermission(permissions.ChannelRead) {
		return nil, permissions.ErrDenied
	}

	return s.channelRepo.ListChannels(ctx, ident.OrgID)
}

type SaveWooCommerceChannelParams struct {
	Name           string
	URL            string
	ConsumerKey    string
	ConsumerSecret string
}

type WooCommerceCredentials struct {
	URL            string `json:"url"`
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
}

func (s *Service) SaveWooCommerceChannel(
	ctx context.Context,
	ident identity.Identity,
	params SaveWooCommerceChannelParams,
) error {
	if !ident.HasPermission(permissions.ChannelCreate) {
		return permissions.ErrDenied
	}

	credentials, err := json.Marshal(WooCommerceCredentials{
		URL:            params.URL,
		ConsumerKey:    params.ConsumerKey,
		ConsumerSecret: params.ConsumerSecret,
	})
	if err != nil {
		return fmt.Errorf("marshal woocommerce credentials: %w", err)
	}

	_, err = s.channelRepo.CreateChannel(ctx, nil, db.CreateChannelParams{
		ID:             uuid.NewV7(),
		OrganizationID: ident.OrgID,
		Name:           params.Name,
		Provider:       db.ChannelProviderWooCommerce,
		Credentials:    credentials,
		Status:         db.ChannelStatusActive,
	})
	if err != nil {
		return fmt.Errorf("db create channel: %w", err)
	}

	return nil
}

type SaveShopifyChannelParams struct {
	Name  string
	Token string
}

type ShopifyCredentials struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

func (s *Service) SaveShopifyChannel(
	ctx context.Context,
	ident identity.Identity,
	params SaveShopifyChannelParams,
) error {
	return nil
}
