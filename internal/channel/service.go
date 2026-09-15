package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"uuid"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/config"
	"mimokocke/internal/shared/identity"
	"mimokocke/internal/shared/permissions"
)

type Service struct {
	cfg         *config.Config
	channelRepo *repository
}

func NewService(cfg *config.Config, channelRepo *repository) *Service {
	return &Service{
		cfg:         cfg,
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

type ChannelDetails struct {
	Channel                      db.Channel
	MaskedWooCommerceCredentials maskedWooCommerceCredentials
	MaskedShopifyCredentials     maskedShopifyCredentials
}

type maskedWooCommerceCredentials struct {
	StoreURL             string
	MaskedConsumerKey    string
	MaskedConsumerSecret string
}
type maskedShopifyCredentials struct{}

func (s *Service) GetChannelDetails(
	ctx context.Context,
	ident identity.Identity,
	channelID uuid.UUID,
) (ChannelDetails, error) {
	if !ident.HasPermission(permissions.ChannelRead) {
		return ChannelDetails{}, permissions.ErrDenied
	}

	var chnDetails ChannelDetails
	var err error

	chnDetails.Channel, err = s.channelRepo.GetChannel(ctx, db.GetChannelParams{
		OrganizationID: ident.OrgID,
		ChannelID:      channelID,
	})
	if err != nil {
		return ChannelDetails{}, err
	}

	credsByte, err := decrypt(chnDetails.Channel.Credentials, s.cfg.MasterKey)
	if err != nil {
		return ChannelDetails{}, err
	}

	switch chnDetails.Channel.Provider {
	case db.ChannelProviderWooCommerce:
		var wooCreds WooCommerceCredentials

		err := json.Unmarshal(credsByte, &wooCreds)
		if err != nil {
			return ChannelDetails{}, err
		}

		chnDetails.MaskedWooCommerceCredentials = maskedWooCommerceCredentials{
			StoreURL:             wooCreds.StoreURL,
			MaskedConsumerKey:    maskKey(wooCreds.ConsumerKey),
			MaskedConsumerSecret: maskKey(wooCreds.ConsumerSecret),
		}

		return chnDetails, nil
	case db.ChannelProviderShopify:
		// TODO(jozekuhar): shopify
		return chnDetails, nil
	default:
		return ChannelDetails{}, fmt.Errorf("match channel provider")
	}
}

type WooCommerceCredentials struct {
	StoreURL       string `json:"store_url"`
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
}

func (s *Service) VerifyWooCommerceCredentials(
	ctx context.Context,
	creds WooCommerceCredentials,
) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/wp-json/wc/v3/orders?per_page=1", creds.StoreURL),
		nil,
	)
	if err != nil {
		return err
	}
	req.SetBasicAuth(creds.ConsumerKey, creds.ConsumerSecret)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 200:
		return nil
	case 401:
		return fmt.Errorf("forbidden")
	default:
		return fmt.Errorf("error")
	}
}

type SaveWooCommerceChannelParams struct {
	Name        string
	Credentials WooCommerceCredentials
}

func (s *Service) SaveWooCommerceChannel(
	ctx context.Context,
	ident identity.Identity,
	params SaveWooCommerceChannelParams,
) (db.Channel, error) {
	if !ident.HasPermission(permissions.ChannelCreate) {
		return db.Channel{}, permissions.ErrDenied
	}

	creds, err := json.Marshal(WooCommerceCredentials{
		StoreURL:       params.Credentials.StoreURL,
		ConsumerKey:    params.Credentials.ConsumerKey,
		ConsumerSecret: params.Credentials.ConsumerSecret,
	})
	if err != nil {
		return db.Channel{}, fmt.Errorf("marshal woo creds: %w", err)
	}

	encryptedCreds, err := encrypt(creds, s.cfg.MasterKey)
	if err != nil {
		return db.Channel{}, fmt.Errorf("encrypt woo creds: %w", err)
	}

	chn, err := s.channelRepo.CreateChannel(ctx, nil, db.CreateChannelParams{
		ID:             uuid.NewV7(),
		OrganizationID: ident.OrgID,
		Name:           params.Name,
		Provider:       db.ChannelProviderWooCommerce,
		Credentials:    encryptedCreds,
		Status:         db.ChannelStatusActive,
	})
	if err != nil {
		// TODO: wrap err with new error for db.Already exists
		return db.Channel{}, fmt.Errorf("db create channel: %w", err)
	}

	return chn, nil
}

func (s *Service) RemoveChannel(
	ctx context.Context,
	ident identity.Identity,
	channelID uuid.UUID,
) error {
	if !ident.HasPermission(permissions.ChannelDelete) {
		return permissions.ErrDenied
	}

	return s.channelRepo.DeleteChannel(ctx, nil, db.DeleteChannelParams{
		OrganizationID: ident.OrgID,
		ChannelID:      channelID,
	})
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "••••••••"
	}
	return "••••••••" + key[len(key)-4:]
}
