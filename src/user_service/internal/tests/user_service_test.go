package user_test

import (
	"context"
	"testing"
	"user_service/internal/models"
	"user_service/internal/service"
	"user_service/internal/usererrors"

	"github.com/stretchr/testify/assert"
)

// MockUserRepository is a mock implementation for testing
type MockUserRepository struct {
	users map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.Email]; exists {
		return usererrors.ErrEmailExists
	}
	user.ID = len(m.users) + 1
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, usererrors.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, usererrors.ErrUserNotFound
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	m.users[user.Email] = user
	return nil
}

// Tests
func TestRegisterSuccess(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := service.NewUserService(mockRepo)

	user, err := service.Register(context.Background(), "John Doe", "john@example.com", "SecurePass123!")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, models.RoleUser, user.Role)
}

func TestRegisterDuplicateEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := service.NewUserService(mockRepo)

	// First registration
	service.Register(context.Background(), "John Doe", "john@example.com", "SecurePass123!")

	// Try to register with same email
	_, err := service.Register(context.Background(), "Jane Doe", "john@example.com", "SecurePass456!")

	assert.Error(t, err)
	assert.True(t, assert.ErrorIs(t, err, usererrors.ErrEmailExists) ||
		assert.ErrorContains(t, err, "already exists"))
}

func TestLoginSuccess(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := service.NewUserService(mockRepo)

	// Register first
	service.Register(context.Background(), "John Doe", "john@example.com", "SecurePass123!")

	// Login
	user, err := service.Login(context.Background(), "john@example.com", "SecurePass123!")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestLoginInvalidPassword(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := service.NewUserService(mockRepo)

	// Register first
	service.Register(context.Background(), "John Doe", "john@example.com", "SecurePass123!")

	// Try login with wrong password
	_, err := service.Login(context.Background(), "john@example.com", "WrongPassword")

	assert.Error(t, err)
	assert.Equal(t, usererrors.ErrInvalidCredentials, err)
}

func TestLoginInvalidEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := service.NewUserService(mockRepo)

	_, err := service.Login(context.Background(), "nonexistent@example.com", "AnyPassword123!")

	assert.Error(t, err)
	assert.Equal(t, usererrors.ErrInvalidCredentials, err)
}
