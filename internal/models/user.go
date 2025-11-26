package models

import "time"

type User struct {
	ID                      string    `json:"id"`
	Email                   string    `json:"email"`
	Name                    string    `json:"name"`
	Password                string    `json:"-"` // Don't serialize password
	EmailVerified           bool      `json:"email_verified"`
	VerificationCode        string    `json:"-"`
	VerificationCodeExpires time.Time `json:"-"`
	ResetCode               string    `json:"-"`
	ResetCodeExpires        time.Time `json:"-"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type CreateUserDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name"  binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=6"` // Required for signup
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserDTO struct {
	Email *string `json:"email" binding:"omitempty,email"`
	Name  *string `json:"name"  binding:"omitempty,min=2,max=64"`
}

type VerifyEmailDTO struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required,len=6"`
	Password string `json:"password" binding:"required,min=6"`
}

type ResendVerificationDTO struct {
	Email string `json:"email" binding:"required,email"`
}
