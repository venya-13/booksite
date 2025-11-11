package repo

import "context"

func (r *PostgresRepo) CreateBook(ctx context.Context, b Book, catIDs []int) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `
		INSERT INTO books (title, author, description, file_url, cover_path)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, b.Title, b.Author, b.Description, b.FileURL, b.CoverPath).Scan(&id)
	if err != nil {
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
