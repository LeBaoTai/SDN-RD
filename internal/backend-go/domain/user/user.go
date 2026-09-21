package user

import (
	"time"

	"github.com/google/uuid"
)

// UserID is a strongly-typed value object for user identity.
type UserID struct {
	value uuid.UUID
}

func NewUserID() UserID {
	return UserID{value: uuid.New()}
}

func UserIDFrom(raw string) (UserID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return UserID{}, ErrInvalidUserID
	}
	return UserID{value: id}, nil
}

func (id UserID) String() string { return id.value.String() }

// Email ---
type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	if !isValidEmail(raw) {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: raw}, nil
}

func (e Email) String() string { return e.value }

// HashedPassword ---
type HashedPassword struct {
	value string
}

func HashedPasswordFrom(hash string) HashedPassword {
	return HashedPassword{value: hash}
}

func (p HashedPassword) String() string { return p.value }

// Username ---
type Username struct {
	value string
}

func NewUsername(raw string) (Username, error) {
	if raw == "" {
		return Username{}, ErrInvalidUsername
	}
	return Username{value: raw}, nil
}

func (u Username) String() string { return u.value }

// Role ---
type Role struct {
	value string
}

func NewRole(raw string) (Role, error) {
	if raw == "" {
		return Role{}, ErrInvalidRole
	}
	return Role{value: raw}, nil
}

func (r Role) String() string { return r.value }

// User ---
type User struct {
	id             UserID
	email          Email
	username       Username
	hashedPassword HashedPassword
	role           Role
	createdAt      time.Time
	updatedAt      time.Time
}

// NewUser creates a new User aggregate. Password must already be hashed
// by the application layer before calling this constructor.
func NewUser(email Email, hashedPassword HashedPassword, username Username, role Role) *User {
	now := time.Now().UTC()
	return &User{
		id:             NewUserID(),
		email:          email,
		username:       username,
		role:           role,
		hashedPassword: hashedPassword,
		createdAt:      now,
		updatedAt:      now,
	}
}

// Reconstitute rebuilds a User from persistent storage (no domain events raised).
func Reconstitute(id UserID, email Email, hashedPassword HashedPassword, role Role, createdAt, updatedAt time.Time) *User {
	return &User{
		id:             id,
		email:          email,
		hashedPassword: hashedPassword,
		role:           role,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// --- Accessors (read-only outside aggregate) ---

func (u *User) ID() UserID                     { return u.id }
func (u *User) Email() Email                   { return u.email }
func (u *User) HashedPassword() HashedPassword { return u.hashedPassword }
func (u *User) CreatedAt() time.Time           { return u.createdAt }
func (u *User) UpdatedAt() time.Time           { return u.updatedAt }
func (u *User) Username() Username             { return u.username }
func (u *User) Role() Role                     { return u.role }

// --- Domain behaviour ---

// ChangePassword replaces the hashed password and bumps the timestamp.
func (u *User) ChangePassword(newHash HashedPassword) {
	u.hashedPassword = newHash
	u.updatedAt = time.Now().UTC()
}
