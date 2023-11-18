package categories

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"nestpass/pkg/httputils"
)

type service struct {
	repo *repository
}

func NewService(repo *repository) *service {
	return &service{repo: repo}
}

func (s *service) GetAllCategories(ctx context.Context, userID uuid.UUID, page *httputils.Pagination) ([]*Category, error) {
	return s.repo.GetAllCategories(ctx, userID, page)
}

func (s *service) GetCategory(ctx context.Context, userID uuid.UUID, key string) (*Category, error) {
	isUUID := true 
	if _, err := uuid.Parse(key); err != nil {
		isUUID = false
	}
		
	return s.repo.GetCategory(ctx, userID, key, isUUID)
}

func (s *service) CreateCategory(ctx context.Context, category *Category) (*Category, error) {
	tx, err := s.repo.postgres.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	categoryResp := New(category.Name, category.Description, category.UserID)
	if err := s.repo.CreateCategory(ctx, tx, categoryResp); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Str("location", "CreateCategory").Msgf("%v: %v", category.UserID, err)
		return nil, err
	}

	return categoryResp, nil
}

func (s *service) UpdateCategory(ctx context.Context, category *Category) (*Category, error) {
	tx, err := s.repo.postgres.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.repo.UpdateCategory(ctx, tx, category); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Str("location", "UpdateCategory").Msgf("%v: %v", category.UserID, err)
		return nil, err
	}

	return category, nil
}

func (s *service) DeleteCategory(ctx context.Context, categoryID, userID uuid.UUID) error {
	tx, err := s.repo.postgres.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.repo.DeleteCategory(ctx, tx, categoryID, userID); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Str("location", "DeleteCategory").Msgf("%v: %v", userID, err)
		return err
	}

	return nil
}
