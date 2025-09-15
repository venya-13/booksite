package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

type PostgresConfig struct {
	DSN string
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

func (r *PostgresRepo) SaveOrUpdate(user map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	googleID := user["id"].(string)
	email := user["email"].(string)
	name := user["name"].(string)
	picture := user["picture"].(string)

	_, err := r.db.Exec(ctx, `
	INSERT INTO users (google_id, email, name, picture, refresh_token)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (google_id) DO UPDATE
	SET email = EXCLUDED.email,
	    name = EXCLUDED.name,
	    picture = EXCLUDED.picture,
	    refresh_token = EXCLUDED.refresh_token,
	    updated_at = CURRENT_TIMESTAMP;
	`, googleID, email, name, picture, user["refresh_token"])

	return err
}

func (r *PostgresRepo) GetRefreshTokenByGoogleID(googleID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var refreshToken string
	err := r.db.QueryRow(ctx, `
		SELECT refresh_token FROM users WHERE google_id=$1
	`, googleID).Scan(&refreshToken)

	if err != nil {
		return "", fmt.Errorf("get refresh token: %w", err)
	}
	return refreshToken, nil
}
