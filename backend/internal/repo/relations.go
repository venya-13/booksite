package repo

import "context"

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
func (r *PostgresRepo) GetCategoriesWithBooks(ctx context.Context) ([]CategoryWithBooks, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, slug, is_system
FROM categories
ORDER BY id
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []CategoryWithBooks
	for rows.Next() {
		var c CategoryWithBooks
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.IsSystem); err != nil {
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
	return res, nil
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
