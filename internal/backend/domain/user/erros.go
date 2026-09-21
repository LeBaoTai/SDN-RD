package user

import "errors"

// Domain errors — these belong to the domain layer and carry no HTTP semantics.
var (
	// auth
	ErrInvalidUserID      = errors.New("user: invalid user id")
	ErrInvalidEmail       = errors.New("user: invalid email address")
	ErrUserNotFound       = errors.New("user: not found")
	ErrEmailAlreadyExists = errors.New("user: email already registered")
	ErrInvalidCredentials = errors.New("user: invalid credentials")
	ErrInvalidUsername    = errors.New("user: invalid username")
	ErrInvalidRole        = errors.New("user: invalid role")

	// user profile
	ErrUnauthorized = errors.New("user: unauthorized")
)
