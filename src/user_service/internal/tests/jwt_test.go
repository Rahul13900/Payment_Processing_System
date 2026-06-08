package user_test

import (
	"testing"
	"time"
	"user_service/pkg/jwt"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateToken(t *testing.T) {
	manager := jwt.NewJWTManager("test-secret-key", 1*time.Hour)

	// Generate token
	token, err := manager.GenerateToken(1, "test@example.com", "USER")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate token
	claims, err := manager.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, 1, claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "USER", claims.Role)
}

func TestExpiredToken(t *testing.T) {
	// Create manager with very short expiry
	manager := jwt.NewJWTManager("test-secret-key", 1*time.Millisecond)

	// Generate token
	token, err := manager.GenerateToken(1, "test@example.com", "USER")
	assert.NoError(t, err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Try to validate - should fail
	_, err = manager.ValidateToken(token)
	assert.Error(t, err)
}

func TestInvalidToken(t *testing.T) {
	manager := jwt.NewJWTManager("test-secret-key", 1*time.Hour)

	// Try to validate invalid token
	_, err := manager.ValidateToken("invalid.token.string")
	assert.Error(t, err)
}

func TestInvalidSignature(t *testing.T) {
	manager1 := jwt.NewJWTManager("secret-key-1", 1*time.Hour)
	manager2 := jwt.NewJWTManager("secret-key-2", 1*time.Hour)

	// Generate with one key
	token, err := manager1.GenerateToken(1, "test@example.com", "USER")
	assert.NoError(t, err)

	// Try to validate with different key - should fail
	_, err = manager2.ValidateToken(token)
	assert.Error(t, err)
}
