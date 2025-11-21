package bcrypt

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/crafty-ezhik/rocket-factory/platform/pkg/hasher"
)

type BcryptPasswordHasher struct {
	cost int
}

func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptPasswordHasher{cost: cost}
}

func (h *BcryptPasswordHasher) Hash(plainText string) (string, error) {
	if err := h.validatePassword(plainText); err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), h.cost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}
	return string(hash), nil
}

func (h *BcryptPasswordHasher) Verify(hashedPassword, plainText string) error {
	if hashedPassword == "" {
		return hasher.ErrEmptyHash
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainText)); err != nil {
		return hasher.ErrInvalidPassword
	}
	return nil
}

func (h *BcryptPasswordHasher) validatePassword(password string) error {
	if len(password) < 8 {
		return hasher.ErrWeakPassword
	}
	return nil
}
