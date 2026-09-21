// Package hash contains the bcrypt adapter for the PasswordHasher port.
package security

import (
	"github.com/LeBaoTai/SDN-RD/internal/backend/domain/user"
	"golang.org/x/crypto/bcrypt"
)

const defaultCost = bcrypt.DefaultCost // 10

// BcryptHasher implements user.PasswordHasher using bcrypt.
type BcryptHasher struct {
	cost int
}

func NewBcryptHasher() *BcryptHasher                 { return &BcryptHasher{cost: defaultCost} }
func NewBcryptHasherWithCost(cost int) *BcryptHasher { return &BcryptHasher{cost: cost} }

func (h *BcryptHasher) Hash(plaintext string) (user.HashedPassword, error) {
	if len(plaintext) < 8 {
		return user.HashedPassword{}, errWeakPassword
	}
	raw, err := bcrypt.GenerateFromPassword([]byte(plaintext), h.cost)
	if err != nil {
		return user.HashedPassword{}, err
	}
	return user.HashedPasswordFrom(string(raw)), nil
}

func (h *BcryptHasher) Verify(plaintext string, hashed user.HashedPassword) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed.String()), []byte(plaintext))
	if err != nil {
		return user.ErrInvalidCredentials
	}
	return nil
}

// errWeakPassword is a package-level sentinel; it is not a domain error.
var errWeakPassword = &passwordError{"password must be at least 8 characters"}

type passwordError struct{ msg string }

func (e *passwordError) Error() string { return e.msg }
