// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/models"
)

type RepositoryInterface interface {
	SaveEstate(ctx context.Context, estate *models.Estate) error

	GetEstate(ctx context.Context, uuid string) (*models.Estate, error)

	SaveTree(ctx context.Context, tree *models.Tree) error

	GetTreeByCoordinate(ctx context.Context, estateId uint64, x uint16, y uint16) (*models.Tree, error)
}
