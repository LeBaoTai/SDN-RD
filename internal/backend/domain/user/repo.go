package user

import "context"

// Repository is the domain port for User persistence.
// Concrete adapters live in the infrastructure layer.
type Repository interface {
	// Save persists a new or updated User aggregate.
	Save(ctx context.Context, u *User) error

	// FindByID retrieves a User by its identity.
	// Returns ErrUserNotFound when no record exists.
	FindByID(ctx context.Context, id UserID) (*User, error)

	// FindByEmail retrieves a User by email address.
	// Returns ErrUserNotFound when no record exists.
	FindByEmail(ctx context.Context, email Email) (*User, error)

	// ExistsByEmail returns true if a user with the given email already exists.
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
}
