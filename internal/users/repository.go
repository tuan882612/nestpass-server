package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tuan882612/apiutils"
)

type repository struct {
	postgres *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) *repository {
	return &repository{postgres: pg}
}

func (r *repository) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	user := &User{}

	row := r.postgres.QueryRow(ctx, GetUserQuery, userID)

	if err := user.Scan(row); err != nil {
		if err == pgx.ErrNoRows {
			return nil, apiutils.NewErrNotFound("user not found")
		}

		return nil, err
	}

	return user, nil
}
