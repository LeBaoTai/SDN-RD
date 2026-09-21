package auth

import (
	"context"
	"time"

	"github.com/LeBaoTai/SDN-RD/internal/backend/domain/user"
)

// LoginUserCommand carries the raw credentials for the login use case.
type LoginUserCommand struct {
	Email    string
	Password string
}

// LoginUserResult is returned on a successful login.
type LoginUserResult struct {
	Token     string
	UserID    string
	Username  string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// LoginUserHandler orchestrates user authentication.
type LoginUserHandler struct {
	repo        user.Repository
	hasher      user.PasswordHasher
	tokenIssuer user.TokenIssuer
}

func NewLoginUserHandler(
	repo user.Repository,
	hasher user.PasswordHasher,
	tokenIssuer user.TokenIssuer,
) *LoginUserHandler {
	return &LoginUserHandler{repo: repo, hasher: hasher, tokenIssuer: tokenIssuer}
}

// Handle executes the login use case.
//
//  1. Build the Email value object
//  2. Fetch the User aggregate by email
//  3. Verify password against stored hash
//  4. Issue a JWT token
func (h *LoginUserHandler) Handle(ctx context.Context, cmd LoginUserCommand) (*LoginUserResult, error) {
	// 1. Validate and build email value object.
	email, err := user.NewEmail(cmd.Email)
	if err != nil {
		// Return opaque error so we don't reveal "email format wrong vs not found".
		return nil, user.ErrInvalidCredentials
	}

	// 2. Fetch aggregate — translate "not found" to opaque credentials error.
	u, err := h.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, user.ErrInvalidCredentials
	}

	// 3. Verify password.
	if err := h.hasher.Verify(cmd.Password, u.HashedPassword()); err != nil {
		return nil, user.ErrInvalidCredentials
	}

	// 4. Issue token.
	token, err := h.tokenIssuer.Issue(ctx, u)
	if err != nil {
		return nil, err
	}

	return &LoginUserResult{
		Token:     token,
		UserID:    u.ID().String(),
		Username:  u.Username().String(),
		Email:     u.Email().String(),
		Role:      u.Role().String(),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}, nil
}
