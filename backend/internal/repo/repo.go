package repo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

func (r *PostgresRepo) SaveOrUpdate(user map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	googleID, _ := user["id"].(string)
	email, _ := user["email"].(string)
	name, _ := user["name"].(string)
	picture, _ := user["picture"].(string)
	accessToken, _ := user["access_token"].(string)
	refreshToken, _ := user["refresh_token"].(string)
	tokenExpiry, _ := user["token_expiry"].(time.Time)

	slog.Info("Saving or updating user in database",
		slog.String("google_id", googleID),
		slog.String("email", email),
		slog.String("name", name),
	)

	_, err := r.db.Exec(ctx, `
		INSERT INTO users (google_id, email, name, picture, access_token, refresh_token, token_expiry)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (google_id) DO UPDATE
		SET email = EXCLUDED.email,
		    name = EXCLUDED.name,
		    picture = EXCLUDED.picture,
		    access_token = EXCLUDED.access_token,
		    refresh_token = EXCLUDED.refresh_token,
		    token_expiry = EXCLUDED.token_expiry,
		    updated_at = CURRENT_TIMESTAMP;
	`, googleID, email, name, picture, accessToken, refreshToken, tokenExpiry)

	if err != nil {
		slog.Error("DB SaveOrUpdate failed",
			slog.String("google_id", googleID),
			slog.String("error", err.Error()),
		)
		return err
	}

	slog.Info("User saved successfully", slog.String("google_id", googleID))
	return nil
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

func (r *PostgresRepo) GetUserByGoogleID(googleID string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.db.QueryRow(ctx, `
		SELECT google_id, email, name, picture, access_token, refresh_token, token_expiry, is_admin
		FROM users
		WHERE google_id=$1
	`, googleID)

	var (
		googleIDVal  string
		email        string
		name         string
		picture      string
		accessToken  string
		refreshToken string
		tokenExpiry  time.Time
		isAdmin      bool
	)

	err := row.Scan(
		&googleIDVal,
		&email,
		&name,
		&picture,
		&accessToken,
		&refreshToken,
		&tokenExpiry,
		&isAdmin,
	)
	if err != nil {
		return nil, fmt.Errorf("get user by google id: %w", err)
	}

	user := map[string]interface{}{
		"google_id":     googleIDVal,
		"email":         email,
		"name":          name,
		"picture":       picture,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_expiry":  tokenExpiry,
		"is_admin":      isAdmin,
	}

	return user, nil
}

func (r *PostgresRepo) GetAllCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, created_at FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *PostgresRepo) CreateCategory(ctx context.Context, name string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO categories (name, created_at) VALUES ($1, NOW())`, name)
	return err
}

func (r *PostgresRepo) DeleteCategory(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id)
	return err
}

func (r *PostgresRepo) RenameCategory(ctx context.Context, id int, name string) error {
	_, err := r.db.Exec(ctx, `UPDATE categories SET name=$1 WHERE id=$2`, name, id)
	if err != nil {
		slog.Error("Failed to rename category",
			slog.Int("id", id),
			slog.String("name", name),
			slog.String("error", err.Error()),
		)
	}
	return err
}
