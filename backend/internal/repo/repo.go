package repo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google-auth-demo/backend/internal/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Book struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	Description string     `json:"description"`
	FileURL     string     `json:"file_url"`
	CoverPath   string     `json:"cover_path"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Categories  []Category `json:"categories"`
}

type PostgresRepo struct {
	db *pgxpool.Pool
}

type PostgresConfig struct {
	DSN string
}

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	IsSystem  bool      `json:"is_system"`
	CreatedAt time.Time `json:"created_at"`
}

// CategoryWithBooks is used by categories-with-books API.
type CategoryWithBooks struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	IsSystem bool   `json:"is_system"`
	Books    []Book `json:"books"`
}

func NewPostgresRepo(cfg PostgresConfig) (*PostgresRepo, error) {
	slog.Info("repo: starting NewPostgresRepo", slog.String("dsn", cfg.DSN))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("repo: creating pgx pool")
	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		slog.Error("repo: pgxpool.New failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}
	slog.Info("repo: pgx pool created successfully")

	slog.Info("repo: running migrations.Up")
	if err := migrations.Up(pool); err != nil {
		slog.Error("repo: migrations.Up failed", slog.String("error", err.Error()))
		pool.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}
	slog.Info("repo: migrations completed successfully")

	return &PostgresRepo{db: pool}, nil
}
