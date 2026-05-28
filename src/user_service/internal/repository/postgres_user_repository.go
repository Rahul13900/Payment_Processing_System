package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"user_service/internal/models"
)

type postgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *postgresUserRepository {
	return &postgresUserRepository{
		db: db,
	}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (name, email, password, role, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	// Using $1, $2 for parameterized queries (prevents SQL injection

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.Status,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return fmt.Errorf("email already exists")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *postgresUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password, role, status, created_at, updated_at
		FROM users
		WHERE email = $1 AND status != 'DELETED'
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, name, email, password, role, status, created_at, updated_at
		FROM users
		WHERE id = $1 AND status != 'DELETED'
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET name = $1, role = $2, status = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Role,
		user.Status,
		user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
