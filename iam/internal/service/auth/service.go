package auth

import (
	"github.com/crafty-ezhik/rocket-factory/iam/internal/repository"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/hasher"
)

type service struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	hasher      hasher.PasswordHasher
}

func NewService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, hasher hasher.PasswordHasher) *service {
	return &service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		hasher:      hasher,
	}
}
