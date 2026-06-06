package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher provides bcrypt password hashing
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new bcrypt hasher with default cost
func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{
		cost: bcrypt.DefaultCost,
	}
}

// Hash hashes a password using bcrypt
func (h *BcryptHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Verify verifies a password against its hash
func (h *BcryptHasher) Verify(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
