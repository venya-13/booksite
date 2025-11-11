package service

import (
	"context"
	"google-auth-demo/backend/internal/repo"
)

func (s *Service) GetAllCategories(ctx context.Context) ([]repo.Category, error) {
	return s.Repo.GetAllCategories(ctx)
}

func (s *Service) CreateCategory(ctx context.Context, name string) error {
	return s.Repo.CreateCategory(ctx, name)
}

func (s *Service) DeleteCategory(ctx context.Context, id int) error {
	return s.Repo.DeleteCategory(ctx, id)
}

func (s *Service) RenameCategory(ctx context.Context, id int, name string) error {
	return s.Repo.RenameCategory(ctx, id, name)
}
