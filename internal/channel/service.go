package channel

type Service struct {
	channelRepo *repository
}

func NewService(channelRepo *repository) *Service {
	return &Service{
		channelRepo: channelRepo,
	}
}
