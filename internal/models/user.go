package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserDTO struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name"  binding:"required,min=2,max=64"`
}

type UpdateUserDTO struct {
	Email *string `json:"email" binding:"omitempty,email"`
	Name  *string `json:"name"  binding:"omitempty,min=2,max=64"`
}
