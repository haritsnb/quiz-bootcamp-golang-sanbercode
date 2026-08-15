package repositories

import (
	"app/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type BookRepository struct {
	DB *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{DB: db}
}

func (r *BookRepository) FindAll() ([]models.Book, error) {
	rows, err := r.DB.Query("SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at, created_by, modified_at, modified_by FROM books ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Description, &b.ImageURL, &b.ReleaseYear, &b.Price, &b.TotalPage, &b.Thickness, &b.CategoryID, &b.CreatedAt, &b.CreatedBy, &b.ModifiedAt, &b.ModifiedBy); err != nil {
			return nil, err
		}
		list = append(list, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *BookRepository) FindByID(id int) (*models.Book, error) {
	var b models.Book
	query := "SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at, created_by, modified_at, modified_by FROM books WHERE id = $1"
	err := r.DB.QueryRow(query, id).Scan(&b.ID, &b.Title, &b.Description, &b.ImageURL, &b.ReleaseYear, &b.Price, &b.TotalPage, &b.Thickness, &b.CategoryID, &b.CreatedAt, &b.CreatedBy, &b.ModifiedAt, &b.ModifiedBy)
	if err == sql.ErrNoRows {
		return nil, errors.New("data book tidak ditemukan")
	}
	return &b, err
}

func (r *BookRepository) Create(b *models.Book) error {
	query := `INSERT INTO books (title, description, image_url, release_year, price, total_page, thickness, category_id, created_by, modified_by)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, created_at, modified_at`
	return r.DB.QueryRow(query, b.Title, b.Description, b.ImageURL, b.ReleaseYear, b.Price, b.TotalPage, b.Thickness, b.CategoryID, b.CreatedBy, b.ModifiedBy).Scan(&b.ID, &b.CreatedAt, &b.ModifiedAt)
}

func (r *BookRepository) UpdatePartial(id int, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	var setClauses []string
	var args []interface{}
	argIdx := 1

	for col, val := range fields {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, argIdx))
		args = append(args, val)
		argIdx++
	}

	setClauses = append(setClauses, "modified_at = CURRENT_TIMESTAMP")
	query := fmt.Sprintf("UPDATE books SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIdx)
	args = append(args, id)

	res, err := r.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data book tidak ditemukan untuk diupdate")
	}
	return nil
}

func (r *BookRepository) Delete(id int) error {
	res, err := r.DB.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data book tidak ditemukan untuk dihapus")
	}
	return nil
}
