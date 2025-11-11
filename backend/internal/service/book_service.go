package service

import (
	"context"
	"google-auth-demo/backend/internal/repo"
)

func (s *Service) CreateBook(ctx context.Context, b repo.Book, catIDs []int) (int, error) {
	id, err := s.Repo.CreateBook(ctx, b, catIDs)
	if err != nil {
		return 0, err
	}

	for _, catID := range catIDs {
		if catID == 0 {
			continue
		}
		if err := s.Repo.AssignBookToCategory(ctx, id, catID); err != nil {
			return id, err
		}
	}

	return id, nil
}

func (s *Service) GetAllBooks(ctx context.Context) ([]repo.Book, error) {
	return s.Repo.GetAllBooks(ctx)
}

func (s *Service) GetBookByID(ctx context.Context, id int) (*repo.Book, error) {
	b, err := s.Repo.GetBookByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *Service) UpdateBook(ctx context.Context, book repo.Book, categoryIDs []int) error {
	return s.Repo.UpdateBook(ctx, book, categoryIDs)
}

func (s *Service) DeleteBook(ctx context.Context, id int) error {
	return s.Repo.DeleteBook(ctx, id)
}

func (s *Service) GetBooksByCategory(ctx context.Context, categoryID int) ([]repo.Book, error) {
	return s.Repo.GetBooksByCategory(ctx, categoryID)
}

func (s *Service) GetCategoriesWithBooks(ctx context.Context) ([]struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Books []repo.Book `json:"books"`
}, error) {
	return s.Repo.GetCategoriesWithBooks(ctx)
}

func (s *Service) GetUncategorizedBooks(ctx context.Context) ([]repo.Book, error) {
	return s.Repo.GetUncategorizedBooks(ctx)
}

func (s *Service) AssignBookToCategory(ctx context.Context, bookID int, categoryID int) error {
	return s.Repo.AssignBookToCategory(ctx, bookID, categoryID)
}
