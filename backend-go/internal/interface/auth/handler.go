// Package handler wires HTTP concerns to application use-case handlers.
package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"lebaotai.com/backend-go/internal/app/auth"
	"lebaotai.com/backend-go/internal/domain/user"
	"lebaotai.com/backend-go/internal/interface/auth/dto"
)

// AuthHandler exposes register and login endpoints.
type AuthHandler struct {
	// register *auth.RegisterUserHandler
	login *auth.LoginUserHandler
}

func NewAuthHandler(
	// register *auth.RegisterUserHandler,
	login *auth.LoginUserHandler,
) *AuthHandler {
	return &AuthHandler{login: login}
}

// // Register handles POST /auth/register
// func (h *AuthHandler) Register(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	var req dto.RegisterRequest
//
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		writeError(c, http.StatusBadRequest, err.Error())
// 		return
// 	}
//
// 	result, err := h.register.Handle(ctx, auth.RegisterUserCommand{
// 		Email:    req.Email,
// 		Password: req.Password,
// 		Username: req.Username,
// 		Role:     req.Role,
// 	})
// 	if err != nil {
// 		writeError(c, statusFor(err), err.Error())
// 		return
// 	}
//
// 	writeJSON(c, http.StatusCreated, dto.RegisterResponse{
// 		StatusCode: http.StatusCreated,
// 		Message:    "user registered successfully",
// 		Data: dto.RegisterResponseData{
// 			UserID:    result.UserID,
// 			Username:  result.Username,
// 			Email:     result.Email,
// 			Role:      result.Role,
// 			CreatedAt: result.CreatedAt,
// 			UpdatedAt: result.UpdatedAt,
// 		},
// 	})
// }

// Login handles POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.login.Handle(ctx, auth.LoginUserCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(c, statusFor(err), err.Error())
		return
	}

	writeJSON(c, http.StatusOK, dto.LoginResponse{
		StatusCode: http.StatusOK,
		Message:    "login successful",
		Data: dto.LoginResponseData{
			Token:     result.Token,
			Email:     result.Email,
			UserID:    result.UserID,
			Username:  result.Username,
			Role:      result.Role,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		},
	})
}

// --- helpers ---

func writeJSON(c *gin.Context, code int, v any) {
	c.JSON(code, v)
}

func writeError(c *gin.Context, code int, msg string) {
	writeJSON(c, code, dto.ErrorResponse{Error: msg})
}

// statusFor maps domain errors to HTTP status codes.
// HTTP semantics stay in the interface layer — the domain knows nothing of HTTP.
func statusFor(err error) int {
	switch {
	case errors.Is(err, user.ErrEmailAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, user.ErrInvalidEmail):
		return http.StatusUnprocessableEntity
	case errors.Is(err, user.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, user.ErrUserNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
