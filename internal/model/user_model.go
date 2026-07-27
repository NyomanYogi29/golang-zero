package model

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRegisterRequestSchema struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequestSchema struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponseSchema struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}
