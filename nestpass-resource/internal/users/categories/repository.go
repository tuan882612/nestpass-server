package categories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"

	"nestpass/pkg/httputils"
)

type repository struct {
	postgres *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) *repository {
	return &repository{postgres: pg}
}

func (r *repository) GetAllCategories(ctx context.Context, userID uuid.UUID, params *httputils.Pagination) ([]*Category, error) {
	categories := []*Category{}

	rows, err := r.postgres.Query(ctx, GetAllCategoriesQuery, userID, params.Index, params.Limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		category := &Category{}
		if err := category.Scan(rows); err != nil {
			log.Error().Str("location", "GetAllCategories").Msgf("%v: %v", userID, err)
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (r *repository) GetCategory(ctx context.Context, userID uuid.UUID, key string, isUUID bool) (*Category, error) {
	category := &Category{}

	query := GetNameCategoryQuery
	if isUUID {
		query = GetUUIDCategoryQuery
	}

	row := r.postgres.QueryRow(ctx, query, key, userID)
	if err := category.Scan(row); err != nil {
		if err == pgx.ErrNoRows {
			return nil, apiutils.NewErrNotFound("category not found")
		}

		log.Error().Str("location", "GetCategory").Msgf("%v: %v", userID, err)
		return nil, err
	}

	return category, nil
}

func (r *repository) CreateCategory(ctx context.Context, tx pgx.Tx, category *Category) error {
	_, err := tx.Exec(ctx, InsertCategoryQuery,
		&category.CategoryID,
		&category.UserID,
		&category.Name,
		&category.Description,
	)

	if err != nil {
		var pgErr *pgconn.PgError = nil
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Error().Str("location", "CreateCategory").Msgf("%v: create category failed conflict", category.UserID)
			return apiutils.NewErrConflict("category already exists")
		}

		log.Error().Str("location", "CreateCategory").Msgf("%v: create category failed internal error - %v", category.UserID, err)
		return err
	}

	return nil
}

func (r *repository) UpdateCategory(ctx context.Context, tx pgx.Tx, category *Category) error {
	_, err := tx.Exec(ctx, UpdateCategoryQuery,
		&category.Name,
		&category.Description,
		&category.CategoryID,
		&category.UserID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Info().Str("location", "UpdateCategory").Msgf("%v: failed to update category not found", category.UserID)
			return apiutils.NewErrNotFound("category not found")
		}

		log.Error().Str("location", "UpdateCategory").Msgf("%v: %v", category.UserID, err)
		return err
	}

	return nil
}

func (r *repository) DeleteCategory(ctx context.Context, tx pgx.Tx, categoryID, userID uuid.UUID) error {
	_, err := tx.Exec(ctx, DeleteCategoryQuery, categoryID, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Info().Str("location", "DeleteCategory").Msgf("%v: failed to delete category not found", userID)
			return apiutils.NewErrNotFound("category not found")
		}

		log.Error().Str("location", "DeleteCategory").Msgf("%v: %v", userID, err)
		return err
	}

	return nil
}
