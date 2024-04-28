package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/models"
)

func (r *Repository) EstatePersist(ctx context.Context, estate *models.Estate) (*models.Estate, error) {
	_, err := r.Db.NewInsert().
		Model(estate).
		ExcludeColumn("id").
		Returning("uuid").
		On("CONFLICT (uuid) DO UPDATE").
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return estate, nil
}
