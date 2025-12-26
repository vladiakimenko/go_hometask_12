package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	if err := ValidateEmail(r.Email); err != nil {
		return fmt.Errorf("email validation failed: %w", err)
	}
	if err := ValidateUsername(r.Username); err != nil {
		return fmt.Errorf("username validation failed: %w", err)
	}
	if err := ValidatePassword(r.Password); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return fmt.Errorf("email must not be empty")
	}
	if strings.TrimSpace(r.Password) == "" {
		return fmt.Errorf("password must not be empty")
	}
	return nil
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
