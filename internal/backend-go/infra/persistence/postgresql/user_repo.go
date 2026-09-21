package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"lebaotai.com/backend-go/domain/user"
)

// PostgresUserRepository is a thread-safe in-memory adapter.
// Swap for a Postgres/MySQL adapter without touching domain or application code.
type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		pool: pool,
	}
}

func (r *PostgresUserRepository) Save(ctx context.Context, u *user.User) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure the transaction is rolled back if not committed.
	// if the transaction is already committed, Rollback will return an error which we can ignore.
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := "INSERT INTO users (id, email, username, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = tx.Exec(ctx, query, u.ID().String(), u.Email().String(), u.Username().String(), u.HashedPassword(), u.CreatedAt(), u.UpdatedAt())
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return tx.Commit(ctx)
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure the transaction is rolled back if not committed.
	// if the transaction is already committed, Rollback will return an error which we can ignore.
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := "SELECT id, email, username, password_hash, created_at, updated_at FROM users WHERE id = $1"
	row := tx.QueryRow(ctx, query, id.String())

	var u *user.User
	err = row.Scan(u.ID(), u.Email(), u.Username(), u.HashedPassword(), u.CreatedAt(), u.UpdatedAt())
	if err != nil {
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	return u, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure the transaction is rolled back if not committed.
	// if the transaction is already committed, Rollback will return an error which we can ignore.
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := "SELECT id, email, username, password_hash, created_at, updated_at FROM users WHERE email = $1"
	row := tx.QueryRow(ctx, query, email.String())

	var u *user.User
	err = row.Scan(u.ID(), u.Email(), u.Username(), u.HashedPassword(), u.CreatedAt(), u.UpdatedAt())
	if err != nil {
		return nil, fmt.Errorf("failed to find user by Email: %w", err)
	}

	return u, nil
}

func (r *PostgresUserRepository) ExistsByEmail(ctx context.Context, email user.Email) (bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure the transaction is rolled back if not committed.
	// if the transaction is already committed, Rollback will return an error which we can ignore.
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	var exists bool
	err = tx.QueryRow(ctx, query, email.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists by email: %w", err)
	}
	return exists, nil
}
