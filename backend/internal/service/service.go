package service

import (
	"context"
	"google-auth-demo/backend/internal/repo"
)

type OAuth interface {
	GetAuthURL() string
	ExchangeCode(string) (*TokenData, error)
	RefreshAccessToken(string) (*TokenData, error)
	FetchProfile(string) (map[string]interface{}, error)
}

type Repository interface {
	SaveOrUpdate(user map[string]interface{}) error
	GetUserByGoogleID(googleID string) (map[string]interface{}, error)

	GetAllCategories(ctx context.Context) ([]repo.Category, error)
	CreateCategory(ctx context.Context, name string) error
	DeleteCategory(ctx context.Context, id int) error

	RenameCategory(ctx context.Context, id int, name string) error

	CreateBook(ctx context.Context, b repo.Book, categoryIDs []int) (int, error)
	GetAllBooks(ctx context.Context) ([]repo.Book, error)
	GetBookByID(ctx context.Context, id int) (repo.Book, error)
	UpdateBook(ctx context.Context, b repo.Book, categoryIDs []int) error
	DeleteBook(ctx context.Context, id int) error
	GetBooksByCategory(ctx context.Context, categoryID int) ([]repo.Book, error)

	GetCategoriesWithBooks(ctx context.Context) ([]struct {
		ID    int         `json:"id"`
		Name  string      `json:"name"`
		Books []repo.Book `json:"books"`
	}, error)

	GetUncategorizedBooks(ctx context.Context) ([]repo.Book, error)
	AssignBookToCategory(ctx context.Context, bookID int, categoryID int) error
}

type Config struct {
	FrontendURL string `env:"FRONTEND_URL"`
}
