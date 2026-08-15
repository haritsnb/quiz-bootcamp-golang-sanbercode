package models

import "time"

type Book struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	ImageURL    *string    `json:"image_url"`
	ReleaseYear int        `json:"release_year"`
	Price       int        `json:"price"`
	TotalPage   int        `json:"total_page"`
	Thickness   string     `json:"thickness"`
	CategoryID  int        `json:"category_id"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   *string    `json:"created_by"`
	ModifiedAt  *time.Time `json:"modified_at"`
	ModifiedBy  *string    `json:"modified_by"`
}

type BookCreateRequest struct {
	Title       string  `form:"title" binding:"required"`
	Description *string `form:"description"`
	ReleaseYear int     `form:"release_year" binding:"required,min=1980,max=2024"`
	Price       int     `form:"price" binding:"min=0"`
	TotalPage   int     `form:"total_page" binding:"required,min=1"`
	CategoryID  int     `form:"category_id" binding:"required"`
	CreatedBy   *string `form:"created_by"`
	ModifiedBy  *string `form:"modified_by"`
}

type BookUpdateRequest struct {
	Title       *string `form:"title"`
	Description *string `form:"description"`
	ReleaseYear *int    `form:"release_year"`
	Price       *int    `form:"price"`
	TotalPage   *int    `form:"total_page"`
	CategoryID  *int    `form:"category_id"`
	CreatedBy   *string `form:"created_by"`
	ModifiedBy  *string `form:"modified_by"`
}
