package user

import "context"

// PasswordHasher is a domain service port for hashing and verifying passwords.
// The concrete bcrypt adapter lives in infrastructure.
type PasswordHasher interface {
	Hash(plaintext string) (HashedPassword, error)
	Verify(plaintext string, hashed HashedPassword) error
}

// TokenIssuer is a domain service port for issuing authentication tokens.
type TokenIssuer interface {
	Issue(ctx context.Context, u *User) (string, error)
}
