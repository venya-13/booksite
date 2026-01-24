package repo

import (
	"context"
	"log/slog"
)

func (r *PostgresRepo) GetAllCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.db.Query(ctx, `
	SELECT id, name, slug, is_system, created_at
	FROM categories
	ORDER BY id
`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.IsSystem, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *PostgresRepo) CreateCategory(ctx context.Context, name string) error {
	_, err := r.db.Exec(ctx, `
	INSERT INTO categories (name, slug, is_system, created_at)
	VALUES ($1, lower($1), FALSE, NOW())
`, name)
	return err
}

func (r *PostgresRepo) DeleteCategory(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `
	DELETE FROM categories
	WHERE id = $1 AND is_system = FALSE
`, id)
	return err
}

func (r *PostgresRepo) RenameCategory(ctx context.Context, id int, name string) error {
	_, err := r.db.Exec(ctx, `
	UPDATE categories
	SET name = $1
	WHERE id = $2 AND is_system = FALSE
`, name, id)

	if err != nil {
		slog.Error("Failed to rename category",
			slog.Int("id", id),
			slog.String("name", name),
			slog.String("error", err.Error()),
		)
	}
	return err
}
