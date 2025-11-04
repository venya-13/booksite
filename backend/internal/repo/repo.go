package repo

import (
	"context"
	"fmt"
	"log/slog"
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

func (r *PostgresRepo) CreateBook(ctx context.Context, b Book, categoryIDs []int) (int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var id int
	err = tx.QueryRow(ctx, `
        INSERT INTO books (title, author, description, file_url, cover_path)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, b.Title, b.Author, b.Description, b.FileURL, b.CoverPath).Scan(&id)
	if err != nil {
		return 0, err
	}

	for _, cid := range categoryIDs {
		if _, err := tx.Exec(ctx, `INSERT INTO book_categories (book_id, category_id) VALUES ($1, $2)`, id, cid); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresRepo) GetAllBooks(ctx context.Context) ([]Book, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, author, description, file_url, cover_path, created_at, updated_at FROM books ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.FileURL, &b.CoverPath, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	return res, nil
}

func (r *PostgresRepo) GetBookByID(ctx context.Context, id int) (Book, error) {
	var b Book
	err := r.db.QueryRow(ctx, `SELECT id, title, author, description, file_url, cover_path, created_at, updated_at FROM books WHERE id=$1`, id).
		Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.FileURL, &b.CoverPath, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return Book{}, err
	}
	return b, nil
}

func (r *PostgresRepo) UpdateBook(ctx context.Context, b Book, categoryIDs []int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
        UPDATE books SET title=$1, author=$2, description=$3, file_url=$4, cover_path=$5, updated_at = now()
        WHERE id=$6
    `, b.Title, b.Author, b.Description, b.FileURL, b.CoverPath, b.ID)
	if err != nil {
		return err
	}

	// clear existing categories and add new ones
	if _, err := tx.Exec(ctx, `DELETE FROM book_categories WHERE book_id=$1`, b.ID); err != nil {
		return err
	}
	for _, cid := range categoryIDs {
		if _, err := tx.Exec(ctx, `INSERT INTO book_categories (book_id, category_id) VALUES ($1, $2)`, b.ID, cid); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepo) DeleteBook(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM books WHERE id=$1`, id)
	return err
}

func (r *PostgresRepo) GetBooksByCategory(ctx context.Context, categoryID int) ([]Book, error) {
	rows, err := r.db.Query(ctx, `
        SELECT b.id, b.title, b.author, b.description, b.file_url, b.cover_path, b.created_at, b.updated_at
        FROM books b
        JOIN book_categories bc ON bc.book_id = b.id
        WHERE bc.category_id = $1
        ORDER BY b.id
    `, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.FileURL, &b.CoverPath, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	return res, nil
}

// GetCategoriesWithBooks returns categories where each category has Books field (inline struct)
func (r *PostgresRepo) GetCategoriesWithBooks(ctx context.Context) ([]struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Books []Book `json:"books"`
}, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type CatWithBooks struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Books []Book `json:"books"`
	}

	var res []CatWithBooks
	for rows.Next() {
		var c CatWithBooks
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		// load books for category
		brows, err := r.db.Query(ctx, `
			SELECT b.id, b.title, b.author, b.description, b.file_url, b.cover_path, b.created_at, b.updated_at
			FROM books b
			JOIN book_categories bc ON bc.book_id = b.id
			WHERE bc.category_id = $1
			ORDER BY b.id
		`, c.ID)
		if err != nil {
			return nil, err
		}
		for brows.Next() {
			var b Book
			if err := brows.Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.FileURL, &b.CoverPath, &b.CreatedAt, &b.UpdatedAt); err != nil {
				brows.Close()
				return nil, err
			}
			c.Books = append(c.Books, b)
		}
		brows.Close()
		res = append(res, c)
	}
	// convert type to anonymous expected (or change signature upstream)
	// build return slice
	out := make([]struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Books []Book `json:"books"`
	}, len(res))
	for i := range res {
		out[i].ID = res[i].ID
		out[i].Name = res[i].Name
		out[i].Books = res[i].Books
	}
	return out, nil
}

// GetUncategorizedBooks returns books which have no entry in book_categories
func (r *PostgresRepo) GetUncategorizedBooks(ctx context.Context) ([]Book, error) {
	rows, err := r.db.Query(ctx, `
		SELECT b.id, b.title, b.author, b.description, b.file_url, b.cover_path, b.created_at, b.updated_at
		FROM books b
		LEFT JOIN book_categories bc ON bc.book_id = b.id
		WHERE bc.category_id IS NULL
		ORDER BY b.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.FileURL, &b.CoverPath, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

// AssignBookToCategory assigns a book to one category without removing existing ones
func (r *PostgresRepo) AssignBookToCategory(ctx context.Context, bookID int, categoryID int) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO book_categories (book_id, category_id)
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING
	`, bookID, categoryID)

	return err
}
