package memory

import (
	"context"
	"sync"

	"lebaotai.com/backend-go/domain/user"
)

// InMemoryUserRepository is a thread-safe in-memory adapter.
// Swap for a Postgres/MySQL adapter without touching domain or application code.
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*user.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*user.User),
	}
}

func (r *InMemoryUserRepository) Save(ctx context.Context, u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[u.ID().String()] = u

	return nil
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id.String()]

	if !ok {
		return nil, user.ErrUserNotFound
	}

	return user.Reconstitute(
		u.ID(),
		u.Email(),
		u.HashedPassword(),
		u.Role(),
		u.CreatedAt(),
		u.UpdatedAt(),
	), nil
}

func (r *InMemoryUserRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email().String() == email.String() {
			return user.Reconstitute(
				u.ID(),
				u.Email(),
				u.HashedPassword(),
				u.Role(),
				u.CreatedAt(),
				u.UpdatedAt(),
			), nil
		}
	}

	return nil, user.ErrUserNotFound
}

func (r *InMemoryUserRepository) ExistsByEmail(ctx context.Context, email user.Email) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email().String() == email.String() {
			return true, user.ErrEmailAlreadyExists
		}
	}
	return false, nil
}
