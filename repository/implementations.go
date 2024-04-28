package repository

import (
	"context"
	"database/sql"

	"github.com/SawitProRecruitment/UserService/models"
)

func (r *Repository) SaveEstate(ctx context.Context, estate *models.Estate) error {
	_, err := r.Db.NewInsert().
		Model(estate).
		ExcludeColumn("id").
		Returning("uuid").
		On("CONFLICT (uuid) DO UPDATE").
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetEstate(ctx context.Context, uuid string) (*models.Estate, error) {
	var estate models.Estate

	err := r.Db.NewSelect().Model(&estate).Where("uuid = ?", uuid).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &estate, nil
}

func (r *Repository) SaveTree(ctx context.Context, tree *models.Tree) error {
	tx, err := r.Db.Begin()
	if err != nil {
		return err
	}

	_, err = r.Db.NewInsert().
		Model(tree).
		ExcludeColumn("id").
		Returning("uuid").
		On("CONFLICT (uuid) DO UPDATE").
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	estate := tree.Estate
	_, err = r.Db.NewUpdate().Model(&estate).Where("id = ?", estate.ID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (r *Repository) GetTreeByCoordinate(ctx context.Context, estateId uint64, x uint16, y uint16) (*models.Tree, error) {
	var tree models.Tree
	err := r.Db.NewSelect().Model(&tree).
		Where("estate_id = ?", estateId).
		Where("x = ?", x).
		Where("y = ?", y).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &tree, nil
}
