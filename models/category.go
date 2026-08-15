package models

import "time"

type Category struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	CreatedAt  *time.Time `json:"created_at"`
	CreatedBy  *string    `json:"created_by"`
	ModifiedAt *time.Time `json:"modified_at"`
	ModifiedBy *string    `json:"modified_by"`
}

type CategoryCreateRequest struct {
	Name       string  `json:"name" binding:"required"`
	CreatedBy  *string `json:"created_by"`
	ModifiedBy *string `json:"modified_by"`
}

type CategoryUpdateRequest struct {
	Name       *string `json:"name"`
	CreatedBy  *string `json:"created_by"`
	ModifiedBy *string `json:"modified_by"`
}
