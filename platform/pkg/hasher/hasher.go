package hasher

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrWeakPassword    = errors.New("password must be at least 8 characters")
	ErrEmptyHash       = errors.New("hash is empty")
)

type PasswordHasher interface {
	// Hash генерирует хеш из plaintext-пароля.
	// Может валидировать сложность пароля.
	Hash(plainText string) (string, error)

	// Verify проверяет, совпадает ли plaintext-пароль с хешем.
	Verify(hashedPassword, plainText string) error
}
