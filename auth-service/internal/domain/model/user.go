package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID        string     `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"-" db:"password"`
	Salt      string     `json:"-" db:"salt"`
	FirstName string     `json:"first_name" db:"first_name"`
	LastName  string     `json:"last_name" db:"last_name"`
	BirthDate *time.Time     `json:"birth_date" db:"birth_date"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" db:"expires_at"`
}

type VerifyRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type VerifyResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message,omitempty"`
	User    *UserInfo `json:"user,omitempty"`
	Token   string    `json:"token,omitempty"` // JWT для gateway
}

type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  bool   `json:"is_active"`
}

type AuthClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	jwt.RegisteredClaims
}