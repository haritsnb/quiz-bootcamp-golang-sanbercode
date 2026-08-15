package repositories

import (
	"app/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	rows, err := r.DB.Query("SELECT id, username, created_at, created_by, modified_at, modified_by FROM users ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.CreatedAt, &u.CreatedBy, &u.ModifiedAt, &u.ModifiedBy); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	var u models.User
	query := "SELECT id, username, created_at, created_by, modified_at, modified_by FROM users WHERE id = $1"
	err := r.DB.QueryRow(query, id).Scan(&u.ID, &u.Username, &u.CreatedAt, &u.CreatedBy, &u.ModifiedAt, &u.ModifiedBy)
	if err == sql.ErrNoRows {
		return nil, errors.New("data user tidak ditemukan")
	}
	return &u, err
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var u models.User
	query := "SELECT id, username, password, created_at, created_by, modified_at, modified_by FROM users WHERE username = $1"
	err := r.DB.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.CreatedBy, &u.ModifiedAt, &u.ModifiedBy)
	if err == sql.ErrNoRows {
		return nil, errors.New("user tidak ditemukan")
	}
	return &u, err
}

func (r *UserRepository) Create(u *models.User) error {
	query := `INSERT INTO users (username, password, created_by, modified_by) VALUES ($1, $2, $3, $4) RETURNING id, created_at, modified_at`
	return r.DB.QueryRow(query, u.Username, u.Password, u.CreatedBy, u.ModifiedBy).Scan(&u.ID, &u.CreatedAt, &u.ModifiedAt)
}

func (r *UserRepository) UpdatePartial(id int, fields map[string]interface{}) error {
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
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIdx)
	args = append(args, id)

	res, err := r.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data user tidak ditemukan untuk diupdate")
	}
	return nil
}

func (r *UserRepository) Delete(id int) error {
	res, err := r.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data user tidak ditemukan untuk dihapus")
	}
	return nil
}
