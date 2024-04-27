// This file contains the repository implementation layer.
package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Repository struct {
	Db *bun.DB
}

type NewRepositoryOptions struct {
	Dsn string
}

func NewRepository(opts NewRepositoryOptions) *Repository {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithDSN(opts.Dsn),
	)

	sqldb := sql.OpenDB(pgconn)
	db := bun.NewDB(sqldb, pgdialect.New())

	return &Repository{
		Db: db,
	}
}
