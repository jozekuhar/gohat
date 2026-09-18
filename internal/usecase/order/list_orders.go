package orderuc

import (
	"context"
	"uuid"

	"mimokocke/internal/provider/db/sqlc"
)

type ListOrdersUseCase struct {
	repo *repository
}

func NewListOrdersUseCase(repo *repository) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		repo: repo,
	}
}

func (c *ListOrdersUseCase) Execute(
	ctx context.Context,
	orgID uuid.UUID,
) ([]sqlc.ListOrdersWithDetailsRow, error) {
	return c.repo.ListOrdersWithDetails(ctx, orgID)
}
