package channel

import (
	"context"

	"mimokocke/internal/shared/authz"

	"github.com/goforj/godump"
)

type Service struct {
	channelRepo *repository
}

func NewService(channelRepo *repository) *Service {
	return &Service{
		channelRepo: channelRepo,
	}
}

func (s *Service) GetChannels(ctx context.Context, identity authz.Identity) ([]Channel, error) {
	// Permissions READ
	godump.Dump(identity)

	channels, err := s.channelRepo.ListChannels(ctx, identity.OrganizationID)
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (s *Service) CreateChannel(ctx context.Context) (*Channel, error) {
	// permissions WRITE
	panic("unimplemented")
}

func (s *Service) DeleteChannel(ctx context.Context) (*Channel, error) {
	// permissions DELETE
	panic("unimplemented")
}
