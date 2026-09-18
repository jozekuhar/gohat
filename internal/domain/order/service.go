package order

type Service struct {
	repo *repository
}

func NewService(repo *repository) *Service {
	return &Service{
		repo: repo,
	}
}
