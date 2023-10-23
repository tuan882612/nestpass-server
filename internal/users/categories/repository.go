package categories

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	postgres *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) *repository {
	return &repository{postgres: pg}
}
