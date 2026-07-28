package model

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResult struct {
	User         *User
	AccessToken  string
	RefreshToken string
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

type RefreshToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}
