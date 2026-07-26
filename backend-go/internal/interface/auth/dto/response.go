package dto

import "time"

type LoginResponseData struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginResponse struct {
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	Data       LoginResponseData `json:"data,omitempty"`
}

type RegisterResponseData struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterResponse struct {
	StatusCode int                  `json:"status_code"`
	Message    string               `json:"message"`
	Data       RegisterResponseData `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
