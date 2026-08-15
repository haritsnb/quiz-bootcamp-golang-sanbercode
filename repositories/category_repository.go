package repositories

import (
	"app/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	rows, err := r.DB.Query("SELECT id, name, created_at, created_by, modified_at, modified_by FROM categories ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.CreatedBy, &c.ModifiedAt, &c.ModifiedBy); err != nil {
			return nil, err
		}
		list = append(list, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *CategoryRepository) FindByID(id int) (*models.Category, error) {
	var c models.Category
	query := "SELECT id, name, created_at, created_by, modified_at, modified_by FROM categories WHERE id = $1"
	err := r.DB.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.CreatedBy, &c.ModifiedAt, &c.ModifiedBy)
	if err == sql.ErrNoRows {
		return nil, errors.New("data category tidak ditemukan")
	}
	return &c, err
}

func (r *CategoryRepository) Create(c *models.Category) error {
	query := `INSERT INTO categories (name, created_by, modified_by) VALUES ($1, $2, $3) RETURNING id, created_at, modified_at`
	return r.DB.QueryRow(query, c.Name, c.CreatedBy, c.ModifiedBy).Scan(&c.ID, &c.CreatedAt, &c.ModifiedAt)
}

func (r *CategoryRepository) UpdatePartial(id int, fields map[string]interface{}) error {
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
	query := fmt.Sprintf("UPDATE categories SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIdx)
	args = append(args, id)

	res, err := r.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data category tidak ditemukan untuk diupdate")
	}
	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	res, err := r.DB.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data category tidak ditemukan untuk dihapus")
	}
	return nil
}
