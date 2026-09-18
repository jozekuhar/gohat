package courier

import "mimokocke/internal/provider/db"

type repository struct {
	*db.BaseRepo
}

func NewRepository(repo *db.BaseRepo) *repository {
	return &repository{
		repo,
	}
}
