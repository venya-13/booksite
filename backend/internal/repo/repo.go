package repo

import (
	"context"
	"fmt"
	"time"

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
	CreatedAt time.Time `json:"created_at"`
}

func NewPostgresRepo(cfg PostgresConfig) (*PostgresRepo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	return &PostgresRepo{db: pool}, nil
}
