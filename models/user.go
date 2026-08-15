package models

import "time"

type User struct {
	ID         int        `json:"id"`
	Username   string     `json:"username"`
	Password   string     `json:"password,omitempty"`
	CreatedAt  *time.Time `json:"created_at"`
	CreatedBy  *string    `json:"created_by"`
	ModifiedAt *time.Time `json:"modified_at"`
	ModifiedBy *string    `json:"modified_by"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserCreateRequest struct {
	Username   string  `json:"username" binding:"required"`
	Password   string  `json:"password" binding:"required,min=6"`
	CreatedBy  *string `json:"created_by"`
	ModifiedBy *string `json:"modified_by"`
}

type UserUpdateRequest struct {
	Username   *string `json:"username"`
	Password   *string `json:"password"`
	CreatedBy  *string `json:"created_by"`
	ModifiedBy *string `json:"modified_by"`
}
