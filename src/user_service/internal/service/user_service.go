package service

import (
	"context"
	"fmt"
	"strings"
	"user_service/internal/models"
	"user_service/internal/repository"
	"user_service/internal/usererrors"
	"user_service/pkg/hash"
)

type UserService struct {
	repo repository.UserRepository
	hash *hash.BcryptHasher
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
		hash: hash.NewBcryptHasher(),
	}
}

func (s *UserService) Register(ctx context.Context, name, email, passwrod string) (*models.User, error) {
	// Normalize input
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	// Hash  Password
	hashedPassword, err := s.hash.Hash(passwrod)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     models.RoleUser,
		Status:   models.StatusActive,
	}

	// save to database
	if err := s.repo.Create(ctx, user); err != nil {
		// check for duplicate email
		if strings.Contains(err.Error(), "already exists") {
			return nil, usererrors.ErrEmailExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// Login

func (s *UserService) Login(ctx context.Context, email, password string) (*models.User, error) {
	// normalize input
	email = strings.TrimSpace(email)

	// Find user by email
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, usererrors.ErrInvalidCredentials
	}
	// check if the user is active
	if user.Status != models.StatusActive {
		return nil, usererrors.ErrUserSuspended
	}
	// verify password
	if err := s.hash.Verify(user.Password, password); err != nil {
		return nil, usererrors.ErrInvalidCredentials
	}
	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, usererrors.ErrUserNotFound
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, usererrors.ErrUserNotFound
	}

	return user, nil
}
